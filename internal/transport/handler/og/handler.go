package og

import (
	"errors"
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"

	postssvc "github.com/eriscoo/blog-backend/internal/application/posts"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

var reTags = regexp.MustCompile(`<[^>]+>`)

type OGHandler struct {
	svc     *postssvc.Service
	siteURL string
}

func New(svc *postssvc.Service, siteURL string) *OGHandler {
	return &OGHandler{svc: svc, siteURL: siteURL}
}

func (h *OGHandler) HandleOG(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.String(http.StatusNotFound, "not found")
		return
	}

	post, err := h.svc.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.String(http.StatusInternalServerError, "internal error")
		return
	}

	title := html.EscapeString(post.Title)
	description := stripAndTruncate(post.Body, 200)
	description = html.EscapeString(description)
	url := h.siteURL + "/" + slug
	imageURL := ""
	if post.ImageURL != "" {
		imageURL = h.siteURL + post.ImageURL
	}
	tags := ""
	if post.TagNames != "" {
		tags = html.EscapeString(post.TagNames)
	}
	publishedAt := ""
	if post.PublishedAt != nil {
		publishedAt = post.PublishedAt.Format(time.RFC3339)
	} else {
		publishedAt = post.CreatedAt.Format(time.RFC3339)
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>` + title + ` | Eriscoo</title>
<meta name="description" content="` + description + `">
<meta property="og:site_name" content="Eriscoo">
<meta property="og:type" content="article">
<meta property="og:title" content="` + title + `">
<meta property="og:description" content="` + description + `">
<meta property="og:url" content="` + url + `">
<meta property="og:image" content="` + imageURL + `">
<meta property="og:image:width" content="1200">
<meta property="og:image:height" content="630">
` + ogImageAlt(imageURL, title) + `
<meta property="article:published_time" content="` + publishedAt + `">
` + ogTags(tags) + `
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="` + title + `">
<meta name="twitter:description" content="` + description + `">
<meta name="twitter:image" content="` + imageURL + `">
<link rel="canonical" href="` + url + `">
</head>
<body></body>
</html>`))
}

func stripAndTruncate(body string, maxLen int) string {
	if body == "" {
		return ""
	}
	text := reTags.ReplaceAllString(body, "")
	text = strings.TrimSpace(text)
	text = strings.Join(strings.Fields(text), " ")
	if len([]rune(text)) > maxLen {
		runes := []rune(text)
		return string(runes[:maxLen]) + "..."
	}
	return text
}

func ogImageAlt(imageURL, title string) string {
	if imageURL == "" {
		return ""
	}
	return `<meta property="og:image:alt" content="` + html.EscapeString(title) + `">`
}

func ogTags(tags string) string {
	if tags == "" {
		return ""
	}
	parts := strings.Split(tags, ",")
	var lines []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			lines = append(lines, `<meta property="article:tag" content="`+html.EscapeString(t)+`">`)
		}
	}
	return strings.Join(lines, "\n")
}
