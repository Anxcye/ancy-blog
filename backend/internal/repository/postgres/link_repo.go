// File: link_repo.go
// Purpose: Implement PostgreSQL repository methods for friend-link submission and review.
// Module: backend/internal/repository/postgres, link persistence layer.
// Related: helpers.go and service link flows.
package postgres

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

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
			return domain.Link{}, apperr.ErrLinkNotFound
		}
		return domain.Link{}, err
	}
	return l, nil
}
