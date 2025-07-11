// File: internal/studentprofile/model.go

package studentprofile

import (
	"time"

	"github.com/lib/pq" // Import pq package for PostgreSQL arrays
)

type Certificate struct {
	ID               string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StudentProfileID string `json:"student_profile_id,omitempty"`
	Name             string `json:"name" binding:"required"`
	File             string `json:"file" binding:"required"`
	IssueDate        string `json:"issue_date" binding:"required"`
}

// TableName specifies the database table name for Certificate
func (Certificate) TableName() string {
	return "certificates"
}

type StudentProfile struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID       string         `json:"user_id" binding:"required"`
	Name         string         `json:"name" binding:"required"`
	Email        string         `json:"email" binding:"required,email"`
	Location     string         `json:"location"`
	PhoneNumber  string         `json:"phone_number"`
	ProfilePhoto string         `json:"profile_photo"`
	Resume       string         `json:"resume"`
	Certificates []Certificate  `gorm:"foreignKey:StudentProfileID" json:"certificates"`
	Skills       pq.StringArray `gorm:"type:text[]" json:"skills"` // Change to pq.StringArray
	Experience   float64        `json:"experience"`
	Education    string         `json:"education"`
	Portfolio    string         `json:"portfolio"`
	Linkedin     string         `json:"linkedin"`
	Github       string         `json:"github"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
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
	PhoneNumber  string         `json:"phone_number,omitempty"`
	ProfilePhoto string         `json:"profile_photo,omitempty"`
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
	IssueDate string `json:"issue_date" binding:"required"`
}
