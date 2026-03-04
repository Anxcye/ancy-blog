// File: upload.go
// Purpose: Provide admin upload endpoint for image files via configured object storage.
// Module: backend/internal/handler, upload HTTP presentation layer.
// Related: storage uploader implementation and admin route group.
package handler

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/anxcye/ancy-blog/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploader           storage.Uploader
	integrationService *service.IntegrationService
}

func NewUploadHandler(uploader storage.Uploader, integrationService *service.IntegrationService) *UploadHandler {
	return &UploadHandler{uploader: uploader, integrationService: integrationService}
}

func (h *UploadHandler) resolveUploader() (storage.Uploader, error) {
	if h.uploader != nil {
		return h.uploader, nil
	}
	if h.integrationService == nil {
		return nil, fmt.Errorf("image storage is not configured")
	}
	provider, ok := h.integrationService.GetIntegrationProviderForRuntime("cloudflare_r2")
	if !ok || !provider.Enabled {
		return nil, fmt.Errorf("cloudflare_r2 provider is disabled or missing")
	}
	uploader, err := storage.NewR2UploaderFromConfigJSON(provider.ConfigJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid cloudflare_r2 configuration: %w", err)
	}
	return uploader, nil
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	uploader, resolveErr := h.resolveUploader()
	if resolveErr != nil {
		response.JSON(c, http.StatusNotImplemented, response.Envelope{Code: "UPLOAD_NOT_CONFIGURED", Message: resolveErr.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", "file is required")
		return
	}

	if file.Size > 10*1024*1024 {
		badRequest(c, "VALIDATION_ERROR", "file size must be <= 10MB")
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
		ext = ".bin"
	}
	objectKey := fmt.Sprintf("uploads/images/%s/%s%s", time.Now().UTC().Format("200601"), uuid.NewString(), ext)
	url, err := uploader.Upload(c.Request.Context(), objectKey, src, contentType)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "UPLOAD_FAILED", Message: err.Error()})
		return
	}

	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"key": objectKey, "url": url}})
}
