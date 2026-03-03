// File: public.go
// Purpose: Implement public read APIs for articles, moments, links, site, taxonomy, and timeline.
// Module: backend/internal/handler, public HTTP presentation layer.
// Related: service.ContentService and API contract public route group.
package handler

import (
	"net/http"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	contentService *service.ContentService
}

func NewPublicHandler(contentService *service.ContentService) *PublicHandler {
	return &PublicHandler{contentService: contentService}
}

func (h *PublicHandler) Articles(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	category := c.Query("category")
	tag := c.Query("tag")
	contentKind := c.DefaultQuery("contentKind", "post")
	rows, total := h.contentService.ListPublishedArticles(page, pageSize, category, tag, contentKind)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Article]{Total: total, Rows: rows}})
}

func (h *PublicHandler) ArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")
	article, ok := h.contentService.GetPublishedArticleBySlug(slug)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "ARTICLE_NOT_FOUND", Message: "article not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: article})
}

func (h *PublicHandler) ArticleByCategory(c *gin.Context) {
	categorySlug := c.Param("categorySlug")
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.contentService.ListPublishedArticles(page, pageSize, categorySlug, "", "post")
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Article]{Total: total, Rows: rows}})
}

func (h *PublicHandler) Moments(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.contentService.ListPublishedMoments(page, pageSize)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Moment]{Total: total, Rows: rows}})
}

func (h *PublicHandler) CommentByArticle(c *gin.Context) {
	articleID := c.Param("articleId")
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.contentService.ListArticleComments(articleID, page, pageSize)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Comment]{Total: total, Rows: rows}})
}

func (h *PublicHandler) CommentChildren(c *gin.Context) {
	parentID := c.Param("id")
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.contentService.ListCommentChildren(parentID, page, pageSize)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Comment]{Total: total, Rows: rows}})
}

func (h *PublicHandler) CommentArticleTotal(c *gin.Context) {
	articleID := c.Param("articleId")
	total, err := h.contentService.CountArticleComments(articleID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "INTERNAL_ERROR", Message: "failed to count comments"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: total})
}

func (h *PublicHandler) AddComment(c *gin.Context) {
	var req domain.Comment
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	req.IP = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")
	comment, err := h.contentService.CreateComment(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": comment.ID}})
}

func (h *PublicHandler) SubmitLink(c *gin.Context) {
	var req domain.Link
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	req.SubmittedIP = c.ClientIP()
	req.SubmittedUserAgent = c.GetHeader("User-Agent")
	link, err := h.contentService.SubmitLink(req)
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": link.ID}})
}

func (h *PublicHandler) ApprovedLinks(c *gin.Context) {
	rows := h.contentService.ListApprovedLinks()
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

func (h *PublicHandler) Categories(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.contentService.ListCategories()})
}

func (h *PublicHandler) Tags(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.contentService.ListTags()})
}

func (h *PublicHandler) SiteSettings(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.contentService.GetSiteSettings()})
}

func (h *PublicHandler) SiteFooter(c *gin.Context) {
	items := h.contentService.ListFooterItems()
	grouped := map[int][]domain.FooterItem{}
	for _, item := range items {
		grouped[item.RowNum] = append(grouped[item.RowNum], item)
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: grouped})
}

func (h *PublicHandler) SiteSocialLinks(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.contentService.ListSocialLinks()})
}

func (h *PublicHandler) SiteNav(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.contentService.ListNavItems()})
}

func (h *PublicHandler) SiteSlotContent(c *gin.Context) {
	slotKey := c.Param("slotKey")
	limit := getIntQuery(c, "limit", 0)
	rows, ok := h.contentService.ListSlotContent(slotKey, limit)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SLOT_NOT_FOUND", Message: "slot not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

func (h *PublicHandler) Timeline(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.contentService.ListTimeline(page, pageSize)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.TimelineItem]{Total: total, Rows: rows}})
}
