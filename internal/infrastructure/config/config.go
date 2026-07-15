package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	JWTSecret     string
	Port          string
	UploadDir     string
	AppEnv        string
	AllowOrigins  []string
	SiteURL       string
}

func Load() *Config {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://eriscoo:Eris1234%21%40%23%24@localhost:5432/eriscoodb?sslmode=disable&timezone=Asia/Jakarta"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-change-in-production"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	allowOriginsRaw := os.Getenv("ALLOW_ORIGINS")
	allowOrigins := []string{"http://localhost:5173"}
	if allowOriginsRaw != "" {
		allowOrigins = strings.Split(allowOriginsRaw, ",")
		for i := range allowOrigins {
			allowOrigins[i] = strings.TrimSpace(allowOrigins[i])
		}
	}

	siteURL := os.Getenv("SITE_URL")
	if siteURL == "" {
		siteURL = "http://localhost:5173"
	}
	siteURL = strings.TrimRight(siteURL, "/")

	return &Config{
		DatabaseURL:  dbURL,
		JWTSecret:    jwtSecret,
		Port:         port,
		UploadDir:    uploadDir,
		AppEnv:       appEnv,
		AllowOrigins: allowOrigins,
		SiteURL:      siteURL,
	}
}
