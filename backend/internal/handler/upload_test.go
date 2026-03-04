// File: upload_test.go
// Purpose: Verify upload handler validation and successful upload response.
// Module: backend/internal/handler, upload HTTP test layer.
// Related: upload.go and storage uploader abstraction.
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/gin-gonic/gin"
)

type uploaderMock struct {
	uploadFunc func(ctx context.Context, objectKey string, body io.Reader, contentType string) (string, error)
}

func (m *uploaderMock) Upload(ctx context.Context, objectKey string, body io.Reader, contentType string) (string, error) {
	if m.uploadFunc != nil {
		return m.uploadFunc(ctx, objectKey, body, contentType)
	}
	return "", nil
}

func TestUploadImageNotConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUploadHandler(nil, nil)
	r := gin.New()
	r.POST("/upload", h.UploadImage)

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}

func TestUploadImageSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUploadHandler(&uploaderMock{uploadFunc: func(_ context.Context, objectKey string, _ io.Reader, contentType string) (string, error) {
		if !strings.HasPrefix(objectKey, "uploads/images/") {
			t.Fatalf("unexpected object key: %s", objectKey)
		}
		if contentType != "image/png" {
			t.Fatalf("unexpected content type: %s", contentType)
		}
		return "https://cdn.example.com/" + objectKey, nil
	}}, nil)

	r := gin.New()
	r.POST("/upload", h.UploadImage)

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Disposition", `form-data; name="file"; filename="a.png"`)
	headers.Set("Content-Type", "image/png")
	part, err := mw.CreatePart(headers)
	if err != nil {
		t.Fatalf("CreatePart failed: %v", err)
	}
	if _, err := part.Write([]byte("png-data")); err != nil {
		t.Fatalf("write file content failed: %v", err)
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if env.Code != "OK" {
		t.Fatalf("unexpected code: %s", env.Code)
	}
}
