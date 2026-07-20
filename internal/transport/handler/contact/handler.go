package contact

import (
	"errors"
	"net/http"

	contactsvc "github.com/eriscoo/blog-backend/internal/application/contact"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *contactsvc.Service
}

func New(svc *contactsvc.Service) *Handler {
	return &Handler{svc: svc}
}

type contactReq struct {
	Name                  string `json:"name"`
	Email                 string `json:"email" binding:"required,email"`
	Subject               string `json:"subject"`
	Phone                 string `json:"phone"`
	Message               string `json:"message"`
	CfTurnstileResponse   string `json:"cf-turnstile-response" binding:"required"`
}

// Submit godoc
// @Summary      Submit contact form
// @Description  Send a message via the contact form (protected by Turnstile)
// @Tags         contact
// @Accept       json
// @Produce      json
// @Param        body body contactReq true "Contact message"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /contact [post]
func (h *Handler) Submit(c *gin.Context) {
	var req contactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &domain.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Phone:   req.Phone,
		Message: req.Message,
	}

	if err := h.svc.Submit(msg, req.CfTurnstileResponse); err != nil {
		if errors.Is(err, domain.ErrSpamDetected) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "spam detected"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent successfully"})
}
