// File: gallery_repository.go
// Purpose: Declare persistence contracts for gallery photo and tag domains.
// Module: backend/internal/repository, repository abstraction layer.
// Related: gallery Postgres implementation, gallery service orchestration.
package repository

import (
	"github.com/anxcye/ancy-blog/backend/internal/domain"
)

// GalleryRepository defines data access operations for gallery photos and tags.
type GalleryRepository interface {
	// Photos
	CreatePhoto(photo domain.GalleryPhoto) (domain.GalleryPhoto, error)
	UpdatePhoto(id string, photo domain.GalleryPhoto) (domain.GalleryPhoto, error)
	DeletePhoto(id string) (bool, error)
	GetPhotoByID(id string) (domain.GalleryPhoto, bool)
	GetPhotoBySlug(slug string) (domain.GalleryPhoto, bool)
	ListPhotos(page, pageSize int, status, tag, keyword string) ([]domain.GalleryPhoto, int)
	ListPublishedPhotos(page, pageSize int, tag string) ([]domain.GalleryPhoto, int)
	GetPublishedPhotoBySlug(slug string) (domain.GalleryPhoto, bool)
	BatchUpdatePhotoStatus(ids []string, status string) int
	UpdatePhotoProcessingStatus(id, status, errorMsg string) error
	UpdatePhotoAssets(id, placeholderData, displayURL, largeURL string, width, height int) error
	PhotoSlugExists(slug string) bool

	// Tags
	CreateGalleryTag(tag domain.GalleryTag) (domain.GalleryTag, error)
	DeleteGalleryTag(id string) bool
	ListGalleryTags() []domain.GalleryTag
	GetGalleryTagBySlug(slug string) (domain.GalleryTag, bool)

	// Photo-Tag associations
	SetPhotoTags(photoID string, tagIDs []string) error
}
