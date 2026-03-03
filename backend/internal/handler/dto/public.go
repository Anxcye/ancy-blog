// File: public.go
// Purpose: Define public API DTOs for visitor-facing request payloads.
// Module: backend/internal/handler/dto, public transport DTO layer.
// Related: handler/public.go and related domain mapping.
package dto

type CreateCommentRequest struct {
	ArticleID   string `json:"articleId"`
	ParentID    string `json:"parentId"`
	RootID      string `json:"rootId"`
	Content     string `json:"content"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	AvatarURL   string `json:"avatarUrl"`
	Source      string `json:"source"`
	ToCommentID string `json:"toCommentId"`
}

type SubmitLinkRequest struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatarUrl"`
	Description  string `json:"description"`
	ContactEmail string `json:"contactEmail"`
}
