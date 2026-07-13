package tags

import (
	"errors"
	"net/http"
	"strconv"

	tagssvc "github.com/eriscoo/blog-backend/internal/application/tags"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type TagsHandler struct {
	svc *tagssvc.Service
}

func New(svc *tagssvc.Service) *TagsHandler {
	return &TagsHandler{svc: svc}
}

type tagResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"golang"`
}

type createTagReq struct {
	Name string `json:"name" binding:"required"`
}

type updateTagReq struct {
	Name string `json:"name" binding:"required"`
}

// used by swaggo
var _ = tagResponse{}

// GetTags godoc
// @Summary      Get all tags
// @Description  Retrieve a list of all tags
// @Tags         tags
// @Produce      json
// @Success      200  {array}  tagResponse
// @Failure      500  {object}  map[string]string
// @Router       /tags [get]
func (h *TagsHandler) GetTags(c *gin.Context) {
	tags, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tags"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// CreateTag godoc
// @Summary      Create a tag
// @Description  Create a new tag (requires auth)
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        body body createTagReq true "Tag name"
// @Success      201  {object}  tagResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /tags [post]
func (h *TagsHandler) CreateTag(c *gin.Context) {
	var req createTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.svc.Create(req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "tag already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// UpdateTag godoc
// @Summary      Update a tag
// @Description  Update an existing tag by ID (requires auth)
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path int true "Tag ID"
// @Param        body body updateTagReq true "New name"
// @Success      200  {object}  tagResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /tags/{id} [put]
func (h *TagsHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(id, req.Name); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		if errors.Is(err, domain.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "tag already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "name": req.Name})
}

// DeleteTag godoc
// @Summary      Delete a tag
// @Description  Delete a tag by ID (requires auth)
// @Tags         tags
// @Produce      json
// @Param        id path int true "Tag ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /tags/{id} [delete]
func (h *TagsHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag deleted"})
}
