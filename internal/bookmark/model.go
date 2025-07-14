package bookmark

import (
	"time"
)

type Bookmark struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid" json:"user_id"`
	JobID     string    `gorm:"type:uuid" json:"job_id"`
	CreatedAt time.Time `json:"created_at"`
}
