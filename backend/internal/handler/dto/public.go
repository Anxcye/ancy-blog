// File: public.go
// Purpose: Define public API DTOs for visitor-facing request payloads.
// Module: backend/internal/handler/dto, public transport DTO layer.
// Related: handler/public.go and related domain mapping.
package dto

import "time"

type CreateCommentRequest struct {
	ArticleID   string `json:"articleId"`
	ContentType string `json:"contentType"`
	ContentID   string `json:"contentId"`
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

type PublicComment struct {
	ID                string          `json:"id"`
	ArticleID         string          `json:"articleId"`
	ContentType       string          `json:"contentType"`
	ContentID         string          `json:"contentId"`
	ParentID          string          `json:"parentId,omitempty"`
	RootID            string          `json:"rootId,omitempty"`
	Content           string          `json:"content"`
	Status            string          `json:"status"`
	IsPinned          bool            `json:"isPinned"`
	IsAuthor          bool            `json:"isAuthor"`
	LikeCount         int             `json:"likeCount"`
	ReplyCount        int             `json:"replyCount"`
	Nickname          string          `json:"nickname"`
	Website           string          `json:"website,omitempty"`
	AvatarURL         string          `json:"avatarUrl,omitempty"`
	ToCommentID       string          `json:"toCommentId,omitempty"`
	ToCommentNickname string          `json:"toCommentNickname,omitempty"`
	ToCommentIsAuthor bool            `json:"toCommentIsAuthor"`
	CreatedAt         time.Time       `json:"createdAt"`
	Children          []PublicComment `json:"children,omitempty"`
}

type SubmitLinkRequest struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatarUrl"`
	Description  string `json:"description"`
	ContactEmail string `json:"contactEmail"`
}
