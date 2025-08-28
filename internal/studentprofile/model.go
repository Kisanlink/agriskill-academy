// File: internal/studentprofile/model.go

package studentprofile

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"gorm.io/gorm"
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
	base.BaseModel
	StudentProfileID string `gorm:"type:varchar(255)" json:"student_profile_id,omitempty"`
	Name             string `json:"name" binding:"required"`
	FileKey          string `json:"file_key,omitempty"`
	FileName         string `json:"file_name,omitempty"`
	FileType         string `json:"file_type,omitempty"`
	FileSize         int64  `json:"file_size,omitempty"`
	IssueDate        string `json:"issue_date" binding:"required"`
}

// TableName specifies the database table name for Certificate
func (Certificate) TableName() string {
	return "certificates"
}

// NewCertificate creates a new Certificate with proper initialization
func NewCertificate() *Certificate {
	return &Certificate{
		BaseModel: *base.NewBaseModel("CERT", hash.Small),
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (c *Certificate) BeforeCreateGORM(tx *gorm.DB) error {
	return c.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (c *Certificate) BeforeUpdateGORM(tx *gorm.DB) error {
	return c.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (c *Certificate) BeforeDeleteGORM(tx *gorm.DB) error {
	return c.BeforeDelete()
}

type StudentProfile struct {
	base.BaseModel
	UserID string `gorm:"type:varchar(255);not null" json:"user_id" binding:"required"`

	// Required basic information
	Name  string `gorm:"not null" json:"name" binding:"required"`
	Email string `gorm:"not null" json:"email" binding:"required"`

	// Optional contact and location
	Location    string `json:"location,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`

	// S3 file keys for profile media
	ProfilePhotoKey string `json:"profile_photo_key,omitempty"`
	ResumeKey       string `json:"resume_key,omitempty"`

	// Optional professional information
	Certificates []Certificate       `gorm:"foreignKey:StudentProfileID" json:"certificates,omitempty"`
	Skills       PostgreSQLTextArray `gorm:"type:text[]" json:"skills,omitempty"`
	Experience   float64             `json:"experience,omitempty"`
	Education    string              `json:"education,omitempty"`
	Portfolio    string              `json:"portfolio,omitempty"`
	Linkedin     string              `json:"linkedin,omitempty"`
	Github       string              `json:"github,omitempty"`
}

// TableName specifies the database table name for StudentProfile
func (StudentProfile) TableName() string {
	return "student_profiles"
}

// NewStudentProfile creates a new StudentProfile with proper initialization
func NewStudentProfile() *StudentProfile {
	return &StudentProfile{
		BaseModel: *base.NewBaseModel("STUD", hash.Medium),
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (s *StudentProfile) BeforeCreateGORM(tx *gorm.DB) error {
	return s.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (s *StudentProfile) BeforeUpdateGORM(tx *gorm.DB) error {
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
		s.Skills = PostgreSQLTextArray(filteredSkills)
	}

	return s.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (s *StudentProfile) BeforeDeleteGORM(tx *gorm.DB) error {
	return s.BeforeDelete()
}

// UpdateStudentProfileRequest - For profile updates (all fields optional except validation)
type UpdateStudentProfileRequest struct {
	UserID          string              `json:"user_id,omitempty"`
	Name            string              `json:"name,omitempty"`
	Email           string              `json:"email,omitempty"`
	Location        string              `json:"location,omitempty"`
	PhoneNumber     string              `json:"phone_number,omitempty"`
	ProfilePhotoKey string              `json:"profile_photo_key,omitempty"`
	ResumeKey       string              `json:"resume_key,omitempty"`
	Skills          PostgreSQLTextArray `json:"skills,omitempty"`
	Experience      *float64            `json:"experience,omitempty"`
	Education       string              `json:"education,omitempty"`
	Portfolio       string              `json:"portfolio,omitempty"`
	Linkedin        string              `json:"linkedin,omitempty"`
	Github          string              `json:"github,omitempty"`
	Certificates    []Certificate       `json:"certificates,omitempty"`
}

// UpdateCertificateRequest - For certificate updates
// (If you want to keep file upload for certificates, keep FileName/FileType/FileSize, else just S3 key)
type UpdateCertificateRequest struct {
	Name      string `json:"name" binding:"required"`
	FileKey   string `json:"file_key,omitempty"`
	IssueDate string `json:"issue_date" binding:"required"`
}
