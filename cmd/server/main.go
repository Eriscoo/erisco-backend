package main

import (
	"log"
	"os"
	"time"

	_ "github.com/eriscoo/blog-backend/docs"
	authsvc "github.com/eriscoo/blog-backend/internal/application/auth"
	catsvc "github.com/eriscoo/blog-backend/internal/application/categories"
	contactsvc "github.com/eriscoo/blog-backend/internal/application/contact"
	postssvc "github.com/eriscoo/blog-backend/internal/application/posts"
	profilesvc "github.com/eriscoo/blog-backend/internal/application/profile"
	tagssvc "github.com/eriscoo/blog-backend/internal/application/tags"
	infraauth "github.com/eriscoo/blog-backend/internal/infrastructure/auth"
	"github.com/eriscoo/blog-backend/internal/infrastructure/config"
	"github.com/eriscoo/blog-backend/internal/infrastructure/email"
	"github.com/eriscoo/blog-backend/internal/infrastructure/persistence"
	"github.com/eriscoo/blog-backend/internal/infrastructure/turnstile"
	"github.com/eriscoo/blog-backend/internal/transport/router"
	oghandler "github.com/eriscoo/blog-backend/internal/transport/handler/og"
	uploadHandler "github.com/eriscoo/blog-backend/internal/transport/handler/upload"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Erisco Blog API
// @version         1.0
// @description     Blog backend API
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	os.Setenv("TZ", "Asia/Jakarta")

	cfg := config.Load()

	db, err := persistence.OpenDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	tokens := infraauth.NewJWTService(cfg.JWTSecret)
	userRepo := persistence.NewUserRepository(db)
	tagRepo := persistence.NewTagRepository(db)
	catRepo := persistence.NewCategoryRepository(db)
	postRepo := persistence.NewPostRepository(db)
	profileRepo := persistence.NewUserProfileRepository(db)
	contactRepo := persistence.NewContactRepository(db)

	authSvc := authsvc.New(userRepo, tokens)
	tagsSvc := tagssvc.New(tagRepo)
	catsSvc := catsvc.New(catRepo)
	postsSvc := postssvc.New(postRepo)
	profileSvc := profilesvc.New(profileRepo)
	uploadH := uploadHandler.New(cfg.UploadDir)

	emailSvc := email.New(email.Config{
		Host:      cfg.SMTPHost,
		Port:      cfg.SMTPPort,
		User:      cfg.SMTPUser,
		Pass:      cfg.SMTPPass,
		FromName:  cfg.SMTPFromName,
		FromEmail: cfg.SMTPFromEmail,
		Secure:    cfg.SMTPSecure,
	})
	turnstileSvc := turnstile.New(cfg.TurnstileSecretKey)
	contactSvc := contactsvc.New(contactRepo, emailSvc, turnstileSvc, cfg.SMTPAdminEmail)

	pages, err := oghandler.LoadPages(cfg.OGPagesConfig)
	if err != nil {
		log.Fatalf("og pages config: %v", err)
	}
	ogH := oghandler.New(postsSvc, cfg.SiteURL, pages)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", cfg.UploadDir)

	router.Setup(r, authSvc, tagsSvc, catsSvc, postsSvc, profileSvc, contactSvc, uploadH, ogH, tokens)

	if cfg.AppEnv != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Printf("server running on :%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
