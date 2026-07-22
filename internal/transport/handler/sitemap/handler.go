package sitemap

import (
	"net/http"
	"strings"
	"time"

	postssvc "github.com/eriscoo/blog-backend/internal/application/posts"
	oghandler "github.com/eriscoo/blog-backend/internal/transport/handler/og"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc     *postssvc.Service
	siteURL string
	pages   map[string]oghandler.PageMeta
}

func New(svc *postssvc.Service, siteURL string, pages map[string]oghandler.PageMeta) *Handler {
	return &Handler{svc: svc, siteURL: siteURL, pages: pages}
}

// GetSitemap godoc
// @Summary      Get sitemap
// @Description  Returns XML sitemap for all published posts and static pages
// @Tags         sitemap
// @Produce      xml
// @Success      200  {string}  string  "XML sitemap"
// @Router       /sitemap.xml [get]
func (h *Handler) GetSitemap(c *gin.Context) {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	// Home
	h.writeURL(&sb, h.siteURL, "", "weekly", "1.0")

	// Static pages
	staticPaths := []string{"about", "portfolio", "posts", "contact"}
	for _, p := range staticPaths {
		if _, ok := h.pages[p]; ok {
			h.writeURL(&sb, h.siteURL+"/"+p, "", "monthly", "0.7")
		}
	}

	// Terms & Privacy
	for _, p := range []string{"terms-and-conditions", "privacy-policy"} {
		if _, ok := h.pages[p]; ok {
			h.writeURL(&sb, h.siteURL+"/"+p, "", "monthly", "0.5")
		}
	}

	// Published posts
	posts, err := h.svc.GetAllPublished()
	if err == nil {
		for _, p := range posts {
			lastmod := ""
			if p.PublishedAt != nil {
				lastmod = p.PublishedAt.Format(time.DateOnly)
			} else {
				lastmod = p.CreatedAt.Format(time.DateOnly)
			}
			h.writeURL(&sb, h.siteURL+"/"+p.Slug, lastmod, "monthly", "0.8")
		}
	}

	sb.WriteString(`</urlset>` + "\n")

	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(sb.String()))
}

func (h *Handler) writeURL(sb *strings.Builder, loc, lastmod, changefreq, priority string) {
	sb.WriteString("  <url>\n")
	sb.WriteString("    <loc>" + loc + "</loc>\n")
	if lastmod != "" {
		sb.WriteString("    <lastmod>" + lastmod + "</lastmod>\n")
	}
	if changefreq != "" {
		sb.WriteString("    <changefreq>" + changefreq + "</changefreq>\n")
	}
	if priority != "" {
		sb.WriteString("    <priority>" + priority + "</priority>\n")
	}
	sb.WriteString("  </url>\n")
}
