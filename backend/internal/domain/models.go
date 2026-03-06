// File: models.go
// Purpose: Define core domain entities and lightweight view models used across services.
// Module: backend/internal/domain, domain model layer.
// Related: repository interfaces, service logic, and HTTP handlers.
package domain

import (
	"encoding/json"
	"time"
)

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	IsAdmin     bool   `json:"isAdmin"`
}

type Article struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	ContentKind   string    `json:"contentKind"`
	Summary       string    `json:"summary"`
	Content       string    `json:"content"`
	Status        string    `json:"status"`
	Visibility    string    `json:"visibility"`
	AllowComment  bool      `json:"allowComment"`
	IsPinned      bool      `json:"isPinned"`
	IsFeatured    bool      `json:"isFeatured"`
	OriginType    string    `json:"originType"`
	SourceURL     string    `json:"sourceUrl,omitempty"`
	AIAssistLevel string    `json:"aiAssistLevel"`
	CoverImage    string    `json:"coverImage,omitempty"`
	CategorySlug  string    `json:"categorySlug,omitempty"`
	TagSlugs      []string  `json:"tagSlugs,omitempty"`
	ViewCount     int64     `json:"viewCount"`
	PublishedAt   time.Time `json:"publishedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Moment struct {
	ID           string    `json:"id"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	AllowComment bool      `json:"allowComment"`
	CommentCount int       `json:"commentCount"`
	PublishedAt  time.Time `json:"publishedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Comment struct {
	ID                string    `json:"id"`
	ArticleID         string    `json:"articleId,omitempty"`
	ContentType       string    `json:"contentType"`
	ContentID         string    `json:"contentId"`
	ParentID          string    `json:"parentId,omitempty"`
	RootID            string    `json:"rootId,omitempty"`
	Content           string    `json:"content"`
	Status            string    `json:"status"`
	IsPinned          string    `json:"isPinned"`
	LikeCount         int       `json:"likeCount"`
	ReplyCount        int       `json:"replyCount"`
	Nickname          string    `json:"nickname"`
	Email             string    `json:"email,omitempty"`
	Website           string    `json:"website,omitempty"`
	AvatarURL         string    `json:"avatarUrl,omitempty"`
	Source            string    `json:"source"`
	IP                string    `json:"ip,omitempty"`
	UserAgent         string    `json:"userAgent,omitempty"`
	ToCommentID       string    `json:"toCommentId,omitempty"`
	ToCommentNickname string    `json:"toCommentNickname,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CommentNode struct {
	Comment
	Children []CommentNode `json:"children,omitempty"`
}

type Link struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	URL                string    `json:"url"`
	AvatarURL          string    `json:"avatarUrl,omitempty"`
	Description        string    `json:"description,omitempty"`
	ContactEmail       string    `json:"contactEmail,omitempty"`
	ReviewStatus       string    `json:"reviewStatus"`
	ReviewNote         string    `json:"reviewNote,omitempty"`
	RelatedArticleID   string    `json:"relatedArticleId,omitempty"`
	SubmittedIP        string    `json:"submittedIp,omitempty"`
	SubmittedUserAgent string    `json:"submittedUserAgent,omitempty"`
	ApprovedAt         time.Time `json:"approvedAt,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SiteSettings struct {
	SiteName               string `json:"siteName"`
	AvatarURL              string `json:"avatarUrl,omitempty"`
	HeroIntroMD            string `json:"heroIntroMd,omitempty"`
	DefaultLocale          string `json:"defaultLocale"`
	CommentEnabled         bool   `json:"commentEnabled"`
	CommentRequireApproval bool   `json:"commentRequireApproval"`
	LinkSubmissionEnabled  bool   `json:"linkSubmissionEnabled"`
	SiteDescription        string `json:"siteDescription,omitempty"`
	SeoKeywords            string `json:"seoKeywords,omitempty"`
	OgImageURL             string `json:"ogImageUrl,omitempty"`
}

// TranslationPolicy holds the global auto-translation configuration stored in site_settings.
type TranslationPolicy struct {
	Enabled       bool     `json:"enabled"`
	TargetLocales []string `json:"targetLocales"`
	ProviderKey   string   `json:"providerKey"`
	AutoPublish   bool     `json:"autoPublish"`
}

type FooterItem struct {
	ID                  string `json:"id"`
	Label               string `json:"label"`
	LinkType            string `json:"linkType"`
	InternalArticleSlug string `json:"internalArticleSlug,omitempty"`
	ExternalURL         string `json:"externalUrl,omitempty"`
	RowNum              int    `json:"rowNum"`
	OrderNum            int    `json:"orderNum"`
	Enabled             bool   `json:"enabled"`
}

type SocialLink struct {
	ID       string `json:"id"`
	Platform string `json:"platform"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	IconKey  string `json:"iconKey,omitempty"`
	OrderNum int    `json:"orderNum"`
	Enabled  bool   `json:"enabled"`
}

type NavItem struct {
	ID          string    `json:"id"`
	ParentID    string    `json:"parentId,omitempty"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Type        string    `json:"type"`
	TargetType  string    `json:"targetType"`
	TargetValue string    `json:"targetValue,omitempty"`
	OrderNum    int       `json:"orderNum"`
	Enabled     bool      `json:"enabled"`
	Children    []NavItem `json:"children,omitempty"`
}

type ContentSlot struct {
	ID          string `json:"id"`
	SlotKey     string `json:"slotKey"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled"`
}

type SlotItem struct {
	ID          string `json:"id"`
	SlotKey     string `json:"slotKey"`
	ContentType string `json:"contentType"`
	ContentID   string `json:"contentId"`
	OrderNum    int    `json:"orderNum"`
	Enabled     bool   `json:"enabled"`
}

type TimelineItem struct {
	ContentType  string    `json:"contentType"`
	ID           string    `json:"id"`
	Title        string    `json:"title,omitempty"`
	Summary      string    `json:"summary,omitempty"`
	Slug         string    `json:"slug,omitempty"`
	CategorySlug string    `json:"categorySlug,omitempty"`
	CategoryName string    `json:"categoryName,omitempty"`
	Content      string    `json:"content,omitempty"`
	PublishedAt  time.Time `json:"publishedAt"`
}

type SlotContentItem struct {
	ContentType string `json:"contentType"`
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Content     string `json:"content,omitempty"`
}

type IntegrationProvider struct {
	ID           string          `json:"id"`
	ProviderType string          `json:"providerType"`
	ProviderKey  string          `json:"providerKey"`
	Name         string          `json:"name"`
	Enabled      bool            `json:"enabled"`
	ConfigJSON   json.RawMessage `json:"configJson"`
	MetaJSON     json.RawMessage `json:"metaJson,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type ProviderTestResult struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	LatencyMS int64  `json:"latencyMs"`
}

type TranslationJob struct {
	ID           string    `json:"id"`
	SourceType   string    `json:"sourceType"`
	SourceID     string    `json:"sourceId"`
	SourceLocale string    `json:"sourceLocale"`
	TargetLocale string    `json:"targetLocale"`
	ProviderKey  string    `json:"providerKey"`
	ModelName    string    `json:"modelName"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
	ResultText   string    `json:"resultText,omitempty"`
	RequestedBy  string    `json:"requestedBy,omitempty"`
	RetryCount   int       `json:"retryCount"`
	MaxRetries   int       `json:"maxRetries"`
	NextRetryAt  time.Time `json:"nextRetryAt,omitempty"`
	AutoPublish  bool      `json:"autoPublish"`
	PublishAt    time.Time `json:"publishAt,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	FinishedAt   time.Time `json:"finishedAt,omitempty"`
}

type TranslationContent struct {
	SourceType        string    `json:"sourceType"`
	SourceID          string    `json:"sourceId"`
	Locale            string    `json:"locale"`
	Title             string    `json:"title,omitempty"`
	Summary           string    `json:"summary,omitempty"`
	Content           string    `json:"content"`
	Status            string    `json:"status"`
	PublishedAt       time.Time `json:"publishedAt,omitempty"`
	TranslatedByJobID string    `json:"translatedByJobId,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
