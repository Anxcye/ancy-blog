// File: gallery.go
// Purpose: Implement admin and public HTTP endpoints for gallery photos and tags.
// Module: backend/internal/handler, gallery HTTP presentation layer.
// Related: gallery service, gallery DTOs, and route registration.
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/handler/dto"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GalleryHandler serves both admin and public gallery endpoints.
type GalleryHandler struct {
	galleryService *service.GalleryService
}

// NewGalleryHandler creates a gallery handler.
func NewGalleryHandler(galleryService *service.GalleryService) *GalleryHandler {
	return &GalleryHandler{galleryService: galleryService}
}

// ============ Admin Endpoints ============

// CreatePhoto handles POST /admin/gallery/photos
func (h *GalleryHandler) CreatePhoto(c *gin.Context) {
	var req dto.GalleryPhotoUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	photo, err := h.galleryService.CreatePhoto(domain.GalleryPhoto{
		Title:           req.Title,
		Slug:            req.Slug,
		Description:     req.Description,
		Status:          req.Status,
		LocationName:    req.LocationName,
		LocationCity:    req.LocationCity,
		LocationState:   req.LocationState,
		LocationCountry: req.LocationCountry,
		TakenAt:         req.TakenAt,
		CameraMake:      req.CameraMake,
		CameraModel:     req.CameraModel,
		LensModel:       req.LensModel,
		FocalLength:     req.FocalLength,
		Aperture:        req.Aperture,
		ShutterSpeed:    req.ShutterSpeed,
		ISO:             req.ISO,
		TakenAtDisplay:  req.TakenAtDisplay,
		CameraDisplay:   req.CameraDisplay,
		LocationDisplay: req.LocationDisplay,
		ExifDisplay:     req.ExifDisplay,
		TagsDisplay:     req.TagsDisplay,
		SortOrder:       req.SortOrder,
		TagSlugs:        req.TagSlugs,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrSlugAlreadyExists) {
			response.JSON(c, http.StatusConflict, response.Envelope{Code: "SLUG_CONFLICT", Message: err.Error()})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: photo})
}

// UpdatePhoto handles PUT /admin/gallery/photos/:id
func (h *GalleryHandler) UpdatePhoto(c *gin.Context) {
	id := c.Param("id")
	var req dto.GalleryPhotoUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	photo, err := h.galleryService.UpdatePhoto(id, domain.GalleryPhoto{
		Title:           req.Title,
		Slug:            req.Slug,
		Description:     req.Description,
		Status:          req.Status,
		LocationName:    req.LocationName,
		LocationCity:    req.LocationCity,
		LocationState:   req.LocationState,
		LocationCountry: req.LocationCountry,
		TakenAt:         req.TakenAt,
		CameraMake:      req.CameraMake,
		CameraModel:     req.CameraModel,
		LensModel:       req.LensModel,
		FocalLength:     req.FocalLength,
		Aperture:        req.Aperture,
		ShutterSpeed:    req.ShutterSpeed,
		ISO:             req.ISO,
		TakenAtDisplay:  req.TakenAtDisplay,
		CameraDisplay:   req.CameraDisplay,
		LocationDisplay: req.LocationDisplay,
		ExifDisplay:     req.ExifDisplay,
		TagsDisplay:     req.TagsDisplay,
		SortOrder:       req.SortOrder,
		TagSlugs:        req.TagSlugs,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrGalleryPhotoNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PHOTO_NOT_FOUND", Message: err.Error()})
			return
		}
		if errors.Is(err, apperr.ErrSlugAlreadyExists) {
			response.JSON(c, http.StatusConflict, response.Envelope{Code: "SLUG_CONFLICT", Message: err.Error()})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: photo})
}

// DeletePhoto handles DELETE /admin/gallery/photos/:id
func (h *GalleryHandler) DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	ok, err := h.galleryService.DeletePhoto(id)
	if err != nil {
		if errors.Is(err, apperr.ErrPhotoHasReferences) {
			response.JSON(c, http.StatusConflict, response.Envelope{Code: "PHOTO_HAS_REFERENCES", Message: err.Error()})
			return
		}
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "DELETE_FAILED", Message: err.Error()})
		return
	}
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PHOTO_NOT_FOUND", Message: "photo not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "deleted"})
}

// PhotoDetail handles GET /admin/gallery/photos/:id
func (h *GalleryHandler) PhotoDetail(c *gin.Context) {
	id := c.Param("id")
	photo, ok := h.galleryService.GetPhotoByID(id)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PHOTO_NOT_FOUND", Message: "photo not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: photo})
}

// ListPhotos handles GET /admin/gallery/photos
func (h *GalleryHandler) ListPhotos(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 20)
	status := c.Query("status")
	tag := c.Query("tag")
	keyword := c.Query("keyword")
	photos, total := h.galleryService.ListPhotos(page, pageSize, status, tag, keyword)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.GalleryPhoto]{Total: total, Rows: photos}})
}

// BatchUpdatePhotoStatus handles POST /admin/gallery/photos/batch-status
func (h *GalleryHandler) BatchUpdatePhotoStatus(c *gin.Context) {
	var req dto.GalleryPhotoBatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	n, err := h.galleryService.BatchUpdatePhotoStatus(req.IDs, req.Status)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]int{"updated": n}})
}

// UploadPhoto handles POST /admin/gallery/photos/upload
// Accepts multipart file, creates a draft photo record and triggers processing.
func (h *GalleryHandler) UploadPhoto(c *gin.Context) {
	uploader, err := h.galleryService.GetUploader()
	if err != nil {
		response.JSON(c, http.StatusNotImplemented, response.Envelope{Code: "UPLOAD_NOT_CONFIGURED", Message: err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", "file is required")
		return
	}
	if file.Size > 50*1024*1024 {
		badRequest(c, "VALIDATION_ERROR", "file size must be <= 50MB")
		return
	}
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		badRequest(c, "VALIDATION_ERROR", "only image uploads are supported")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "UPLOAD_IO_ERROR", Message: "failed to open uploaded file"})
		return
	}
	defer src.Close()

	ext := path.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	photoID := uuid.NewString()
	objectKey := fmt.Sprintf("gallery/large/%s/%s%s", time.Now().UTC().Format("200601"), photoID, ext)

	url, err := uploader.Upload(c.Request.Context(), objectKey, src, contentType)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "UPLOAD_FAILED", Message: err.Error()})
		return
	}

	// Generate a slug from filename
	baseName := strings.TrimSuffix(file.Filename, ext)
	slug := slugify(baseName) + "-" + photoID[:8]

	photo, createErr := h.galleryService.CreatePhoto(domain.GalleryPhoto{
		Slug:             slug,
		Status:           "draft",
		LargeURL:         url,
		ProcessingStatus: "pending",
	})
	if createErr != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "CREATE_FAILED", Message: createErr.Error()})
		return
	}

	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: photo})
}

// CreateGalleryTag handles POST /admin/gallery/tags
func (h *GalleryHandler) CreateGalleryTag(c *gin.Context) {
	var req dto.GalleryTagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	tag, err := h.galleryService.CreateGalleryTag(domain.GalleryTag{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrSlugAlreadyExists) {
			response.JSON(c, http.StatusConflict, response.Envelope{Code: "SLUG_CONFLICT", Message: err.Error()})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: tag})
}

// DeleteGalleryTag handles DELETE /admin/gallery/tags/:id
func (h *GalleryHandler) DeleteGalleryTag(c *gin.Context) {
	id := c.Param("id")
	if !h.galleryService.DeleteGalleryTag(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "TAG_NOT_FOUND", Message: "gallery tag not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "deleted"})
}

// ListGalleryTags handles GET /admin/gallery/tags
func (h *GalleryHandler) ListGalleryTags(c *gin.Context) {
	tags := h.galleryService.ListGalleryTags()
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: tags})
}

// ============ Public Endpoints ============

// PublicPhotos handles GET /public/gallery/photos
func (h *GalleryHandler) PublicPhotos(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 20)
	tag := c.Query("tag")
	photos, total := h.galleryService.ListPublishedPhotos(page, pageSize, tag)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.GalleryPhotoPublic]{Total: total, Rows: photos}})
}

// PublicPhotoBySlug handles GET /public/gallery/photos/:slug
func (h *GalleryHandler) PublicPhotoBySlug(c *gin.Context) {
	slug := c.Param("slug")
	photo, ok := h.galleryService.GetPublishedPhotoBySlug(slug)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PHOTO_NOT_FOUND", Message: "photo not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: photo})
}

// PublicGalleryTags handles GET /public/gallery/tags
func (h *GalleryHandler) PublicGalleryTags(c *gin.Context) {
	tags := h.galleryService.ListGalleryTags()
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: tags})
}

// ============ Helpers ============

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		if r == ' ' || r == '.' {
			return '-'
		}
		return -1
	}, s)
	// collapse multiple dashes
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	if s == "" {
		s = "photo"
	}
	return s
}
