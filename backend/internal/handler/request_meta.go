// File: request_meta.go
// Purpose: Provide shared request metadata helpers for client IP and user-agent inference.
// Module: backend/internal/handler, HTTP metadata utility layer.
// Related: public/admin handlers and analytics ingestion flows.
package handler

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func inferUAProfile(userAgent string) (deviceType, browserName, osName string, isBot bool) {
	ua := strings.ToLower(strings.TrimSpace(userAgent))
	if ua == "" {
		return "unknown", "unknown", "unknown", false
	}

	botMarkers := []string{"bot", "crawler", "spider", "slurp", "bingpreview", "headless", "facebookexternalhit", "curl", "wget"}
	for _, marker := range botMarkers {
		if strings.Contains(ua, marker) {
			isBot = true
			break
		}
	}

	deviceType = "desktop"
	switch {
	case isBot:
		deviceType = "bot"
	case strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet"):
		deviceType = "tablet"
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "iphone") || strings.Contains(ua, "android"):
		deviceType = "mobile"
	}

	switch {
	case strings.Contains(ua, "edg/"):
		browserName = "Edge"
	case strings.Contains(ua, "opr/") || strings.Contains(ua, "opera"):
		browserName = "Opera"
	case strings.Contains(ua, "chrome/") && !strings.Contains(ua, "edg/"):
		browserName = "Chrome"
	case strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/"):
		browserName = "Safari"
	case strings.Contains(ua, "firefox/"):
		browserName = "Firefox"
	case strings.Contains(ua, "micromessenger"):
		browserName = "WeChat"
	default:
		browserName = "Other"
	}

	switch {
	case strings.Contains(ua, "windows"):
		osName = "Windows"
	case strings.Contains(ua, "android"):
		osName = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ios"):
		osName = "iOS"
	case strings.Contains(ua, "mac os x") || strings.Contains(ua, "macintosh"):
		osName = "macOS"
	case strings.Contains(ua, "linux"):
		osName = "Linux"
	default:
		osName = "Other"
	}

	return deviceType, browserName, osName, isBot
}
