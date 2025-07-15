package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB                *gorm.DB
	AAAServiceBaseURL string
	AAASecret         string
	// New gRPC configuration
	AAAGrpcEndpoint string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AAAServiceBaseURL = os.Getenv("AAA_SERVICE_URL")
	AAASecret = os.Getenv("SECRET_KEY")
	// Load gRPC endpoint
	AAAGrpcEndpoint = os.Getenv("AAA_GRPC_ENDPOINT")
	if AAAGrpcEndpoint == "" {
		AAAGrpcEndpoint = "localhost:50052" // Default gRPC endpoint
	}
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
