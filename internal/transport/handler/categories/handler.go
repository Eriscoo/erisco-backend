package categories

import (
	"errors"
	"net/http"
	"strconv"

	catsvc "github.com/eriscoo/blog-backend/internal/application/categories"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type CategoriesHandler struct {
	svc *catsvc.Service
}

func New(svc *catsvc.Service) *CategoriesHandler {
	return &CategoriesHandler{svc: svc}
}

type categoryResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"technology"`
}

type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

type updateCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

// used by swaggo
var _ = categoryResponse{}

// GetCategories godoc
// @Summary      Get all categories
// @Description  Retrieve a list of all categories
// @Tags         categories
// @Produce      json
// @Success      200  {array}  categoryResponse
// @Failure      500  {object}  map[string]string
// @Router       /categories [get]
func (h *CategoriesHandler) GetCategories(c *gin.Context) {
	categories, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// CreateCategory godoc
// @Summary      Create a category
// @Description  Create a new category (requires auth)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        body body createCategoryReq true "Category name"
// @Success      201  {object}  categoryResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /categories [post]
func (h *CategoriesHandler) CreateCategory(c *gin.Context) {
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.svc.Create(req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "category already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Update an existing category by ID (requires auth)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path int true "Category ID"
// @Param        body body updateCategoryReq true "New name"
// @Success      200  {object}  categoryResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /categories/{id} [put]
func (h *CategoriesHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Update(id, req.Name); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		if errors.Is(err, domain.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, gin.H{"error": "category already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "name": req.Name})
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete a category by ID (requires auth)
// @Tags         categories
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /categories/{id} [delete]
func (h *CategoriesHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
