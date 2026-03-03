// File: integration_repo.go
// Purpose: Implement PostgreSQL repository methods for integration providers and translation jobs.
// Module: backend/internal/repository/postgres, integration persistence layer.
// Related: helpers.go and service integration/translation flows.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

func (r *Repository) ListIntegrationProviders(providerType string) []domain.IntegrationProvider {
	query := `
SELECT id::text, provider_type, provider_key, name, enabled,
       config_json::text, COALESCE(meta_json::text, '{}'), created_at, updated_at
FROM integration_providers`
	args := []any{}
	if strings.TrimSpace(providerType) != "" {
		query += ` WHERE provider_type=$1`
		args = append(args, providerType)
	}
	query += ` ORDER BY provider_type ASC, provider_key ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []domain.IntegrationProvider{}
	}
	defer rows.Close()

	items := make([]domain.IntegrationProvider, 0)
	for rows.Next() {
		var p domain.IntegrationProvider
		var configText, metaText string
		if err := rows.Scan(&p.ID, &p.ProviderType, &p.ProviderKey, &p.Name, &p.Enabled, &configText, &metaText, &p.CreatedAt, &p.UpdatedAt); err == nil {
			p.ConfigJSON = ensureJSONText(configText)
			p.MetaJSON = ensureJSONText(metaText)
			items = append(items, p)
		}
	}
	return items
}

func (r *Repository) GetIntegrationProvider(providerKey string) (domain.IntegrationProvider, bool) {
	var p domain.IntegrationProvider
	var configText, metaText string
	err := r.db.QueryRow(`
SELECT id::text, provider_type, provider_key, name, enabled,
       config_json::text, COALESCE(meta_json::text, '{}'), created_at, updated_at
FROM integration_providers
WHERE provider_key=$1
`, providerKey).Scan(&p.ID, &p.ProviderType, &p.ProviderKey, &p.Name, &p.Enabled, &configText, &metaText, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return domain.IntegrationProvider{}, false
	}
	p.ConfigJSON = ensureJSONText(configText)
	p.MetaJSON = ensureJSONText(metaText)
	return p, true
}

func (r *Repository) UpdateIntegrationProvider(providerKey string, enabled bool, configJSON, metaJSON []byte) (domain.IntegrationProvider, error) {
	if len(configJSON) == 0 {
		configJSON = []byte("{}")
	}
	if len(metaJSON) == 0 {
		metaJSON = []byte("{}")
	}
	var p domain.IntegrationProvider
	var configText, metaText string
	err := r.db.QueryRow(`
UPDATE integration_providers
SET enabled=$2, config_json=$3::jsonb, meta_json=$4::jsonb, updated_at=NOW()
WHERE provider_key=$1
RETURNING id::text, provider_type, provider_key, name, enabled, config_json::text, COALESCE(meta_json::text, '{}'), created_at, updated_at
`, providerKey, enabled, string(configJSON), string(metaJSON)).
		Scan(&p.ID, &p.ProviderType, &p.ProviderKey, &p.Name, &p.Enabled, &configText, &metaText, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.IntegrationProvider{}, apperr.ErrProviderNotFound
		}
		return domain.IntegrationProvider{}, err
	}
	p.ConfigJSON = ensureJSONText(configText)
	p.MetaJSON = ensureJSONText(metaText)
	return p, nil
}

func (r *Repository) CreateTranslationJob(job domain.TranslationJob) (domain.TranslationJob, error) {
	var finishedAt sql.NullTime
	var errorMessage sql.NullString
	err := r.db.QueryRow(`
INSERT INTO translation_jobs (source_type, source_id, source_locale, target_locale, provider_key, model_name, status, error_message, requested_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING id::text, created_at, updated_at, finished_at, error_message
`, job.SourceType, job.SourceID, job.SourceLocale, job.TargetLocale, job.ProviderKey, job.ModelName, job.Status, nullableString(job.ErrorMessage), nullableUUID(job.RequestedBy)).
		Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt, &finishedAt, &errorMessage)
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.TranslationJob{}, apperr.ErrProviderNotFound
		}
		return domain.TranslationJob{}, err
	}
	if finishedAt.Valid {
		job.FinishedAt = finishedAt.Time
	}
	if errorMessage.Valid {
		job.ErrorMessage = errorMessage.String
	}
	return job, nil
}

func (r *Repository) ListTranslationJobs(page, pageSize int, status, sourceType, sourceID string) ([]domain.TranslationJob, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	conds := []string{"1=1"}
	args := make([]any, 0)
	countArgs := make([]any, 0)
	idx := 1
	if strings.TrimSpace(status) != "" {
		conds = append(conds, fmt.Sprintf("status=$%d", idx))
		args = append(args, status)
		countArgs = append(countArgs, status)
		idx++
	}
	if strings.TrimSpace(sourceType) != "" {
		conds = append(conds, fmt.Sprintf("source_type=$%d", idx))
		args = append(args, sourceType)
		countArgs = append(countArgs, sourceType)
		idx++
	}
	if strings.TrimSpace(sourceID) != "" {
		conds = append(conds, fmt.Sprintf("source_id=$%d", idx))
		args = append(args, sourceID)
		countArgs = append(countArgs, sourceID)
		idx++
	}

	where := strings.Join(conds, " AND ")
	var total int
	countQuery := "SELECT COUNT(*) FROM translation_jobs WHERE " + where
	if err := r.db.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		return []domain.TranslationJob{}, 0
	}

	args = append(args, pageSize, offset)
	query := fmt.Sprintf(`
SELECT id::text, source_type, source_id::text, source_locale, target_locale, provider_key, model_name,
       status, COALESCE(error_message,''), COALESCE(requested_by::text,''), created_at, updated_at, finished_at
FROM translation_jobs
WHERE %s
ORDER BY created_at DESC
LIMIT $%d OFFSET $%d
`, where, idx, idx+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []domain.TranslationJob{}, total
	}
	defer rows.Close()

	items := make([]domain.TranslationJob, 0)
	for rows.Next() {
		var j domain.TranslationJob
		var finishedAt sql.NullTime
		if err := rows.Scan(&j.ID, &j.SourceType, &j.SourceID, &j.SourceLocale, &j.TargetLocale, &j.ProviderKey, &j.ModelName,
			&j.Status, &j.ErrorMessage, &j.RequestedBy, &j.CreatedAt, &j.UpdatedAt, &finishedAt); err == nil {
			if finishedAt.Valid {
				j.FinishedAt = finishedAt.Time
			}
			items = append(items, j)
		}
	}
	return items, total
}

func (r *Repository) GetTranslationJobByID(id string) (domain.TranslationJob, bool) {
	var j domain.TranslationJob
	var finishedAt sql.NullTime
	err := r.db.QueryRow(`
SELECT id::text, source_type, source_id::text, source_locale, target_locale, provider_key, model_name,
       status, COALESCE(error_message,''), COALESCE(requested_by::text,''), created_at, updated_at, finished_at
FROM translation_jobs
WHERE id=$1
`, id).Scan(&j.ID, &j.SourceType, &j.SourceID, &j.SourceLocale, &j.TargetLocale, &j.ProviderKey, &j.ModelName,
		&j.Status, &j.ErrorMessage, &j.RequestedBy, &j.CreatedAt, &j.UpdatedAt, &finishedAt)
	if err != nil {
		return domain.TranslationJob{}, false
	}
	if finishedAt.Valid {
		j.FinishedAt = finishedAt.Time
	}
	return j, true
}
