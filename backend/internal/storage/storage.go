// File: storage.go
// Purpose: Define storage abstraction for uploading binary objects.
// Module: backend/internal/storage, infrastructure abstraction layer.
// Related: R2 storage implementation and admin upload handlers.
package storage

import (
	"context"
	"io"
)

type Uploader interface {
	Upload(ctx context.Context, objectKey string, body io.Reader, contentType string) (string, error)
}
