// File: content_test.go
// Purpose: Verify integration-center and translation-job business rules in content service.
// Module: backend/internal/service, content unit test layer.
// Related: content.go and repository contracts.
package service

import (
	"encoding/json"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
)

type contentRepoStub struct {
	repository.ContentRepository

	getIntegrationProviderFunc    func(providerKey string) (domain.IntegrationProvider, bool)
	updateIntegrationProviderFunc func(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error)
	createTranslationJobFunc      func(job domain.TranslationJob) (domain.TranslationJob, error)
}

func (s *contentRepoStub) GetIntegrationProvider(providerKey string) (domain.IntegrationProvider, bool) {
	if s.getIntegrationProviderFunc != nil {
		return s.getIntegrationProviderFunc(providerKey)
	}
	return domain.IntegrationProvider{}, false
}

func (s *contentRepoStub) UpdateIntegrationProvider(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
	if s.updateIntegrationProviderFunc != nil {
		return s.updateIntegrationProviderFunc(providerKey, enabled, configJSON, metaJSON)
	}
	return domain.IntegrationProvider{}, nil
}

func (s *contentRepoStub) CreateTranslationJob(job domain.TranslationJob) (domain.TranslationJob, error) {
	if s.createTranslationJobFunc != nil {
		return s.createTranslationJobFunc(job)
	}
	return job, nil
}

func TestUpdateIntegrationProviderMasksSecrets(t *testing.T) {
	repo := &contentRepoStub{
		updateIntegrationProviderFunc: func(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
			return domain.IntegrationProvider{
				ProviderKey: providerKey,
				Enabled:     enabled,
				ConfigJSON:  configJSON,
				MetaJSON:    metaJSON,
			}, nil
		},
	}
	svc := NewContentService(repo, nil)

	config := []byte(`{"access_key_id":"abc","secret_access_key":"def","public_base_url":"https://cdn.example.com"}`)
	meta := []byte(`{"health":"ok"}`)
	got, err := svc.UpdateIntegrationProvider("cloudflare_r2", true, config, meta)
	if err != nil {
		t.Fatalf("expected update success, got error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(got.ConfigJSON, &payload); err != nil {
		t.Fatalf("failed to parse masked config json: %v", err)
	}
	if payload["access_key_id"] != "******" || payload["secret_access_key"] != "******" {
		t.Fatalf("expected secret keys to be masked, got: %#v", payload)
	}
	if payload["public_base_url"] != "https://cdn.example.com" {
		t.Fatalf("expected non-secret key to remain unchanged, got: %#v", payload)
	}
}

func TestUpdateIntegrationProviderInvalidConfigJSON(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	if _, err := svc.UpdateIntegrationProvider("cloudflare_r2", true, []byte("{"), nil); err == nil {
		t.Fatalf("expected validation error for invalid config json")
	}
}

func TestTestIntegrationProviderValidation(t *testing.T) {
	cases := []struct {
		name    string
		provide func(string) (domain.IntegrationProvider, bool)
		wantErr bool
	}{
		{
			name: "provider not found",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{}, false
			},
			wantErr: true,
		},
		{
			name: "provider disabled",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: false, ConfigJSON: []byte(`{}`)}, true
			},
			wantErr: true,
		},
		{
			name: "missing required config",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{ProviderKey: "openai_compatible", Enabled: true, ConfigJSON: []byte(`{"base_url":"https://example.com"}`)}, true
			},
			wantErr: true,
		},
		{
			name: "valid openai compatible config",
			provide: func(string) (domain.IntegrationProvider, bool) {
				return domain.IntegrationProvider{
					ProviderKey: "openai_compatible",
					Enabled:     true,
					ConfigJSON:  []byte(`{"base_url":"https://example.com","api_key":"k","model":"gpt-4.1-mini"}`),
				}, true
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := NewContentService(&contentRepoStub{getIntegrationProviderFunc: tc.provide}, nil)
			_, err := svc.TestIntegrationProvider("openai_compatible")
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected success, got error: %v", err)
			}
		})
	}
}

func TestCreateTranslationJobValidation(t *testing.T) {
	svc := NewContentService(&contentRepoStub{
		getIntegrationProviderFunc: func(string) (domain.IntegrationProvider, bool) {
			return domain.IntegrationProvider{ProviderKey: "openai_compatible", ProviderType: "llm", Enabled: true}, true
		},
		createTranslationJobFunc: func(job domain.TranslationJob) (domain.TranslationJob, error) {
			job.ID = "job-1"
			return job, nil
		},
	}, nil)

	_, err := svc.CreateTranslationJob(domain.TranslationJob{
		SourceType:   "article",
		SourceID:     "a1",
		SourceLocale: "zh-CN",
		TargetLocale: "en-US",
		ProviderKey:  "openai_compatible",
		ModelName:    "gpt-4.1-mini",
	})
	if err != nil {
		t.Fatalf("expected translation job creation success, got error: %v", err)
	}
}

func TestCreateTranslationJobRejectsSameLocale(t *testing.T) {
	svc := NewContentService(&contentRepoStub{}, nil)
	_, err := svc.CreateTranslationJob(domain.TranslationJob{
		SourceType:   "article",
		SourceID:     "a1",
		SourceLocale: "en-US",
		TargetLocale: "en-US",
		ProviderKey:  "openai_compatible",
		ModelName:    "gpt-4.1-mini",
	})
	if err == nil {
		t.Fatalf("expected error when sourceLocale equals targetLocale")
	}
}
