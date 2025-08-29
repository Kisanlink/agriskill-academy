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

	// Email configuration
	MailFrom string
	MailHost string
	MailPort string
	MailPass string

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
	rateLimitWindowStr := getEnvWithDefault("RATE_LIMIT_WINDOW", "1m")
	rateLimitWindow, err := time.ParseDuration(rateLimitWindowStr)
	if err != nil {
		log.Printf("Warning: Invalid RATE_LIMIT_WINDOW, using default: %v", err)
		rateLimitWindow = time.Minute
	}

	// Parse health check timeout
	healthCheckTimeoutStr := getEnvWithDefault("HEALTH_CHECK_TIMEOUT", "30s")
	healthCheckTimeout, err := time.ParseDuration(healthCheckTimeoutStr)
	if err != nil {
		log.Printf("Warning: Invalid HEALTH_CHECK_TIMEOUT, using default: %v", err)
		healthCheckTimeout = 30 * time.Second
	}

	// Parse allowed origins
	allowedOriginsStr := getEnvWithDefault("CORS_ALLOWED_ORIGINS", "*")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	// Parse numeric values
	rateLimitRequests, _ := strconv.Atoi(getEnvWithDefault("RATE_LIMIT_REQUESTS", "100"))
	maxRequestSize, _ := strconv.ParseInt(getEnvWithDefault("MAX_REQUEST_SIZE", "10485760"), 10, 64) // 10MB default
	maxFileSize, _ := strconv.ParseInt(getEnvWithDefault("MAX_FILE_SIZE", "5242880"), 10, 64)        // 5MB default
	redisDB, _ := strconv.Atoi(getEnvWithDefault("REDIS_DB", "0"))

	// Parse boolean values
	enableCORS := getEnvWithDefault("ENABLE_CORS", "true") == "true"
	awsS3ForcePathStyle := getEnvWithDefault("AWS_S3_FORCE_PATH_STYLE", "false") == "true"
	awsS3DisableSSL := getEnvWithDefault("AWS_S3_DISABLE_SSL", "false") == "true"
	logDevelopment := getEnvWithDefault("LOG_DEVELOPMENT", "false") == "true"

	return &Config{
		// Database configuration
		DBHost:     getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:     getEnvWithDefault("DB_PORT", "5432"),
		DBUser:     getEnvWithDefault("POSTGRESS_USER", "postgres"),
		DBPassword: getEnvWithDefault("POSTGRESS_PASS", ""),
		DBName:     getEnvWithDefault("DB_NAME", "agrijobs"),
		DBSSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),

		// Authentication
		JWTSecret: getEnvWithDefault("JWT_SECRET", "your-secret-key"),

		// Server configuration
		ServerPort: getEnvWithDefault("SERVER_PORT", "8080"),

		// Email configuration
		MailFrom: getEnvWithDefault("MAIL_FROM", ""),
		MailHost: getEnvWithDefault("MAIL_HOST", ""),
		MailPort: getEnvWithDefault("MAIL_PORT", "587"),
		MailPass: getEnvWithDefault("MAIL_PASS", ""),

		// AWS S3 configuration
		AWSRegion:           getEnvWithDefault("AWS_REGION", ""),
		AWSS3Bucket:         getEnvWithDefault("AWS_S3_BUCKET", ""),
		AWSAccessKeyID:      getEnvWithDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey:        getEnvWithDefault("AWS_SECRET_ACCESS_KEY", ""),
		AWSS3Endpoint:       getEnvWithDefault("AWS_S3_ENDPOINT", ""),
		AWSS3ForcePathStyle: awsS3ForcePathStyle,
		AWSS3DisableSSL:     awsS3DisableSSL,

		// File upload configuration
		MaxFileSize: maxFileSize,
		UploadDir:   getEnvWithDefault("UPLOAD_DIR", "./uploads"),

		// Application configuration
		AppName:    getEnvWithDefault("APP_NAME", "AgriJobs"),
		AppVersion: getEnvWithDefault("APP_VERSION", "1.0.0"),
		AppEnv:     getEnvWithDefault("APP_ENV", "development"),
		ASABaseURL: getEnvWithDefault("ASA_BASE_URL", ""),

		// Security configuration
		RateLimitRequests: rateLimitRequests,
		RateLimitWindow:   rateLimitWindow,
		AllowedOrigins:    allowedOrigins,
		MaxRequestSize:    maxRequestSize,
		EnableCORS:        enableCORS,

		// Logging configuration
		LogLevel:       getEnvWithDefault("LOG_LEVEL", "info"),
		LogOutputPath:  getEnvWithDefault("LOG_FILE", ""),
		LogFormat:      getEnvWithDefault("LOG_FORMAT", "json"),
		LogDevelopment: logDevelopment,

		// Redis configuration
		RedisAddr:     getEnvWithDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnvWithDefault("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// Health check configuration
		HealthCheckTimeout: healthCheckTimeout,

		// Job queue configuration
		JobMaxRetries: getEnvAsIntWithDefault("JOB_MAX_RETRIES", 3),
	}
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntWithDefault gets an environment variable as int with a default value
func getEnvAsIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDefaultMaxRetries returns the default max retries from environment or config
func GetDefaultMaxRetries() int {
	return getEnvAsIntWithDefault("JOB_MAX_RETRIES", 3)
}

func InitDB() (*gorm.DB, error) {
	LoadEnv()
	fmt.Println("Env variables:")
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("POSTGRESS_USER:", os.Getenv("POSTGRESS_USER"))
	fmt.Println("POSTGRESS_PASS:", os.Getenv("POSTGRESS_PASS"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_SSLMODE:", os.Getenv("DB_SSLMODE"))

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("POSTGRESS_USER")
	password := os.Getenv("POSTGRESS_PASS")
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
