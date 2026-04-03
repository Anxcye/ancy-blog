// File: errors.go
// Purpose: Define typed application errors for stable cross-layer error handling.
// Module: backend/internal/apperr, shared error taxonomy layer.
// Related: repository, service, and handler error mapping.
package apperr

import "errors"

var (
	ErrValidation = errors.New("validation error")

	ErrArticleNotFound        = errors.New("article not found")
	ErrMomentNotFound         = errors.New("moment not found")
	ErrCommentNotFound        = errors.New("comment not found")
	ErrLinkNotFound           = errors.New("link not found")
	ErrLinkSubmissionDisabled = errors.New("link submission disabled")
	ErrFooterItemNotFound     = errors.New("footer item not found")
	ErrSocialLinkNotFound     = errors.New("social link not found")
	ErrNavItemNotFound        = errors.New("nav item not found")
	ErrSlotNotFound           = errors.New("slot not found")
	ErrSlotItemNotFound       = errors.New("slot item not found")
	ErrProviderNotFound       = errors.New("provider not found")
	ErrTranslationJobNotFound = errors.New("translation job not found")

	ErrSlugAlreadyExists = errors.New("slug already exists")

	ErrGalleryPhotoNotFound = errors.New("gallery photo not found")
	ErrGalleryTagNotFound   = errors.New("gallery tag not found")
	ErrPhotoHasReferences   = errors.New("photo is referenced by articles and cannot be deleted")
)
