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

const commentSelectColumns = `
SELECT id::text, COALESCE(article_id::text,''), content_type, content_id::text, COALESCE(parent_id::text,''), COALESCE(root_id::text,''),
       content, status, is_pinned, like_count, reply_count, nickname, COALESCE(email,''), COALESCE(website,''),
       COALESCE(avatar_url,''), source, COALESCE(ip,''), COALESCE(user_agent,''), risk_score, approved_at, approved_by::text, created_at, updated_at
FROM comments`

func (r *Repository) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if comment.Status == "" {
		comment.Status = "approved"
	}
	if comment.Source == "" {
		comment.Source = "web"
	}
	if comment.ContentType == "" && comment.ArticleID != "" {
		comment.ContentType = "article"
	}
	if comment.ContentID == "" && comment.ArticleID != "" {
		comment.ContentID = comment.ArticleID
	}
	if comment.ContentType == "article" {
		comment.ArticleID = comment.ContentID
	}
	if comment.Status == "approved" && comment.ApprovedAt.IsZero() {
		comment.ApprovedAt = time.Now().UTC()
	}

	var id string
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(`
INSERT INTO comments (article_id, content_type, content_id, parent_id, root_id, content, status, is_pinned, nickname, email, website, avatar_url, source, ip, user_agent, approved_at, approved_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
RETURNING id::text, created_at, updated_at
`, nullableUUID(comment.ArticleID), comment.ContentType, nullableUUID(comment.ContentID), nullableUUID(comment.ParentID), nullableUUID(comment.RootID), comment.Content, comment.Status, false, comment.Nickname,
		nullableString(comment.Email), nullableString(comment.Website), nullableString(comment.AvatarURL), comment.Source, nullableString(comment.IP), nullableString(comment.UserAgent), nullableTime(comment.ApprovedAt), nullableUUID(comment.ApprovedBy)).
		Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return domain.Comment{}, err
	}
	comment.ID = id
	comment.IsPinned = "0"
	comment.CreatedAt = createdAt
	comment.UpdatedAt = updatedAt
	r.bumpReplyCounters(comment.ParentID, comment.RootID)
	return comment, nil
}

func (r *Repository) ListArticleComments(articleID string, page, pageSize int) ([]domain.Comment, int) {
	return r.ListContentComments("article", articleID, page, pageSize)
}

func (r *Repository) ListContentComments(contentType, contentID string, page, pageSize int) ([]domain.Comment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	var total int
	_ = r.db.QueryRow(`
SELECT COUNT(*) FROM comments
WHERE content_type=$1 AND content_id=$2 AND status='approved' AND deleted_at IS NULL AND parent_id IS NULL
`, contentType, contentID).Scan(&total)

	rows, err := r.db.Query(commentSelectColumns+`
WHERE content_type=$1 AND content_id=$2 AND status='approved' AND deleted_at IS NULL AND parent_id IS NULL
ORDER BY is_pinned DESC, created_at DESC
LIMIT $3 OFFSET $4
`, contentType, contentID, pageSize, offset)
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

	rows, err := r.db.Query(commentSelectColumns+`
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

	query := commentSelectColumns + fmt.Sprintf(`
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
	return r.CountContentComments("article", articleID)
}

func (r *Repository) CountContentComments(contentType, contentID string) (int, error) {
	var total int
	err := r.db.QueryRow(`
SELECT COUNT(*) FROM comments
WHERE content_type=$1 AND content_id=$2 AND status='approved' AND deleted_at IS NULL
`, contentType, contentID).Scan(&total)
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

	query := commentSelectColumns + `
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

func (r *Repository) GetCommentByID(id string) (domain.Comment, bool) {
	row := r.db.QueryRow(commentSelectColumns+`
WHERE id=$1 AND deleted_at IS NULL
`, id)
	comment, err := scanComment(row)
	if err != nil {
		return domain.Comment{}, false
	}
	return comment, true
}

func (r *Repository) UpdateCommentAdmin(id string, status, isPinned string) (domain.Comment, error) {
	pinned := isPinned == "1" || strings.EqualFold(isPinned, "true")
	row := r.db.QueryRow(`
UPDATE comments
SET status=$2,
    is_pinned=$3,
    approved_at=CASE
      WHEN $2='approved' THEN COALESCE(approved_at, NOW())
      ELSE NULL
    END,
    updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
RETURNING id::text, COALESCE(article_id::text,''), content_type, content_id::text, COALESCE(parent_id::text,''), COALESCE(root_id::text,''),
          content, status, is_pinned, like_count, reply_count, nickname, COALESCE(email,''), COALESCE(website,''),
          COALESCE(avatar_url,''), source, COALESCE(ip,''), COALESCE(user_agent,''), risk_score, approved_at, approved_by::text, created_at, updated_at
`, id, status, pinned)
	c, err := scanComment(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Comment{}, apperr.ErrCommentNotFound
		}
		return domain.Comment{}, err
	}
	return c, nil
}

func (r *Repository) bumpReplyCounters(parentID, rootID string) {
	if strings.TrimSpace(parentID) == "" {
		return
	}
	_, _ = r.db.Exec(`UPDATE comments SET reply_count = reply_count + 1, updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, parentID)
	if strings.TrimSpace(rootID) != "" && rootID != parentID {
		_, _ = r.db.Exec(`UPDATE comments SET reply_count = reply_count + 1, updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, rootID)
	}
}

type commentRowScanner interface {
	Scan(dest ...any) error
}

func scanComment(scanner commentRowScanner) (domain.Comment, error) {
	var c domain.Comment
	var pinned bool
	var approvedAt sql.NullTime
	var approvedBy sql.NullString
	err := scanner.Scan(&c.ID, &c.ArticleID, &c.ContentType, &c.ContentID, &c.ParentID, &c.RootID, &c.Content, &c.Status, &pinned, &c.LikeCount, &c.ReplyCount,
		&c.Nickname, &c.Email, &c.Website, &c.AvatarURL, &c.Source, &c.IP, &c.UserAgent, &c.RiskScore, &approvedAt, &approvedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return domain.Comment{}, err
	}
	if pinned {
		c.IsPinned = "1"
	} else {
		c.IsPinned = "0"
	}
	if approvedAt.Valid {
		c.ApprovedAt = approvedAt.Time
	}
	if approvedBy.Valid {
		c.ApprovedBy = approvedBy.String
	}
	return c, nil
}

func scanCommentRows(rows *sql.Rows) []domain.Comment {
	defer rows.Close()

	items := make([]domain.Comment, 0)
	for rows.Next() {
		if c, err := scanComment(rows); err == nil {
			items = append(items, c)
		}
	}
	return items
}
