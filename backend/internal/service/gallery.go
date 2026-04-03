// File: gallery.go
// Purpose: Orchestrate gallery photo and tag business operations including validation and public projection.
// Module: backend/internal/service, gallery service layer.
// Related: gallery repository, gallery handler, image processing worker.
package service

import (
	"fmt"
	"strings"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
	"github.com/anxcye/ancy-blog/backend/internal/storage"
)

// GalleryService handles gallery photo lifecycle and tag management.
type GalleryService struct {
	repo       repository.GalleryRepository
	uploader   func() (storage.Uploader, error)
}

// NewGalleryService creates a gallery service with the given repository.
func NewGalleryService(repo repository.GalleryRepository) *GalleryService {
	return &GalleryService{repo: repo}
}

// WithUploaderFactory injects a lazy uploader resolver (R2) for image asset uploads.
func (s *GalleryService) WithUploaderFactory(factory func() (storage.Uploader, error)) *GalleryService {
	s.uploader = factory
	return s
}

// --------------- Photos ---------------

func (s *GalleryService) CreatePhoto(photo domain.GalleryPhoto) (domain.GalleryPhoto, error) {
	if strings.TrimSpace(photo.Slug) == "" {
		return domain.GalleryPhoto{}, fmt.Errorf("%w: slug is required", apperr.ErrValidation)
	}
	if photo.Status == "" {
		photo.Status = "draft"
	}
	validStatuses := map[string]bool{"draft": true, "published": true, "hidden": true}
	if !validStatuses[photo.Status] {
		return domain.GalleryPhoto{}, fmt.Errorf("%w: invalid status %q", apperr.ErrValidation, photo.Status)
	}
	return s.repo.CreatePhoto(photo)
}

func (s *GalleryService) UpdatePhoto(id string, photo domain.GalleryPhoto) (domain.GalleryPhoto, error) {
	if strings.TrimSpace(id) == "" {
		return domain.GalleryPhoto{}, fmt.Errorf("%w: photo id is required", apperr.ErrValidation)
	}
	if strings.TrimSpace(photo.Slug) == "" {
		return domain.GalleryPhoto{}, fmt.Errorf("%w: slug is required", apperr.ErrValidation)
	}
	validStatuses := map[string]bool{"draft": true, "published": true, "hidden": true}
	if photo.Status != "" && !validStatuses[photo.Status] {
		return domain.GalleryPhoto{}, fmt.Errorf("%w: invalid status %q", apperr.ErrValidation, photo.Status)
	}
	return s.repo.UpdatePhoto(id, photo)
}

func (s *GalleryService) DeletePhoto(id string) (bool, error) {
	return s.repo.DeletePhoto(id)
}

func (s *GalleryService) GetPhotoByID(id string) (domain.GalleryPhoto, bool) {
	return s.repo.GetPhotoByID(id)
}

func (s *GalleryService) ListPhotos(page, pageSize int, status, tag, keyword string) ([]domain.GalleryPhoto, int) {
	return s.repo.ListPhotos(page, pageSize, status, tag, keyword)
}

func (s *GalleryService) BatchUpdatePhotoStatus(ids []string, status string) (int, error) {
	validStatuses := map[string]bool{"draft": true, "published": true, "hidden": true}
	if !validStatuses[status] {
		return 0, fmt.Errorf("%w: invalid status %q", apperr.ErrValidation, status)
	}
	return s.repo.BatchUpdatePhotoStatus(ids, status), nil
}

func (s *GalleryService) PhotoSlugExists(slug string) bool {
	return s.repo.PhotoSlugExists(slug)
}

func (s *GalleryService) UpdatePhotoProcessingStatus(id, status, errorMsg string) error {
	return s.repo.UpdatePhotoProcessingStatus(id, status, errorMsg)
}

func (s *GalleryService) UpdatePhotoAssets(id, placeholderData, displayURL, largeURL string, width, height int) error {
	return s.repo.UpdatePhotoAssets(id, placeholderData, displayURL, largeURL, width, height)
}

// --------------- Public ---------------

func (s *GalleryService) ListPublishedPhotos(page, pageSize int, tag string) ([]domain.GalleryPhotoPublic, int) {
	photos, total := s.repo.ListPublishedPhotos(page, pageSize, tag)
	result := make([]domain.GalleryPhotoPublic, 0, len(photos))
	for _, p := range photos {
		result = append(result, toPublicPhoto(p))
	}
	return result, total
}

func (s *GalleryService) GetPublishedPhotoBySlug(slug string) (domain.GalleryPhotoPublic, bool) {
	p, ok := s.repo.GetPublishedPhotoBySlug(slug)
	if !ok {
		return domain.GalleryPhotoPublic{}, false
	}
	return toPublicPhoto(p), true
}

// toPublicPhoto applies display-switch rules: only include metadata fields that have values AND enabled display switches.
func toPublicPhoto(p domain.GalleryPhoto) domain.GalleryPhotoPublic {
	pub := domain.GalleryPhotoPublic{
		ID:              p.ID,
		Title:           p.Title,
		Slug:            p.Slug,
		Description:     p.Description,
		FileSizeBytes:   p.FileSizeBytes,
		Width:           p.Width,
		Height:          p.Height,
		PlaceholderData: p.PlaceholderData,
		DisplayURL:      p.DisplayURL,
		LargeURL:        p.LargeURL,
		CreatedAt:       p.CreatedAt,
	}
	if p.TakenAtDisplay && !p.TakenAt.IsZero() {
		pub.TakenAt = p.TakenAt
	}
	if p.LocationDisplay {
		pub.LocationName = p.LocationName
		pub.LocationCity = p.LocationCity
		pub.LocationCountry = p.LocationCountry
	}
	if p.CameraDisplay {
		pub.CameraMake = p.CameraMake
		pub.CameraModel = p.CameraModel
	}
	if p.ExifDisplay {
		pub.LensModel = p.LensModel
		pub.FocalLength = p.FocalLength
		pub.Aperture = p.Aperture
		pub.ShutterSpeed = p.ShutterSpeed
		pub.ISO = p.ISO
	}
	if p.TagsDisplay {
		pub.TagSlugs = p.TagSlugs
	}
	return pub
}

// --------------- Tags ---------------

func (s *GalleryService) CreateGalleryTag(tag domain.GalleryTag) (domain.GalleryTag, error) {
	if strings.TrimSpace(tag.Name) == "" || strings.TrimSpace(tag.Slug) == "" {
		return domain.GalleryTag{}, fmt.Errorf("%w: name and slug are required", apperr.ErrValidation)
	}
	return s.repo.CreateGalleryTag(tag)
}

func (s *GalleryService) DeleteGalleryTag(id string) bool {
	return s.repo.DeleteGalleryTag(id)
}

func (s *GalleryService) ListGalleryTags() []domain.GalleryTag {
	return s.repo.ListGalleryTags()
}

// GetUploader resolves the storage uploader for image asset operations.
func (s *GalleryService) GetUploader() (storage.Uploader, error) {
	if s.uploader == nil {
		return nil, fmt.Errorf("image storage is not configured")
	}
	return s.uploader()
}
