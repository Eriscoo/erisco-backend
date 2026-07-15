package main

import (
	"log"
	"os"
	"time"

	_ "github.com/eriscoo/blog-backend/docs"
	authsvc "github.com/eriscoo/blog-backend/internal/application/auth"
	catsvc "github.com/eriscoo/blog-backend/internal/application/categories"
	postssvc "github.com/eriscoo/blog-backend/internal/application/posts"
	profilesvc "github.com/eriscoo/blog-backend/internal/application/profile"
	tagssvc "github.com/eriscoo/blog-backend/internal/application/tags"
	infraauth "github.com/eriscoo/blog-backend/internal/infrastructure/auth"
	"github.com/eriscoo/blog-backend/internal/infrastructure/config"
	"github.com/eriscoo/blog-backend/internal/infrastructure/persistence"
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

	authSvc := authsvc.New(userRepo, tokens)
	tagsSvc := tagssvc.New(tagRepo)
	catsSvc := catsvc.New(catRepo)
	postsSvc := postssvc.New(postRepo)
	profileSvc := profilesvc.New(profileRepo)
	uploadH := uploadHandler.New(cfg.UploadDir)
	ogH := oghandler.New(postsSvc, cfg.SiteURL)

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

	router.Setup(r, authSvc, tagsSvc, catsSvc, postsSvc, profileSvc, uploadH, ogH, tokens)

	if cfg.AppEnv != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Printf("server running on :%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
