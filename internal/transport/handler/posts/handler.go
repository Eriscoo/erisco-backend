package posts

import (
	"errors"
	"net/http"
	"strconv"

	postssvc "github.com/eriscoo/blog-backend/internal/application/posts"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type PostsHandler struct {
	svc *postssvc.Service
}

func New(svc *postssvc.Service) *PostsHandler {
	return &PostsHandler{svc: svc}
}

type postResponse struct {
	ID          int     `json:"id" example:"1"`
	Title       string  `json:"title" example:"Hello World"`
	Slug        string  `json:"slug" example:"hello-world"`
	Body        string  `json:"body" example:"Post content here"`
	ImageURL    string  `json:"image_url" example:"/uploads/image.jpg"`
	Categories  string  `json:"categories" example:"1,3"`
	Tags        string  `json:"tags" example:"2,5"`
	CreatedBy   int     `json:"created_by" example:"1"`
	Status      string  `json:"status" example:"draft"`
	PublishedAt *string `json:"published_at" example:"2026-07-12T00:00:00Z"`
	CreatedAt   string  `json:"created_at" example:"2026-07-12T00:00:00Z"`
	UpdatedAt   string  `json:"updated_at" example:"2026-07-12T00:00:00Z"`
}

type createPostReq struct {
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug"`
	Body       string `json:"body"`
	ImageURL   string `json:"image_url"`
	Categories string `json:"categories"`
	Tags       string `json:"tags"`
	Status     string `json:"status"`
}

type updatePostReq struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Body       string `json:"body"`
	ImageURL   string `json:"image_url"`
	Categories string `json:"categories"`
	Tags       string `json:"tags"`
	Status     string `json:"status"`
}

var _ = postResponse{}

// GetPostBySlug godoc
// @Summary      Get a published post by slug
// @Description  Retrieve a published post by its slug (public)
// @Tags         posts
// @Produce      json
// @Param        slug path string true "Post slug"
// @Success      200  {object}  postResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /public/posts/{slug} [get]
func (h *PostsHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	post, err := h.svc.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetAllPublishedPosts godoc
// @Summary      Get all published posts
// @Description  Retrieve a list of all published posts (public)
// @Tags         posts
// @Produce      json
// @Success      200  {array}  postResponse
// @Failure      500  {object}  map[string]string
// @Router       /public/posts/all [get]
func (h *PostsHandler) GetAllPublishedPosts(c *gin.Context) {
	posts, err := h.svc.GetAllPublished()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetPost godoc
// @Summary      Get a post by ID
// @Description  Retrieve a single post by its ID
// @Tags         posts
// @Produce      json
// @Param        id path int true "Post ID"
// @Success      200  {object}  postResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /posts/{id} [get]
func (h *PostsHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	post, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetPosts godoc
// @Summary      Get all posts
// @Description  Retrieve a list of all posts
// @Tags         posts
// @Produce      json
// @Success      200  {array}  postResponse
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /posts [get]
func (h *PostsHandler) GetPosts(c *gin.Context) {
	posts, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// CreatePost godoc
// @Summary      Create a post
// @Description  Create a new post (requires auth)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        body body createPostReq true "Post data"
// @Success      201  {object}  postResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /posts [post]
func (h *PostsHandler) CreatePost(c *gin.Context) {
	var req createPostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	post := domain.Post{
		Title:      req.Title,
		Slug:       req.Slug,
		Body:       req.Body,
		ImageURL:   req.ImageURL,
		Categories: req.Categories,
		Tags:       req.Tags,
		CreatedBy:  userID.(int),
		Status:     req.Status,
	}

	if err := h.svc.Create(&post); err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// UpdatePost godoc
// @Summary      Update a post
// @Description  Update an existing post by ID (requires auth)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id   path int true "Post ID"
// @Param        body body updatePostReq true "Post data"
// @Success      200  {object}  postResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /posts/{id} [put]
func (h *PostsHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.Body != "" {
		post.Body = req.Body
	}
	if req.ImageURL != "" {
		post.ImageURL = req.ImageURL
	}
	if req.Categories != "" {
		post.Categories = req.Categories
	}
	if req.Tags != "" {
		post.Tags = req.Tags
	}
	if req.Status != "" {
		post.Status = req.Status
	}

	if err := h.svc.Update(post); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		if errors.Is(err, domain.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost godoc
// @Summary      Delete a post
// @Description  Delete a post by ID (requires auth)
// @Tags         posts
// @Produce      json
// @Param        id path int true "Post ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /posts/{id} [delete]
func (h *PostsHandler) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
