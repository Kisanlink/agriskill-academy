package bookmark

import (
	"time"

	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bookmark struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid" json:"user_id"`
	JobID     string    `gorm:"type:uuid" json:"job_id"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate is a GORM hook that generates UUID for ID if it's empty and validates if not empty
func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	} else {
		if _, err := uuid.Parse(b.ID); err != nil {
			return fmt.Errorf("invalid UUID format for Bookmark ID: %w", err)
		}
	}
	return nil
}
