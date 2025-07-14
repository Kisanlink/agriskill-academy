package bookmark

import (
	"time"
)

type Bookmark struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	UserID    string    `gorm:"type:varchar(255)" json:"user_id"`
	JobID     string    `gorm:"type:varchar(255)" json:"job_id"`
	CreatedAt time.Time `json:"created_at"`
}
