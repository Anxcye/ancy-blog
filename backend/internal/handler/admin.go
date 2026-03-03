// File: admin.go
// Purpose: Implement authenticated admin APIs for content and site management.
// Module: backend/internal/handler, admin HTTP presentation layer.
// Related: service.ContentService and auth middleware.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/middleware"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	contentService *service.ContentService
}

func NewAdminHandler(contentService *service.ContentService) *AdminHandler {
	return &AdminHandler{contentService: contentService}
}

func (h *AdminHandler) CreateArticle(c *gin.Context) {
	var req domain.Article
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	article, err := h.contentService.CreateArticle(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": article.ID}})
}

func (h *AdminHandler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var req domain.Article
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	article, err := h.contentService.UpdateArticle(id, req)
	if err != nil {
		if err.Error() == "article not found" {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "ARTICLE_NOT_FOUND", Message: "article not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": article.ID}})
}

func (h *AdminHandler) CreateMoment(c *gin.Context) {
	var req domain.Moment
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	moment, err := h.contentService.CreateMoment(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": moment.ID}})
}

func (h *AdminHandler) CommentPage(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	status := c.Query("status")
	rows, total := h.contentService.ListCommentPage(page, pageSize, status)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Comment]{Total: total, Rows: rows}})
}

type commentUpdateRequest struct {
	Status   string `json:"status"`
	IsPinned string `json:"isPinned"`
}

func (h *AdminHandler) CommentUpdate(c *gin.Context) {
	id := c.Param("id")
	var req commentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	comment, err := h.contentService.UpdateCommentAdmin(id, req.Status, req.IsPinned)
	if err != nil {
		if err.Error() == "comment not found" {
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
	rows, total := h.contentService.ListLinkSubmissions(page, pageSize, reviewStatus)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Link]{Total: total, Rows: rows}})
}

type reviewLinkRequest struct {
	ReviewStatus     string `json:"reviewStatus"`
	ReviewNote       string `json:"reviewNote"`
	RelatedArticleID string `json:"relatedArticleId"`
}

func (h *AdminHandler) ReviewLink(c *gin.Context) {
	id := c.Param("id")
	var req reviewLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	_, err := h.contentService.ReviewLink(id, req.ReviewStatus, req.ReviewNote, req.RelatedArticleID)
	if err != nil {
		if err.Error() == "link not found" {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "LINK_NOT_FOUND", Message: "link not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) UpdateSiteSettings(c *gin.Context) {
	var req domain.SiteSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	updated := h.contentService.UpdateSiteSettings(req)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: updated})
}

func (h *AdminHandler) CreateFooterItem(c *gin.Context) {
	var req domain.FooterItem
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.contentService.CreateFooterItem(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) UpdateFooterItem(c *gin.Context) {
	id := c.Param("id")
	var req domain.FooterItem
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.contentService.UpdateFooterItem(id, req)
	if err != nil {
		if err.Error() == "footer item not found" {
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
	if !h.contentService.DeleteFooterItem(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "FOOTER_ITEM_NOT_FOUND", Message: "footer item not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) CreateSocialLink(c *gin.Context) {
	var req domain.SocialLink
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.contentService.CreateSocialLink(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) UpdateSocialLink(c *gin.Context) {
	id := c.Param("id")
	var req domain.SocialLink
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.contentService.UpdateSocialLink(id, req)
	if err != nil {
		if err.Error() == "social link not found" {
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
	if !h.contentService.DeleteSocialLink(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SOCIAL_LINK_NOT_FOUND", Message: "social link not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) CreateNavItem(c *gin.Context) {
	var req domain.NavItem
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.contentService.CreateNavItem(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) UpdateNavItem(c *gin.Context) {
	id := c.Param("id")
	var req domain.NavItem
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.contentService.UpdateNavItem(id, req)
	if err != nil {
		if err.Error() == "nav item not found" {
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
	if !h.contentService.DeleteNavItem(id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "NAV_ITEM_NOT_FOUND", Message: "nav item not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) CreateSlot(c *gin.Context) {
	var req domain.ContentSlot
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	slot, err := h.contentService.CreateContentSlot(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": slot.ID}})
}

func (h *AdminHandler) CreateSlotItem(c *gin.Context) {
	slotKey := c.Param("slotKey")
	var req domain.SlotItem
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	item, err := h.contentService.CreateSlotItem(slotKey, req)
	if err != nil {
		if err.Error() == "slot not found" {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SLOT_NOT_FOUND", Message: "slot not found"})
			return
		}
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": item.ID}})
}

func (h *AdminHandler) DeleteSlotItem(c *gin.Context) {
	slotKey := c.Param("slotKey")
	id := c.Param("id")
	if !h.contentService.DeleteSlotItem(slotKey, id) {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SLOT_ITEM_NOT_FOUND", Message: "slot item not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: true})
}

func (h *AdminHandler) ListIntegrations(c *gin.Context) {
	providerType := c.Query("providerType")
	rows := h.contentService.ListIntegrationProviders(providerType)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

type integrationUpdateRequest struct {
	Enabled    bool           `json:"enabled"`
	ConfigJSON map[string]any `json:"configJson"`
	MetaJSON   map[string]any `json:"metaJson"`
}

func (h *AdminHandler) UpdateIntegration(c *gin.Context) {
	providerKey := c.Param("providerKey")
	var req integrationUpdateRequest
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
	_, err = h.contentService.UpdateIntegrationProvider(providerKey, req.Enabled, configJSON, metaJSON)
	if err != nil {
		if err.Error() == "provider not found" {
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
	result, err := h.contentService.TestIntegrationProvider(providerKey)
	if err != nil {
		if err.Error() == "provider not found" {
			response.JSON(c, http.StatusNotFound, response.Envelope{Code: "PROVIDER_NOT_FOUND", Message: "provider not found"})
			return
		}
		response.JSON(c, http.StatusBadRequest, response.Envelope{Code: "PROVIDER_TEST_FAILED", Message: err.Error()})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: result})
}

type createTranslationJobRequest struct {
	SourceType   string `json:"sourceType"`
	SourceID     string `json:"sourceId"`
	SourceLocale string `json:"sourceLocale"`
	TargetLocale string `json:"targetLocale"`
	ProviderKey  string `json:"providerKey"`
	ModelName    string `json:"modelName"`
}

func (h *AdminHandler) CreateTranslationJob(c *gin.Context) {
	var req createTranslationJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	user := middleware.MustUser(c)
	job, err := h.contentService.CreateTranslationJob(domain.TranslationJob{
		SourceType:   req.SourceType,
		SourceID:     req.SourceID,
		SourceLocale: req.SourceLocale,
		TargetLocale: req.TargetLocale,
		ProviderKey:  req.ProviderKey,
		ModelName:    req.ModelName,
		RequestedBy:  user.ID,
	})
	if err != nil {
		if err.Error() == "provider not found" {
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
	rows, total := h.contentService.ListTranslationJobs(page, pageSize, status, sourceType, sourceID)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.TranslationJob]{Total: total, Rows: rows}})
}

func (h *AdminHandler) TranslationJobDetail(c *gin.Context) {
	id := c.Param("id")
	job, ok := h.contentService.GetTranslationJobByID(id)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "TRANSLATION_JOB_NOT_FOUND", Message: "translation job not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: job})
}
