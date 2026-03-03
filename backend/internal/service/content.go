// File: content.go
// Purpose: Orchestrate content, site configuration, and timeline business operations.
// Module: backend/internal/service, content service layer.
// Related: content repository implementations and HTTP handlers.
package service

import (
	"errors"
	"net/url"
	"strings"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
)

type ContentService struct {
	repo repository.ContentRepository
}

func NewContentService(repo repository.ContentRepository) *ContentService {
	return &ContentService{repo: repo}
}

func (s *ContentService) CreateArticle(article domain.Article) (domain.Article, error) {
	if strings.TrimSpace(article.Title) == "" || strings.TrimSpace(article.Slug) == "" {
		return domain.Article{}, errors.New("title and slug are required")
	}
	if article.ContentKind == "" {
		article.ContentKind = "post"
	}
	if article.Status == "" {
		article.Status = "draft"
	}
	if article.Visibility == "" {
		article.Visibility = "public"
	}
	if article.OriginType == "" {
		article.OriginType = "original"
	}
	if article.AIAssistLevel == "" {
		article.AIAssistLevel = "none"
	}
	return s.repo.CreateArticle(article)
}

func (s *ContentService) UpdateArticle(id string, article domain.Article) (domain.Article, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Article{}, errors.New("article id is required")
	}
	if strings.TrimSpace(article.Title) == "" || strings.TrimSpace(article.Slug) == "" {
		return domain.Article{}, errors.New("title and slug are required")
	}
	if article.ContentKind == "" {
		article.ContentKind = "post"
	}
	return s.repo.UpdateArticle(id, article)
}

func (s *ContentService) ListPublishedArticles(page, pageSize int, category, tag, contentKind string) ([]domain.Article, int) {
	if contentKind == "" {
		contentKind = "post"
	}
	return s.repo.ListPublishedArticles(page, pageSize, category, tag, contentKind)
}

func (s *ContentService) GetPublishedArticleBySlug(slug string) (domain.Article, bool) {
	return s.repo.GetPublishedArticleBySlug(slug)
}

func (s *ContentService) CreateMoment(moment domain.Moment) (domain.Moment, error) {
	if strings.TrimSpace(moment.Content) == "" {
		return domain.Moment{}, errors.New("content is required")
	}
	if moment.Status == "" {
		moment.Status = "draft"
	}
	return s.repo.CreateMoment(moment)
}

func (s *ContentService) ListPublishedMoments(page, pageSize int) ([]domain.Moment, int) {
	return s.repo.ListPublishedMoments(page, pageSize)
}

func (s *ContentService) SubmitLink(link domain.Link) (domain.Link, error) {
	if strings.TrimSpace(link.Name) == "" || strings.TrimSpace(link.URL) == "" {
		return domain.Link{}, errors.New("name and url are required")
	}
	if _, err := url.ParseRequestURI(link.URL); err != nil {
		return domain.Link{}, errors.New("url is invalid")
	}
	return s.repo.SubmitLink(link)
}

func (s *ContentService) ListApprovedLinks() []domain.Link {
	return s.repo.ListApprovedLinks()
}

func (s *ContentService) ListLinkSubmissions(page, pageSize int, reviewStatus string) ([]domain.Link, int) {
	return s.repo.ListLinkSubmissions(page, pageSize, reviewStatus)
}

func (s *ContentService) ReviewLink(id, reviewStatus, reviewNote, relatedArticleID string) (domain.Link, error) {
	if reviewStatus != "pending" && reviewStatus != "approved" && reviewStatus != "rejected" {
		return domain.Link{}, errors.New("invalid review status")
	}
	if relatedArticleID != "" {
		if _, ok := s.repo.GetArticleByID(relatedArticleID); !ok {
			return domain.Link{}, errors.New("related article not found")
		}
	}
	return s.repo.ReviewLink(id, reviewStatus, reviewNote, relatedArticleID)
}

func (s *ContentService) ListCategories() []domain.Category {
	return s.repo.ListCategories()
}

func (s *ContentService) ListTags() []domain.Tag {
	return s.repo.ListTags()
}

func (s *ContentService) GetSiteSettings() domain.SiteSettings {
	return s.repo.GetSiteSettings()
}

func (s *ContentService) UpdateSiteSettings(settings domain.SiteSettings) domain.SiteSettings {
	if settings.DefaultLocale == "" {
		settings.DefaultLocale = "en"
	}
	return s.repo.UpdateSiteSettings(settings)
}

func (s *ContentService) CreateFooterItem(item domain.FooterItem) (domain.FooterItem, error) {
	if strings.TrimSpace(item.Label) == "" {
		return domain.FooterItem{}, errors.New("label is required")
	}
	if item.RowNum < 1 || item.RowNum > 3 {
		return domain.FooterItem{}, errors.New("rowNum must be between 1 and 3")
	}
	if item.LinkType == "internal" {
		a, ok := s.repo.GetPublishedArticleBySlug(item.InternalArticleSlug)
		if !ok || a.ContentKind != "page" {
			return domain.FooterItem{}, errors.New("internal article slug must point to a published page")
		}
	}
	if item.Enabled == false && item.OrderNum == 0 {
		item.OrderNum = 1
	}
	return s.repo.CreateFooterItem(item)
}

func (s *ContentService) UpdateFooterItem(id string, item domain.FooterItem) (domain.FooterItem, error) {
	if item.RowNum < 1 || item.RowNum > 3 {
		return domain.FooterItem{}, errors.New("rowNum must be between 1 and 3")
	}
	if item.LinkType == "internal" {
		a, ok := s.repo.GetPublishedArticleBySlug(item.InternalArticleSlug)
		if !ok || a.ContentKind != "page" {
			return domain.FooterItem{}, errors.New("internal article slug must point to a published page")
		}
	}
	return s.repo.UpdateFooterItem(id, item)
}

func (s *ContentService) DeleteFooterItem(id string) bool {
	return s.repo.DeleteFooterItem(id)
}

func (s *ContentService) ListFooterItems() []domain.FooterItem {
	return s.repo.ListFooterItems()
}

func (s *ContentService) CreateSocialLink(item domain.SocialLink) (domain.SocialLink, error) {
	if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.URL) == "" {
		return domain.SocialLink{}, errors.New("title and url are required")
	}
	return s.repo.CreateSocialLink(item)
}

func (s *ContentService) UpdateSocialLink(id string, item domain.SocialLink) (domain.SocialLink, error) {
	if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.URL) == "" {
		return domain.SocialLink{}, errors.New("title and url are required")
	}
	return s.repo.UpdateSocialLink(id, item)
}

func (s *ContentService) DeleteSocialLink(id string) bool {
	return s.repo.DeleteSocialLink(id)
}

func (s *ContentService) ListSocialLinks() []domain.SocialLink {
	return s.repo.ListSocialLinks()
}

func (s *ContentService) CreateNavItem(item domain.NavItem) (domain.NavItem, error) {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Key) == "" {
		return domain.NavItem{}, errors.New("name and key are required")
	}
	return s.repo.CreateNavItem(item)
}

func (s *ContentService) UpdateNavItem(id string, item domain.NavItem) (domain.NavItem, error) {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Key) == "" {
		return domain.NavItem{}, errors.New("name and key are required")
	}
	return s.repo.UpdateNavItem(id, item)
}

func (s *ContentService) DeleteNavItem(id string) bool {
	return s.repo.DeleteNavItem(id)
}

func (s *ContentService) ListNavItems() []domain.NavItem {
	return s.repo.ListNavItems()
}

func (s *ContentService) CreateContentSlot(slot domain.ContentSlot) (domain.ContentSlot, error) {
	if strings.TrimSpace(slot.SlotKey) == "" || strings.TrimSpace(slot.Name) == "" {
		return domain.ContentSlot{}, errors.New("slotKey and name are required")
	}
	return s.repo.CreateContentSlot(slot)
}

func (s *ContentService) CreateSlotItem(slotKey string, item domain.SlotItem) (domain.SlotItem, error) {
	if item.ContentType != "article" && item.ContentType != "moment" {
		return domain.SlotItem{}, errors.New("contentType must be article or moment")
	}
	if item.ContentType == "article" {
		if _, ok := s.repo.GetArticleByID(item.ContentID); !ok {
			return domain.SlotItem{}, errors.New("article not found")
		}
	}
	return s.repo.CreateSlotItem(slotKey, item)
}

func (s *ContentService) DeleteSlotItem(slotKey, itemID string) bool {
	return s.repo.DeleteSlotItem(slotKey, itemID)
}

func (s *ContentService) ListSlotContent(slotKey string, limit int) ([]domain.SlotContentItem, bool) {
	return s.repo.ListSlotContent(slotKey, limit)
}

func (s *ContentService) ListTimeline(page, pageSize int) ([]domain.TimelineItem, int) {
	return s.repo.ListTimeline(page, pageSize)
}
