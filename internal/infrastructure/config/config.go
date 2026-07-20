package config

import (
	"os"
	"strconv"
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
	OGPagesConfig string

	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPass      string
	SMTPFromName  string
	SMTPFromEmail string
	SMTPAdminEmail string
	SMTPSecure    bool

	TurnstileSecretKey string
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

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpSecure, _ := strconv.ParseBool(os.Getenv("SMTP_SECURE"))

	return &Config{
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		Port:          port,
		UploadDir:     uploadDir,
		AppEnv:        appEnv,
		AllowOrigins:  allowOrigins,
		SiteURL:       siteURL,
		OGPagesConfig: ogPagesConfig,

		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       smtpPort,
		SMTPUser:       os.Getenv("SMTP_USER"),
		SMTPPass:       os.Getenv("SMTP_PASS"),
		SMTPFromName:   os.Getenv("SMTP_FROM_NAME"),
		SMTPFromEmail:  os.Getenv("SMTP_FROM_EMAIL"),
		SMTPAdminEmail: os.Getenv("SMTP_ADMIN_EMAIL"),
		SMTPSecure:     smtpSecure,

		TurnstileSecretKey: os.Getenv("TURNSTILE_SECRET_KEY"),
	}
}
