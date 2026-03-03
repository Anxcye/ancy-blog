// File: admin.go
// Purpose: Define admin API DTOs for content and site management payloads.
// Module: backend/internal/handler/dto, admin transport DTO layer.
// Related: handler/admin.go and service modules.
package dto

import "time"

type ArticleUpsertRequest struct {
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	ContentKind   string    `json:"contentKind"`
	Summary       string    `json:"summary"`
	Content       string    `json:"content"`
	Status        string    `json:"status"`
	Visibility    string    `json:"visibility"`
	AllowComment  bool      `json:"allowComment"`
	OriginType    string    `json:"originType"`
	SourceURL     string    `json:"sourceUrl"`
	AIAssistLevel string    `json:"aiAssistLevel"`
	CoverImage    string    `json:"coverImage"`
	CategorySlug  string    `json:"categorySlug"`
	TagSlugs      []string  `json:"tagSlugs"`
	PublishedAt   time.Time `json:"publishedAt"`
}

type MomentCreateRequest struct {
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	AllowComment bool      `json:"allowComment"`
	PublishedAt  time.Time `json:"publishedAt"`
}

type SiteSettingsUpdateRequest struct {
	SiteName      string `json:"siteName"`
	AvatarURL     string `json:"avatarUrl"`
	HeroIntroMD   string `json:"heroIntroMd"`
	DefaultLocale string `json:"defaultLocale"`
}

type FooterItemUpsertRequest struct {
	Label               string `json:"label"`
	LinkType            string `json:"linkType"`
	InternalArticleSlug string `json:"internalArticleSlug"`
	ExternalURL         string `json:"externalUrl"`
	RowNum              int    `json:"rowNum"`
	OrderNum            int    `json:"orderNum"`
	Enabled             bool   `json:"enabled"`
}

type SocialLinkUpsertRequest struct {
	Platform string `json:"platform"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	IconKey  string `json:"iconKey"`
	OrderNum int    `json:"orderNum"`
	Enabled  bool   `json:"enabled"`
}

type NavItemUpsertRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Type        string `json:"type"`
	TargetType  string `json:"targetType"`
	TargetValue string `json:"targetValue"`
	OrderNum    int    `json:"orderNum"`
	Enabled     bool   `json:"enabled"`
}

type SlotCreateRequest struct {
	SlotKey     string `json:"slotKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type SlotItemCreateRequest struct {
	ContentType string `json:"contentType"`
	ContentID   string `json:"contentId"`
	OrderNum    int    `json:"orderNum"`
	Enabled     bool   `json:"enabled"`
}

type CommentUpdateRequest struct {
	Status   string `json:"status"`
	IsPinned string `json:"isPinned"`
}

type ReviewLinkRequest struct {
	ReviewStatus     string `json:"reviewStatus"`
	ReviewNote       string `json:"reviewNote"`
	RelatedArticleID string `json:"relatedArticleId"`
}

type IntegrationUpdateRequest struct {
	Enabled    bool           `json:"enabled"`
	ConfigJSON map[string]any `json:"configJson"`
	MetaJSON   map[string]any `json:"metaJson"`
}

type CreateTranslationJobRequest struct {
	SourceType   string `json:"sourceType"`
	SourceID     string `json:"sourceId"`
	SourceLocale string `json:"sourceLocale"`
	TargetLocale string `json:"targetLocale"`
	ProviderKey  string `json:"providerKey"`
	ModelName    string `json:"modelName"`
}

type UpsertTranslationContentRequest struct {
	SourceType        string `json:"sourceType"`
	SourceID          string `json:"sourceId"`
	Locale            string `json:"locale"`
	Content           string `json:"content"`
	TranslatedByJobID string `json:"translatedByJobId"`
}
