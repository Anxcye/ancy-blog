// File: comment_repo.go
// Purpose: Implement PostgreSQL repository methods for comment creation, queries, and moderation.
// Module: backend/internal/repository/postgres, comment persistence layer.
// Related: helpers.go and service comment flows.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

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
	return scanCommentRows(rows), total
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
	return scanCommentRows(rows), total
}

func (r *Repository) ListCommentDescendants(rootIDs []string) []domain.Comment {
	if len(rootIDs) == 0 {
		return []domain.Comment{}
	}

	args := make([]any, 0, len(rootIDs))
	placeholders := make([]string, 0, len(rootIDs))
	for idx, id := range rootIDs {
		args = append(args, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx+1))
	}

	query := fmt.Sprintf(`
SELECT id::text, article_id::text, COALESCE(parent_id::text,''), COALESCE(root_id::text,''), content, status, is_pinned,
       like_count, reply_count, nickname, COALESCE(email,''), COALESCE(website,''), COALESCE(avatar_url,''),
       source, COALESCE(ip,''), COALESCE(user_agent,''), created_at, updated_at
FROM comments
WHERE status='approved' AND deleted_at IS NULL AND parent_id IS NOT NULL AND root_id IN (%s)
ORDER BY root_id, is_pinned DESC, created_at ASC
`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return []domain.Comment{}
	}
	return scanCommentRows(rows)
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
	return scanCommentRows(rows), total
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
			return domain.Comment{}, apperr.ErrCommentNotFound
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

func scanCommentRows(rows *sql.Rows) []domain.Comment {
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
	return items
}
