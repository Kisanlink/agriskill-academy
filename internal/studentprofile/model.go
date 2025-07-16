// File: internal/studentprofile/model.go

package studentprofile

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// PostgreSQLTextArray is a custom type for PostgreSQL text arrays
type PostgreSQLTextArray []string

// Value implements the driver.Valuer interface
func (a PostgreSQLTextArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}

	// Convert to PostgreSQL array format: {"item1","item2","item3"}
	quoted := make([]string, len(a))
	for i, item := range a {
		// Escape quotes and wrap in quotes
		escaped := strings.ReplaceAll(item, `"`, `\"`)
		quoted[i] = `"` + escaped + `"`
	}
	return "{" + strings.Join(quoted, ",") + "}", nil
}

// Scan implements the sql.Scanner interface
func (a *PostgreSQLTextArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("cannot scan %T into PostgreSQLTextArray", value)
	}

	// Handle empty array
	if str == "{}" {
		*a = PostgreSQLTextArray{}
		return nil
	}

	// Remove outer braces and split by comma
	str = strings.Trim(str, "{}")
	if str == "" {
		*a = PostgreSQLTextArray{}
		return nil
	}

	// Split by comma and unquote each item
	parts := strings.Split(str, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		// Remove quotes if present
		if strings.HasPrefix(part, `"`) && strings.HasSuffix(part, `"`) {
			part = part[1 : len(part)-1]
			// Unescape quotes
			part = strings.ReplaceAll(part, `\"`, `"`)
		}
		result[i] = part
	}

	*a = PostgreSQLTextArray(result)
	return nil
}

// MarshalJSON implements json.Marshaler
func (a PostgreSQLTextArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(a))
}

// UnmarshalJSON implements json.Unmarshaler
func (a *PostgreSQLTextArray) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*a = PostgreSQLTextArray(arr)
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
	Certificates []Certificate       `gorm:"foreignKey:StudentProfileID" json:"certificates,omitempty"`
	Skills       PostgreSQLTextArray `gorm:"type:text[]" json:"skills,omitempty"`
	Experience   float64             `json:"experience,omitempty"`
	Education    string              `json:"education,omitempty"`
	Portfolio    string              `json:"portfolio,omitempty"`
	Linkedin     string              `json:"linkedin,omitempty"`
	Github       string              `json:"github,omitempty"`

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
	UserID           string              `json:"user_id,omitempty"`
	Name             string              `json:"name,omitempty"`
	Email            string              `json:"email,omitempty"`
	Location         string              `json:"location,omitempty"`
	PhoneNumber      string              `json:"phone_number,omitempty"`
	ProfilePhoto     []byte              `json:"profile_photo,omitempty"`
	ProfilePhotoName string              `json:"profile_photo_name,omitempty"`
	ProfilePhotoType string              `json:"profile_photo_type,omitempty"`
	ProfilePhotoSize int64               `json:"profile_photo_size,omitempty"`
	Resume           []byte              `json:"resume,omitempty"`
	ResumeName       string              `json:"resume_name,omitempty"`
	ResumeType       string              `json:"resume_type,omitempty"`
	ResumeSize       int64               `json:"resume_size,omitempty"`
	Skills           PostgreSQLTextArray `json:"skills,omitempty"`
	Experience       *float64            `json:"experience,omitempty"`
	Education        string              `json:"education,omitempty"`
	Portfolio        string              `json:"portfolio,omitempty"`
	Linkedin         string              `json:"linkedin,omitempty"`
	Github           string              `json:"github,omitempty"`
	Certificates     []Certificate       `json:"certificates,omitempty"`
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
