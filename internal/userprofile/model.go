// File: internal/userprofile/model.go

package userprofile

import (
	"time"

	"github.com/lib/pq" // Import pq package for PostgreSQL arrays
)

type Certificate struct {
	ID            string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserProfileID string `json:"userProfileId,omitempty"`
	Name          string `json:"name" binding:"required"`
	File          string `json:"file" binding:"required"`
	IssueDate     string `json:"issueDate" binding:"required"`
}

type UserProfile struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID       string         `json:"userId" binding:"required"`
	Name         string         `json:"name" binding:"required"`
	Email        string         `json:"email" binding:"required,email"`
	Location     string         `json:"location"`
	ProfilePhoto string         `json:"profilePhoto"`
	Resume       string         `json:"resume"`
	Certificates []Certificate  `gorm:"foreignKey:UserProfileID" json:"certificates"`
	Skills       pq.StringArray `gorm:"type:text[]" json:"skills"` // Change to pq.StringArray
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

// UpdateUserProfileRequest - For profile updates (all fields optional except validation)
type UpdateUserProfileRequest struct {
	Name         string         `json:"name,omitempty"`
	Email        string         `json:"email,omitempty"`
	Location     string         `json:"location,omitempty"`
	ProfilePhoto string         `json:"profilePhoto,omitempty"`
	Resume       string         `json:"resume,omitempty"`
	Skills       pq.StringArray `json:"skills,omitempty"`
	Certificates []Certificate  `json:"certificates,omitempty"`
}

// UpdateCertificateRequest - For certificate updates
type UpdateCertificateRequest struct {
	Name      string `json:"name" binding:"required"`
	File      string `json:"file" binding:"required"`
	IssueDate string `json:"issueDate" binding:"required"`
}
