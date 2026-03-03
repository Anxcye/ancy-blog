// File: repository.go
// Purpose: Provide PostgreSQL-backed repository implementation for content domains.
// Module: backend/internal/repository/postgres, persistence implementation layer.
// Related: repository contracts, service layer, and schema_v1.sql.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/config"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	db *sql.DB
}

func New(ctx context.Context, cfg config.DBConfig) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Close() error { return r.db.Close() }

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
	var id string
	var createdAt, updatedAt time.Time
	var publishedAt sql.NullTime
	err := r.db.QueryRow(`
INSERT INTO articles (title, slug, content_kind, summary, content, status, visibility, allow_comment, origin_type, source_url, ai_assist_level, cover_image, published_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
RETURNING id::text, created_at, updated_at, published_at
`, article.Title, article.Slug, article.ContentKind, article.Summary, article.Content, article.Status, article.Visibility, article.AllowComment, article.OriginType, nullableString(article.SourceURL), article.AIAssistLevel, nullableString(article.CoverImage), nullableTime(article.PublishedAt)).
		Scan(&id, &createdAt, &updatedAt, &publishedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Article{}, errors.New("slug already exists")
		}
		return domain.Article{}, err
	}
	article.ID = id
	article.CreatedAt = createdAt
	article.UpdatedAt = updatedAt
	if publishedAt.Valid {
		article.PublishedAt = publishedAt.Time
	}
	return article, nil
}

func (r *Repository) UpdateArticle(id string, article domain.Article) (domain.Article, error) {
	var createdAt time.Time
	err := r.db.QueryRow(`SELECT created_at FROM articles WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Article{}, errors.New("article not found")
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

	var updatedAt time.Time
	var publishedAt sql.NullTime
	err = r.db.QueryRow(`
UPDATE articles
SET title=$2, slug=$3, content_kind=$4, summary=$5, content=$6, status=$7, visibility=$8,
    allow_comment=$9, origin_type=$10, source_url=$11, ai_assist_level=$12, cover_image=$13,
    published_at=$14, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
RETURNING updated_at, published_at
`, id, article.Title, article.Slug, article.ContentKind, article.Summary, article.Content, article.Status,
		article.Visibility, article.AllowComment, article.OriginType, nullableString(article.SourceURL), article.AIAssistLevel,
		nullableString(article.CoverImage), nullableTime(article.PublishedAt)).Scan(&updatedAt, &publishedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Article{}, errors.New("slug already exists")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Article{}, errors.New("article not found")
		}
		return domain.Article{}, err
	}
	article.ID = id
	article.CreatedAt = createdAt
	article.UpdatedAt = updatedAt
	if publishedAt.Valid {
		article.PublishedAt = publishedAt.Time
	}
	return article, nil
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

func (r *Repository) GetArticleByID(id string) (domain.Article, bool) {
	var a domain.Article
	err := r.db.QueryRow(`
SELECT id::text, title, slug, content_kind, COALESCE(summary,''), COALESCE(content,''), status, visibility,
       allow_comment, origin_type, COALESCE(source_url,''), ai_assist_level, COALESCE(cover_image,''),
       COALESCE(published_at, created_at), created_at, updated_at
FROM articles
WHERE id=$1 AND deleted_at IS NULL
`, id).Scan(&a.ID, &a.Title, &a.Slug, &a.ContentKind, &a.Summary, &a.Content, &a.Status, &a.Visibility,
		&a.AllowComment, &a.OriginType, &a.SourceURL, &a.AIAssistLevel, &a.CoverImage, &a.PublishedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return domain.Article{}, false
	}
	return a, true
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

func (r *Repository) ListPublishedMoments(page, pageSize int) ([]domain.Moment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM moments WHERE status='published' AND deleted_at IS NULL`).Scan(&total); err != nil {
		return []domain.Moment{}, 0
	}
	rows, err := r.db.Query(`
SELECT id::text, content, status, allow_comment, COALESCE(published_at, created_at), created_at, updated_at
FROM moments
WHERE status='published' AND deleted_at IS NULL
ORDER BY published_at DESC NULLS LAST, created_at DESC
LIMIT $1 OFFSET $2
`, pageSize, offset)
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

func (r *Repository) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if comment.Status == "" {
		comment.Status = "approved"
	}
	if comment.Source == "" {
		comment.Source = "web"
	}
	var id string
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(`
INSERT INTO comments (article_id, parent_id, root_id, content, status, is_pinned, nickname, email, website, avatar_url, source, ip, user_agent)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
RETURNING id::text, created_at, updated_at
`, comment.ArticleID, nullableUUID(comment.ParentID), nullableUUID(comment.RootID), comment.Content, comment.Status, false, comment.Nickname,
		nullableString(comment.Email), nullableString(comment.Website), nullableString(comment.AvatarURL), comment.Source, nullableString(comment.IP), nullableString(comment.UserAgent)).
		Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return domain.Comment{}, err
	}
	comment.ID = id
	comment.IsPinned = "0"
	comment.CreatedAt = createdAt
	comment.UpdatedAt = updatedAt
	return comment, nil
}

func (r *Repository) ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	var total int
	_ = r.db.QueryRow(`
SELECT COUNT(*) FROM comments
WHERE article_id=$1 AND status='approved' AND deleted_at IS NULL AND parent_id IS NULL
`, articleID).Scan(&total)

	rows, err := r.db.Query(`
SELECT id::text, article_id::text, COALESCE(parent_id::text,''), COALESCE(root_id::text,''), content, status, is_pinned,
       like_count, reply_count, nickname, COALESCE(email,''), COALESCE(website,''), COALESCE(avatar_url,''),
       source, COALESCE(ip,''), COALESCE(user_agent,''), created_at, updated_at
FROM comments
WHERE article_id=$1 AND status='approved' AND deleted_at IS NULL AND parent_id IS NULL
ORDER BY is_pinned DESC, created_at DESC
LIMIT $2 OFFSET $3
`, articleID, pageSize, offset)
	if err != nil {
		return []domain.Comment{}, total
	}
	defer rows.Close()
	items := make([]domain.Comment, 0)
	for rows.Next() {
		var c domain.Comment
		var pinned bool
		if err := rows.Scan(&c.ID, &c.ArticleID, &c.ParentID, &c.RootID, &c.Content, &c.Status, &pinned, &c.LikeCount, &c.ReplyCount,
			&c.Nickname, &c.Email, &c.Website, &c.AvatarURL, &c.Source, &c.IP, &c.UserAgent, &c.CreatedAt, &c.UpdatedAt); err == nil {
			if pinned {
				c.IsPinned = "1"
			} else {
				c.IsPinned = "0"
			}
			items = append(items, c)
		}
	}
	return items, total
}

func (r *Repository) ListCommentChildren(parentID string, page, pageSize int) ([]domain.Comment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	var total int
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE parent_id=$1 AND status='approved' AND deleted_at IS NULL`, parentID).Scan(&total)

	rows, err := r.db.Query(`
SELECT id::text, article_id::text, COALESCE(parent_id::text,''), COALESCE(root_id::text,''), content, status, is_pinned,
       like_count, reply_count, nickname, COALESCE(email,''), COALESCE(website,''), COALESCE(avatar_url,''),
       source, COALESCE(ip,''), COALESCE(user_agent,''), created_at, updated_at
FROM comments
WHERE parent_id=$1 AND status='approved' AND deleted_at IS NULL
ORDER BY is_pinned DESC, created_at ASC
LIMIT $2 OFFSET $3
`, parentID, pageSize, offset)
	if err != nil {
		return []domain.Comment{}, total
	}
	defer rows.Close()
	items := make([]domain.Comment, 0)
	for rows.Next() {
		var c domain.Comment
		var pinned bool
		if err := rows.Scan(&c.ID, &c.ArticleID, &c.ParentID, &c.RootID, &c.Content, &c.Status, &pinned, &c.LikeCount, &c.ReplyCount,
			&c.Nickname, &c.Email, &c.Website, &c.AvatarURL, &c.Source, &c.IP, &c.UserAgent, &c.CreatedAt, &c.UpdatedAt); err == nil {
			if pinned {
				c.IsPinned = "1"
			} else {
				c.IsPinned = "0"
			}
			items = append(items, c)
		}
	}
	return items, total
}

func (r *Repository) CountArticleComments(articleID string) (int, error) {
	var total int
	err := r.db.QueryRow(`
SELECT COUNT(*) FROM comments
WHERE article_id=$1 AND status='approved' AND deleted_at IS NULL
`, articleID).Scan(&total)
	return total, err
}

func (r *Repository) ListCommentPage(page, pageSize int, status string) ([]domain.Comment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	countSQL := `SELECT COUNT(*) FROM comments WHERE deleted_at IS NULL`
	countArgs := []any{}
	if status != "" {
		countSQL += ` AND status=$1`
		countArgs = append(countArgs, status)
	}
	var total int
	_ = r.db.QueryRow(countSQL, countArgs...).Scan(&total)

	query := `
SELECT id::text, article_id::text, COALESCE(parent_id::text,''), COALESCE(root_id::text,''), content, status, is_pinned,
       like_count, reply_count, nickname, COALESCE(email,''), COALESCE(website,''), COALESCE(avatar_url,''),
       source, COALESCE(ip,''), COALESCE(user_agent,''), created_at, updated_at
FROM comments
WHERE deleted_at IS NULL`
	args := []any{}
	if status != "" {
		query += ` AND status=$1`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	if status == "" {
		query = strings.Replace(query, "$2", "$1", 1)
		query = strings.Replace(query, "$3", "$2", 1)
		args = []any{pageSize, offset}
	} else {
		args = append(args, pageSize, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []domain.Comment{}, total
	}
	defer rows.Close()
	items := make([]domain.Comment, 0)
	for rows.Next() {
		var c domain.Comment
		var pinned bool
		if err := rows.Scan(&c.ID, &c.ArticleID, &c.ParentID, &c.RootID, &c.Content, &c.Status, &pinned, &c.LikeCount, &c.ReplyCount,
			&c.Nickname, &c.Email, &c.Website, &c.AvatarURL, &c.Source, &c.IP, &c.UserAgent, &c.CreatedAt, &c.UpdatedAt); err == nil {
			if pinned {
				c.IsPinned = "1"
			} else {
				c.IsPinned = "0"
			}
			items = append(items, c)
		}
	}
	return items, total
}

func (r *Repository) UpdateCommentAdmin(id string, status, isPinned string) (domain.Comment, error) {
	pinned := isPinned == "1" || strings.EqualFold(isPinned, "true")
	var c domain.Comment
	var b bool
	err := r.db.QueryRow(`
UPDATE comments
SET status=$2, is_pinned=$3, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
RETURNING id::text, article_id::text, COALESCE(parent_id::text,''), COALESCE(root_id::text,''), content, status, is_pinned,
          like_count, reply_count, nickname, COALESCE(email,''), COALESCE(website,''), COALESCE(avatar_url,''),
          source, COALESCE(ip,''), COALESCE(user_agent,''), created_at, updated_at
`, id, status, pinned).Scan(&c.ID, &c.ArticleID, &c.ParentID, &c.RootID, &c.Content, &c.Status, &b, &c.LikeCount, &c.ReplyCount,
		&c.Nickname, &c.Email, &c.Website, &c.AvatarURL, &c.Source, &c.IP, &c.UserAgent, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Comment{}, errors.New("comment not found")
		}
		return domain.Comment{}, err
	}
	if b {
		c.IsPinned = "1"
	} else {
		c.IsPinned = "0"
	}
	return c, nil
}

func (r *Repository) SubmitLink(link domain.Link) (domain.Link, error) {
	var id string
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(`
INSERT INTO links (name, url, avatar_url, description, contact_email, review_status, submitted_ip, submitted_user_agent)
VALUES ($1,$2,$3,$4,$5,'pending',$6,$7)
RETURNING id::text, created_at, updated_at
`, link.Name, link.URL, nullableString(link.AvatarURL), nullableString(link.Description), nullableString(link.ContactEmail), nullableString(link.SubmittedIP), nullableString(link.SubmittedUserAgent)).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return domain.Link{}, err
	}
	link.ID = id
	link.ReviewStatus = "pending"
	link.CreatedAt = createdAt
	link.UpdatedAt = updatedAt
	return link, nil
}

func (r *Repository) ListApprovedLinks() []domain.Link {
	rows, err := r.db.Query(`
SELECT id::text, name, url, COALESCE(avatar_url,''), COALESCE(description,''), COALESCE(contact_email,''), review_status,
       COALESCE(review_note,''), COALESCE(related_article_id::text,''), created_at, updated_at
FROM links
WHERE review_status='approved' AND deleted_at IS NULL
ORDER BY created_at DESC
`)
	if err != nil {
		return []domain.Link{}
	}
	defer rows.Close()
	items := make([]domain.Link, 0)
	for rows.Next() {
		var l domain.Link
		if err := rows.Scan(&l.ID, &l.Name, &l.URL, &l.AvatarURL, &l.Description, &l.ContactEmail, &l.ReviewStatus, &l.ReviewNote, &l.RelatedArticleID, &l.CreatedAt, &l.UpdatedAt); err == nil {
			items = append(items, l)
		}
	}
	return items
}

func (r *Repository) ListLinkSubmissions(page, pageSize int, reviewStatus string) ([]domain.Link, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	var total int
	if reviewStatus == "" {
		_ = r.db.QueryRow(`SELECT COUNT(*) FROM links WHERE deleted_at IS NULL`).Scan(&total)
	} else {
		_ = r.db.QueryRow(`SELECT COUNT(*) FROM links WHERE review_status=$1 AND deleted_at IS NULL`, reviewStatus).Scan(&total)
	}

	query := `
SELECT id::text, name, url, COALESCE(avatar_url,''), COALESCE(description,''), COALESCE(contact_email,''), review_status,
       COALESCE(review_note,''), COALESCE(related_article_id::text,''), created_at, updated_at
FROM links
WHERE deleted_at IS NULL`
	args := []any{}
	if reviewStatus != "" {
		query += ` AND review_status=$1`
		args = append(args, reviewStatus)
	}
	query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	if reviewStatus == "" {
		query = strings.Replace(query, "$2", "$1", 1)
		query = strings.Replace(query, "$3", "$2", 1)
		args = []any{pageSize, offset}
	} else {
		args = append(args, pageSize, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []domain.Link{}, total
	}
	defer rows.Close()
	items := make([]domain.Link, 0)
	for rows.Next() {
		var l domain.Link
		if err := rows.Scan(&l.ID, &l.Name, &l.URL, &l.AvatarURL, &l.Description, &l.ContactEmail, &l.ReviewStatus, &l.ReviewNote, &l.RelatedArticleID, &l.CreatedAt, &l.UpdatedAt); err == nil {
			items = append(items, l)
		}
	}
	return items, total
}

func (r *Repository) ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error) {
	var approvedAt sql.NullTime
	if reviewStatus == "approved" {
		approvedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	}
	var l domain.Link
	err := r.db.QueryRow(`
UPDATE links
SET review_status=$2, review_note=$3, related_article_id=$4, approved_at=$5, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
RETURNING id::text, name, url, COALESCE(avatar_url,''), COALESCE(description,''), COALESCE(contact_email,''),
          review_status, COALESCE(review_note,''), COALESCE(related_article_id::text,''), created_at, updated_at
`, id, reviewStatus, nullableString(reviewNote), nullableUUID(relatedArticleID), approvedAt).Scan(&l.ID, &l.Name, &l.URL, &l.AvatarURL, &l.Description, &l.ContactEmail, &l.ReviewStatus, &l.ReviewNote, &l.RelatedArticleID, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, errors.New("link not found")
		}
		return domain.Link{}, err
	}
	return l, nil
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

func (r *Repository) GetSiteSettings() domain.SiteSettings {
	var s domain.SiteSettings
	err := r.db.QueryRow(`SELECT site_name, COALESCE(avatar_url,''), COALESCE(hero_intro_md,''), default_locale FROM site_settings ORDER BY created_at ASC LIMIT 1`).
		Scan(&s.SiteName, &s.AvatarURL, &s.HeroIntroMD, &s.DefaultLocale)
	if err != nil {
		return domain.SiteSettings{SiteName: "Ancy Blog", DefaultLocale: "en"}
	}
	return s
}

func (r *Repository) UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings {
	var id string
	err := r.db.QueryRow(`SELECT id::text FROM site_settings ORDER BY created_at ASC LIMIT 1`).Scan(&id)
	if err != nil {
		_ = r.db.QueryRow(`
INSERT INTO site_settings (site_name, avatar_url, hero_intro_md, default_locale)
VALUES ($1,$2,$3,$4) RETURNING id::text
`, settings.SiteName, nullableString(settings.AvatarURL), nullableString(settings.HeroIntroMD), settings.DefaultLocale).Scan(&id)
		return settings
	}
	_, _ = r.db.Exec(`
UPDATE site_settings
SET site_name=$2, avatar_url=$3, hero_intro_md=$4, default_locale=$5, updated_at=NOW()
WHERE id=$1
`, id, settings.SiteName, nullableString(settings.AvatarURL), nullableString(settings.HeroIntroMD), settings.DefaultLocale)
	return settings
}

func (r *Repository) CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO footer_items (label, link_type, internal_article_slug, external_url, row_num, order_num, enabled)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id::text
`, item.Label, item.LinkType, nullableString(item.InternalArticleSlug), nullableString(item.ExternalURL), item.RowNum, item.OrderNum, item.Enabled).Scan(&id)
	if err != nil {
		return domain.FooterItem{}, err
	}
	item.ID = id
	return item, nil
}

func (r *Repository) UpdateFooterItem(id string, item domain.FooterItem) (domain.FooterItem, error) {
	res, err := r.db.Exec(`
UPDATE footer_items
SET label=$2, link_type=$3, internal_article_slug=$4, external_url=$5, row_num=$6, order_num=$7, enabled=$8, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
`, id, item.Label, item.LinkType, nullableString(item.InternalArticleSlug), nullableString(item.ExternalURL), item.RowNum, item.OrderNum, item.Enabled)
	if err != nil {
		return domain.FooterItem{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.FooterItem{}, errors.New("footer item not found")
	}
	item.ID = id
	return item, nil
}

func (r *Repository) DeleteFooterItem(id string) bool {
	res, err := r.db.Exec(`UPDATE footer_items SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListFooterItems() []domain.FooterItem {
	rows, err := r.db.Query(`
SELECT id::text, label, link_type, COALESCE(internal_article_slug,''), COALESCE(external_url,''), row_num, order_num, enabled
FROM footer_items
WHERE enabled=TRUE AND deleted_at IS NULL
ORDER BY row_num ASC, order_num ASC
`)
	if err != nil {
		return []domain.FooterItem{}
	}
	defer rows.Close()
	items := make([]domain.FooterItem, 0)
	for rows.Next() {
		var f domain.FooterItem
		if err := rows.Scan(&f.ID, &f.Label, &f.LinkType, &f.InternalArticleSlug, &f.ExternalURL, &f.RowNum, &f.OrderNum, &f.Enabled); err == nil {
			items = append(items, f)
		}
	}
	return items
}

func (r *Repository) CreateSocialLink(item domain.SocialLink) (domain.SocialLink, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO social_links (platform, title, url, icon_key, order_num, enabled)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id::text
`, item.Platform, item.Title, item.URL, nullableString(item.IconKey), item.OrderNum, item.Enabled).Scan(&id)
	if err != nil {
		return domain.SocialLink{}, err
	}
	item.ID = id
	return item, nil
}

func (r *Repository) UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error) {
	res, err := r.db.Exec(`
UPDATE social_links
SET platform=$2, title=$3, url=$4, icon_key=$5, order_num=$6, enabled=$7, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
`, id, item.Platform, item.Title, item.URL, nullableString(item.IconKey), item.OrderNum, item.Enabled)
	if err != nil {
		return domain.SocialLink{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.SocialLink{}, errors.New("social link not found")
	}
	item.ID = id
	return item, nil
}

func (r *Repository) DeleteSocialLink(id string) bool {
	res, err := r.db.Exec(`UPDATE social_links SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListSocialLinks() []domain.SocialLink {
	rows, err := r.db.Query(`
SELECT id::text, platform, title, url, COALESCE(icon_key,''), order_num, enabled
FROM social_links
WHERE enabled=TRUE AND deleted_at IS NULL
ORDER BY order_num ASC
`)
	if err != nil {
		return []domain.SocialLink{}
	}
	defer rows.Close()
	items := make([]domain.SocialLink, 0)
	for rows.Next() {
		var s domain.SocialLink
		if err := rows.Scan(&s.ID, &s.Platform, &s.Title, &s.URL, &s.IconKey, &s.OrderNum, &s.Enabled); err == nil {
			items = append(items, s)
		}
	}
	return items
}

func (r *Repository) CreateNavItem(item domain.NavItem) (domain.NavItem, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO nav_items (name, key, type, target_type, target_value, order_num, enabled)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id::text
`, item.Name, item.Key, item.Type, item.TargetType, nullableString(item.TargetValue), item.OrderNum, item.Enabled).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.NavItem{}, errors.New("name and key are required")
		}
		return domain.NavItem{}, err
	}
	item.ID = id
	return item, nil
}

func (r *Repository) UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error) {
	res, err := r.db.Exec(`
UPDATE nav_items
SET name=$2, key=$3, type=$4, target_type=$5, target_value=$6, order_num=$7, enabled=$8, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
`, id, item.Name, item.Key, item.Type, item.TargetType, nullableString(item.TargetValue), item.OrderNum, item.Enabled)
	if err != nil {
		return domain.NavItem{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.NavItem{}, errors.New("nav item not found")
	}
	item.ID = id
	return item, nil
}

func (r *Repository) DeleteNavItem(id string) bool {
	res, err := r.db.Exec(`UPDATE nav_items SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListNavItems() []domain.NavItem {
	rows, err := r.db.Query(`
SELECT id::text, name, key, type, target_type, COALESCE(target_value,''), order_num, enabled
FROM nav_items
WHERE enabled=TRUE AND deleted_at IS NULL
ORDER BY order_num ASC
`)
	if err != nil {
		return []domain.NavItem{}
	}
	defer rows.Close()
	items := make([]domain.NavItem, 0)
	for rows.Next() {
		var n domain.NavItem
		if err := rows.Scan(&n.ID, &n.Name, &n.Key, &n.Type, &n.TargetType, &n.TargetValue, &n.OrderNum, &n.Enabled); err == nil {
			items = append(items, n)
		}
	}
	return items
}

func (r *Repository) CreateContentSlot(slot domain.ContentSlot) (domain.ContentSlot, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO content_slots (slot_key, name, description, enabled)
VALUES ($1,$2,$3,$4)
RETURNING id::text
`, slot.SlotKey, slot.Name, nullableString(slot.Description), slot.Enabled).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ContentSlot{}, errors.New("slot already exists")
		}
		return domain.ContentSlot{}, err
	}
	slot.ID = id
	return slot, nil
}

func (r *Repository) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO content_slot_items (slot_key, content_type, content_id, order_num, enabled)
VALUES ($1,$2,$3,$4,$5)
RETURNING id::text
`, slotKey, item.ContentType, item.ContentID, item.OrderNum, item.Enabled).Scan(&id)
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.SlotItem{}, errors.New("slot not found")
		}
		return domain.SlotItem{}, err
	}
	item.ID = id
	item.SlotKey = slotKey
	return item, nil
}

func (r *Repository) DeleteSlotItem(slotKey, itemID string) bool {
	res, err := r.db.Exec(`DELETE FROM content_slot_items WHERE slot_key=$1 AND id=$2`, slotKey, itemID)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListSlotContent(slotKey string, limit int) ([]domain.SlotContentItem, bool) {
	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM content_slots WHERE slot_key=$1)`, slotKey).Scan(&exists); err != nil || !exists {
		return nil, false
	}
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.Query(`
SELECT csi.content_type, csi.content_id::text,
       COALESCE(a.title,''), COALESCE(a.slug,''), COALESCE(a.summary,''),
       COALESCE(m.content,'')
FROM content_slot_items csi
LEFT JOIN articles a ON csi.content_type='article' AND a.id=csi.content_id AND a.status='published' AND a.deleted_at IS NULL
LEFT JOIN moments m ON csi.content_type='moment' AND m.id=csi.content_id AND m.status='published' AND m.deleted_at IS NULL
WHERE csi.slot_key=$1 AND csi.enabled=TRUE
ORDER BY csi.order_num ASC
LIMIT $2
`, slotKey, limit)
	if err != nil {
		return []domain.SlotContentItem{}, true
	}
	defer rows.Close()
	items := make([]domain.SlotContentItem, 0)
	for rows.Next() {
		var it domain.SlotContentItem
		if err := rows.Scan(&it.ContentType, &it.ID, &it.Title, &it.Slug, &it.Summary, &it.Content); err == nil {
			if it.ContentType == "article" && it.Slug == "" {
				continue
			}
			if it.ContentType == "moment" && it.Content == "" {
				continue
			}
			items = append(items, it)
		}
	}
	return items, true
}

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
			return domain.IntegrationProvider{}, errors.New("provider not found")
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
			return domain.TranslationJob{}, errors.New("provider not found")
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

func (r *Repository) ListTimeline(page, pageSize int) ([]domain.TimelineItem, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	var total int
	_ = r.db.QueryRow(`
SELECT
  (SELECT COUNT(*) FROM articles WHERE status='published' AND deleted_at IS NULL)
+ (SELECT COUNT(*) FROM moments WHERE status='published' AND deleted_at IS NULL)
`).Scan(&total)

	rows, err := r.db.Query(`
SELECT content_type, id, title, summary, slug, content, published_at
FROM (
  SELECT 'article' AS content_type,
         a.id::text AS id,
         a.title,
         COALESCE(a.summary,'') AS summary,
         a.slug,
         '' AS content,
         COALESCE(a.published_at, a.created_at) AS published_at
  FROM articles a
  WHERE a.status='published' AND a.deleted_at IS NULL
  UNION ALL
  SELECT 'moment' AS content_type,
         m.id::text AS id,
         '' AS title,
         '' AS summary,
         '' AS slug,
         m.content,
         COALESCE(m.published_at, m.created_at) AS published_at
  FROM moments m
  WHERE m.status='published' AND m.deleted_at IS NULL
) t
ORDER BY published_at DESC
LIMIT $1 OFFSET $2
`, pageSize, offset)
	if err != nil {
		return []domain.TimelineItem{}, total
	}
	defer rows.Close()
	items := make([]domain.TimelineItem, 0)
	for rows.Next() {
		var t domain.TimelineItem
		if err := rows.Scan(&t.ContentType, &t.ID, &t.Title, &t.Summary, &t.Slug, &t.Content, &t.PublishedAt); err == nil {
			items = append(items, t)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].PublishedAt.After(items[j].PublishedAt) })
	return items, total
}

func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func nullableString(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func nullableUUID(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func nullableTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func isForeignKeyViolation(err error) bool {
	return strings.Contains(err.Error(), "violates foreign key constraint")
}

func ensureJSONText(v string) []byte {
	if strings.TrimSpace(v) == "" {
		return []byte("{}")
	}
	raw := []byte(v)
	if json.Valid(raw) {
		return raw
	}
	return []byte("{}")
}
