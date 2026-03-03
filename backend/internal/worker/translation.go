// File: translation.go
// Purpose: Execute queued translation jobs asynchronously and persist job state transitions.
// Module: backend/internal/worker, background worker layer.
// Related: service translation/integration modules and integration provider configs.
package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/service"
)

type TranslationWorker struct {
	logger             *slog.Logger
	translationService *service.TranslationService
	integrationService *service.IntegrationService
	pollInterval       time.Duration
	backoffBase        time.Duration
	backoffMax         time.Duration
	httpClient         *http.Client
}

func NewTranslationWorker(
	logger *slog.Logger,
	translationService *service.TranslationService,
	integrationService *service.IntegrationService,
	pollInterval time.Duration,
	backoffBase time.Duration,
	backoffMax time.Duration,
) *TranslationWorker {
	if pollInterval <= 0 {
		pollInterval = 3 * time.Second
	}
	if backoffBase <= 0 {
		backoffBase = 3 * time.Second
	}
	if backoffMax <= 0 {
		backoffMax = 60 * time.Second
	}
	if backoffMax < backoffBase {
		backoffMax = backoffBase
	}
	return &TranslationWorker{
		logger:             logger,
		translationService: translationService,
		integrationService: integrationService,
		pollInterval:       pollInterval,
		backoffBase:        backoffBase,
		backoffMax:         backoffMax,
		httpClient:         &http.Client{Timeout: 45 * time.Second},
	}
}

func (w *TranslationWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	w.logger.Info("translation worker started", "poll_interval", w.pollInterval.String())
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("translation worker stopped")
			return
		case <-ticker.C:
			if err := w.processOnce(ctx); err != nil {
				w.logger.Error("translation worker tick failed", "error", err)
			}
		}
	}
}

func (w *TranslationWorker) processOnce(ctx context.Context) error {
	job, ok, err := w.translationService.ClaimNextQueuedTranslationJob()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	sourceText, sourceOK, err := w.translationService.GetTranslationSourceText(job.SourceType, job.SourceID)
	if err != nil {
		w.handleFailure(job, fmt.Sprintf("load source failed: %v", err))
		return nil
	}
	if !sourceOK || strings.TrimSpace(sourceText) == "" {
		w.handleFailure(job, "source not found or empty")
		return nil
	}

	provider, providerOK := w.integrationService.GetIntegrationProviderForRuntime(job.ProviderKey)
	if !providerOK {
		w.handleFailure(job, "provider not found")
		return nil
	}
	if !provider.Enabled {
		w.handleFailure(job, "provider is disabled")
		return nil
	}

	status, publishAt := w.resolvePublishPolicy(job)
	translatedTitle := ""
	translatedSummary := ""
	translatedContent := ""
	if job.SourceType == "article" {
		var err error
		translatedTitle, translatedSummary, translatedContent, err = w.translateArticleWithProvider(ctx, provider.ConfigJSON, job.ModelName, job.SourceLocale, job.TargetLocale, sourceText)
		if err != nil {
			w.handleFailure(job, err.Error())
			return nil
		}
	} else {
		translated, err := w.translateWithProvider(ctx, provider.ConfigJSON, job.ModelName, job.SourceLocale, job.TargetLocale, sourceText)
		if err != nil {
			w.handleFailure(job, err.Error())
			return nil
		}
		translatedContent = translated
	}
	if strings.TrimSpace(translatedContent) == "" {
		w.handleFailure(job, "empty translation output")
		return nil
	}

	if err := w.translationService.UpsertTranslationResult(job.SourceType, job.SourceID, job.TargetLocale, translatedTitle, translatedSummary, translatedContent, status, publishAt, job.ID); err != nil {
		w.handleFailure(job, fmt.Sprintf("persist translation failed: %v", err))
		return nil
	}

	if err := w.translationService.MarkTranslationJobSucceeded(job.ID, translatedContent); err != nil {
		return err
	}
	w.logger.Info("translation job succeeded", "job_id", job.ID, "provider_key", job.ProviderKey)
	return nil
}

func (w *TranslationWorker) resolvePublishPolicy(job domain.TranslationJob) (string, time.Time) {
	if !job.AutoPublish {
		return "draft", time.Time{}
	}
	if !job.PublishAt.IsZero() {
		return "published", job.PublishAt.UTC()
	}
	return "published", time.Now().UTC()
}

func (w *TranslationWorker) handleFailure(job domain.TranslationJob, errorMessage string) {
	if job.RetryCount < job.MaxRetries {
		delay := w.computeBackoff(job.RetryCount)
		nextRetryAt := time.Now().UTC().Add(delay)
		if err := w.translationService.ScheduleTranslationJobRetry(job.ID, errorMessage, nextRetryAt); err != nil {
			w.logger.Error("schedule translation retry failed", "job_id", job.ID, "error", err)
			_ = w.translationService.MarkTranslationJobFailed(job.ID, errorMessage)
			return
		}
		w.logger.Warn("translation job scheduled for retry",
			"job_id", job.ID,
			"retry_count", job.RetryCount+1,
			"max_retries", job.MaxRetries,
			"next_retry_at", nextRetryAt.Format(time.RFC3339),
		)
		return
	}
	_ = w.translationService.MarkTranslationJobFailed(job.ID, errorMessage)
	w.logger.Warn("translation job marked failed", "job_id", job.ID, "retry_count", job.RetryCount, "max_retries", job.MaxRetries)
}

func (w *TranslationWorker) computeBackoff(retryCount int) time.Duration {
	if retryCount < 0 {
		retryCount = 0
	}
	delay := w.backoffBase
	for i := 0; i < retryCount; i++ {
		delay *= 2
		if delay >= w.backoffMax {
			return w.backoffMax
		}
	}
	if delay > w.backoffMax {
		return w.backoffMax
	}
	return delay
}

func (w *TranslationWorker) translateWithProvider(ctx context.Context, configJSON []byte, modelName, sourceLocale, targetLocale, sourceText string) (string, error) {
	var cfg struct {
		BaseURL string `json:"base_url"`
		APIKey  string `json:"api_key"`
		Model   string `json:"model"`
	}
	if err := json.Unmarshal(configJSON, &cfg); err != nil {
		return "", fmt.Errorf("provider config parse failed: %w", err)
	}
	if strings.TrimSpace(cfg.BaseURL) == "" || strings.TrimSpace(cfg.APIKey) == "" {
		return "", fmt.Errorf("provider config missing base_url or api_key")
	}
	if strings.TrimSpace(modelName) == "" {
		modelName = cfg.Model
	}
	if strings.TrimSpace(modelName) == "" {
		return "", fmt.Errorf("model is empty")
	}
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	endpoint := baseURL + "/chat/completions"

	systemPrompt := "You are a professional translator. Preserve markdown structure and meaning. Return only translated text."
	userPrompt := fmt.Sprintf("Translate the following text from %s to %s:\n\n%s", sourceLocale, targetLocale, sourceText)

	payload := map[string]any{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.2,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm request failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("llm response status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", fmt.Errorf("llm response decode failed: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("llm response has no choices")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

func (w *TranslationWorker) translateArticleWithProvider(ctx context.Context, configJSON []byte, modelName, sourceLocale, targetLocale, sourceText string) (string, string, string, error) {
	var cfg struct {
		BaseURL string `json:"base_url"`
		APIKey  string `json:"api_key"`
		Model   string `json:"model"`
	}
	if err := json.Unmarshal(configJSON, &cfg); err != nil {
		return "", "", "", fmt.Errorf("provider config parse failed: %w", err)
	}
	if strings.TrimSpace(cfg.BaseURL) == "" || strings.TrimSpace(cfg.APIKey) == "" {
		return "", "", "", fmt.Errorf("provider config missing base_url or api_key")
	}
	if strings.TrimSpace(modelName) == "" {
		modelName = cfg.Model
	}
	if strings.TrimSpace(modelName) == "" {
		return "", "", "", fmt.Errorf("model is empty")
	}
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	endpoint := baseURL + "/chat/completions"

	systemPrompt := "You are a professional translator. Preserve meaning and markdown structure."
	userPrompt := fmt.Sprintf(`Translate the following article source from %s to %s.
Return STRICT JSON only with fields: title, summary, content.
Do not add markdown code fences.

Source:
%s`, sourceLocale, targetLocale, sourceText)

	payload := map[string]any{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.2,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("llm request failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", "", fmt.Errorf("llm response status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", "", "", fmt.Errorf("llm response decode failed: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", "", "", fmt.Errorf("llm response has no choices")
	}
	raw := strings.TrimSpace(out.Choices[0].Message.Content)
	if raw == "" {
		return "", "", "", fmt.Errorf("empty translation output")
	}
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start >= 0 && end > start {
		raw = raw[start : end+1]
	}
	var parsed struct {
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return "", "", "", fmt.Errorf("article translation json decode failed: %w", err)
	}
	return strings.TrimSpace(parsed.Title), strings.TrimSpace(parsed.Summary), strings.TrimSpace(parsed.Content), nil
}
