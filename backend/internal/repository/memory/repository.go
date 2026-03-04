// File: repository.go
// Purpose: Provide an in-memory repository implementation for rapid local development.
// Module: backend/internal/repository/memory, infrastructure persistence layer.
// Related: repository contracts and content service APIs.
package memory

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/google/uuid"
)

type Repository struct {
	mu sync.RWMutex

	articles map[string]domain.Article
	moments  map[string]domain.Moment
	links    map[string]domain.Link

	articleTranslations map[string]map[string]translationRecord
	momentTranslations  map[string]map[string]translationRecord

	categories []domain.Category
	tags       []domain.Tag

	siteSettings      domain.SiteSettings
	adminPasswordHash string
	footerItems       map[string]domain.FooterItem
	socialLinks       map[string]domain.SocialLink
	navItems          map[string]domain.NavItem
	slots             map[string]domain.ContentSlot
	slotItems         map[string]map[string]domain.SlotItem
}

type translationRecord struct {
	title             string
	summary           string
	content           string
	status            string
	publishedAt       time.Time
	translatedByJobID string
	createdAt         time.Time
	updatedAt         time.Time
}

func NewRepository() *Repository {
	now := time.Now().UTC()
	r := &Repository{
		articles:            make(map[string]domain.Article),
		moments:             make(map[string]domain.Moment),
		links:               make(map[string]domain.Link),
		articleTranslations: make(map[string]map[string]translationRecord),
		momentTranslations:  make(map[string]map[string]translationRecord),
		categories: []domain.Category{
			{ID: uuid.NewString(), Name: "Tech", Slug: "tech"},
			{ID: uuid.NewString(), Name: "Life", Slug: "life"},
		},
		tags: []domain.Tag{
			{ID: uuid.NewString(), Name: "Go", Slug: "go"},
			{ID: uuid.NewString(), Name: "Vue", Slug: "vue"},
		},
		siteSettings: domain.SiteSettings{
			SiteName:       "Ancy Blog",
			AvatarURL:      "",
			HeroIntroMD:    "Hi, I build things.",
			DefaultLocale:  "en",
			CommentEnabled: true,
		},
		footerItems: make(map[string]domain.FooterItem),
		socialLinks: make(map[string]domain.SocialLink),
		navItems:    make(map[string]domain.NavItem),
		slots:       make(map[string]domain.ContentSlot),
		slotItems:   make(map[string]map[string]domain.SlotItem),
	}

	about := domain.Article{
		ID:            uuid.NewString(),
		Title:         "About Me",
		Slug:          "about-me",
		ContentKind:   "page",
		Summary:       "About me page",
		Content:       "# About Me\n\nThis is a starter page.",
		Status:        "published",
		Visibility:    "public",
		AllowComment:  false,
		OriginType:    "original",
		AIAssistLevel: "none",
		PublishedAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	r.articles[about.ID] = about

	seedPost := domain.Article{
		ID:            uuid.NewString(),
		Title:         "Welcome to the Rewrite",
		Slug:          "welcome-rewrite",
		ContentKind:   "post",
		Summary:       "Initial post for the new backend.",
		Content:       "# Welcome\n\nThe rewrite starts now.",
		Status:        "published",
		Visibility:    "public",
		AllowComment:  true,
		OriginType:    "original",
		AIAssistLevel: "assisted",
		CategorySlug:  "tech",
		TagSlugs:      []string{"go"},
		PublishedAt:   now.Add(-1 * time.Hour),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	r.articles[seedPost.ID] = seedPost

	moment := domain.Moment{
		ID:           uuid.NewString(),
		Content:      "Started rebuilding the blog backend with Go.",
		Status:       "published",
		AllowComment: true,
		PublishedAt:  now.Add(-30 * time.Minute),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	r.moments[moment.ID] = moment

	r.navItems["home"] = domain.NavItem{ID: "home", Name: "Home", Key: "home", Type: "menu", TargetType: "route", TargetValue: "/", OrderNum: 1, Enabled: true}
	r.navItems["articles"] = domain.NavItem{ID: "articles", Name: "Articles", Key: "articles", Type: "dropdown", TargetType: "category", TargetValue: "", OrderNum: 2, Enabled: true}
	r.navItems["moments"] = domain.NavItem{ID: "moments", Name: "Moments", Key: "moments", Type: "menu", TargetType: "route", TargetValue: "/moments", OrderNum: 3, Enabled: true}
	r.navItems["timeline"] = domain.NavItem{ID: "timeline", Name: "Timeline", Key: "timeline", Type: "menu", TargetType: "route", TargetValue: "/timeline", OrderNum: 4, Enabled: true}
	r.navItems["links"] = domain.NavItem{ID: "links", Name: "Links", Key: "links", Type: "menu", TargetType: "route", TargetValue: "/links", OrderNum: 5, Enabled: true}

	r.slots["home_featured"] = domain.ContentSlot{ID: uuid.NewString(), SlotKey: "home_featured", Name: "Home Featured", Enabled: true}
	r.slotItems["home_featured"] = map[string]domain.SlotItem{}
	r.slotItems["home_featured"][uuid.NewString()] = domain.SlotItem{ID: uuid.NewString(), SlotKey: "home_featured", ContentType: "article", ContentID: seedPost.ID, OrderNum: 1, Enabled: true}
	r.slots["home_about"] = domain.ContentSlot{ID: uuid.NewString(), SlotKey: "home_about", Name: "Home About", Enabled: true}
	r.slotItems["home_about"] = map[string]domain.SlotItem{}
	r.slotItems["home_about"][uuid.NewString()] = domain.SlotItem{ID: uuid.NewString(), SlotKey: "home_about", ContentType: "article", ContentID: about.ID, OrderNum: 1, Enabled: true}

	r.footerItems["f1"] = domain.FooterItem{ID: "f1", Label: "About", LinkType: "internal", InternalArticleSlug: "about-me", RowNum: 1, OrderNum: 1, Enabled: true}
	r.footerItems["f2"] = domain.FooterItem{ID: "f2", Label: "ICP 00000000", LinkType: "none", RowNum: 2, OrderNum: 1, Enabled: true}
	r.socialLinks["s1"] = domain.SocialLink{ID: "s1", Platform: "github", Title: "GitHub", URL: "https://github.com", OrderNum: 1, Enabled: true}

	return r
}

func (r *Repository) CreateArticle(article domain.Article) (domain.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.slugExists(article.Slug, "") {
		return domain.Article{}, apperr.ErrSlugAlreadyExists
	}
	now := time.Now().UTC()
	article.ID = uuid.NewString()
	article.CreatedAt = now
	article.UpdatedAt = now
	if article.Status == "published" && article.PublishedAt.IsZero() {
		article.PublishedAt = now
	}
	r.articles[article.ID] = article
	return article, nil
}

func (r *Repository) UpdateArticle(id string, article domain.Article) (domain.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.articles[id]
	if !ok {
		return domain.Article{}, apperr.ErrArticleNotFound
	}
	if r.slugExists(article.Slug, id) {
		return domain.Article{}, apperr.ErrSlugAlreadyExists
	}
	article.ID = id
	article.CreatedAt = current.CreatedAt
	article.UpdatedAt = time.Now().UTC()
	if article.Status == "published" {
		if article.PublishedAt.IsZero() {
			if current.PublishedAt.IsZero() {
				article.PublishedAt = article.UpdatedAt
			} else {
				article.PublishedAt = current.PublishedAt
			}
		}
	}
	r.articles[id] = article
	return article, nil
}

func (r *Repository) DeleteArticle(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.articles[id]; !ok {
		return false
	}
	delete(r.articles, id)
	return true
}

func (r *Repository) BatchUpdateArticleStatus(ids []string, status string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if status == "" {
		return 0
	}
	affected := 0
	now := time.Now().UTC()
	for _, id := range ids {
		article, ok := r.articles[id]
		if !ok {
			continue
		}
		article.Status = status
		article.UpdatedAt = now
		if status == "published" && article.PublishedAt.IsZero() {
			article.PublishedAt = now
		}
		r.articles[id] = article
		affected++
	}
	return affected
}

func (r *Repository) ListArticles(page, pageSize int, status, contentKind, keyword string) ([]domain.Article, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.Article, 0)
	keyword = strings.TrimSpace(strings.ToLower(keyword))
	for _, a := range r.articles {
		if status != "" && a.Status != status {
			continue
		}
		if contentKind != "" && a.ContentKind != contentKind {
			continue
		}
		if keyword != "" {
			searchText := strings.ToLower(a.Title + " " + a.Slug)
			if !strings.Contains(searchText, keyword) {
				continue
			}
		}
		items = append(items, a)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return paginateArticles(items, page, pageSize)
}

func (r *Repository) ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.Article, 0)
	for _, a := range r.articles {
		if a.Status != "published" {
			continue
		}
		if contentKind != "" && a.ContentKind != contentKind {
			continue
		}
		if category != "" && a.CategorySlug != category {
			continue
		}
		if tag != "" && !contains(a.TagSlugs, tag) {
			continue
		}
		items = append(items, a)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].PublishedAt.After(items[j].PublishedAt) })
	return paginateArticles(items, page, pageSize)
}

func (r *Repository) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.articles {
		if a.Slug == slug && a.Status == "published" {
			return a, true
		}
	}
	return domain.Article{}, false
}

func (r *Repository) GetPublishedArticleBySlugWithLocale(slug, locale string) (domain.Article, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, a := range r.articles {
		if a.Slug == slug && a.Status == "published" {
			if strings.TrimSpace(locale) != "" {
				if tr, ok := r.getArticleTranslation(a.ID, locale); ok && translationVisible(tr) {
					if strings.TrimSpace(tr.title) != "" {
						a.Title = tr.title
					}
					if strings.TrimSpace(tr.summary) != "" {
						a.Summary = tr.summary
					}
					if strings.TrimSpace(tr.content) != "" {
						a.Content = tr.content
					}
				}
			}
			return a, true
		}
	}
	return domain.Article{}, false
}

func (r *Repository) GetArticleByID(id string) (domain.Article, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.articles[id]
	return a, ok
}

func (r *Repository) SlugExists(slug string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.slugExists(slug, "")
}

func (r *Repository) CreateMoment(moment domain.Moment) (domain.Moment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	moment.ID = uuid.NewString()
	moment.CreatedAt = now
	moment.UpdatedAt = now
	if moment.Status == "published" && moment.PublishedAt.IsZero() {
		moment.PublishedAt = now
	}
	r.moments[moment.ID] = moment
	return moment, nil
}

func (r *Repository) UpdateMoment(id string, moment domain.Moment) (domain.Moment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.moments[id]
	if !ok {
		return domain.Moment{}, apperr.ErrMomentNotFound
	}
	moment.ID = id
	moment.CreatedAt = current.CreatedAt
	moment.UpdatedAt = time.Now().UTC()
	if moment.Status == "published" {
		if moment.PublishedAt.IsZero() {
			if current.PublishedAt.IsZero() {
				moment.PublishedAt = moment.UpdatedAt
			} else {
				moment.PublishedAt = current.PublishedAt
			}
		}
	}
	r.moments[id] = moment
	return moment, nil
}

func (r *Repository) DeleteMoment(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.moments[id]; !ok {
		return false
	}
	delete(r.moments, id)
	return true
}

func (r *Repository) BatchUpdateMomentStatus(ids []string, status string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if status == "" {
		return 0
	}
	affected := 0
	now := time.Now().UTC()
	for _, id := range ids {
		moment, ok := r.moments[id]
		if !ok {
			continue
		}
		moment.Status = status
		moment.UpdatedAt = now
		if status == "published" && moment.PublishedAt.IsZero() {
			moment.PublishedAt = now
		}
		r.moments[id] = moment
		affected++
	}
	return affected
}

func (r *Repository) ListMoments(page, pageSize int, status string) ([]domain.Moment, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.Moment, 0)
	for _, m := range r.moments {
		if status != "" && m.Status != status {
			continue
		}
		items = append(items, m)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return paginateMoments(items, page, pageSize)
}

func (r *Repository) ListPublishedMoments(page, pageSize int, locale string) ([]domain.Moment, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.Moment, 0)
	for _, m := range r.moments {
		if m.Status == "published" {
			if strings.TrimSpace(locale) != "" {
				if tr, ok := r.getMomentTranslation(m.ID, locale); ok && translationVisible(tr) {
					if strings.TrimSpace(tr.content) != "" {
						m.Content = tr.content
					}
				}
			}
			items = append(items, m)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].PublishedAt.After(items[j].PublishedAt) })
	return paginateMoments(items, page, pageSize)
}

func (r *Repository) SubmitLink(link domain.Link) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	link.ID = uuid.NewString()
	link.ReviewStatus = "pending"
	link.CreatedAt = now
	link.UpdatedAt = now
	r.links[link.ID] = link
	return link, nil
}

func (r *Repository) ListApprovedLinks() []domain.Link {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.Link, 0)
	for _, l := range r.links {
		if l.ReviewStatus == "approved" {
			items = append(items, l)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items
}

func (r *Repository) ListLinkSubmissions(page, pageSize int, reviewStatus string) ([]domain.Link, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.Link, 0)
	for _, l := range r.links {
		if reviewStatus != "" && l.ReviewStatus != reviewStatus {
			continue
		}
		items = append(items, l)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return paginateLinks(items, page, pageSize)
}

func (r *Repository) ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	link, ok := r.links[id]
	if !ok {
		return domain.Link{}, apperr.ErrLinkNotFound
	}
	link.ReviewStatus = reviewStatus
	link.ReviewNote = reviewNote
	link.RelatedArticleID = relatedArticleID
	if reviewStatus == "approved" {
		link.ApprovedAt = time.Now().UTC()
	}
	link.UpdatedAt = time.Now().UTC()
	r.links[id] = link
	return link, nil
}

func (r *Repository) ListCategories() []domain.Category {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Category, len(r.categories))
	copy(out, r.categories)
	return out
}

func (r *Repository) ListTags() []domain.Tag {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Tag, len(r.tags))
	copy(out, r.tags)
	return out
}

func (r *Repository) GetSiteSettings() domain.SiteSettings {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.siteSettings
}

func (r *Repository) UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings {
	r.mu.Lock()
	defer r.mu.Unlock()
	if settings.SiteName == "" {
		settings.SiteName = r.siteSettings.SiteName
	}
	if settings.DefaultLocale == "" {
		settings.DefaultLocale = r.siteSettings.DefaultLocale
	}
	r.siteSettings = settings
	return r.siteSettings
}

func (r *Repository) GetTranslationPolicy() domain.TranslationPolicy {
	return domain.TranslationPolicy{TargetLocales: []string{}}
}

func (r *Repository) UpdateTranslationPolicy(_ domain.TranslationPolicy) error {
	return nil
}

func (r *Repository) GetAdminPasswordHash() (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.adminPasswordHash == "" {
		return "", false
	}
	return r.adminPasswordHash, true
}

func (r *Repository) SetAdminPasswordHash(hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adminPasswordHash = hash
	return nil
}

func (r *Repository) CreateCategory(category domain.Category) (domain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	category.ID = uuid.NewString()
	r.categories = append(r.categories, category)
	return category, nil
}

func (r *Repository) DeleteCategory(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, c := range r.categories {
		if c.ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return true
		}
	}
	return false
}

func (r *Repository) CreateTag(tag domain.Tag) (domain.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tag.ID = uuid.NewString()
	r.tags = append(r.tags, tag)
	return tag, nil
}

func (r *Repository) DeleteTag(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, t := range r.tags {
		if t.ID == id {
			r.tags = append(r.tags[:i], r.tags[i+1:]...)
			return true
		}
	}
	return false
}

func (r *Repository) CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = uuid.NewString()
	r.footerItems[item.ID] = item
	return item, nil
}

func (r *Repository) UpdateFooterItem(id string, item domain.FooterItem) (domain.FooterItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.footerItems[id]; !ok {
		return domain.FooterItem{}, apperr.ErrFooterItemNotFound
	}
	item.ID = id
	r.footerItems[id] = item
	return item, nil
}

func (r *Repository) DeleteFooterItem(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.footerItems[id]; !ok {
		return false
	}
	delete(r.footerItems, id)
	return true
}

func (r *Repository) ListFooterItems() []domain.FooterItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.FooterItem, 0)
	for _, f := range r.footerItems {
		if f.Enabled {
			items = append(items, f)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].RowNum == items[j].RowNum {
			return items[i].OrderNum < items[j].OrderNum
		}
		return items[i].RowNum < items[j].RowNum
	})
	return items
}

func (r *Repository) CreateSocialLink(item domain.SocialLink) (domain.SocialLink, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = uuid.NewString()
	r.socialLinks[item.ID] = item
	return item, nil
}

func (r *Repository) UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.socialLinks[id]; !ok {
		return domain.SocialLink{}, apperr.ErrSocialLinkNotFound
	}
	item.ID = id
	r.socialLinks[id] = item
	return item, nil
}

func (r *Repository) DeleteSocialLink(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.socialLinks[id]; !ok {
		return false
	}
	delete(r.socialLinks, id)
	return true
}

func (r *Repository) ListSocialLinks() []domain.SocialLink {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.SocialLink, 0)
	for _, s := range r.socialLinks {
		if s.Enabled {
			items = append(items, s)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].OrderNum < items[j].OrderNum })
	return items
}

func (r *Repository) CreateNavItem(item domain.NavItem) (domain.NavItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = uuid.NewString()
	r.navItems[item.ID] = item
	return item, nil
}

func (r *Repository) UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.navItems[id]; !ok {
		return domain.NavItem{}, apperr.ErrNavItemNotFound
	}
	item.ID = id
	r.navItems[id] = item
	return item, nil
}

func (r *Repository) DeleteNavItem(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.navItems[id]; !ok {
		return false
	}
	delete(r.navItems, id)
	return true
}

func (r *Repository) ListNavItems() []domain.NavItem {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.NavItem, 0)
	for _, n := range r.navItems {
		if n.Enabled {
			items = append(items, n)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].OrderNum < items[j].OrderNum })
	return items
}

func (r *Repository) CreateContentSlot(slot domain.ContentSlot) (domain.ContentSlot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.slots[slot.SlotKey]; ok {
		return domain.ContentSlot{}, apperr.ErrValidation
	}
	slot.ID = uuid.NewString()
	r.slots[slot.SlotKey] = slot
	r.slotItems[slot.SlotKey] = map[string]domain.SlotItem{}
	return slot, nil
}

func (r *Repository) ListContentSlots() []domain.ContentSlot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.ContentSlot, 0, len(r.slots))
	for _, slot := range r.slots {
		items = append(items, slot)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].SlotKey < items[j].SlotKey })
	return items
}

func (r *Repository) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.slots[slotKey]; !ok {
		return domain.SlotItem{}, apperr.ErrSlotNotFound
	}
	item.ID = uuid.NewString()
	item.SlotKey = slotKey
	if _, ok := r.slotItems[slotKey]; !ok {
		r.slotItems[slotKey] = map[string]domain.SlotItem{}
	}
	r.slotItems[slotKey][item.ID] = item
	return item, nil
}

func (r *Repository) ListSlotItems(slotKey string) ([]domain.SlotItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.slots[slotKey]; !ok {
		return nil, false
	}
	items := make([]domain.SlotItem, 0)
	for _, it := range r.slotItems[slotKey] {
		items = append(items, it)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].OrderNum == items[j].OrderNum {
			return items[i].ID < items[j].ID
		}
		return items[i].OrderNum < items[j].OrderNum
	})
	return items, true
}

func (r *Repository) DeleteSlotItem(slotKey, itemID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	items, ok := r.slotItems[slotKey]
	if !ok {
		return false
	}
	if _, ok := items[itemID]; !ok {
		return false
	}
	delete(items, itemID)
	return true
}

func (r *Repository) ListSlotContent(slotKey string, limit int) ([]domain.SlotContentItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.slots[slotKey]; !ok {
		return nil, false
	}
	items := make([]domain.SlotItem, 0)
	for _, it := range r.slotItems[slotKey] {
		if it.Enabled {
			items = append(items, it)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].OrderNum < items[j].OrderNum })
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	out := make([]domain.SlotContentItem, 0, len(items))
	for _, it := range items {
		switch it.ContentType {
		case "article":
			if a, ok := r.articles[it.ContentID]; ok && a.Status == "published" {
				out = append(out, domain.SlotContentItem{ContentType: "article", ID: a.ID, Title: a.Title, Slug: a.Slug, Summary: a.Summary})
			}
		case "moment":
			if m, ok := r.moments[it.ContentID]; ok && m.Status == "published" {
				out = append(out, domain.SlotContentItem{ContentType: "moment", ID: m.ID, Content: m.Content})
			}
		}
	}
	return out, true
}

func (r *Repository) ListTimeline(page, pageSize int, locale string) ([]domain.TimelineItem, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]domain.TimelineItem, 0)
	for _, a := range r.articles {
		if a.Status == "published" {
			content := ""
			if strings.TrimSpace(locale) != "" {
				if tr, ok := r.getArticleTranslation(a.ID, locale); ok && translationVisible(tr) {
					if strings.TrimSpace(tr.title) != "" {
						a.Title = tr.title
					}
					if strings.TrimSpace(tr.summary) != "" {
						a.Summary = tr.summary
					}
					content = tr.content
				}
			}
			items = append(items, domain.TimelineItem{
				ContentType: "article",
				ID:          a.ID,
				Title:       a.Title,
				Summary:     a.Summary,
				Slug:        a.Slug,
				Content:     content,
				PublishedAt: a.PublishedAt,
			})
		}
	}
	for _, m := range r.moments {
		if m.Status == "published" {
			content := m.Content
			if strings.TrimSpace(locale) != "" {
				if tr, ok := r.getMomentTranslation(m.ID, locale); ok && translationVisible(tr) {
					content = tr.content
				}
			}
			items = append(items, domain.TimelineItem{ContentType: "moment", ID: m.ID, Content: content, PublishedAt: m.PublishedAt})
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].PublishedAt.After(items[j].PublishedAt) })
	return paginateTimeline(items, page, pageSize)
}

func (r *Repository) ClaimNextQueuedTranslationJob() (domain.TranslationJob, bool, error) {
	return domain.TranslationJob{}, false, nil
}

func (r *Repository) MarkTranslationJobRunning(id string) error {
	_ = id
	return nil
}

func (r *Repository) MarkTranslationJobSucceeded(id, resultText string) error {
	_ = id
	_ = resultText
	return nil
}

func (r *Repository) MarkTranslationJobFailed(id, errorMessage string) error {
	_ = id
	_ = errorMessage
	return nil
}

func (r *Repository) GetTranslationSourceText(sourceType, sourceID string) (string, bool, error) {
	switch sourceType {
	case "article":
		a, ok := r.GetArticleByID(sourceID)
		if !ok {
			return "", false, nil
		}
		return strings.TrimSpace(a.Title + "\n\n" + a.Summary + "\n\n" + a.Content), true, nil
	case "moment":
		r.mu.RLock()
		defer r.mu.RUnlock()
		m, ok := r.moments[sourceID]
		if !ok {
			return "", false, nil
		}
		return m.Content, true, nil
	default:
		return "", false, nil
	}
}

func (r *Repository) UpsertArticleTranslation(articleID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if status == "" {
		status = "draft"
	}
	if status == "published" && publishedAt.IsZero() {
		publishedAt = time.Now().UTC()
	}
	if _, ok := r.articleTranslations[articleID]; !ok {
		r.articleTranslations[articleID] = map[string]translationRecord{}
	}
	now := time.Now().UTC()
	current, exists := r.articleTranslations[articleID][locale]
	if !exists {
		current.createdAt = now
	}
	current.title = title
	current.summary = summary
	current.content = content
	current.status = status
	current.publishedAt = publishedAt
	current.translatedByJobID = translatedByJobID
	current.updatedAt = now
	r.articleTranslations[articleID][locale] = current
	return nil
}

func (r *Repository) UpsertMomentTranslation(momentID, locale, content, status string, publishedAt time.Time, translatedByJobID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if status == "" {
		status = "draft"
	}
	if status == "published" && publishedAt.IsZero() {
		publishedAt = time.Now().UTC()
	}
	if _, ok := r.momentTranslations[momentID]; !ok {
		r.momentTranslations[momentID] = map[string]translationRecord{}
	}
	now := time.Now().UTC()
	current, exists := r.momentTranslations[momentID][locale]
	if !exists {
		current.createdAt = now
	}
	current.content = content
	current.status = status
	current.publishedAt = publishedAt
	current.translatedByJobID = translatedByJobID
	current.updatedAt = now
	r.momentTranslations[momentID][locale] = current
	return nil
}

func (r *Repository) getArticleTranslation(articleID, locale string) (translationRecord, bool) {
	translations, ok := r.articleTranslations[articleID]
	if !ok {
		return translationRecord{}, false
	}
	row, ok := translations[locale]
	return row, ok
}

func (r *Repository) getMomentTranslation(momentID, locale string) (translationRecord, bool) {
	translations, ok := r.momentTranslations[momentID]
	if !ok {
		return translationRecord{}, false
	}
	row, ok := translations[locale]
	return row, ok
}

func (r *Repository) ListTranslationContents(page, pageSize int, sourceType, sourceID, locale string) ([]domain.TranslationContent, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows := make([]domain.TranslationContent, 0)
	switch sourceType {
	case "article":
		rows = r.listArticleTranslationContents(sourceID, locale)
	case "moment":
		rows = r.listMomentTranslationContents(sourceID, locale)
	default:
		return []domain.TranslationContent{}, 0
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].UpdatedAt.After(rows[j].UpdatedAt) })
	return paginateTranslationContents(rows, page, pageSize)
}

func (r *Repository) GetTranslationContent(sourceType, sourceID, locale string) (domain.TranslationContent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	switch sourceType {
	case "article":
		translations, ok := r.articleTranslations[sourceID]
		if !ok {
			return domain.TranslationContent{}, false
		}
		row, ok := translations[locale]
		if !ok {
			return domain.TranslationContent{}, false
		}
		return domain.TranslationContent{
			SourceType:        "article",
			SourceID:          sourceID,
			Locale:            locale,
			Title:             row.title,
			Summary:           row.summary,
			Content:           row.content,
			Status:            row.status,
			PublishedAt:       row.publishedAt,
			TranslatedByJobID: row.translatedByJobID,
			CreatedAt:         row.createdAt,
			UpdatedAt:         row.updatedAt,
		}, true
	case "moment":
		translations, ok := r.momentTranslations[sourceID]
		if !ok {
			return domain.TranslationContent{}, false
		}
		row, ok := translations[locale]
		if !ok {
			return domain.TranslationContent{}, false
		}
		return domain.TranslationContent{
			SourceType:        "moment",
			SourceID:          sourceID,
			Locale:            locale,
			Content:           row.content,
			Status:            row.status,
			PublishedAt:       row.publishedAt,
			TranslatedByJobID: row.translatedByJobID,
			CreatedAt:         row.createdAt,
			UpdatedAt:         row.updatedAt,
		}, true
	default:
		return domain.TranslationContent{}, false
	}
}

func (r *Repository) UpsertTranslationContent(sourceType, sourceID, locale, title, summary, content, status string, publishedAt time.Time, translatedByJobID string) (domain.TranslationContent, error) {
	switch sourceType {
	case "article":
		if err := r.UpsertArticleTranslation(sourceID, locale, title, summary, content, status, publishedAt, translatedByJobID); err != nil {
			return domain.TranslationContent{}, err
		}
	case "moment":
		if err := r.UpsertMomentTranslation(sourceID, locale, content, status, publishedAt, translatedByJobID); err != nil {
			return domain.TranslationContent{}, err
		}
	default:
		return domain.TranslationContent{}, apperr.ErrValidation
	}
	row, _ := r.GetTranslationContent(sourceType, sourceID, locale)
	return row, nil
}

func (r *Repository) listArticleTranslationContents(sourceID, locale string) []domain.TranslationContent {
	rows := make([]domain.TranslationContent, 0)
	for articleID, translations := range r.articleTranslations {
		if strings.TrimSpace(sourceID) != "" && articleID != sourceID {
			continue
		}
		for lc, rec := range translations {
			if strings.TrimSpace(locale) != "" && lc != locale {
				continue
			}
			rows = append(rows, domain.TranslationContent{
				SourceType:        "article",
				SourceID:          articleID,
				Locale:            lc,
				Title:             rec.title,
				Summary:           rec.summary,
				Content:           rec.content,
				Status:            rec.status,
				PublishedAt:       rec.publishedAt,
				TranslatedByJobID: rec.translatedByJobID,
				CreatedAt:         rec.createdAt,
				UpdatedAt:         rec.updatedAt,
			})
		}
	}
	return rows
}

func (r *Repository) listMomentTranslationContents(sourceID, locale string) []domain.TranslationContent {
	rows := make([]domain.TranslationContent, 0)
	for momentID, translations := range r.momentTranslations {
		if strings.TrimSpace(sourceID) != "" && momentID != sourceID {
			continue
		}
		for lc, rec := range translations {
			if strings.TrimSpace(locale) != "" && lc != locale {
				continue
			}
			rows = append(rows, domain.TranslationContent{
				SourceType:        "moment",
				SourceID:          momentID,
				Locale:            lc,
				Content:           rec.content,
				Status:            rec.status,
				PublishedAt:       rec.publishedAt,
				TranslatedByJobID: rec.translatedByJobID,
				CreatedAt:         rec.createdAt,
				UpdatedAt:         rec.updatedAt,
			})
		}
	}
	return rows
}

func (r *Repository) slugExists(slug, excludedID string) bool {
	for id, a := range r.articles {
		if strings.EqualFold(a.Slug, slug) && id != excludedID {
			return true
		}
	}
	return false
}

func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func translationVisible(rec translationRecord) bool {
	if rec.status != "published" {
		return false
	}
	if rec.publishedAt.IsZero() {
		return true
	}
	return !rec.publishedAt.After(time.Now().UTC())
}

func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func paginateArticles(items []domain.Article, page, pageSize int) ([]domain.Article, int) {
	page, pageSize = normalizePagination(page, pageSize)
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []domain.Article{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}

func paginateMoments(items []domain.Moment, page, pageSize int) ([]domain.Moment, int) {
	page, pageSize = normalizePagination(page, pageSize)
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []domain.Moment{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}

func paginateLinks(items []domain.Link, page, pageSize int) ([]domain.Link, int) {
	page, pageSize = normalizePagination(page, pageSize)
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []domain.Link{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}

func paginateTimeline(items []domain.TimelineItem, page, pageSize int) ([]domain.TimelineItem, int) {
	page, pageSize = normalizePagination(page, pageSize)
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []domain.TimelineItem{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}

func paginateTranslationContents(items []domain.TranslationContent, page, pageSize int) ([]domain.TranslationContent, int) {
	page, pageSize = normalizePagination(page, pageSize)
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []domain.TranslationContent{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}
