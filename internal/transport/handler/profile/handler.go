package profile

import (
	"net/http"
	"strconv"

	profilesvc "github.com/eriscoo/blog-backend/internal/application/profile"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	svc *profilesvc.Service
}

func New(svc *profilesvc.Service) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

type profileResponse struct {
	UserID    int    `json:"user_id" example:"1"`
	Bio       string `json:"bio" example:"Full-stack developer"`
	AvatarURL string `json:"avatar_url" example:"/uploads/avatar.jpg"`
	Website   string `json:"website" example:"https://erisco.dev"`
	Location  string `json:"location" example:"Jakarta, Indonesia"`
	Phone     string `json:"phone" example:"+6281234567890"`
}

type updateProfileReq struct {
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	Website   string `json:"website"`
	Location  string `json:"location"`
	Phone     string `json:"phone"`
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get profile by user ID (requires auth)
// @Tags         profile
// @Produce      json
// @Param        user_id path int true "User ID"
// @Success      200  {object}  profileResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /profile/{user_id} [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	profile, err := h.svc.Get(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Create or update user profile (requires auth, only own profile)
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param        user_id path int true "User ID"
// @Param        body body updateProfileReq true "Profile data"
// @Success      200  {object}  domain.UserProfile
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /profile/{user_id} [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	authUserID := c.GetInt("user_id")
	if authUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own profile"})
		return
	}

	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile := &domain.UserProfile{
		UserID:    userID,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		Website:   req.Website,
		Location:  req.Location,
		Phone:     req.Phone,
	}

	if err := h.svc.Update(profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	updated, _ := h.svc.Get(userID)
	c.JSON(http.StatusOK, updated)
}
