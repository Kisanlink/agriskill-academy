package bookmark

import (
	"time"
)

type Bookmark struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `json:"user_id"`
	JobID     string    `json:"job_id"`
	CreatedAt time.Time `json:"created_at"`
}
