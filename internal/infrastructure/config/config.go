package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	UploadDir      string
	AppEnv         string
	AllowOrigins   []string
	SiteURL        string
	OGPagesConfig  string
}

func Load() *Config {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	uploadDir := os.Getenv("UPLOAD_DIR")
	appEnv := os.Getenv("APP_ENV")

	allowOriginsRaw := os.Getenv("ALLOW_ORIGINS")
	var allowOrigins []string
	if allowOriginsRaw != "" {
		allowOrigins = strings.Split(allowOriginsRaw, ",")
		for i := range allowOrigins {
			allowOrigins[i] = strings.TrimSpace(allowOrigins[i])
		}
	}

	siteURL := strings.TrimRight(os.Getenv("SITE_URL"), "/")
	ogPagesConfig := os.Getenv("OG_PAGES_CONFIG")

	return &Config{
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		Port:          port,
		UploadDir:     uploadDir,
		AppEnv:        appEnv,
		AllowOrigins:  allowOrigins,
		SiteURL:       siteURL,
		OGPagesConfig: ogPagesConfig,
	}
}
