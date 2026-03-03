// File: timeline_repo.go
// Purpose: Implement PostgreSQL repository methods for mixed timeline querying.
// Module: backend/internal/repository/postgres, timeline persistence layer.
// Related: service timeline flows.
package postgres

import (
	"sort"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

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
