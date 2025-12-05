package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB        *gorm.DB
	JWTSecret string
)

// Config holds all application configuration
type Config struct {
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Authentication
	JWTSecret string

	// Server configuration
	ServerPort string
	GinMode    string

	// SMTP Email configuration (for notifications: job alerts, updates, etc.)
	// NOT used for authentication emails (Firebase handles those)
	EmailNotificationEnabled bool
	MailFrom                 string
	MailHost                 string
	MailPort                 int
	MailPass                 string

	// AWS S3 configuration
	AWSRegion           string
	AWSS3Bucket         string
	AWSAccessKeyID      string
	AWSSecretKey        string
	AWSS3Endpoint       string
	AWSS3ForcePathStyle bool
	AWSS3DisableSSL     bool

	// File upload configuration
	MaxFileSize int64
	UploadDir   string

	// Application configuration
	AppName    string
	AppVersion string
	AppEnv     string
	ASABaseURL string

	// Security configuration
	RateLimitRequests int
	RateLimitWindow   time.Duration
	AllowedOrigins    []string
	MaxRequestSize    int64
	EnableCORS        bool

	// Logging configuration
	LogLevel       string
	LogOutputPath  string
	LogFormat      string
	LogDevelopment bool

	// Redis configuration (for job queue)
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Health check configuration
	HealthCheckTimeout time.Duration

	// Job queue configuration
	JobMaxRetries int

	// Firebase configuration (for email sending only)
	FirebaseProjectID       string
	FirebaseCredentialsPath string
	FirebaseCredentialsJSON string
	FirebaseWebAPIKey       string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	JWTSecret = os.Getenv("JWT_SECRET")
}

// LoadConfig loads all configuration from environment variables
func LoadConfig() *Config {
	LoadEnv()

	// Parse rate limit window
	rateLimitWindowStr := os.Getenv("RATE_LIMIT_WINDOW")
	var rateLimitWindow time.Duration
	if rateLimitWindowStr != "" {
		if parsed, err := time.ParseDuration(rateLimitWindowStr); err == nil {
			rateLimitWindow = parsed
		} else {
			log.Printf("Warning: Invalid RATE_LIMIT_WINDOW, using default: %v", err)
			rateLimitWindow = time.Minute
		}
	} else {
		rateLimitWindow = time.Minute
	}

	// Parse health check timeout
	healthCheckTimeoutStr := os.Getenv("HEALTH_CHECK_TIMEOUT")
	var healthCheckTimeout time.Duration
	if healthCheckTimeoutStr != "" {
		if parsed, err := time.ParseDuration(healthCheckTimeoutStr); err == nil {
			healthCheckTimeout = parsed
		} else {
			log.Printf("Warning: Invalid HEALTH_CHECK_TIMEOUT, using default: %v", err)
			healthCheckTimeout = 30 * time.Second
		}
	} else {
		healthCheckTimeout = 30 * time.Second
	}

	// Parse allowed origins
	allowedOriginsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	} else {
		allowedOrigins = []string{"*"}
	}

	// Parse numeric values
	rateLimitRequestsStr := os.Getenv("RATE_LIMIT_REQUESTS")
	rateLimitRequests := 100
	if rateLimitRequestsStr != "" {
		if val, err := strconv.Atoi(rateLimitRequestsStr); err == nil {
			rateLimitRequests = val
		}
	}

	maxRequestSizeStr := os.Getenv("MAX_REQUEST_SIZE")
	maxRequestSize := int64(10485760) // 10MB
	if maxRequestSizeStr != "" {
		if val, err := strconv.ParseInt(maxRequestSizeStr, 10, 64); err == nil {
			maxRequestSize = val
		}
	}

	maxFileSizeStr := os.Getenv("MAX_FILE_SIZE")
	maxFileSize := int64(5242880) // 5MB
	if maxFileSizeStr != "" {
		if val, err := strconv.ParseInt(maxFileSizeStr, 10, 64); err == nil {
			maxFileSize = val
		}
	}

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB := 0
	if redisDBStr != "" {
		if val, err := strconv.Atoi(redisDBStr); err == nil {
			redisDB = val
		}
	}

	// Parse boolean values
	enableCORS := os.Getenv("ENABLE_CORS") == "true"
	awsS3ForcePathStyle := os.Getenv("AWS_S3_FORCE_PATH_STYLE") == "true"
	awsS3DisableSSL := os.Getenv("AWS_S3_DISABLE_SSL") == "true"
	logDevelopment := os.Getenv("LOG_DEVELOPMENT") == "true"
	emailNotificationEnabled := os.Getenv("EMAIL_NOTIFICATION") == "true"

	// Parse SMTP port with default
	mailPort := 587
	if mailPortStr := os.Getenv("MAIL_PORT"); mailPortStr != "" {
		if val, err := strconv.Atoi(mailPortStr); err == nil {
			mailPort = val
		} else {
			log.Printf("Warning: Invalid MAIL_PORT value '%s', using default: 587", mailPortStr)
		}
	}

	return &Config{
		// Database configuration
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASS"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		// Authentication
		JWTSecret: os.Getenv("JWT_SECRET"),

		// Server configuration
		ServerPort: os.Getenv("SERVER_PORT"),
		GinMode:    os.Getenv("GIN_MODE"),

		// Email configuration
		EmailNotificationEnabled: emailNotificationEnabled,
		MailFrom:                 os.Getenv("MAIL_FROM"),
		MailHost:                 os.Getenv("MAIL_HOST"),
		MailPort:                 mailPort,
		MailPass:                 os.Getenv("MAIL_PASS"),

		// AWS S3 configuration
		AWSRegion:           os.Getenv("AWS_REGION"),
		AWSS3Bucket:         os.Getenv("AWS_S3_BUCKET"),
		AWSAccessKeyID:      os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey:        os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSS3Endpoint:       os.Getenv("AWS_S3_ENDPOINT"),
		AWSS3ForcePathStyle: awsS3ForcePathStyle,
		AWSS3DisableSSL:     awsS3DisableSSL,

		// File upload configuration
		MaxFileSize: maxFileSize,
		UploadDir:   os.Getenv("UPLOAD_DIR"),

		// Application configuration
		AppName:    os.Getenv("APP_NAME"),
		AppVersion: os.Getenv("APP_VERSION"),
		AppEnv:     os.Getenv("APP_ENV"),
		ASABaseURL: os.Getenv("ASA_BASE_URL"),

		// Security configuration
		RateLimitRequests: rateLimitRequests,
		RateLimitWindow:   rateLimitWindow,
		AllowedOrigins:    allowedOrigins,
		MaxRequestSize:    maxRequestSize,
		EnableCORS:        enableCORS,

		// Logging configuration
		LogLevel:       os.Getenv("LOG_LEVEL"),
		LogOutputPath:  os.Getenv("LOG_FILE"),
		LogFormat:      os.Getenv("LOG_FORMAT"),
		LogDevelopment: logDevelopment,

		// Redis configuration
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,

		// Health check configuration
		HealthCheckTimeout: healthCheckTimeout,

		// Job queue configuration
		JobMaxRetries: GetDefaultMaxRetries(),

		// Firebase configuration
		FirebaseProjectID:       os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseCredentialsPath: os.Getenv("FIREBASE_CREDENTIALS_PATH"),
		FirebaseCredentialsJSON: os.Getenv("FIREBASE_CREDENTIALS_JSON"),
		FirebaseWebAPIKey:       os.Getenv("FIREBASE_WEB_API_KEY"),
	}
}

// GetDefaultMaxRetries returns the default max retries from environment or config
func GetDefaultMaxRetries() int {
	if value := os.Getenv("JOB_MAX_RETRIES"); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return 3
}

func InitDB() (*gorm.DB, error) {
	LoadEnv()
	fmt.Println("Env variables:")
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("POSTGRES_USER:", os.Getenv("POSTGRES_USER"))
	fmt.Println("POSTGRES_PASS:", os.Getenv("POSTGRES_PASS"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_SSLMODE:", os.Getenv("DB_SSLMODE"))

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASS")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Check required variables
	if host == "" || port == "" || user == "" || password == "" || dbname == "" || sslmode == "" {
		return nil, fmt.Errorf("missing one or more required DB environment variables")
	}

	// Construct the DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	// Configure GORM logger based on GIN_MODE
	var gormLogger logger.Interface
	if os.Getenv("GIN_MODE") == "debug" {
		// In debug mode: show all SQL queries with colors
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		// In production: only log warnings and errors (no SQL queries)
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return db, nil
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}
