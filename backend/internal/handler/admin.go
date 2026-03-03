// File: admin.go
// Purpose: Implement authenticated admin APIs for content and site management.
// Module: backend/internal/handler, admin HTTP presentation layer.
// Related: service.ContentService and auth middleware.
package handler

import (
	"net/http"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
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
