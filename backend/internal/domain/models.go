// File: models.go
// Purpose: Define core domain entities and lightweight view models used across services.
// Module: backend/internal/domain, domain model layer.
// Related: repository interfaces, service logic, and HTTP handlers.
package domain

import "time"

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
	OriginType    string    `json:"originType"`
	SourceURL     string    `json:"sourceUrl,omitempty"`
	AIAssistLevel string    `json:"aiAssistLevel"`
	CoverImage    string    `json:"coverImage,omitempty"`
	CategorySlug  string    `json:"categorySlug,omitempty"`
	TagSlugs      []string  `json:"tagSlugs,omitempty"`
	PublishedAt   time.Time `json:"publishedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Moment struct {
	ID           string    `json:"id"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	AllowComment bool      `json:"allowComment"`
	PublishedAt  time.Time `json:"publishedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
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
	SiteName      string `json:"siteName"`
	AvatarURL     string `json:"avatarUrl,omitempty"`
	HeroIntroMD   string `json:"heroIntroMd,omitempty"`
	DefaultLocale string `json:"defaultLocale"`
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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Type        string `json:"type"`
	TargetType  string `json:"targetType"`
	TargetValue string `json:"targetValue,omitempty"`
	OrderNum    int    `json:"orderNum"`
	Enabled     bool   `json:"enabled"`
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
	ContentType string    `json:"contentType"`
	ID          string    `json:"id"`
	Title       string    `json:"title,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Content     string    `json:"content,omitempty"`
	PublishedAt time.Time `json:"publishedAt"`
}

type SlotContentItem struct {
	ContentType string `json:"contentType"`
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Content     string `json:"content,omitempty"`
}
