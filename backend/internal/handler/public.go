// File: public.go
// Purpose: Implement public read APIs for articles, moments, links, site, taxonomy, and timeline.
// Module: backend/internal/handler, public HTTP presentation layer.
// Related: service module facades and API contract public route group.
package handler

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/handler/dto"
	"github.com/anxcye/ancy-blog/backend/internal/response"
	"github.com/anxcye/ancy-blog/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	articleService  *service.ArticleService
	commentService  *service.CommentService
	linkService     *service.LinkService
	siteService     *service.SiteService
	timelineService *service.TimelineService
}

func NewPublicHandler(
	articleService *service.ArticleService,
	commentService *service.CommentService,
	linkService *service.LinkService,
	siteService *service.SiteService,
	timelineService *service.TimelineService,
) *PublicHandler {
	return &PublicHandler{
		articleService:  articleService,
		commentService:  commentService,
		linkService:     linkService,
		siteService:     siteService,
		timelineService: timelineService,
	}
}

func (h *PublicHandler) Articles(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	category := c.Query("category")
	tag := c.Query("tag")
	contentKind := c.DefaultQuery("contentKind", "post")
	rows, total := h.articleService.ListPublishedArticles(page, pageSize, category, tag, contentKind)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Article]{Total: total, Rows: rows}})
}

func (h *PublicHandler) ArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")
	locale := c.Query("locale")
	article, ok := h.articleService.GetPublishedArticleBySlugWithLocale(slug, locale)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "ARTICLE_NOT_FOUND", Message: "article not found"})
		return
	}
	// Record view asynchronously: fingerprint = SHA-256(ip + ua + date)
	go func(articleID, ip, ua string) {
		date := time.Now().UTC().Format("2006-01-02")
		raw := fmt.Sprintf("%s|%s|%s", ip, ua, date)
		sum := sha256.Sum256([]byte(raw))
		key := fmt.Sprintf("%x", sum)
		_, _ = h.articleService.RecordView(articleID, key)
	}(article.ID, c.ClientIP(), c.GetHeader("User-Agent"))
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: article})
}

func (h *PublicHandler) ArticleByCategory(c *gin.Context) {
	categorySlug := c.Param("categorySlug")
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.articleService.ListPublishedArticles(page, pageSize, categorySlug, "", "post")
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Article]{Total: total, Rows: rows}})
}

func (h *PublicHandler) Moments(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	locale := c.Query("locale")
	rows, total := h.articleService.ListPublishedMoments(page, pageSize, locale)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Moment]{Total: total, Rows: rows}})
}

func (h *PublicHandler) CommentByArticle(c *gin.Context) {
	articleID := c.Param("articleId")
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.commentService.ListArticleComments(articleID, page, pageSize)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Comment]{Total: total, Rows: rows}})
}

func (h *PublicHandler) CommentChildren(c *gin.Context) {
	parentID := c.Param("id")
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.commentService.ListCommentChildren(parentID, page, pageSize)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.Comment]{Total: total, Rows: rows}})
}

func (h *PublicHandler) CommentArticleTotal(c *gin.Context) {
	articleID := c.Param("articleId")
	total, err := h.commentService.CountArticleComments(articleID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "INTERNAL_ERROR", Message: "failed to count comments"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: total})
}

func (h *PublicHandler) AddComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	comment, err := h.commentService.CreateComment(domain.Comment{
		ArticleID:   req.ArticleID,
		ParentID:    req.ParentID,
		RootID:      req.RootID,
		Content:     req.Content,
		Nickname:    req.Nickname,
		Email:       req.Email,
		Website:     req.Website,
		AvatarURL:   req.AvatarURL,
		Source:      req.Source,
		ToCommentID: req.ToCommentID,
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": comment.ID}})
}

func (h *PublicHandler) SubmitLink(c *gin.Context) {
	var req dto.SubmitLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	link, err := h.linkService.SubmitLink(domain.Link{
		Name:               req.Name,
		URL:                req.URL,
		AvatarURL:          req.AvatarURL,
		Description:        req.Description,
		ContactEmail:       req.ContactEmail,
		SubmittedIP:        c.ClientIP(),
		SubmittedUserAgent: c.GetHeader("User-Agent"),
	})
	if err != nil {
		badRequest(c, "VALIDATION_ERROR", err.Error())
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: map[string]string{"id": link.ID}})
}

func (h *PublicHandler) ApprovedLinks(c *gin.Context) {
	rows := h.linkService.ListApprovedLinks()
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

func (h *PublicHandler) Categories(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.articleService.ListCategories()})
}

func (h *PublicHandler) Tags(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.articleService.ListTags()})
}

func (h *PublicHandler) SiteSettings(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.siteService.GetSiteSettings()})
}

func (h *PublicHandler) SiteFooter(c *gin.Context) {
	items := h.siteService.ListFooterItems()
	grouped := map[int][]domain.FooterItem{}
	for _, item := range items {
		grouped[item.RowNum] = append(grouped[item.RowNum], item)
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: grouped})
}

func (h *PublicHandler) SiteSocialLinks(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.siteService.ListSocialLinks()})
}

func (h *PublicHandler) SiteNav(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: h.siteService.ListNavItems()})
}

func (h *PublicHandler) SiteSlotContent(c *gin.Context) {
	slotKey := c.Param("slotKey")
	limit := getIntQuery(c, "limit", 0)
	rows, ok := h.siteService.ListSlotContent(slotKey, limit)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "SLOT_NOT_FOUND", Message: "slot not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: rows})
}

func (h *PublicHandler) Timeline(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	locale := c.Query("locale")
	rows, total := h.timelineService.ListTimeline(page, pageSize, locale)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[domain.TimelineItem]{Total: total, Rows: rows}})
}
