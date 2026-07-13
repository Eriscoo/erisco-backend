package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	UploadDir string
}

func New(uploadDir string) *UploadHandler {
	return &UploadHandler{UploadDir: uploadDir}
}

// Upload godoc
// @Summary      Upload file
// @Description  Upload an image, returns public URL
// @Tags         upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Image (jpg, png, webp, gif)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	ext := filepath.Ext(file.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file type %s not allowed", ext)})
		return
	}

	subdir := c.DefaultQuery("dir", "profile")
	if subdir == "" {
		subdir = "profile"
	}

	dir := filepath.Join(h.UploadDir, subdir)
	os.MkdirAll(dir, 0755)

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dest := filepath.Join(dir, filename)

	if err := c.SaveUploadedFile(file, dest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": fmt.Sprintf("/uploads/%s/%s", subdir, filename)})
}
