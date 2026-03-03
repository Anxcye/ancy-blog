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
INSERT INTO translation_jobs (source_type, source_id, source_locale, target_locale, provider_key, model_name, status, error_message, result_text, requested_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING id::text, created_at, updated_at, finished_at, error_message, COALESCE(result_text,'')
`, job.SourceType, job.SourceID, job.SourceLocale, job.TargetLocale, job.ProviderKey, job.ModelName, job.Status, nullableString(job.ErrorMessage), nullableUUID(job.RequestedBy)).
		Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt, &finishedAt, &errorMessage, &job.ResultText)
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
       status, COALESCE(error_message,''), COALESCE(result_text,''), COALESCE(requested_by::text,''), created_at, updated_at, finished_at
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
			&j.Status, &j.ErrorMessage, &j.ResultText, &j.RequestedBy, &j.CreatedAt, &j.UpdatedAt, &finishedAt); err == nil {
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
       status, COALESCE(error_message,''), COALESCE(result_text,''), COALESCE(requested_by::text,''), created_at, updated_at, finished_at
FROM translation_jobs
WHERE id=$1
`, id).Scan(&j.ID, &j.SourceType, &j.SourceID, &j.SourceLocale, &j.TargetLocale, &j.ProviderKey, &j.ModelName,
		&j.Status, &j.ErrorMessage, &j.ResultText, &j.RequestedBy, &j.CreatedAt, &j.UpdatedAt, &finishedAt)
	if err != nil {
		return domain.TranslationJob{}, false
	}
	if finishedAt.Valid {
		j.FinishedAt = finishedAt.Time
	}
	return j, true
}

func (r *Repository) ClaimNextQueuedTranslationJob() (domain.TranslationJob, bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return domain.TranslationJob{}, false, err
	}
	defer func() { _ = tx.Rollback() }()

	var id string
	err = tx.QueryRow(`
SELECT id::text
FROM translation_jobs
WHERE status='queued'
ORDER BY created_at ASC
FOR UPDATE SKIP LOCKED
LIMIT 1
`).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TranslationJob{}, false, nil
		}
		return domain.TranslationJob{}, false, err
	}

	var j domain.TranslationJob
	var finishedAt sql.NullTime
	err = tx.QueryRow(`
UPDATE translation_jobs
SET status='running', updated_at=NOW()
WHERE id=$1
RETURNING id::text, source_type, source_id::text, source_locale, target_locale, provider_key, model_name,
          status, COALESCE(error_message,''), COALESCE(result_text,''), COALESCE(requested_by::text,''), created_at, updated_at, finished_at
`, id).Scan(&j.ID, &j.SourceType, &j.SourceID, &j.SourceLocale, &j.TargetLocale, &j.ProviderKey, &j.ModelName,
		&j.Status, &j.ErrorMessage, &j.ResultText, &j.RequestedBy, &j.CreatedAt, &j.UpdatedAt, &finishedAt)
	if err != nil {
		return domain.TranslationJob{}, false, err
	}
	if finishedAt.Valid {
		j.FinishedAt = finishedAt.Time
	}
	if err := tx.Commit(); err != nil {
		return domain.TranslationJob{}, false, err
	}
	return j, true, nil
}

func (r *Repository) MarkTranslationJobRunning(id string) error {
	_, err := r.db.Exec(`UPDATE translation_jobs SET status='running', updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *Repository) MarkTranslationJobSucceeded(id, resultText string) error {
	_, err := r.db.Exec(`
UPDATE translation_jobs
SET status='succeeded', result_text=$2, error_message=NULL, updated_at=NOW(), finished_at=NOW()
WHERE id=$1
`, id, nullableString(resultText))
	return err
}

func (r *Repository) MarkTranslationJobFailed(id, errorMessage string) error {
	_, err := r.db.Exec(`
UPDATE translation_jobs
SET status='failed', error_message=$2, updated_at=NOW(), finished_at=NOW()
WHERE id=$1
`, id, nullableString(errorMessage))
	return err
}

func (r *Repository) GetTranslationSourceText(sourceType, sourceID string) (string, bool, error) {
	switch sourceType {
	case "article":
		var title, summary, content string
		err := r.db.QueryRow(`
SELECT COALESCE(title,''), COALESCE(summary,''), COALESCE(content,'')
FROM articles
WHERE id=$1 AND deleted_at IS NULL
`, sourceID).Scan(&title, &summary, &content)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", false, nil
			}
			return "", false, err
		}
		payload := fmt.Sprintf("# %s\n\n%s\n\n%s", title, summary, content)
		return strings.TrimSpace(payload), true, nil
	case "moment":
		var content string
		err := r.db.QueryRow(`
SELECT COALESCE(content,'')
FROM moments
WHERE id=$1 AND deleted_at IS NULL
`, sourceID).Scan(&content)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", false, nil
			}
			return "", false, err
		}
		return content, true, nil
	default:
		return "", false, nil
	}
}

func (r *Repository) UpsertArticleTranslation(articleID, locale, content, translatedByJobID string) error {
	_, err := r.db.Exec(`
INSERT INTO article_translations (article_id, locale, content, translated_by_job_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (article_id, locale)
DO UPDATE SET content = EXCLUDED.content, translated_by_job_id = EXCLUDED.translated_by_job_id, updated_at = NOW()
`, articleID, locale, nullableString(content), nullableUUID(translatedByJobID))
	return err
}

func (r *Repository) UpsertMomentTranslation(momentID, locale, content, translatedByJobID string) error {
	_, err := r.db.Exec(`
INSERT INTO moment_translations (moment_id, locale, content, translated_by_job_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (moment_id, locale)
DO UPDATE SET content = EXCLUDED.content, translated_by_job_id = EXCLUDED.translated_by_job_id, updated_at = NOW()
`, momentID, locale, nullableString(content), nullableUUID(translatedByJobID))
	return err
}

func (r *Repository) ListTranslationContents(page, pageSize int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	where := []string{"1=1"}
	args := make([]any, 0, 4)
	idx := 1
	if strings.TrimSpace(sourceID) != "" {
		where = append(where, fmt.Sprintf("source_id=$%d", idx))
		args = append(args, sourceID)
		idx++
	}
	if strings.TrimSpace(locale) != "" {
		where = append(where, fmt.Sprintf("locale=$%d", idx))
		args = append(args, locale)
		idx++
	}
	cond := strings.Join(where, " AND ")

	var baseQuery string
	switch sourceType {
	case "article":
		baseQuery = `
SELECT article_id::text AS source_id, locale, COALESCE(content,''), COALESCE(translated_by_job_id::text,''), created_at, updated_at
FROM article_translations`
	case "moment":
		baseQuery = `
SELECT moment_id::text AS source_id, locale, COALESCE(content,''), COALESCE(translated_by_job_id::text,''), created_at, updated_at
FROM moment_translations`
	default:
		return []domain.TranslationContent{}, 0
	}

	var total int
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) t WHERE %s", baseQuery, cond)
	if err := r.db.QueryRow(countSQL, args...).Scan(&total); err != nil {
		return []domain.TranslationContent{}, 0
	}

	listArgs := append(args, pageSize, offset)
	listSQL := fmt.Sprintf(`
SELECT source_id, locale, content, translated_by_job_id, created_at, updated_at
FROM (%s) t
WHERE %s
ORDER BY updated_at DESC
LIMIT $%d OFFSET $%d
`, baseQuery, cond, idx, idx+1)
	rows, err := r.db.Query(listSQL, listArgs...)
	if err != nil {
		return []domain.TranslationContent{}, total
	}
	defer rows.Close()

	out := make([]domain.TranslationContent, 0)
	for rows.Next() {
		var item domain.TranslationContent
		item.SourceType = sourceType
		if err := rows.Scan(&item.SourceID, &item.Locale, &item.Content, &item.TranslatedByJobID, &item.CreatedAt, &item.UpdatedAt); err == nil {
			out = append(out, item)
		}
	}
	return out, total
}

func (r *Repository) GetTranslationContent(sourceType, sourceID, locale string) (domain.TranslationContent, bool) {
	row := domain.TranslationContent{SourceType: sourceType}
	var err error
	switch sourceType {
	case "article":
		err = r.db.QueryRow(`
SELECT article_id::text, locale, COALESCE(content,''), COALESCE(translated_by_job_id::text,''), created_at, updated_at
FROM article_translations
WHERE article_id=$1 AND locale=$2
`, sourceID, locale).Scan(&row.SourceID, &row.Locale, &row.Content, &row.TranslatedByJobID, &row.CreatedAt, &row.UpdatedAt)
	case "moment":
		err = r.db.QueryRow(`
SELECT moment_id::text, locale, COALESCE(content,''), COALESCE(translated_by_job_id::text,''), created_at, updated_at
FROM moment_translations
WHERE moment_id=$1 AND locale=$2
`, sourceID, locale).Scan(&row.SourceID, &row.Locale, &row.Content, &row.TranslatedByJobID, &row.CreatedAt, &row.UpdatedAt)
	default:
		return domain.TranslationContent{}, false
	}
	if err != nil {
		return domain.TranslationContent{}, false
	}
	return row, true
}

func (r *Repository) UpsertTranslationContent(sourceType, sourceID, locale, content, translatedByJobID string) (domain.TranslationContent, error) {
	switch sourceType {
	case "article":
		if err := r.UpsertArticleTranslation(sourceID, locale, content, translatedByJobID); err != nil {
			return domain.TranslationContent{}, err
		}
	case "moment":
		if err := r.UpsertMomentTranslation(sourceID, locale, content, translatedByJobID); err != nil {
			return domain.TranslationContent{}, err
		}
	default:
		return domain.TranslationContent{}, apperr.ErrValidation
	}
	row, ok := r.GetTranslationContent(sourceType, sourceID, locale)
	if !ok {
		return domain.TranslationContent{}, sql.ErrNoRows
	}
	return row, nil
}
