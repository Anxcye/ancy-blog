// File: errors.go
// Purpose: Provide handler-level helpers for mapping typed service errors to HTTP responses.
// Module: backend/internal/handler, transport error mapping layer.
// Related: internal/apperr and response envelope.
package handler

import (
	"errors"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
)

func isValidationError(err error) bool {
	return errors.Is(err, apperr.ErrValidation)
}
