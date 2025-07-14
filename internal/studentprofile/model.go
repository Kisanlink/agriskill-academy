// File: internal/studentprofile/model.go

package studentprofile

import (
	"encoding/json"
	"time"

	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq" // Import pq package for PostgreSQL arrays
	"gorm.io/gorm"
)

// Skills is a custom type to handle JSON marshaling/unmarshaling for PostgreSQL string arrays
type Skills []string

// MarshalJSON implements json.Marshaler
func (s Skills) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(s))
}

// UnmarshalJSON implements json.Unmarshaler
func (s *Skills) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as an array
	var skills []string
	if err := json.Unmarshal(data, &skills); err == nil {
		*s = Skills(skills)
		return nil
	}

	// If that fails, try to unmarshal as a single string
	var singleSkill string
	if err := json.Unmarshal(data, &singleSkill); err == nil {
		if singleSkill != "" {
			*s = Skills([]string{singleSkill})
		} else {
			*s = nil
		}
		return nil
	}

	return fmt.Errorf("skills must be either a string or an array of strings")
}

// Value implements driver.Valuer for database storage
func (s Skills) Value() (interface{}, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}

	// Filter out empty strings
	var filteredSkills []string
	for _, skill := range s {
		if skill != "" {
			filteredSkills = append(filteredSkills, skill)
		}
	}

	if len(filteredSkills) == 0 {
		return nil, nil
	}

	return pq.StringArray(filteredSkills), nil
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

// BeforeCreate is a GORM hook that generates UUID for ID if it's empty and validates if not empty
func (c *Certificate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	} else {
		if _, err := uuid.Parse(c.ID); err != nil {
			return fmt.Errorf("invalid UUID format for Certificate ID: %w", err)
		}
	}
	return nil
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

// BeforeCreate is a GORM hook that generates UUID for ID if it's empty and validates if not empty
func (s *StudentProfile) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	} else {
		if _, err := uuid.Parse(s.ID); err != nil {
			return fmt.Errorf("invalid UUID format for StudentProfile ID: %w", err)
		}
	}
	return nil
}

// BeforeUpdate is a GORM hook to handle Skills field conversion
func (s *StudentProfile) BeforeUpdate(tx *gorm.DB) error {
	// Ensure Skills is properly formatted as an array
	// This handles cases where the frontend might send a single string
	if len(s.Skills) == 1 && s.Skills[0] == "" {
		// If it's an empty string, set to nil/empty array
		s.Skills = nil
	}

	// Filter out empty strings from skills
	if len(s.Skills) > 0 {
		var filteredSkills []string
		for _, skill := range s.Skills {
			if skill != "" {
				filteredSkills = append(filteredSkills, skill)
			}
		}
		s.Skills = Skills(filteredSkills)
	}

	return nil
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
