package configs

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUsername string
	DBPassword string

	// App
	AppPort                string
	AppBaseURL             string // เช่น "https://myapp.com" (ใช้ใน Reset Password Link)
	InstanceConnectionName string // สำหรับเชื่อมต่อ Cloud SQL

	// Auth Expiry
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	ResetTokenExpiry   time.Duration
	OTPExpiry          time.Duration
	OTPRateLimit       time.Duration

	// JWT
	JWTSecret string

	// SMTP (Gmail)
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string // Gmail App Password (ไม่ใช่ Password จริง)
	SMTPFrom     string // เช่น "MyApp <noreply@gmail.com>"

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string // เช่น "http://localhost:4001/api/auth/google/callback"

	// Facebook OAuth
	FacebookClientID     string
	FacebookClientSecret string
	FacebookRedirectURL  string

	// Twilio (SMS OTP)
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string // เช่น "+12025551234"
}

func LoadConfig() *Config {
	// โหลด .env file (ถ้าไม่มีไฟล์ .env ไม่เป็นไร จะข้ามไปอ่านจาก environment ปกติ)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it, falling back to environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "clean_arch"),
		DBUsername: getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),

		AppPort:                getEnv("APP_PORT", "4001"),
		AppBaseURL:             getEnv("APP_BASE_URL", "http://localhost:4001"),
		InstanceConnectionName: getEnv("INSTANCE_CONNECTION_NAME", ""),

		AccessTokenExpiry:  getEnvAsDuration("ACCESS_TOKEN_EXPIRY", "15m"),
		RefreshTokenExpiry: getEnvAsDuration("REFRESH_TOKEN_EXPIRY", "720h"), // 30 days
		ResetTokenExpiry:   getEnvAsDuration("RESET_TOKEN_EXPIRY", "15m"),
		OTPExpiry:          getEnvAsDuration("OTP_EXPIRY", "5m"),
		OTPRateLimit:       getEnvAsDuration("OTP_RATE_LIMIT", "60s"),

		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),

		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:4001/api/auth/google/callback"),

		FacebookClientID:     getEnv("FACEBOOK_CLIENT_ID", ""),
		FacebookClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""),
		FacebookRedirectURL:  getEnv("FACEBOOK_REDIRECT_URL", "http://localhost:4001/api/auth/facebook/callback"),

		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsDuration(key, fallback string) time.Duration {
	valStr := getEnv(key, fallback)
	duration, err := time.ParseDuration(valStr)
	if err != nil {
		log.Printf("Invalid duration for %s: %v. Using fallback %s", key, valStr, fallback)
		duration, _ = time.ParseDuration(fallback)
	}
	return duration
}
