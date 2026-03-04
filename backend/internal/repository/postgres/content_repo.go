// File: content_repo.go
// Purpose: Implement PostgreSQL repository methods for articles, moments, and taxonomy.
// Module: backend/internal/repository/postgres, content persistence layer.
// Related: repository.go, helpers.go, and service article flows.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

func (r *Repository) CreateArticle(article domain.Article) (domain.Article, error) {
	if article.ContentKind == "" {
		article.ContentKind = "post"
	}
	if article.Status == "" {
		article.Status = "draft"
	}
	if article.Visibility == "" {
		article.Visibility = "public"
	}
	if article.OriginType == "" {
		article.OriginType = "original"
	}
	if article.AIAssistLevel == "" {
		article.AIAssistLevel = "none"
	}
	categoryID := r.resolveCategoryID(article.CategorySlug)
	var id string
	var createdAt, updatedAt time.Time
	var publishedAt sql.NullTime
	err := r.db.QueryRow(`
INSERT INTO articles (title, slug, content_kind, summary, content, status, visibility, allow_comment, origin_type, source_url, ai_assist_level, cover_image, published_at, category_id, is_pinned, is_featured)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
RETURNING id::text, created_at, updated_at, published_at
`, article.Title, article.Slug, article.ContentKind, article.Summary, article.Content, article.Status, article.Visibility, article.AllowComment, article.OriginType, nullableString(article.SourceURL), article.AIAssistLevel, nullableString(article.CoverImage), nullableTime(article.PublishedAt), nullableString(categoryID), article.IsPinned, article.IsFeatured).
		Scan(&id, &createdAt, &updatedAt, &publishedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Article{}, apperr.ErrSlugAlreadyExists
		}
		return domain.Article{}, err
	}
	article.ID = id
	article.CreatedAt = createdAt
	article.UpdatedAt = updatedAt
	if publishedAt.Valid {
		article.PublishedAt = publishedAt.Time
	}
	r.syncArticleTags(id, r.resolveTagIDs(article.TagSlugs))
	return article, nil
}

func (r *Repository) UpdateArticle(id string, article domain.Article) (domain.Article, error) {
	var createdAt time.Time
	err := r.db.QueryRow(`SELECT created_at FROM articles WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Article{}, apperr.ErrArticleNotFound
		}
		return domain.Article{}, err
	}

	if article.ContentKind == "" {
		article.ContentKind = "post"
	}
	if article.Visibility == "" {
		article.Visibility = "public"
	}
	if article.OriginType == "" {
		article.OriginType = "original"
	}
	if article.AIAssistLevel == "" {
		article.AIAssistLevel = "none"
	}

	categoryID := r.resolveCategoryID(article.CategorySlug)
	var updatedAt time.Time
	var publishedAt sql.NullTime
	err = r.db.QueryRow(`
UPDATE articles
SET title=$2, slug=$3, content_kind=$4, summary=$5, content=$6, status=$7, visibility=$8,
    allow_comment=$9, origin_type=$10, source_url=$11, ai_assist_level=$12, cover_image=$13,
    published_at=$14, category_id=$15, is_pinned=$16, is_featured=$17, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
RETURNING updated_at, published_at
`, id, article.Title, article.Slug, article.ContentKind, article.Summary, article.Content, article.Status,
		article.Visibility, article.AllowComment, article.OriginType, nullableString(article.SourceURL), article.AIAssistLevel,
		nullableString(article.CoverImage), nullableTime(article.PublishedAt), nullableString(categoryID), article.IsPinned, article.IsFeatured).Scan(&updatedAt, &publishedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Article{}, apperr.ErrSlugAlreadyExists
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Article{}, apperr.ErrArticleNotFound
		}
		return domain.Article{}, err
	}
	article.ID = id
	article.CreatedAt = createdAt
	article.UpdatedAt = updatedAt
	if publishedAt.Valid {
		article.PublishedAt = publishedAt.Time
	}
	r.syncArticleTags(id, r.resolveTagIDs(article.TagSlugs))
	return article, nil
}

func (r *Repository) ListArticles(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	conditions := []string{"deleted_at IS NULL"}
	args := make([]any, 0, 5)
	if status != "" {
		args = append(args, status)
		conditions = append(conditions, "status = $"+strconv.Itoa(len(args)))
	}
	if contentKind != "" {
		args = append(args, contentKind)
		conditions = append(conditions, "content_kind = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(keyword) != "" {
		args = append(args, "%"+strings.TrimSpace(keyword)+"%")
		conditions = append(conditions, "(title ILIKE $"+strconv.Itoa(len(args))+" OR slug ILIKE $"+strconv.Itoa(len(args))+")")
	}

	whereClause := strings.Join(conditions, " AND ")

	var total int
	countQuery := "SELECT COUNT(*) FROM articles WHERE " + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return []domain.Article{}, 0
	}

	listArgs := append(args, pageSize, offset)
	query := `
SELECT id::text, title, slug, content_kind, COALESCE(summary,''), COALESCE(content,''), status, visibility,
       allow_comment, origin_type, COALESCE(source_url,''), ai_assist_level, COALESCE(cover_image,''),
       COALESCE(published_at, created_at), created_at, updated_at
FROM articles
WHERE ` + whereClause + `
ORDER BY updated_at DESC, created_at DESC
LIMIT $` + strconv.Itoa(len(listArgs)-1) + ` OFFSET $` + strconv.Itoa(len(listArgs))
	rows, err := r.db.Query(query, listArgs...)
	if err != nil {
		return []domain.Article{}, total
	}
	defer rows.Close()

	items := make([]domain.Article, 0)
	for rows.Next() {
		var a domain.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Slug, &a.ContentKind, &a.Summary, &a.Content, &a.Status, &a.Visibility,
			&a.AllowComment, &a.OriginType, &a.SourceURL, &a.AIAssistLevel, &a.CoverImage, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt); err == nil {
			items = append(items, a)
		}
	}
	return items, total
}

func (r *Repository) DeleteArticle(id string) bool {
	res, err := r.db.Exec(`UPDATE articles SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) BatchUpdateArticleStatus(ids []string, status string) int {
	if len(ids) == 0 || strings.TrimSpace(status) == "" {
		return 0
	}
	placeholders := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids)+1)
	args = append(args, status)
	for i, id := range ids {
		args = append(args, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
	}
	query := `
UPDATE articles
SET status=$1,
    published_at=CASE WHEN $1='published' AND published_at IS NULL THEN NOW() ELSE published_at END,
    updated_at=NOW()
WHERE id IN (` + strings.Join(placeholders, ",") + `) AND deleted_at IS NULL`
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0
	}
	n, _ := res.RowsAffected()
	return int(n)
}

func (r *Repository) ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int) {
	_ = category
	_ = tag
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM articles WHERE status='published' AND deleted_at IS NULL`
	args := []any{}
	if contentKind != "" {
		countQuery += ` AND content_kind=$1`
		args = append(args, contentKind)
	}
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return []domain.Article{}, 0
	}

	query := `
SELECT id::text, title, slug, content_kind, COALESCE(summary,''), COALESCE(content,''), status, visibility,
       allow_comment, origin_type, COALESCE(source_url,''), ai_assist_level, COALESCE(cover_image,''),
       COALESCE(published_at, created_at), created_at, updated_at
FROM articles
WHERE status='published' AND deleted_at IS NULL`
	if contentKind != "" {
		query += ` AND content_kind=$1`
	}
	query += ` ORDER BY published_at DESC NULLS LAST, created_at DESC LIMIT $2 OFFSET $3`
	if contentKind == "" {
		query = strings.Replace(query, "$2", "$1", 1)
		query = strings.Replace(query, "$3", "$2", 1)
		args = []any{pageSize, offset}
	} else {
		args = append(args, pageSize, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []domain.Article{}, total
	}
	defer rows.Close()
	items := make([]domain.Article, 0)
	for rows.Next() {
		var a domain.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.Slug, &a.ContentKind, &a.Summary, &a.Content, &a.Status, &a.Visibility,
			&a.AllowComment, &a.OriginType, &a.SourceURL, &a.AIAssistLevel, &a.CoverImage, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt); err == nil {
			items = append(items, a)
		}
	}
	return items, total
}

func (r *Repository) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	var a domain.Article
	err := r.db.QueryRow(`
SELECT id::text, title, slug, content_kind, COALESCE(summary,''), COALESCE(content,''), status, visibility,
       allow_comment, origin_type, COALESCE(source_url,''), ai_assist_level, COALESCE(cover_image,''),
       COALESCE(published_at, created_at), created_at, updated_at
FROM articles
WHERE slug=$1 AND status='published' AND deleted_at IS NULL
`, slug).Scan(&a.ID, &a.Title, &a.Slug, &a.ContentKind, &a.Summary, &a.Content, &a.Status, &a.Visibility,
		&a.AllowComment, &a.OriginType, &a.SourceURL, &a.AIAssistLevel, &a.CoverImage, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return domain.Article{}, false
	}
	return a, true
}

func (r *Repository) GetPublishedArticleBySlugWithLocale(slug, locale string) (domain.Article, bool) {
	var a domain.Article
	err := r.db.QueryRow(`
SELECT a.id::text, COALESCE(at.title, a.title), a.slug, a.content_kind, COALESCE(at.summary, a.summary, ''), COALESCE(at.content, a.content), a.status, a.visibility,
       a.allow_comment, a.origin_type, COALESCE(a.source_url,''), a.ai_assist_level, COALESCE(a.cover_image,''),
       COALESCE(a.published_at, a.created_at), a.created_at, a.updated_at
FROM articles a
LEFT JOIN article_translations at
  ON at.article_id = a.id
 AND at.locale = $2
 AND at.status = 'published'
 AND (at.published_at IS NULL OR at.published_at <= NOW())
WHERE a.slug=$1 AND a.status='published' AND a.deleted_at IS NULL
`, slug, locale).Scan(&a.ID, &a.Title, &a.Slug, &a.ContentKind, &a.Summary, &a.Content, &a.Status, &a.Visibility,
		&a.AllowComment, &a.OriginType, &a.SourceURL, &a.AIAssistLevel, &a.CoverImage, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return domain.Article{}, false
	}
	return a, true
}

func (r *Repository) GetArticleByID(id string) (domain.Article, bool) {
	var a domain.Article
	var categorySlug sql.NullString
	err := r.db.QueryRow(`
SELECT a.id::text, a.title, a.slug, a.content_kind, COALESCE(a.summary,''), COALESCE(a.content,''), a.status, a.visibility,
       a.allow_comment, a.is_pinned, a.is_featured, a.origin_type, COALESCE(a.source_url,''), a.ai_assist_level, COALESCE(a.cover_image,''),
       COALESCE(a.published_at, a.created_at), a.created_at, a.updated_at,
       c.slug
FROM articles a
LEFT JOIN categories c ON c.id=a.category_id AND c.deleted_at IS NULL
WHERE a.id=$1 AND a.deleted_at IS NULL
`, id).Scan(&a.ID, &a.Title, &a.Slug, &a.ContentKind, &a.Summary, &a.Content, &a.Status, &a.Visibility,
		&a.AllowComment, &a.IsPinned, &a.IsFeatured, &a.OriginType, &a.SourceURL, &a.AIAssistLevel, &a.CoverImage, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt,
		&categorySlug)
	if err != nil {
		return domain.Article{}, false
	}
	if categorySlug.Valid {
		a.CategorySlug = categorySlug.String
	}
	a.TagSlugs = r.articleTagSlugs(id)
	return a, true
}

func (r *Repository) SlugExists(slug string) bool {
	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM articles WHERE slug=$1 AND deleted_at IS NULL)`, slug).Scan(&exists); err != nil {
		return false
	}
	return exists
}

func (r *Repository) CreateMoment(moment domain.Moment) (domain.Moment, error) {
	var id string
	var createdAt, updatedAt time.Time
	var publishedAt sql.NullTime
	err := r.db.QueryRow(`
INSERT INTO moments (content, status, allow_comment, published_at)
VALUES ($1,$2,$3,$4)
RETURNING id::text, created_at, updated_at, published_at
`, moment.Content, moment.Status, moment.AllowComment, nullableTime(moment.PublishedAt)).Scan(&id, &createdAt, &updatedAt, &publishedAt)
	if err != nil {
		return domain.Moment{}, err
	}
	moment.ID = id
	moment.CreatedAt = createdAt
	moment.UpdatedAt = updatedAt
	if publishedAt.Valid {
		moment.PublishedAt = publishedAt.Time
	}
	return moment, nil
}

func (r *Repository) UpdateMoment(id string, moment domain.Moment) (domain.Moment, error) {
	var createdAt time.Time
	err := r.db.QueryRow(`SELECT created_at FROM moments WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Moment{}, apperr.ErrMomentNotFound
		}
		return domain.Moment{}, err
	}

	var updatedAt time.Time
	var publishedAt sql.NullTime
	err = r.db.QueryRow(`
UPDATE moments
SET content=$2, status=$3, allow_comment=$4, published_at=$5, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
RETURNING updated_at, published_at
`, id, moment.Content, moment.Status, moment.AllowComment, nullableTime(moment.PublishedAt)).
		Scan(&updatedAt, &publishedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Moment{}, apperr.ErrMomentNotFound
		}
		return domain.Moment{}, err
	}

	moment.ID = id
	moment.CreatedAt = createdAt
	moment.UpdatedAt = updatedAt
	if publishedAt.Valid {
		moment.PublishedAt = publishedAt.Time
	}
	return moment, nil
}

func (r *Repository) DeleteMoment(id string) bool {
	res, err := r.db.Exec(`UPDATE moments SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) BatchUpdateMomentStatus(ids []string, status string) int {
	if len(ids) == 0 || strings.TrimSpace(status) == "" {
		return 0
	}
	placeholders := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids)+1)
	args = append(args, status)
	for i, id := range ids {
		args = append(args, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+2))
	}
	query := `
UPDATE moments
SET status=$1,
    published_at=CASE WHEN $1='published' AND published_at IS NULL THEN NOW() ELSE published_at END,
    updated_at=NOW()
WHERE id IN (` + strings.Join(placeholders, ",") + `) AND deleted_at IS NULL`
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0
	}
	n, _ := res.RowsAffected()
	return int(n)
}

func (r *Repository) ListMoments(page, pageSize int, status string) ([]domain.Moment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	whereClause := "deleted_at IS NULL"
	args := make([]any, 0, 3)
	if status != "" {
		whereClause += " AND status=$1"
		args = append(args, status)
	}

	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM moments WHERE `+whereClause, args...).Scan(&total); err != nil {
		return []domain.Moment{}, 0
	}

	listQuery := `
SELECT id::text, content, status, allow_comment, COALESCE(published_at, created_at), created_at, updated_at
FROM moments
WHERE ` + whereClause + `
ORDER BY updated_at DESC, created_at DESC
LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	listArgs := append(args, pageSize, offset)
	rows, err := r.db.Query(listQuery, listArgs...)
	if err != nil {
		return []domain.Moment{}, total
	}
	defer rows.Close()

	items := make([]domain.Moment, 0)
	for rows.Next() {
		var m domain.Moment
		if err := rows.Scan(&m.ID, &m.Content, &m.Status, &m.AllowComment, &m.PublishedAt, &m.CreatedAt, &m.UpdatedAt); err == nil {
			items = append(items, m)
		}
	}
	return items, total
}

func (r *Repository) ListPublishedMoments(page, pageSize int, locale string) ([]domain.Moment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM moments WHERE status='published' AND deleted_at IS NULL`).Scan(&total); err != nil {
		return []domain.Moment{}, 0
	}
	rows, err := r.db.Query(`
SELECT m.id::text, COALESCE(mt.content, m.content), m.status, m.allow_comment,
       COALESCE(m.published_at, m.created_at), m.created_at, m.updated_at
FROM moments m
LEFT JOIN moment_translations mt
  ON mt.moment_id = m.id
 AND mt.locale = $1
 AND mt.status = 'published'
 AND (mt.published_at IS NULL OR mt.published_at <= NOW())
WHERE m.status='published' AND m.deleted_at IS NULL
ORDER BY m.published_at DESC NULLS LAST, m.created_at DESC
LIMIT $2 OFFSET $3
`, locale, pageSize, offset)
	if err != nil {
		return []domain.Moment{}, total
	}
	defer rows.Close()
	items := make([]domain.Moment, 0)
	for rows.Next() {
		var m domain.Moment
		if err := rows.Scan(&m.ID, &m.Content, &m.Status, &m.AllowComment, &m.PublishedAt, &m.CreatedAt, &m.UpdatedAt); err == nil {
			items = append(items, m)
		}
	}
	return items, total
}

func (r *Repository) ListCategories() []domain.Category {
	rows, err := r.db.Query(`SELECT id::text, name, slug FROM categories WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		return []domain.Category{}
	}
	defer rows.Close()
	items := make([]domain.Category, 0)
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err == nil {
			items = append(items, c)
		}
	}
	return items
}

func (r *Repository) CreateCategory(category domain.Category) (domain.Category, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO categories (name, slug) VALUES ($1, $2) RETURNING id::text`,
		category.Name, category.Slug,
	).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Category{}, apperr.ErrValidation
		}
		return domain.Category{}, err
	}
	category.ID = id
	return category, nil
}

func (r *Repository) DeleteCategory(id string) bool {
	res, err := r.db.Exec(`UPDATE categories SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListTags() []domain.Tag {
	rows, err := r.db.Query(`SELECT id::text, name, slug FROM tags WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		return []domain.Tag{}
	}
	defer rows.Close()
	items := make([]domain.Tag, 0)
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err == nil {
			items = append(items, t)
		}
	}
	return items
}

func (r *Repository) CreateTag(tag domain.Tag) (domain.Tag, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO tags (name, slug) VALUES ($1, $2) RETURNING id::text`,
		tag.Name, tag.Slug,
	).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Tag{}, apperr.ErrValidation
		}
		return domain.Tag{}, err
	}
	tag.ID = id
	return tag, nil
}

func (r *Repository) DeleteTag(id string) bool {
	res, err := r.db.Exec(`UPDATE tags SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

// resolveCategoryID returns the UUID for a category slug, or empty string if slug is blank.
func (r *Repository) resolveCategoryID(slug string) string {
	if strings.TrimSpace(slug) == "" {
		return ""
	}
	var id string
	if err := r.db.QueryRow(`SELECT id::text FROM categories WHERE slug=$1 AND deleted_at IS NULL`, slug).Scan(&id); err != nil {
		return ""
	}
	return id
}

// resolveTagIDs converts a list of tag slugs to UUIDs, silently skipping unknown slugs.
func (r *Repository) resolveTagIDs(slugs []string) []string {
	ids := make([]string, 0, len(slugs))
	for _, s := range slugs {
		if strings.TrimSpace(s) == "" {
			continue
		}
		var id string
		if err := r.db.QueryRow(`SELECT id::text FROM tags WHERE slug=$1 AND deleted_at IS NULL`, s).Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// syncArticleTags replaces article_tags rows for the given article.
func (r *Repository) syncArticleTags(articleID string, tagIDs []string) {
	_, _ = r.db.Exec(`DELETE FROM article_tags WHERE article_id=$1`, articleID)
	for _, tid := range tagIDs {
		_, _ = r.db.Exec(`INSERT INTO article_tags (article_id, tag_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, articleID, tid)
	}
}

// articleTagSlugs returns the tag slugs associated with an article.
func (r *Repository) articleTagSlugs(articleID string) []string {
	rows, err := r.db.Query(`
SELECT t.slug FROM tags t
JOIN article_tags at ON at.tag_id=t.id
WHERE at.article_id=$1 AND t.deleted_at IS NULL
ORDER BY t.name ASC`, articleID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()
	slugs := make([]string, 0)
	for rows.Next() {
		var s string
		if rows.Scan(&s) == nil {
			slugs = append(slugs, s)
		}
	}
	return slugs
}
