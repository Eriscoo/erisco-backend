# Crawler Meta Tags (`/og/:slug` & `/og/page/:page`)

## Overview

The `/og/:slug` endpoint generates a static HTML page containing **Open Graph** and **Twitter Card** meta tags for each blog post. Unlike other API endpoints, this route does **not** return JSON — it returns a full HTML document (`text/html; charset=utf-8`) so social media crawlers (Facebook, Twitter/X, LinkedIn, WhatsApp, Discord, etc.) can read it directly without executing JavaScript.

## Why It's Needed

The Erisco frontend is a React SPA where all content is rendered client-side. When a post link is shared on social media, crawlers only read the raw HTML from the server without running JavaScript. Without this endpoint, crawlers would see only the default meta tags from `index.html` (`"Developer thoughts, project showcases..."`) for every page, rather than the specific meta tags for the shared post.

The solution: **server-side rendering** of meta tags via `/og/:slug`, with the frontend pointing crawlers to this endpoint through the `<meta property="og:url">` tag.

## How It Works

```
User shares https://eriscoo.com/post-slug
       │
       ▼
Crawler reads <meta property="og:url" content="...">
       │
       ▼
Crawler fetches https://eriscoo.com/og/post-slug
       │
       ▼
Backend fetches post data from the database by slug
       │
       ▼
Backend generates HTML with post-specific meta tags
       │
       ▼
Crawler renders preview card with post title, description, and image
```

## Generated Meta Tags

| Tag | Value | Data Source |
|-----|-------|-------------|
| `<title>` | `{Post Title} \| Eriscoo` | `post.Title` |
| `meta name="description"` | First 200 characters of body (HTML stripped) | `post.Body` |
| `og:site_name` | `Eriscoo` | Hardcoded |
| `og:type` | `article` | Hardcoded (article pages only) |
| `og:title` | `{Post Title}` | `post.Title` |
| `og:description` | First 200 characters of body (HTML stripped, escaped) | `post.Body` |
| `og:url` | `{SITE_URL}/{slug}` | `siteURL` + `slug` |
| `og:image` | `{SITE_URL}/uploads/post-cover/{filename}` | `post.ImageURL` |
| `og:image:width` | `1200` | Hardcoded |
| `og:image:height` | `630` | Hardcoded |
| `og:image:alt` | `{Post Title}` | `post.Title` (if image exists) |
| `article:published_time` | RFC3339 | `post.PublishedAt` or `post.CreatedAt` |
| `article:tag` | One tag per meta element, comma-separated | `post.TagNames` |
| `twitter:card` | `summary_large_image` | Hardcoded |
| `twitter:title` | `{Post Title}` | `post.Title` |
| `twitter:description` | First 200 characters of body | `post.Body` |
| `twitter:image` | `{SITE_URL}/uploads/post-cover/{filename}` | `post.ImageURL` |
| `<link rel="canonical">` | `{SITE_URL}/{slug}` | `siteURL` + `slug` |

## Supported Features

- **Social media preview**: Facebook, Twitter/X, LinkedIn, WhatsApp, Telegram, Discord — all platforms that read Open Graph and/or Twitter Card tags.
- **SEO**: Search engines read `<title>`, `<meta name="description">`, and `<link rel="canonical">` from this page.
- **Article metadata**: `article:published_time` and `article:tag` help with content classification on Facebook.

## Pages That Use This

- **Post detail pages** (`/post/:slug`) use `GET /og/:slug` — meta tags are dynamically generated from the database.
- **Static pages** (Portfolio, About, Terms, Privacy) use `GET /og/page/:page` — meta tags are loaded from a JSON config file.

### Frontend Integration

**For articles**, post pages should set `og:url` to the crawler endpoint:

```html
<meta property="og:url" content="https://eriscoo.com/og/post-slug" />
```

**For static pages** (portfolio, about, etc.), set `og:url` accordingly:

```html
<meta property="og:url" content="https://eriscoo.com/og/page/portfolio" />
```

So when crawlers visit the SPA page, they are redirected to the correct OG endpoint to get the proper meta tags.

## Static Pages Configuration (`og-pages.json`)

Static page meta tags are configured in a JSON file (default: `./og-pages.json`, overridable via `OG_PAGES_CONFIG` env var):

```json
{
  "portfolio": {
    "title": "Portfolio | Eriscoo",
    "description": "My projects, work, and side builds.",
    "image": "/og-image.jpg"
  },
  "about": {
    "title": "About | Eriscoo",
    "description": "Learn more about me and what I do.",
    "image": "/og-image.jpg"
  }
}
```

| Field | Description |
|-------|-------------|
| `title` | Page title (include `\| Eriscoo` suffix if desired) |
| `description` | Meta description (plain text, no HTML) — will be HTML-escaped |
| `image` | Image path relative to `SITE_URL` (e.g., `/og-image.jpg` resolves to `{SITE_URL}/og-image.jpg`) |

To add a new static page, just add a new entry to this JSON file and restart the server. No code changes needed.

### `og:type` for Static Pages

The handler always sets `og:type` to `website` for static pages (as opposed to `article` for blog posts). No `article:published_time` or `article:tag` tags are emitted.

## Technical Specs

| Property | `/og/:slug` | `/og/page/:page` |
|----------|-------------|-------------------|
| **Handler file** | `internal/transport/handler/og/handler.go` | Same |
| **Auth** | None (public) | None (public) |
| **Response Content-Type** | `text/html; charset=utf-8` | Same |
| **Error responses** | `404` (empty/unknown slug), `500` | `404` (unknown page) |
| **Listed in Swagger** | No (not a JSON API endpoint) | No |
| **`og:type`** | `article` | `website` |
| **Config source** | Database (posts table) | `og-pages.json` file |

### Description Processing (Strip & Truncate)

The description is extracted from `post.Body` (HTML content), then:
1. Strip all HTML tags using regex
2. Trim whitespace and normalize spaces
3. Truncate to a maximum of **200 characters** (runes, not bytes — safe for non-ASCII characters)
4. Append `...` if truncated
5. HTML-escape the final result

## Related Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /og/:slug` | Generate HTML meta tags for blog post articles |
| `GET /og/page/:page` | Generate HTML meta tags for static pages (portfolio, about, etc.) |
| `GET /api/v1/public/posts/:slug` | Get post details as JSON (for the frontend) |
| `GET /api/v1/public/posts/all` | Get all published posts (for listing pages) |
