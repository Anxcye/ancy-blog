// File: ai_assist.go
// Purpose: Provide AI-assisted utility capabilities for admin authoring workflows.
// Module: backend/internal/service, AI assist service layer.
// Related: article/integration services and admin AI endpoints.
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type AIAssistService struct {
	articleService     *ArticleService
	integrationService *IntegrationService
	httpClient         *http.Client
}

func NewAIAssistService(articleService *ArticleService, integrationService *IntegrationService) *AIAssistService {
	return &AIAssistService{
		articleService:     articleService,
		integrationService: integrationService,
		httpClient:         &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *AIAssistService) GenerateSummary(ctx context.Context, title, content, providerKey, modelName string, maxLength int) (string, bool, string, error) {
	if maxLength <= 0 {
		maxLength = 180
	}
	if providerKey == "" {
		providerKey = "openai_compatible"
	}

	provider, ok := s.integrationService.GetIntegrationProviderForRuntime(providerKey)
	if !ok {
		return fallbackSummary(title, content, maxLength), true, fmt.Sprintf("provider %q not found", providerKey), nil
	}
	if !provider.Enabled {
		return fallbackSummary(title, content, maxLength), true, fmt.Sprintf("provider %q is disabled", providerKey), nil
	}

	prompt := fmt.Sprintf("Generate a concise blog summary in <= %d characters. Return plain text only.\n\nTitle: %s\n\nContent:\n%s", maxLength, title, content)
	result, err := callOpenAICompatible(ctx, s.httpClient, provider.ConfigJSON, modelName, prompt)
	if err != nil {
		return fallbackSummary(title, content, maxLength), true, err.Error(), nil
	}
	clean := strings.TrimSpace(result)
	if clean == "" {
		return fallbackSummary(title, content, maxLength), true, "LLM returned empty response", nil
	}
	if len([]rune(clean)) > maxLength {
		runes := []rune(clean)
		clean = string(runes[:maxLength])
	}
	return clean, false, "", nil
}

func (s *AIAssistService) SuggestSlug(ctx context.Context, title, providerKey, modelName string) (string, bool, string, error) {
	if providerKey == "" {
		providerKey = "openai_compatible"
	}

	baseSlug := slugify(title)
	fallback := s.ensureUniqueSlug(baseSlug)

	provider, ok := s.integrationService.GetIntegrationProviderForRuntime(providerKey)
	if !ok {
		return fallback, true, fmt.Sprintf("provider %q not found", providerKey), nil
	}
	if !provider.Enabled {
		return fallback, true, fmt.Sprintf("provider %q is disabled", providerKey), nil
	}

	prompt := fmt.Sprintf("Generate a short URL slug for this title. Use lowercase letters, numbers, and hyphens only. Return only slug text.\n\nTitle: %s", title)
	result, err := callOpenAICompatible(ctx, s.httpClient, provider.ConfigJSON, modelName, prompt)
	if err != nil {
		return fallback, true, err.Error(), nil
	}
	clean := slugify(result)
	if clean == "" {
		return fallback, true, "LLM returned empty or invalid slug", nil
	}
	return s.ensureUniqueSlug(clean), false, "", nil
}

func (s *AIAssistService) ensureUniqueSlug(base string) string {
	if strings.TrimSpace(base) == "" {
		base = fmt.Sprintf("post-%d", time.Now().Unix())
	}
	slug := base
	idx := 2
	for s.articleService.SlugExists(slug) {
		slug = fmt.Sprintf("%s-%d", base, idx)
		idx++
	}
	return slug
}

func fallbackSummary(title, content string, maxLength int) string {
	base := strings.TrimSpace(content)
	if base == "" {
		base = strings.TrimSpace(title)
	}
	if base == "" {
		return ""
	}
	runes := []rune(base)
	if len(runes) <= maxLength {
		return base
	}
	return string(runes[:maxLength])
}

var slugCleanRe = regexp.MustCompile(`[^a-z0-9-]+`)
var slugDashRe = regexp.MustCompile(`-+`)

func slugify(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	s = slugCleanRe.ReplaceAllString(s, "-")
	s = slugDashRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return ""
	}
	if len(s) > 96 {
		s = s[:96]
		s = strings.Trim(s, "-")
	}
	return s
}

func callOpenAICompatible(ctx context.Context, httpClient *http.Client, configJSON []byte, modelName, userPrompt string) (string, error) {
	var cfg struct {
		BaseURL string `json:"base_url"`
		APIKey  string `json:"api_key"`
		Model   string `json:"model"`
	}
	if err := json.Unmarshal(configJSON, &cfg); err != nil {
		return "", err
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

	endpoint := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	payload := map[string]any{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a careful writing assistant. Return only final text."},
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

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
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
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("llm response has no choices")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}
