// File: internal/bookmark/model.go

package bookmark

import (
	"time"
)

type Bookmark struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `json:"userId"`
	JobID     string    `json:"jobId"`
	CreatedAt time.Time `json:"createdAt"`
}
