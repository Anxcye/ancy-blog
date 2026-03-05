// File: admin.go
// Purpose: Implement authenticated admin APIs for content and site management.
// Module: backend/internal/handler, admin HTTP presentation layer.
// Related: service module facades and auth middleware.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/handler/dto"
	"github.com/anxcye/ancy-blog/backend/internal/middleware"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	articleService     *service.ArticleService
	commentService     *service.CommentService
	linkService        *service.LinkService
	siteService        *service.SiteService
	integrationService *service.IntegrationService
	translationService *service.TranslationService
	aiAssistService    *service.AIAssistService
	authService        *service.AuthService
	tmdbService        *service.TMDBService
}

func NewAdminHandler(
	articleService *service.ArticleService,
	commentService *service.CommentService,
	linkService *service.LinkService,
	siteService *service.SiteService,
	integrationService *service.IntegrationService,
	translationService *service.TranslationService,
	aiAssistService *service.AIAssistService,
	authService *service.AuthService,
	tmdbService *service.TMDBService,
) *AdminHandler {
	return &AdminHandler{
		articleService:     articleService,
		commentService:     commentService,
		linkService:        linkService,
		siteService:        siteService,
		integrationService: integrationService,
		translationService: translationService,
		aiAssistService:    aiAssistService,
		authService:        authService,
		tmdbService:        tmdbService,
	}
}

func (h *AdminHandler) CreateArticle(c *gin.Context) {
	var req dto.ArticleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	article, err := h.articleService.CreateArticle(domain.Article{
		Title:         req.Title,
		Slug:          req.Slug,
		ContentKind:   req.ContentKind,
		Summary:       req.Summary,
		Content:       req.Content,
		Status:        req.Status,
		Visibility:    req.Visibility,
		AllowComment:  req.AllowComment,
		IsPinned:      req.IsPinned,
		IsFeatured:    req.IsFeatured,
		OriginType:    req.OriginType,
		SourceURL:     req.SourceURL,
		AIAssistLevel: req.AIAssistLevel,
		CoverImage:    req.CoverImage,
		CategorySlug:  req.CategorySlug,
		TagSlugs:      req.TagSlugs,
		PublishedAt:   req.PublishedAt,
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": article.ID}})
}

func (h *AdminHandler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var req dto.ArticleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	article, err := h.articleService.UpdateArticle(id, domain.Article{
		Title:         req.Title,
		Slug:          req.Slug,
		ContentKind:   req.ContentKind,
		Summary:       req.Summary,
		Content:       req.Content,
		Status:        req.Status,
		Visibility:    req.Visibility,
		AllowComment:  req.AllowComment,
		IsPinned:      req.IsPinned,
		IsFeatured:    req.IsFeatured,
		OriginType:    req.OriginType,
		SourceURL:     req.SourceURL,
		AIAssistLevel: req.AIAssistLevel,
		CoverImage:    req.CoverImage,
		CategorySlug:  req.CategorySlug,
		TagSlugs:      req.TagSlugs,
		PublishedAt:   req.PublishedAt,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrArticleNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "ARTICLE_NOT_FOUND", Message: "article not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": article.ID}})
}

func (h *AdminHandler) ListArticles(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	status := c.Query("status")
	contentKind := c.Query("contentKind")
	keyword := c.Query("keyword")
	rows, total := h.articleService.ListArticles(page, pageSize, status, contentKind, keyword)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Article]{Total: total, Rows: rows}})
}

func (h *AdminHandler) ArticleDetail(c *gin.Context) {
	id := c.Param("id")
	article, ok := h.articleService.GetArticleByID(id)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "ARTICLE_NOT_FOUND", Message: "article not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: article})
}

func (h *AdminHandler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if !h.articleService.DeleteArticle(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "ARTICLE_NOT_FOUND", Message: "article not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) BatchUpdateArticleStatus(c *gin.Context) {
	var req dto.ArticleBatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	affected, err := h.articleService.BatchUpdateArticleStatus(req.IDs, req.Status)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]int{"affected": affected}})
}

func (h *AdminHandler) BatchDeleteArticle(c *gin.Context) {
	var req dto.ArticleBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	affected := 0
	for _, id := range req.IDs {
		if h.articleService.DeleteArticle(id) {
			affected++
		}
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]int{"affected": affected}})
}

func (h *AdminHandler) CreateMoment(c *gin.Context) {
	var req dto.MomentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	moment, err := h.articleService.CreateMoment(domain.Moment{
		Content:      req.Content,
		Status:       req.Status,
		AllowComment: req.AllowComment,
		PublishedAt:  req.PublishedAt,
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": moment.ID}})
}

func (h *AdminHandler) ListMoments(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	status := c.Query("status")
	rows, total := h.articleService.ListMoments(page, pageSize, status)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Moment]{Total: total, Rows: rows}})
}

func (h *AdminHandler) UpdateMoment(c *gin.Context) {
	id := c.Param("id")
	var req dto.MomentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	moment, err := h.articleService.UpdateMoment(id, domain.Moment{
		Content:      req.Content,
		Status:       req.Status,
		AllowComment: req.AllowComment,
		PublishedAt:  req.PublishedAt,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrMomentNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "MOMENT_NOT_FOUND", Message: "moment not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: moment})
}

func (h *AdminHandler) DeleteMoment(c *gin.Context) {
	id := c.Param("id")
	if !h.articleService.DeleteMoment(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "MOMENT_NOT_FOUND", Message: "moment not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) BatchUpdateMomentStatus(c *gin.Context) {
	var req dto.MomentBatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	affected, err := h.articleService.BatchUpdateMomentStatus(req.IDs, req.Status)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]int{"affected": affected}})
}

func (h *AdminHandler) BatchDeleteMoment(c *gin.Context) {
	var req dto.MomentBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	affected := 0
	for _, id := range req.IDs {
		if h.articleService.DeleteMoment(id) {
			affected++
		}
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]int{"affected": affected}})
}

func (h *AdminHandler) CommentPage(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	status := c.Query("status")
	rows, total := h.commentService.ListCommentPage(page, pageSize, status)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Comment]{Total: total, Rows: rows}})
}

func (h *AdminHandler) CommentUpdate(c *gin.Context) {
	id := c.Param("id")
	var req dto.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	comment, err := h.commentService.UpdateCommentAdmin(id, req.Status, req.IsPinned)
	if err != nil {
		if errors.Is(err, apperr.ErrCommentNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "COMMENT_NOT_FOUND", Message: "comment not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: comment})
}

func (h *AdminHandler) ListLinkSubmissions(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	reviewStatus := c.Query("reviewStatus")
	rows, total := h.linkService.ListLinkSubmissions(page, pageSize, reviewStatus)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Link]{Total: total, Rows: rows}})
}

func (h *AdminHandler) ReviewLink(c *gin.Context) {
	id := c.Param("id")
	var req dto.ReviewLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	_, err := h.linkService.ReviewLink(id, req.ReviewStatus, req.ReviewNote, req.RelatedArticleID)
	if err != nil {
		if errors.Is(err, apperr.ErrLinkNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "LINK_NOT_FOUND", Message: "link not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) UpdateSiteSettings(c *gin.Context) {
	var req dto.SiteSettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	updated := h.siteService.UpdateSiteSettings(domain.SiteSettings{
		SiteName:               req.SiteName,
		AvatarURL:              req.AvatarURL,
		HeroIntroMD:            req.HeroIntroMD,
		DefaultLocale:          req.DefaultLocale,
		CommentEnabled:         req.CommentEnabled,
		CommentRequireApproval: req.CommentRequireApproval,
		SiteDescription:        req.SiteDescription,
		SeoKeywords:            req.SeoKeywords,
		OgImageURL:             req.OgImageURL,
	})
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: updated})
}

func (h *AdminHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	if err := h.authService.ChangePassword(req.OldPassword, req.NewPassword); err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success"})
}

func (h *AdminHandler) GetTranslationPolicy(c *gin.Context) {
	policy := h.siteService.GetTranslationPolicy()
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: policy})
}

func (h *AdminHandler) UpdateTranslationPolicy(c *gin.Context) {
	var policy domain.TranslationPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	if err := h.siteService.UpdateTranslationPolicy(policy); err != nil {
		badRequest(c, "UPDATE_FAILED", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: policy})
}

func (h *AdminHandler) CreateFooterItem(c *gin.Context) {
	var req dto.FooterItemUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.siteService.CreateFooterItem(domain.FooterItem{
		Label:               req.Label,
		LinkType:            req.LinkType,
		InternalArticleSlug: req.InternalArticleSlug,
		ExternalURL:         req.ExternalURL,
		RowNum:              req.RowNum,
		OrderNum:            req.OrderNum,
		Enabled:             req.Enabled,
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) UpdateFooterItem(c *gin.Context) {
	id := c.Param("id")
	var req dto.FooterItemUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.siteService.UpdateFooterItem(id, domain.FooterItem{
		Label:               req.Label,
		LinkType:            req.LinkType,
		InternalArticleSlug: req.InternalArticleSlug,
		ExternalURL:         req.ExternalURL,
		RowNum:              req.RowNum,
		OrderNum:            req.OrderNum,
		Enabled:             req.Enabled,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrFooterItemNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "FOOTER_ITEM_NOT_FOUND", Message: "footer item not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: item})
}

func (h *AdminHandler) DeleteFooterItem(c *gin.Context) {
	id := c.Param("id")
	if !h.siteService.DeleteFooterItem(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "FOOTER_ITEM_NOT_FOUND", Message: "footer item not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) CreateSocialLink(c *gin.Context) {
	var req dto.SocialLinkUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.siteService.CreateSocialLink(domain.SocialLink{
		Platform: req.Platform,
		Title:    req.Title,
		URL:      req.URL,
		IconKey:  req.IconKey,
		OrderNum: req.OrderNum,
		Enabled:  req.Enabled,
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) UpdateSocialLink(c *gin.Context) {
	id := c.Param("id")
	var req dto.SocialLinkUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.siteService.UpdateSocialLink(id, domain.SocialLink{
		Platform: req.Platform,
		Title:    req.Title,
		URL:      req.URL,
		IconKey:  req.IconKey,
		OrderNum: req.OrderNum,
		Enabled:  req.Enabled,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrSocialLinkNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SOCIAL_LINK_NOT_FOUND", Message: "social link not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: item})
}

func (h *AdminHandler) DeleteSocialLink(c *gin.Context) {
	id := c.Param("id")
	if !h.siteService.DeleteSocialLink(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SOCIAL_LINK_NOT_FOUND", Message: "social link not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) CreateNavItem(c *gin.Context) {
	var req dto.NavItemUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.siteService.CreateNavItem(domain.NavItem{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Key:         req.Key,
		Type:        req.Type,
		TargetType:  req.TargetType,
		TargetValue: req.TargetValue,
		OrderNum:    req.OrderNum,
		Enabled:     req.Enabled,
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) UpdateNavItem(c *gin.Context) {
	id := c.Param("id")
	var req dto.NavItemUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.siteService.UpdateNavItem(id, domain.NavItem{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Key:         req.Key,
		Type:        req.Type,
		TargetType:  req.TargetType,
		TargetValue: req.TargetValue,
		OrderNum:    req.OrderNum,
		Enabled:     req.Enabled,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrNavItemNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "NAV_ITEM_NOT_FOUND", Message: "nav item not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: item})
}

func (h *AdminHandler) DeleteNavItem(c *gin.Context) {
	id := c.Param("id")
	if !h.siteService.DeleteNavItem(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "NAV_ITEM_NOT_FOUND", Message: "nav item not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) CreateSlot(c *gin.Context) {
	var req dto.SlotCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	slot, err := h.siteService.CreateContentSlot(domain.ContentSlot{
		SlotKey:     req.SlotKey,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": slot.ID}})
}

func (h *AdminHandler) ListSlots(c *gin.Context) {
	rows := h.siteService.ListContentSlots()
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

func (h *AdminHandler) CreateSlotItem(c *gin.Context) {
	slotKey := c.Param("slotKey")
	var req dto.SlotItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.siteService.CreateSlotItem(slotKey, domain.SlotItem{
		ContentType: req.ContentType,
		ContentID:   req.ContentID,
		OrderNum:    req.OrderNum,
		Enabled:     req.Enabled,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrSlotNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SLOT_NOT_FOUND", Message: "slot not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) ListSlotItems(c *gin.Context) {
	slotKey := c.Param("slotKey")
	rows, ok := h.siteService.ListSlotItems(slotKey)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SLOT_NOT_FOUND", Message: "slot not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

func (h *AdminHandler) DeleteSlotItem(c *gin.Context) {
	slotKey := c.Param("slotKey")
	id := c.Param("id")
	if !h.siteService.DeleteSlotItem(slotKey, id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SLOT_ITEM_NOT_FOUND", Message: "slot item not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) ListIntegrations(c *gin.Context) {
	providerType := c.Query("providerType")
	rows := h.integrationService.ListIntegrationProviders(providerType)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

func (h *AdminHandler) UpdateIntegration(c *gin.Context) {
	providerKey := c.Param("providerKey")
	var req dto.IntegrationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	configJSON, err := json.Marshal(req.ConfigJSON)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", "configJson must be JSON object")
		return
	}
	metaJSON, err := json.Marshal(req.MetaJSON)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", "metaJson must be JSON object")
		return
	}
	_, err = h.integrationService.UpdateIntegrationProvider(providerKey, req.Enabled, configJSON, metaJSON)
	if err != nil {
		if errors.Is(err, apperr.ErrProviderNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PROVIDER_NOT_FOUND", Message: "provider not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) TestIntegration(c *gin.Context) {
	providerKey := c.Param("providerKey")
	result, err := h.integrationService.TestIntegrationProvider(providerKey)
	if err != nil {
		if errors.Is(err, apperr.ErrProviderNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PROVIDER_NOT_FOUND", Message: "provider not found"})
			return
		}
		response.JSON(c, http.StatusBadRequest, response.Envelope{Code: "PROVIDER_TEST_FAILED", Message: err.Error()})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: result})
}

func (h *AdminHandler) CreateTranslationJob(c *gin.Context) {
	var req dto.CreateTranslationJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	user := middleware.MustUser(c)
	job, err := h.translationService.CreateTranslationJob(domain.TranslationJob{
		SourceType:   req.SourceType,
		SourceID:     req.SourceID,
		SourceLocale: req.SourceLocale,
		TargetLocale: req.TargetLocale,
		ProviderKey:  req.ProviderKey,
		ModelName:    req.ModelName,
		MaxRetries:   req.MaxRetries,
		AutoPublish:  req.AutoPublish,
		PublishAt:    req.PublishAt,
		RequestedBy:  user.ID,
	})
	if err != nil {
		if errors.Is(err, apperr.ErrProviderNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PROVIDER_NOT_FOUND", Message: "provider not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": job.ID}})
}

func (h *AdminHandler) ListTranslationJobs(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	status := c.Query("status")
	sourceType := c.Query("sourceType")
	sourceID := c.Query("sourceId")
	rows, total := h.translationService.ListTranslationJobs(page, pageSize, status, sourceType, sourceID)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.TranslationJob]{Total: total, Rows: rows}})
}

func (h *AdminHandler) TranslationJobDetail(c *gin.Context) {
	id := c.Param("id")
	job, ok := h.translationService.GetTranslationJobByID(id)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "TRANSLATION_JOB_NOT_FOUND", Message: "translation job not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: job})
}

func (h *AdminHandler) RetryTranslationJob(c *gin.Context) {
	id := c.Param("id")
	job, err := h.translationService.RetryTranslationJob(id)
	if err != nil {
		if errors.Is(err, apperr.ErrTranslationJobNotFound) {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "TRANSLATION_JOB_NOT_FOUND", Message: "translation job not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: job})
}

func (h *AdminHandler) ListTranslationContents(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	sourceType := c.Query("sourceType")
	sourceID := c.Query("sourceId")
	locale := c.Query("locale")
	rows, total, err := h.translationService.ListTranslationContents(page, pageSize, sourceType, sourceID, locale)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.TranslationContent]{Total: total, Rows: rows}})
}

func (h *AdminHandler) TranslationContentDetail(c *gin.Context) {
	sourceType := c.Param("sourceType")
	sourceID := c.Param("sourceId")
	locale := c.Param("locale")
	row, ok, err := h.translationService.GetTranslationContent(sourceType, sourceID, locale)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "TRANSLATION_CONTENT_NOT_FOUND", Message: "translation content not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: row})
}

func (h *AdminHandler) UpsertTranslationContent(c *gin.Context) {
	var req dto.UpsertTranslationContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	row, err := h.translationService.UpsertTranslationContent(
		req.SourceType,
		req.SourceID,
		req.Locale,
		req.Title,
		req.Summary,
		req.Content,
		req.Status,
		req.PublishedAt,
		req.TranslatedByJobID,
	)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: row})
}

type aiSummaryRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	ProviderKey string `json:"providerKey"`
	ModelName   string `json:"modelName"`
	MaxLength   int    `json:"maxLength"`
}

func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	cat, err := h.articleService.CreateCategory(domain.Category{Name: req.Name, Slug: req.Slug})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: cat})
}

func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if !h.articleService.DeleteCategory(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "NOT_FOUND", Message: "category not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success"})
}

func (h *AdminHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	tag, err := h.articleService.CreateTag(domain.Tag{Name: req.Name, Slug: req.Slug})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: tag})
}

func (h *AdminHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if !h.articleService.DeleteTag(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "NOT_FOUND", Message: "tag not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success"})
}

func (h *AdminHandler) GenerateSummary(c *gin.Context) {
	var req aiSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	summary, fallbackUsed, fallbackReason, err := h.aiAssistService.GenerateSummary(c.Request.Context(), req.Title, req.Content, req.ProviderKey, req.ModelName, req.MaxLength)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "INTERNAL_ERROR", Message: "failed to generate summary"})
		return
	}
	data := map[string]any{
		"summary":      summary,
		"fallbackUsed": fallbackUsed,
	}
	if fallbackReason != "" {
		data["fallbackReason"] = fallbackReason
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: data})
}

type aiSlugRequest struct {
	Title       string `json:"title"`
	ProviderKey string `json:"providerKey"`
	ModelName   string `json:"modelName"`
}

func (h *AdminHandler) SuggestSlug(c *gin.Context) {
	var req aiSlugRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	slug, fallbackUsed, fallbackReason, err := h.aiAssistService.SuggestSlug(c.Request.Context(), req.Title, req.ProviderKey, req.ModelName)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "INTERNAL_ERROR", Message: "failed to suggest slug"})
		return
	}
	data := map[string]any{
		"slug":         slug,
		"fallbackUsed": fallbackUsed,
	}
	if fallbackReason != "" {
		data["fallbackReason"] = fallbackReason
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: data})
}

// GetTMDBMetadata fetches movie/TV metadata from TMDB API
func (h *AdminHandler) GetTMDBMetadata(c *gin.Context) {
	mediaType := c.Param("type")
	id := c.Param("id")

	if mediaType != "movie" && mediaType != "tv" {
		badRequest(c, "INVALID_TYPE", "type must be movie or tv")
		return
	}

	data, err := h.tmdbService.GetMetadata(mediaType, id)
	if err != nil {
		badRequest(c, "TMDB_ERROR", err.Error())
		return
	}

	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: data})
}
