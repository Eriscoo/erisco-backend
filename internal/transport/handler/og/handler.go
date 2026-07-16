package og

import (
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	postssvc "github.com/eriscoo/blog-backend/internal/application/posts"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

var reTags = regexp.MustCompile(`<[^>]+>`)

type PageMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func LoadPages(path string) (map[string]PageMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pages map[string]PageMeta
	if err := json.Unmarshal(data, &pages); err != nil {
		return nil, err
	}
	return pages, nil
}

type OGHandler struct {
	svc     *postssvc.Service
	siteURL string
	pages   map[string]PageMeta
}

func New(svc *postssvc.Service, siteURL string, pages map[string]PageMeta) *OGHandler {
	return &OGHandler{svc: svc, siteURL: siteURL, pages: pages}
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

	title := html.EscapeString(post.Title) + " | Eriscoo"
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

	extra := ogImageAlt(imageURL, title) + "\n" +
		`<meta property="article:published_time" content="` + publishedAt + `">` + "\n" +
		ogTags(tags)

	h.renderHTML(c, title, description, url, imageURL, "article", extra)
}

func (h *OGHandler) HandleStaticPage(c *gin.Context) {
	page := c.Param("page")
	if page == "" {
		c.String(http.StatusNotFound, "not found")
		return
	}

	meta, ok := h.pages[page]
	if !ok {
		c.String(http.StatusNotFound, "not found")
		return
	}

	title := html.EscapeString(meta.Title)
	description := html.EscapeString(meta.Description)
	url := h.siteURL + "/" + page
	imageURL := ""
	if meta.Image != "" {
		imageURL = h.siteURL + meta.Image
	}

	h.renderHTML(c, title, description, url, imageURL, "website", "")
}

func (h *OGHandler) renderHTML(c *gin.Context, title, description, url, imageURL, ogType, extra string) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>` + title + `</title>
<meta name="description" content="` + description + `">
<meta property="og:site_name" content="Eriscoo">
<meta property="og:type" content="` + ogType + `">
<meta property="og:title" content="` + title + `">
<meta property="og:description" content="` + description + `">
<meta property="og:url" content="` + url + `">
<meta property="og:image" content="` + imageURL + `">
<meta property="og:image:width" content="1200">
<meta property="og:image:height" content="630">
` + extra + `
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
