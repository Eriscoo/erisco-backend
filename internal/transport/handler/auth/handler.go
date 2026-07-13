package auth

import (
	"errors"
	"net/http"

	authsvc "github.com/eriscoo/blog-backend/internal/application/auth"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *authsvc.Service
}

func New(svc *authsvc.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type registerReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account (role: user)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body registerReq true "User info"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body loginReq true "Credentials"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get the authenticated user's ID, name, and email
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Router       /me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := h.svc.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "name": user.Name, "email": user.Email})
}
