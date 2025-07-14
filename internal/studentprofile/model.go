// File: internal/studentprofile/model.go

package studentprofile

import (
	"encoding/json"
	"time"

	"github.com/lib/pq" // Import pq package for PostgreSQL arrays
)

// Skills is a custom type to handle JSON marshaling/unmarshaling for PostgreSQL string arrays
type Skills []string

// MarshalJSON implements json.Marshaler
func (s Skills) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(s))
}

// UnmarshalJSON implements json.Unmarshaler
func (s *Skills) UnmarshalJSON(data []byte) error {
	var skills []string
	if err := json.Unmarshal(data, &skills); err != nil {
		return err
	}
	*s = Skills(skills)
	return nil
}

// Value implements driver.Valuer for database storage
func (s Skills) Value() (interface{}, error) {
	if s == nil {
		return nil, nil
	}
	return pq.StringArray(s), nil
}

// Scan implements sql.Scanner for database retrieval
func (s *Skills) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	var pqArray pq.StringArray
	if err := pqArray.Scan(value); err != nil {
		return err
	}
	*s = Skills(pqArray)
	return nil
}

type Certificate struct {
	ID               string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StudentProfileID string `gorm:"type:uuid" json:"student_profile_id,omitempty"`
	Name             string `json:"name" binding:"required"`
	File             []byte `json:"file" gorm:"type:bytea"`
	FileName         string `json:"file_name"`
	FileType         string `json:"file_type"`
	FileSize         int64  `json:"file_size"`
	IssueDate        string `json:"issue_date" binding:"required"`
}

// TableName specifies the database table name for Certificate
func (Certificate) TableName() string {
	return "certificates"
}

type StudentProfile struct {
	ID     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID string `gorm:"type:uuid;not null" json:"user_id" binding:"required"`

	// Required basic information
	Name  string `gorm:"not null" json:"name" binding:"required"`
	Email string `gorm:"not null" json:"email" binding:"required,email"`

	// Optional contact and location
	Location    string `json:"location,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`

	// Optional profile media
	ProfilePhoto     []byte `gorm:"type:bytea" json:"profile_photo,omitempty"`
	ProfilePhotoName string `json:"profile_photo_name,omitempty"`
	ProfilePhotoType string `json:"profile_photo_type,omitempty"`
	ProfilePhotoSize int64  `json:"profile_photo_size,omitempty"`
	Resume           []byte `gorm:"type:bytea" json:"resume,omitempty"`
	ResumeName       string `json:"resume_name,omitempty"`
	ResumeType       string `json:"resume_type,omitempty"`
	ResumeSize       int64  `json:"resume_size,omitempty"`

	// Optional professional information
	Certificates []Certificate `gorm:"foreignKey:StudentProfileID" json:"certificates,omitempty"`
	Skills       Skills        `gorm:"type:text[]" json:"skills,omitempty"`
	Experience   float64       `json:"experience,omitempty"`
	Education    string        `json:"education,omitempty"`
	Portfolio    string        `json:"portfolio,omitempty"`
	Linkedin     string        `json:"linkedin,omitempty"`
	Github       string        `json:"github,omitempty"`

	// System managed fields
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the database table name for StudentProfile
func (StudentProfile) TableName() string {
	return "student_profiles"
}

// UpdateStudentProfileRequest - For profile updates (all fields optional except validation)
type UpdateStudentProfileRequest struct {
	UserID           string        `json:"user_id,omitempty"`
	Name             string        `json:"name,omitempty"`
	Email            string        `json:"email,omitempty"`
	Location         string        `json:"location,omitempty"`
	PhoneNumber      string        `json:"phone_number,omitempty"`
	ProfilePhoto     []byte        `json:"profile_photo,omitempty"`
	ProfilePhotoName string        `json:"profile_photo_name,omitempty"`
	ProfilePhotoType string        `json:"profile_photo_type,omitempty"`
	ProfilePhotoSize int64         `json:"profile_photo_size,omitempty"`
	Resume           []byte        `json:"resume,omitempty"`
	ResumeName       string        `json:"resume_name,omitempty"`
	ResumeType       string        `json:"resume_type,omitempty"`
	ResumeSize       int64         `json:"resume_size,omitempty"`
	Skills           Skills        `json:"skills,omitempty"`
	Experience       *float64      `json:"experience,omitempty"`
	Education        string        `json:"education,omitempty"`
	Portfolio        string        `json:"portfolio,omitempty"`
	Linkedin         string        `json:"linkedin,omitempty"`
	Github           string        `json:"github,omitempty"`
	Certificates     []Certificate `json:"certificates,omitempty"`
}

// UpdateCertificateRequest - For certificate updates
type UpdateCertificateRequest struct {
	Name      string `json:"name" binding:"required"`
	File      []byte `json:"file" binding:"required"`
	FileName  string `json:"file_name,omitempty"`
	FileType  string `json:"file_type,omitempty"`
	FileSize  int64  `json:"file_size,omitempty"`
	IssueDate string `json:"issue_date" binding:"required"`
}
