// File: internal/userprofile/model.go

package userprofile

import (
	"time"

	"github.com/lib/pq" // Import pq package for PostgreSQL arrays
)

type Certificate struct {
	ID            string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserProfileID string `json:"userProfileId"`
	Name          string `json:"name"`
	File          string `json:"file"`
	IssueDate     string `json:"issueDate"`
}

type UserProfile struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID       string         `json:"userId"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Location     string         `json:"location"`
	ProfilePhoto string         `json:"profilePhoto"`
	Resume       string         `json:"resume"`
	Certificates []Certificate  `gorm:"foreignKey:UserProfileID" json:"certificates"`
	Skills       pq.StringArray `gorm:"type:text[]" json:"skills"` // Change to pq.StringArray
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}
