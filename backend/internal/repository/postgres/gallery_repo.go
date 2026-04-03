// File: gallery_repo.go
// Purpose: Implement PostgreSQL repository methods for gallery photos and tags.
// Module: backend/internal/repository/postgres, gallery persistence layer.
// Related: gallery_repository.go interface, domain models, gallery service.
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

// --------------- Photos ---------------

func (r *Repository) CreatePhoto(photo domain.GalleryPhoto) (domain.GalleryPhoto, error) {
	if photo.Status == "" {
		photo.Status = "draft"
	}
	if photo.ProcessingStatus == "" {
		photo.ProcessingStatus = "pending"
	}

	var id string
	var createdAt, updatedAt time.Time
	var takenAt sql.NullTime
	err := r.db.QueryRow(`
INSERT INTO gallery_photos
  (title, slug, description, status,
   location_name, location_city, location_state, location_country,
   taken_at,
   camera_make, camera_model, lens_model, focal_length, aperture, shutter_speed, iso,
   file_size_bytes, width, height,
   taken_at_display, camera_display, location_display, exif_display, tags_display,
   placeholder_data, display_url, large_url,
   processing_status, processing_error, sort_order)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)
RETURNING id::text, taken_at, created_at, updated_at
`,
		photo.Title, photo.Slug, photo.Description, photo.Status,
		photo.LocationName, photo.LocationCity, photo.LocationState, photo.LocationCountry,
		nullableTime(photo.TakenAt),
		photo.CameraMake, photo.CameraModel, photo.LensModel, photo.FocalLength,
		photo.Aperture, photo.ShutterSpeed, photo.ISO,
		photo.FileSizeBytes, photo.Width, photo.Height,
		photo.TakenAtDisplay, photo.CameraDisplay, photo.LocationDisplay, photo.ExifDisplay, photo.TagsDisplay,
		photo.PlaceholderData, photo.DisplayURL, photo.LargeURL,
		photo.ProcessingStatus, photo.ProcessingError, photo.SortOrder,
	).Scan(&id, &takenAt, &createdAt, &updatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.GalleryPhoto{}, apperr.ErrSlugAlreadyExists
		}
		return domain.GalleryPhoto{}, err
	}
	photo.ID = id
	if takenAt.Valid {
		photo.TakenAt = takenAt.Time
	}
	photo.CreatedAt = createdAt
	photo.UpdatedAt = updatedAt

	if len(photo.TagSlugs) > 0 {
		tagIDs := r.resolveGalleryTagIDs(photo.TagSlugs)
		_ = r.setPhotoTags(id, tagIDs)
	}
	return photo, nil
}

func (r *Repository) UpdatePhoto(id string, photo domain.GalleryPhoto) (domain.GalleryPhoto, error) {
	var createdAt time.Time
	err := r.db.QueryRow(`SELECT created_at FROM gallery_photos WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GalleryPhoto{}, apperr.ErrGalleryPhotoNotFound
		}
		return domain.GalleryPhoto{}, err
	}

	var updatedAt time.Time
	var takenAt sql.NullTime
	err = r.db.QueryRow(`
UPDATE gallery_photos
SET title=$2, slug=$3, description=$4, status=$5,
    location_name=$6, location_city=$7, location_state=$8, location_country=$9,
    taken_at=$10,
    camera_make=$11, camera_model=$12, lens_model=$13, focal_length=$14,
    aperture=$15, shutter_speed=$16, iso=$17,
    file_size_bytes=$18,
    taken_at_display=$19, camera_display=$20, location_display=$21, exif_display=$22, tags_display=$23,
    sort_order=$24, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL
RETURNING taken_at, updated_at
`, id,
		photo.Title, photo.Slug, photo.Description, photo.Status,
		photo.LocationName, photo.LocationCity, photo.LocationState, photo.LocationCountry,
		nullableTime(photo.TakenAt),
		photo.CameraMake, photo.CameraModel, photo.LensModel, photo.FocalLength,
		photo.Aperture, photo.ShutterSpeed, photo.ISO,
		photo.FileSizeBytes,
		photo.TakenAtDisplay, photo.CameraDisplay, photo.LocationDisplay, photo.ExifDisplay, photo.TagsDisplay,
		photo.SortOrder,
	).Scan(&takenAt, &updatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.GalleryPhoto{}, apperr.ErrSlugAlreadyExists
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.GalleryPhoto{}, apperr.ErrGalleryPhotoNotFound
		}
		return domain.GalleryPhoto{}, err
	}
	photo.ID = id
	photo.CreatedAt = createdAt
	photo.UpdatedAt = updatedAt
	if takenAt.Valid {
		photo.TakenAt = takenAt.Time
	}

	if photo.TagSlugs != nil {
		tagIDs := r.resolveGalleryTagIDs(photo.TagSlugs)
		_ = r.setPhotoTags(id, tagIDs)
	}
	return photo, nil
}

func (r *Repository) DeletePhoto(id string) (bool, error) {
	var refCount int
	err := r.db.QueryRow(`SELECT article_ref_count FROM gallery_photos WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&refCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if refCount > 0 {
		return false, apperr.ErrPhotoHasReferences
	}
	res, err := r.db.Exec(`UPDATE gallery_photos SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *Repository) GetPhotoByID(id string) (domain.GalleryPhoto, bool) {
	return r.scanPhoto(`
SELECT `+photoColumns()+`
FROM gallery_photos
WHERE id=$1 AND deleted_at IS NULL`, id)
}

func (r *Repository) GetPhotoBySlug(slug string) (domain.GalleryPhoto, bool) {
	return r.scanPhoto(`
SELECT `+photoColumns()+`
FROM gallery_photos
WHERE slug=$1 AND deleted_at IS NULL`, slug)
}

func (r *Repository) GetPublishedPhotoBySlug(slug string) (domain.GalleryPhoto, bool) {
	return r.scanPhoto(`
SELECT `+photoColumns()+`
FROM gallery_photos
WHERE slug=$1 AND status IN ('published','hidden') AND deleted_at IS NULL`, slug)
}

func (r *Repository) ListPhotos(page, pageSize int, status, tag, keyword string) ([]domain.GalleryPhoto, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	conditions := []string{"gp.deleted_at IS NULL"}
	args := make([]any, 0, 5)
	joinTag := false

	if status != "" {
		args = append(args, status)
		conditions = append(conditions, "gp.status = $"+strconv.Itoa(len(args)))
	}
	if strings.TrimSpace(keyword) != "" {
		args = append(args, "%"+strings.TrimSpace(keyword)+"%")
		conditions = append(conditions, "(gp.title ILIKE $"+strconv.Itoa(len(args))+" OR gp.slug ILIKE $"+strconv.Itoa(len(args))+" OR gp.location_name ILIKE $"+strconv.Itoa(len(args))+")")
	}
	if tag != "" {
		joinTag = true
		args = append(args, tag)
		conditions = append(conditions, "gt.slug = $"+strconv.Itoa(len(args)))
	}

	whereClause := strings.Join(conditions, " AND ")
	fromClause := "gallery_photos gp"
	if joinTag {
		fromClause += " JOIN gallery_photo_tags gpt ON gpt.photo_id = gp.id JOIN gallery_tags gt ON gt.id = gpt.tag_id AND gt.deleted_at IS NULL"
	}

	var total int
	countQuery := "SELECT COUNT(DISTINCT gp.id) FROM " + fromClause + " WHERE " + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return []domain.GalleryPhoto{}, 0
	}

	listArgs := append(args, pageSize, offset)
	query := `
SELECT DISTINCT ON (gp.sort_order, gp.created_at, gp.id)
  ` + photoColumnsAliased("gp") + `
FROM ` + fromClause + `
WHERE ` + whereClause + `
ORDER BY gp.sort_order DESC, gp.created_at DESC, gp.id
LIMIT $` + strconv.Itoa(len(listArgs)-1) + ` OFFSET $` + strconv.Itoa(len(listArgs))

	rows, err := r.db.Query(query, listArgs...)
	if err != nil {
		return []domain.GalleryPhoto{}, total
	}
	defer rows.Close()

	items := r.scanPhotoRows(rows)
	r.loadPhotoTags(items)
	return items, total
}

func (r *Repository) ListPublishedPhotos(page, pageSize int, tag string) ([]domain.GalleryPhoto, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize

	conditions := []string{"gp.deleted_at IS NULL", "gp.status = 'published'"}
	args := make([]any, 0, 3)
	joinTag := false

	if tag != "" {
		joinTag = true
		args = append(args, tag)
		conditions = append(conditions, "gt.slug = $"+strconv.Itoa(len(args)))
	}

	whereClause := strings.Join(conditions, " AND ")
	fromClause := "gallery_photos gp"
	if joinTag {
		fromClause += " JOIN gallery_photo_tags gpt ON gpt.photo_id = gp.id JOIN gallery_tags gt ON gt.id = gpt.tag_id AND gt.deleted_at IS NULL"
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(DISTINCT gp.id) FROM "+fromClause+" WHERE "+whereClause, args...).Scan(&total); err != nil {
		return []domain.GalleryPhoto{}, 0
	}

	listArgs := append(args, pageSize, offset)
	query := `
SELECT DISTINCT ON (gp.sort_order, COALESCE(gp.taken_at, gp.created_at), gp.id)
  ` + photoColumnsAliased("gp") + `
FROM ` + fromClause + `
WHERE ` + whereClause + `
ORDER BY gp.sort_order DESC, COALESCE(gp.taken_at, gp.created_at) DESC, gp.id
LIMIT $` + strconv.Itoa(len(listArgs)-1) + ` OFFSET $` + strconv.Itoa(len(listArgs))

	rows, err := r.db.Query(query, listArgs...)
	if err != nil {
		return []domain.GalleryPhoto{}, total
	}
	defer rows.Close()

	items := r.scanPhotoRows(rows)
	r.loadPhotoTags(items)
	return items, total
}

func (r *Repository) BatchUpdatePhotoStatus(ids []string, status string) int {
	if len(ids) == 0 || strings.TrimSpace(status) == "" {
		return 0
	}
	placeholders := make([]string, len(ids))
	args := make([]any, 0, len(ids)+1)
	args = append(args, status)
	for i, id := range ids {
		args = append(args, id)
		placeholders[i] = "$" + strconv.Itoa(i+2)
	}
	query := fmt.Sprintf(`UPDATE gallery_photos SET status=$1, updated_at=NOW() WHERE id IN (%s) AND deleted_at IS NULL`, strings.Join(placeholders, ","))
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0
	}
	n, _ := res.RowsAffected()
	return int(n)
}

func (r *Repository) UpdatePhotoProcessingStatus(id, status, errorMsg string) error {
	_, err := r.db.Exec(`
UPDATE gallery_photos SET processing_status=$2, processing_error=$3, updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL`, id, status, errorMsg)
	return err
}

func (r *Repository) UpdatePhotoAssets(id, placeholderData, displayURL, largeURL string, width, height int) error {
	_, err := r.db.Exec(`
UPDATE gallery_photos SET placeholder_data=$2, display_url=$3, large_url=$4, width=$5, height=$6,
  processing_status='completed', processing_error='', updated_at=NOW()
WHERE id=$1 AND deleted_at IS NULL`, id, placeholderData, displayURL, largeURL, width, height)
	return err
}

func (r *Repository) PhotoSlugExists(slug string) bool {
	var exists bool
	_ = r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM gallery_photos WHERE slug=$1 AND deleted_at IS NULL)`, slug).Scan(&exists)
	return exists
}

// --------------- Tags ---------------

func (r *Repository) CreateGalleryTag(tag domain.GalleryTag) (domain.GalleryTag, error) {
	var id string
	err := r.db.QueryRow(`
INSERT INTO gallery_tags (name, slug) VALUES ($1, $2)
RETURNING id::text`, tag.Name, tag.Slug).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.GalleryTag{}, apperr.ErrSlugAlreadyExists
		}
		return domain.GalleryTag{}, err
	}
	tag.ID = id
	return tag, nil
}

func (r *Repository) DeleteGalleryTag(id string) bool {
	res, err := r.db.Exec(`UPDATE gallery_tags SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n > 0
}

func (r *Repository) ListGalleryTags() []domain.GalleryTag {
	rows, err := r.db.Query(`SELECT id::text, name, slug FROM gallery_tags WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return []domain.GalleryTag{}
	}
	defer rows.Close()
	tags := make([]domain.GalleryTag, 0)
	for rows.Next() {
		var t domain.GalleryTag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err == nil {
			tags = append(tags, t)
		}
	}
	return tags
}

func (r *Repository) GetGalleryTagBySlug(slug string) (domain.GalleryTag, bool) {
	var t domain.GalleryTag
	err := r.db.QueryRow(`SELECT id::text, name, slug FROM gallery_tags WHERE slug=$1 AND deleted_at IS NULL`, slug).
		Scan(&t.ID, &t.Name, &t.Slug)
	if err != nil {
		return domain.GalleryTag{}, false
	}
	return t, true
}

// --------------- Photo-Tag associations ---------------

func (r *Repository) SetPhotoTags(photoID string, tagIDs []string) error {
	return r.setPhotoTags(photoID, tagIDs)
}

func (r *Repository) setPhotoTags(photoID string, tagIDs []string) error {
	_, _ = r.db.Exec(`DELETE FROM gallery_photo_tags WHERE photo_id=$1`, photoID)
	for _, tagID := range tagIDs {
		_, _ = r.db.Exec(`INSERT INTO gallery_photo_tags (photo_id, tag_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, photoID, tagID)
	}
	return nil
}

// --------------- Helpers ---------------

func (r *Repository) resolveGalleryTagIDs(slugs []string) []string {
	if len(slugs) == 0 {
		return nil
	}
	ids := make([]string, 0, len(slugs))
	for _, s := range slugs {
		var id string
		err := r.db.QueryRow(`SELECT id::text FROM gallery_tags WHERE slug=$1 AND deleted_at IS NULL`, s).Scan(&id)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func photoColumns() string {
	return `id::text, title, slug, description, status,
  location_name, location_city, location_state, location_country,
  taken_at,
  camera_make, camera_model, lens_model, focal_length, aperture, shutter_speed, iso,
  file_size_bytes, width, height,
  taken_at_display, camera_display, location_display, exif_display, tags_display,
  placeholder_data, display_url, large_url,
  processing_status, processing_error, sort_order, article_ref_count,
  created_at, updated_at`
}

func photoColumnsAliased(alias string) string {
	cols := []string{
		"id::text", "title", "slug", "description", "status",
		"location_name", "location_city", "location_state", "location_country",
		"taken_at",
		"camera_make", "camera_model", "lens_model", "focal_length", "aperture", "shutter_speed", "iso",
		"file_size_bytes", "width", "height",
		"taken_at_display", "camera_display", "location_display", "exif_display", "tags_display",
		"placeholder_data", "display_url", "large_url",
		"processing_status", "processing_error", "sort_order", "article_ref_count",
		"created_at", "updated_at",
	}
	for i, c := range cols {
		cols[i] = alias + "." + c
	}
	return strings.Join(cols, ", ")
}

func (r *Repository) scanPhoto(query string, args ...any) (domain.GalleryPhoto, bool) {
	var p domain.GalleryPhoto
	var takenAt sql.NullTime
	err := r.db.QueryRow(query, args...).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Description, &p.Status,
		&p.LocationName, &p.LocationCity, &p.LocationState, &p.LocationCountry,
		&takenAt,
		&p.CameraMake, &p.CameraModel, &p.LensModel, &p.FocalLength, &p.Aperture, &p.ShutterSpeed, &p.ISO,
		&p.FileSizeBytes, &p.Width, &p.Height,
		&p.TakenAtDisplay, &p.CameraDisplay, &p.LocationDisplay, &p.ExifDisplay, &p.TagsDisplay,
		&p.PlaceholderData, &p.DisplayURL, &p.LargeURL,
		&p.ProcessingStatus, &p.ProcessingError, &p.SortOrder, &p.ArticleRefCount,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return domain.GalleryPhoto{}, false
	}
	if takenAt.Valid {
		p.TakenAt = takenAt.Time
	}
	r.loadPhotoTags([]domain.GalleryPhoto{p})
	if len(p.TagSlugs) == 0 {
		// reload from slice
	}
	items := []domain.GalleryPhoto{p}
	r.loadPhotoTags(items)
	p = items[0]
	return p, true
}

func (r *Repository) scanPhotoRows(rows *sql.Rows) []domain.GalleryPhoto {
	items := make([]domain.GalleryPhoto, 0)
	for rows.Next() {
		var p domain.GalleryPhoto
		var takenAt sql.NullTime
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Slug, &p.Description, &p.Status,
			&p.LocationName, &p.LocationCity, &p.LocationState, &p.LocationCountry,
			&takenAt,
			&p.CameraMake, &p.CameraModel, &p.LensModel, &p.FocalLength, &p.Aperture, &p.ShutterSpeed, &p.ISO,
			&p.FileSizeBytes, &p.Width, &p.Height,
			&p.TakenAtDisplay, &p.CameraDisplay, &p.LocationDisplay, &p.ExifDisplay, &p.TagsDisplay,
			&p.PlaceholderData, &p.DisplayURL, &p.LargeURL,
			&p.ProcessingStatus, &p.ProcessingError, &p.SortOrder, &p.ArticleRefCount,
			&p.CreatedAt, &p.UpdatedAt,
		); err == nil {
			if takenAt.Valid {
				p.TakenAt = takenAt.Time
			}
			items = append(items, p)
		}
	}
	return items
}

func (r *Repository) loadPhotoTags(photos []domain.GalleryPhoto) {
	if len(photos) == 0 {
		return
	}
	ids := make([]string, len(photos))
	idxMap := make(map[string]int, len(photos))
	for i, p := range photos {
		ids[i] = p.ID
		idxMap[p.ID] = i
		photos[i].TagSlugs = []string{}
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	query := fmt.Sprintf(`
SELECT gpt.photo_id::text, gt.slug
FROM gallery_photo_tags gpt
JOIN gallery_tags gt ON gt.id = gpt.tag_id AND gt.deleted_at IS NULL
WHERE gpt.photo_id IN (%s)
ORDER BY gt.name`, strings.Join(placeholders, ","))
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var photoID, tagSlug string
		if err := rows.Scan(&photoID, &tagSlug); err == nil {
			if idx, ok := idxMap[photoID]; ok {
				photos[idx].TagSlugs = append(photos[idx].TagSlugs, tagSlug)
			}
		}
	}
}
