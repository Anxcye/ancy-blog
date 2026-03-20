// File: public.go
// Purpose: Implement public read APIs for articles, moments, links, site, taxonomy, and timeline.
// Module: backend/internal/handler, public HTTP presentation layer.
// Related: service module facades and API contract public route group.
package handler

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/apperr"
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

func clientIPFromRequest(c *gin.Context) string {
	if ip := strings.TrimSpace(c.GetHeader("CF-Connecting-IP")); ip != "" {
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}
	if xff := strings.TrimSpace(c.GetHeader("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		for _, part := range parts {
			ip := strings.TrimSpace(part)
			if parsed := net.ParseIP(ip); parsed != nil {
				return parsed.String()
			}
		}
	}
	if ip := strings.TrimSpace(c.GetHeader("X-Real-IP")); ip != "" {
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed.String()
		}
	}
	return c.ClientIP()
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
	}(article.ID, clientIPFromRequest(c), c.GetHeader("User-Agent"))
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

func (h *PublicHandler) MomentByID(c *gin.Context) {
	id := c.Param("id")
	locale := c.Query("locale")
	moment, ok := h.articleService.GetPublishedMomentByID(id, locale)
	if !ok {
		response.JSON(c, http.StatusNotFound, response.Envelope{Code: "MOMENT_NOT_FOUND", Message: "moment not found"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: moment})
}

func (h *PublicHandler) CommentByArticle(c *gin.Context) {
	h.renderCommentThreads(c, "article", c.Param("articleId"))
}

func (h *PublicHandler) CommentByContent(c *gin.Context) {
	contentType, contentID, ok := publicCommentTarget(c)
	if !ok {
		badRequest(c, "VALIDATION_ERROR", "invalid content target")
		return
	}
	h.renderCommentThreads(c, contentType, contentID)
}

func (h *PublicHandler) renderCommentThreads(c *gin.Context, contentType, contentID string) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.commentService.ListContentCommentThreads(contentType, contentID, page, pageSize)
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[dto.PublicComment]{Total: total, Rows: mapPublicComments(rows)}})
}

func (h *PublicHandler) CommentChildren(c *gin.Context) {
	parentID := c.Param("id")
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)
	rows, total := h.commentService.ListCommentChildren(parentID, page, pageSize)
	items := make([]dto.PublicComment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPublicComment(domain.CommentNode{Comment: row}))
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: pageResult[dto.PublicComment]{Total: total, Rows: items}})
}

func (h *PublicHandler) CommentArticleTotal(c *gin.Context) {
	h.renderCommentTotal(c, "article", c.Param("articleId"))
}

func (h *PublicHandler) CommentContentTotal(c *gin.Context) {
	contentType, contentID, ok := publicCommentTarget(c)
	if !ok {
		badRequest(c, "VALIDATION_ERROR", "invalid content target")
		return
	}
	h.renderCommentTotal(c, contentType, contentID)
}

func (h *PublicHandler) renderCommentTotal(c *gin.Context, contentType, contentID string) {
	total, err := h.commentService.CountContentComments(contentType, contentID)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, response.Envelope{Code: "INTERNAL_ERROR", Message: "failed to count comments"})
		return
	}
	response.JSON(c, http.StatusOK, response.Envelope{Code: "OK", Message: "success", Data: total})
}

func publicCommentTarget(c *gin.Context) (string, string, bool) {
	contentType := strings.TrimSpace(c.Param("contentType"))
	contentID := strings.TrimSpace(c.Param("contentId"))
	if contentType != "article" && contentType != "moment" {
		return "", "", false
	}
	if contentID == "" {
		return "", "", false
	}
	return contentType, contentID, true
}

func (h *PublicHandler) AddComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "VALIDATION_ERROR", "invalid request body")
		return
	}
	contentType := req.ContentType
	contentID := req.ContentID
	if contentType == "" && req.ArticleID != "" {
		contentType = "article"
		contentID = req.ArticleID
	}
	clientIP := clientIPFromRequest(c)
	comment, err := h.commentService.CreateComment(domain.Comment{
		ArticleID:   req.ArticleID,
		ContentType: contentType,
		ContentID:   contentID,
		ParentID:    req.ParentID,
		RootID:      req.RootID,
		Content:     req.Content,
		Nickname:    req.Nickname,
		Email:       req.Email,
		Website:     req.Website,
		AvatarURL:   req.AvatarURL,
		Source:      req.Source,
		ToCommentID: req.ToCommentID,
		IP:          clientIP,
		UserAgent:   c.GetHeader("User-Agent"),
	})
	if err != nil {
		if errors.Is(err, apperr.ErrLinkSubmissionDisabled) {
			response.JSON(c, http.StatusForbidden, response.Envelope{Code: "LINK_SUBMISSION_DISABLED", Message: "link submission is disabled"})
			return
		}
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
	clientIP := clientIPFromRequest(c)
	link, err := h.linkService.SubmitLink(domain.Link{
		Name:               req.Name,
		URL:                req.URL,
		AvatarURL:          req.AvatarURL,
		Description:        req.Description,
		ContactEmail:       req.ContactEmail,
		SubmittedIP:        clientIP,
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

func mapPublicComments(rows []domain.CommentNode) []dto.PublicComment {
	items := make([]dto.PublicComment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPublicComment(row))
	}
	return items
}

func mapPublicComment(row domain.CommentNode) dto.PublicComment {
	return dto.PublicComment{
		ID:                row.ID,
		ArticleID:         row.ArticleID,
		ContentType:       row.ContentType,
		ContentID:         row.ContentID,
		ParentID:          row.ParentID,
		RootID:            row.RootID,
		Content:           row.Content,
		Status:            row.Status,
		IsPinned:          row.IsPinned == "1",
		IsAuthor:          isCommentAuthor(row.Comment),
		LikeCount:         row.LikeCount,
		ReplyCount:        row.ReplyCount,
		Nickname:          row.Nickname,
		Website:           row.Website,
		AvatarURL:         row.AvatarURL,
		ToCommentID:       row.ToCommentID,
		ToCommentNickname: row.ToCommentNickname,
		ToCommentIsAuthor: row.ToCommentIsAuthor,
		CreatedAt:         row.CreatedAt,
		Children:          mapPublicComments(row.Children),
	}
}

func isCommentAuthor(comment domain.Comment) bool {
	return strings.EqualFold(comment.Source, "admin") || strings.EqualFold(comment.Nickname, "admin")
}
