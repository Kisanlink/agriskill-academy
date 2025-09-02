package seeding

import (
	"asa/internal/auth"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

type SeedingService struct {
	db *gorm.DB
}

func NewSeedingService(db *gorm.DB) *SeedingService {
	return &SeedingService{db: db}
}

// RunSeeding runs all seeding operations based on environment variables
func (s *SeedingService) RunSeeding() error {
	log.Println("Starting seeding process...")

	// Check if seeding is enabled
	if os.Getenv("ASA_RUN_SEED") != "true" {
		log.Println("Seeding disabled (ASA_RUN_SEED != true)")
		return nil
	}

	// Run admin seeding
	if err := s.seedAdmin(); err != nil {
		return fmt.Errorf("failed to seed admin: %w", err)
	}

	log.Println("Seeding completed successfully!")
	return nil
}

// seedAdmin creates a default admin account if no admin exists
func (s *SeedingService) seedAdmin() error {
	log.Println("Checking for existing admin accounts...")

	// Check if any admin exists
	var adminCount int64
	if err := s.db.Model(&auth.User{}).Where("role = ?", "asa_admin").Count(&adminCount).Error; err != nil {
		return fmt.Errorf("failed to check admin count: %w", err)
	}

	if adminCount > 0 {
		log.Printf("Admin account already exists (%d admins found), skipping admin seeding", adminCount)
		return nil
	}

	log.Println("No admin account found, creating default admin...")

	// Get default admin credentials from environment variables
	defaultEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
	defaultPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	defaultName := os.Getenv("DEFAULT_ADMIN_NAME")
	defaultUsername := os.Getenv("DEFAULT_ADMIN_USERNAME")

	// Validate required environment variables
	if defaultEmail == "" {
		return fmt.Errorf("DEFAULT_ADMIN_EMAIL environment variable is required for admin seeding")
	}
	if defaultPassword == "" {
		return fmt.Errorf("DEFAULT_ADMIN_PASSWORD environment variable is required for admin seeding")
	}
	if defaultName == "" {
		return fmt.Errorf("DEFAULT_ADMIN_NAME environment variable is required for admin seeding")
	}
	if defaultUsername == "" {
		return fmt.Errorf("DEFAULT_ADMIN_USERNAME environment variable is required for admin seeding")
	}

	// Check for existing users with the same email or username before proceeding
	// This prevents unique constraint violations and provides clear error messages
	var emailCount int64
	if err := s.db.Model(&auth.User{}).Where("email = ?", defaultEmail).Count(&emailCount).Error; err != nil {
		return fmt.Errorf("failed to check existing user by email: %w", err)
	}
	if emailCount > 0 {
		return fmt.Errorf("user with email '%s' already exists; please set a different DEFAULT_ADMIN_EMAIL", defaultEmail)
	}

	var usernameCount int64
	if err := s.db.Model(&auth.User{}).Where("username = ?", defaultUsername).Count(&usernameCount).Error; err != nil {
		return fmt.Errorf("failed to check existing user by username: %w", err)
	}
	if usernameCount > 0 {
		return fmt.Errorf("user with username '%s' already exists; please set a different DEFAULT_ADMIN_USERNAME", defaultUsername)
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	// Create admin user within a transaction to handle potential concurrent seeding
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Double-check within transaction to handle concurrent access
		var adminCount int64
		if err := tx.Model(&auth.User{}).Where("role = ?", "asa_admin").Count(&adminCount).Error; err != nil {
			return fmt.Errorf("failed to recheck admin count in transaction: %w", err)
		}
		if adminCount > 0 {
			log.Printf("Admin account was created by another process, skipping")
			return nil // Not an error, just skip
		}

		// Create admin user
		adminUser := auth.NewUser()
		adminUser.Name = defaultName
		adminUser.Username = defaultUsername
		adminUser.Email = defaultEmail
		adminUser.Password = hashedPassword
		adminUser.Role = "asa_admin"

		if err := tx.Create(adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		log.Printf("✅ Default admin account created successfully!")
		log.Printf("   Username: %s", defaultUsername)
		log.Printf("   Email: %s", defaultEmail)
		log.Printf("   Password: %s", defaultPassword)
		log.Printf("   ID: %s", adminUser.ID)
		log.Printf("   ⚠️  IMPORTANT: Change the default password after first login!")

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
