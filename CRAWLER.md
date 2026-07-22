# Crawler Meta Tags (`/og/:slug`, `/og/page/:page`, `/sitemap.xml`)

## Overview

The OG (Open Graph) endpoints generate static HTML pages containing **Open Graph**, **Twitter Card**, and **JSON-LD structured data** meta tags. Unlike other API endpoints, these routes return a full HTML document (`text/html; charset=utf-8`) so crawlers — both social media and search engines — can read them directly without executing JavaScript.

Supported crawlers: Facebook, Twitter/X, WhatsApp, LinkedIn, Slack, Telegram, Discord, Skype, Google, Bing, Yandex, DuckDuckGo, Naver (Yeti), Baidu, Petal (Huawei), and Seznam.

## Why It's Needed

The Erisco frontend is a React SPA where all content is rendered client-side. When a link is shared on social media or crawled by a search engine, crawlers only read the raw HTML from the server without running JavaScript. Without these endpoints, crawlers would see only the default meta tags from `index.html` for every page.

The solution: **server-side rendering** of full HTML pages via dedicated OG endpoints, with **nginx detecting crawlers by User-Agent** and rewriting requests to the backend.

## How It Works

```
Crawler visits https://eriscoo.com/post-slug
       │
       ▼
Nginx detects crawler User-Agent and rewrites to /og/post-slug
       │
       ▼
Backend fetches post data from the database by slug (or static page config)
       │
       ▼
Backend generates full HTML with meta tags, JSON-LD, and body content
       │
       ▼
Crawler renders preview card / indexes the page
```

For **regular users** (non-crawlers), nginx serves the SPA `index.html` as usual — there is no impact on the normal browsing experience.

## Nginx Crawler Detection

Nginx inspects the `User-Agent` header and rewrites requests for known crawlers to the backend OG endpoints.

**Rewrite logic:**

The nginx configuration sets a variable `$og_redirect` through three conditions (all must be true for the rewrite to trigger):

| # | Condition | Appends |
|---|-----------|---------|
| 1 | User-Agent is a known crawler | `"1"` |
| 2 | Request URI is NOT root (`/`) | `"2"` |
| 3 | Request URI is NOT a static file (`.js`, `.css`, `.jpg`, etc.) | `"3"` |

When `$og_redirect = "123"`, the request is rewritten:

- **Static pages** (`about`, `portfolio`, `posts`, `contact`, `terms-and-conditions`, `privacy-policy`) → `/og/page/$1`
- **Everything else** (assumed to be a blog post slug) → `/og/$1`

**Supported crawler User-Agents:**

| Engine | User-Agent token |
|--------|-----------------|
| Facebook/Instagram | `facebookexternalhit` |
| Twitter/X | `Twitterbot` |
| WhatsApp | `WhatsApp` |
| LinkedIn | `LinkedInBot` |
| Slack | `Slackbot` |
| Telegram | `TelegramBot` |
| Discord | `Discordbot` |
| Skype | `SkypeUriPreview` |
| Google | `Googlebot` |
| Bing (also Yahoo, Ecosia, DuckDuckGo) | `Bingbot` |
| Yandex | `YandexBot` |
| DuckDuckGo (own crawler) | `DuckDuckBot` |
| Naver | `Yeti` |
| Baidu | `Baiduspider` |
| Petal (Huawei) | `PetalBot` |
| Seznam | `SeznamBot` |

## Generated Meta Tags

### Article pages (`/og/:slug`)

| Tag | Value | Data Source |
|-----|-------|-------------|
| `<title>` | `{Post Title} \| Eriscoo` | `post.Title` |
| `meta name="description"` | First 200 chars of body (HTML stripped) | `post.Body` |
| `og:site_name` | `Eriscoo` | Hardcoded |
| `og:type` | `article` | Hardcoded |
| `og:title` | `{Post Title} \| Eriscoo` | `post.Title` |
| `og:description` | First 200 chars of body (HTML stripped, escaped) | `post.Body` |
| `og:url` | `{SITE_URL}/{slug}` | `siteURL` + `slug` |
| `og:image` | `{SITE_URL}/uploads/post-cover/{filename}` | `post.ImageURL` |
| `og:image:width` | `1200` | Hardcoded |
| `og:image:height` | `630` | Hardcoded |
| `og:image:alt` | `{Post Title}` | `post.Title` (if image exists) |
| `article:published_time` | RFC3339 | `post.PublishedAt` or `post.CreatedAt` |
| `article:tag` | One tag per meta element | `post.TagNames` (comma-separated) |
| `twitter:card` | `summary_large_image` | Hardcoded |
| `twitter:title` | `{Post Title} \| Eriscoo` | `post.Title` |
| `twitter:description` | First 200 chars of body | `post.Body` |
| `twitter:image` | `{SITE_URL}/uploads/post-cover/{filename}` | `post.ImageURL` |
| `<link rel="canonical">` | `{SITE_URL}/{slug}` | `siteURL` + `slug` |
| **JSON-LD** | `BlogPosting` schema | See [JSON-LD](#json-ld-structured-data) |
| **`<body>`** | Full article HTML (`<h1>`, `<h2>`, `<p>`, etc.) | `post.Body` (raw HTML) |

### Static pages (`/og/page/:page`)

| Tag | Value | Data Source |
|-----|-------|-------------|
| `<title>` | `{Title}` | `og-pages-crawler.json` |
| `meta name="description"` | `{Description}` | `og-pages-crawler.json` |
| `og:site_name` | `Eriscoo` | Hardcoded |
| `og:type` | `website` | Hardcoded |
| `og:title` | `{Title}` | `og-pages-crawler.json` |
| `og:description` | `{Description}` | `og-pages-crawler.json` |
| `og:url` | `{SITE_URL}/{page}` | `siteURL` + `page` |
| `og:image` | `{SITE_URL}/{image}` | `og-pages-crawler.json` |
| `og:image:width` | `1200` | Hardcoded |
| `og:image:height` | `630` | Hardcoded |
| `og:image:alt` | `{Title}` | `og-pages-crawler.json` (if image exists) |
| `twitter:card` | `summary_large_image` | Hardcoded |
| `twitter:title` | `{Title}` | `og-pages-crawler.json` |
| `twitter:description` | `{Description}` | `og-pages-crawler.json` |
| `twitter:image` | `{SITE_URL}/{image}` | `og-pages-crawler.json` |
| `<link rel="canonical">` | `{SITE_URL}/{page}` | `siteURL` + `page` |
| **JSON-LD** | `WebPage` schema | See [JSON-LD](#json-ld-structured-data) |
| **`<body>`** | `<h1>{Title}</h1><p>{Description}</p>` | `og-pages-crawler.json` |

## JSON-LD Structured Data

JSON-LD (JavaScript Object Notation for Linked Data) is embedded in the `<head>` of every OG response. It helps search engines understand the page structure and enables **rich results** (breadcrumbs, enhanced article cards) in search results.

### Article — `BlogPosting` schema

```json
{
  "@context": "https://schema.org",
  "@type": "BlogPosting",
  "headline": "Post Title",
  "description": "First 200 chars of body...",
  "image": "https://eriscoo.com/uploads/post-cover/xxx.png",
  "datePublished": "2026-07-21T01:08:29+07:00",
  "author": { "@type": "Person", "name": "Eriscoo" },
  "url": "https://eriscoo.com/post-slug",
  "articleSection": "Linux",
  "keywords": "KDE, Gnome"
}
```

| JSON-LD field | Source | Effect in Google |
|---|---|---|
| `headline` | `post.Title` | Rich result title |
| `description` | `post.Body` (stripped, 200 chars) | Rich result snippet |
| `image` | `post.ImageURL` | Article thumbnail |
| `datePublished` | `post.PublishedAt` or `post.CreatedAt` | Date display in results |
| `author` | Hardcoded `"Eriscoo"` | Author attribution |
| `url` | `siteURL + slug` | Canonical URL |
| `articleSection` | `post.CategoryNames` | **Breadcrumb** (e.g., `eriscoo.com › Linux › Post Title`) |
| `keywords` | `post.TagNames` | Content classification (not displayed) |

### Static Page — `WebPage` schema

```json
{
  "@context": "https://schema.org",
  "@type": "WebPage",
  "name": "About Me",
  "description": "Learn more about my background...",
  "image": "https://eriscoo.com/og-about.jpg",
  "url": "https://eriscoo.com/about"
}
```

## Body Content (SEO)

Each OG response includes actual HTML content in `<body>` for search engine indexing:

| Page type | Body content | Source |
|-----------|-------------|--------|
| **Article** | Full article HTML — `<h1>` heading + raw `post.Body` (preserves `<h2>`, `<h3>`, `<ul>`, `<ol>`, `<blockquote>`, `<code>`, `<strong>`, `<em>`, etc.) | Database |
| **Static page** | `<h1>{Title}</h1><p>{Description}</p>` | `og-pages-crawler.json` |

This gives search engines rich, structural content to analyse for ranking — not an empty `<body>`.

## Supported Features

- **Social media preview**: Facebook, Twitter/X, LinkedIn, WhatsApp, Telegram, Discord, Skype — all platforms that read Open Graph and/or Twitter Card tags.
- **SEO indexing**: Search engines receive full HTML with `<title>`, `<meta name="description">`, `<link rel="canonical">`, structured heading hierarchy, and body content.
- **JSON-LD structured data**: `BlogPosting` schema for articles (with categories and tags), `WebPage` schema for static pages. Enables rich results in Google.
- **Article metadata**: `article:published_time` and `article:tag` help with content classification.
- **Sitemap auto-generation**: Dynamic XML sitemap at `/sitemap.xml` includes all static pages and published posts with `lastmod`, `changefreq`, and `priority`.

## Pages That Use This

- **Post detail pages** (`/post/:slug`) use `GET /og/:slug` — meta tags are dynamically generated from the database.
- **Static pages** (About, Portfolio, Posts, Contact, Terms, Privacy) use `GET /og/page/:page` — meta tags are loaded from a JSON config file.

## SEO: `robots.txt` & Sitemap

### `robots.txt`

A static `robots.txt` file served from the frontend's `public/` directory:

```
User-agent: *
Content-Signal: search=yes,ai-train=no,use=reference
Allow: /

User-agent: Amazonbot
Disallow: /

# (other AI bots disallowed)

Sitemap: https://eriscoo.com/sitemap.xml
```

Key points:
- All search engine crawlers are allowed (`Allow: /`).
- AI training bots (GPTBot, ClaudeBot, Google-Extended, etc.) are blocked.
- The `Sitemap:` directive helps bots discover the sitemap automatically — no manual submission required.

### `GET /sitemap.xml`

A dynamic endpoint that generates an XML sitemap listing all discoverable URLs:

| URL group | Source | `changefreq` | `priority` |
|-----------|--------|-------------|-----------|
| Home `/` | Hardcoded | `weekly` | `1.0` |
| Static pages (about, portfolio, posts, contact) | `og-pages-crawler.json` | `monthly` | `0.7` |
| Legal pages (terms, privacy) | `og-pages-crawler.json` | `monthly` | `0.5` |
| Published posts | Database (`FindAllPublished`) | `monthly` | `0.8` |

Posts include a `<lastmod>` date from `PublishedAt` (or `CreatedAt` as fallback).

**Handler file:** `internal/transport/handler/sitemap/handler.go`

Submit the sitemap URL to search engine webmaster tools (Google Search Console, Bing Webmaster Tools) for faster discovery.

## Static Pages Configuration (`og-pages-crawler.json`)

Static page meta tags are configured in a JSON file (default: `./og-pages-crawler.json`, overridable via `OG_PAGES_CONFIG` env var):

```json
{
  "portfolio": {
    "title": "Portfolio | Eriscoo",
    "description": "My projects, work, and side builds.",
    "image": "/og-portfolio.jpg"
  },
  "about": {
    "title": "About | Eriscoo",
    "description": "Learn more about me and what I do.",
    "image": "/og-about.jpg"
  }
}
```

| Field | Description |
|-------|-------------|
| `title` | Page title (include `\| Eriscoo` suffix if desired). Also used for `og:image:alt` when an image is set. |
| `description` | Meta description (plain text, no HTML) — used for `<meta>`, OG, Twitter, JSON-LD, and body content. |
| `image` | Image path relative to `SITE_URL` (e.g., `/og-about.jpg` resolves to `{SITE_URL}/og-about.jpg`). OG images should be placed in the frontend's `public/` directory. |

To add a new static page, add a new entry to this JSON file and restart the server. No code changes needed.

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
| **JSON-LD `@type`** | `BlogPosting` | `WebPage` |
| **Config source** | Database (posts table) | `og-pages-crawler.json` file |
| **Body content** | Full article HTML (h1-h3, lists, code, etc.) | h1 + description paragraph |

### Description Processing (Strip & Truncate)

The description is extracted from `post.Body` (HTML content), then:
1. Strip all HTML tags using regex
2. Trim whitespace and normalize spaces
3. Truncate to a maximum of **200 characters** (runes, not bytes — safe for non-ASCII characters)
4. Append `...` if truncated
5. HTML-escape the final result (for meta tags; JSON-LD uses the raw unescaped text)

## Related Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /og/:slug` | Generate HTML meta tags + JSON-LD + body content for blog post articles |
| `GET /og/page/:page` | Generate HTML meta tags + JSON-LD + body content for static pages |
| `GET /sitemap.xml` | Generate XML sitemap for all published posts and static pages |
| `GET /api/v1/public/posts/:slug` | Get post details as JSON (for the frontend) |
| `GET /api/v1/public/posts/all` | Get all published posts (for listing pages) |
