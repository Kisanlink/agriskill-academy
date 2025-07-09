// File: internal/studentprofile/model.go

package studentprofile

import (
	"time"

	"github.com/lib/pq" // Import pq package for PostgreSQL arrays
)

type Certificate struct {
	ID               string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StudentProfileID string `json:"studentProfileId,omitempty"`
	Name             string `json:"name" binding:"required"`
	File             string `json:"file" binding:"required"`
	IssueDate        string `json:"issueDate" binding:"required"`
}

// TableName specifies the database table name for Certificate
func (Certificate) TableName() string {
	return "certificates"
}

type StudentProfile struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID       string         `json:"userId" binding:"required"`
	Name         string         `json:"name" binding:"required"`
	Email        string         `json:"email" binding:"required,email"`
	Location     string         `json:"location"`
	ProfilePhoto string         `json:"profilePhoto"`
	Resume       string         `json:"resume"`
	Certificates []Certificate  `gorm:"foreignKey:StudentProfileID" json:"certificates"`
	Skills       pq.StringArray `gorm:"type:text[]" json:"skills"` // Change to pq.StringArray
	Experience   float64        `json:"experience"`
	Education    string         `json:"education"`
	Portfolio    string         `json:"portfolio"`
	Linkedin     string         `json:"linkedin"`
	Github       string         `json:"github"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

// TableName specifies the database table name for StudentProfile
func (StudentProfile) TableName() string {
	return "student_profiles"
}

// UpdateStudentProfileRequest - For profile updates (all fields optional except validation)
type UpdateStudentProfileRequest struct {
	Name         string         `json:"name,omitempty"`
	Email        string         `json:"email,omitempty"`
	Location     string         `json:"location,omitempty"`
	ProfilePhoto string         `json:"profilePhoto,omitempty"`
	Resume       string         `json:"resume,omitempty"`
	Skills       pq.StringArray `json:"skills,omitempty"`
	Experience   *float64       `json:"experience,omitempty"`
	Education    string         `json:"education,omitempty"`
	Portfolio    string         `json:"portfolio,omitempty"`
	Linkedin     string         `json:"linkedin,omitempty"`
	Github       string         `json:"github,omitempty"`
	Certificates []Certificate  `json:"certificates,omitempty"`
}

// UpdateCertificateRequest - For certificate updates
type UpdateCertificateRequest struct {
	Name      string `json:"name" binding:"required"`
	File      string `json:"file" binding:"required"`
	IssueDate string `json:"issueDate" binding:"required"`
}
