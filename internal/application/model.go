// File: internal/application/model.go

package application

import (
	"time"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"gorm.io/gorm"
)

// Application status constants
const (
	StatusApplied     = "applied"
	StatusReviewing   = "reviewing"
	StatusShortlisted = "shortlisted"
	StatusInterview   = "interview"
	StatusRejected    = "rejected"
	StatusAccepted    = "accepted"
	StatusWithdrawn   = "withdrawn"
)

// Valid statuses for validation
var ValidStatuses = []string{
	StatusApplied,
	StatusReviewing,
	StatusShortlisted,
	StatusInterview,
	StatusRejected,
	StatusAccepted,
	StatusWithdrawn,
}

type Application struct {
	base.BaseModel
	JobID              string    `gorm:"type:varchar(255)" json:"job_id" binding:"required"`
	StudentID          string    `gorm:"type:varchar(255)" json:"student_id" binding:"required"`
	AppliedAt          time.Time `json:"applied_at"`
	Status             string    `json:"status"`
	CoverLetter        string    `json:"cover_letter"`
	ResumeKey          string    `json:"resume_key,omitempty"`
	ResumeFileName     string    `json:"resume_file_name,omitempty"`
	ResumeFileType     string    `json:"resume_file_type,omitempty"`
	ResumeFileSize     int64     `json:"resume_file_size,omitempty"`
	JobTitle           string    `json:"job_title"`
	Company            string    `json:"company"`
	Location           string    `json:"location"`
	JobType            string    `json:"job_type"`
	Experience         string    `json:"experience"`
	StudentPhoneNumber string    `json:"student_phone_number,omitempty"`
}

// TableName specifies the database table name for Application
func (Application) TableName() string {
	return "applications"
}

// NewApplication creates a new Application with proper initialization
func NewApplication() *Application {
	return &Application{
		BaseModel: *base.NewBaseModel("APPL", hash.Large),
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (a *Application) BeforeCreateGORM(tx *gorm.DB) error {
	return a.BaseModel.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (a *Application) BeforeUpdateGORM(tx *gorm.DB) error {
	return a.BaseModel.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (a *Application) BeforeDeleteGORM(tx *gorm.DB) error {
	return a.BaseModel.BeforeDelete()
}

// Request/Response Models
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes,omitempty"`
}

type ApplicationResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	Application  *Application  `json:"application,omitempty"`
	Applications []Application `json:"applications,omitempty"`
}

// Helper function to validate status
func IsValidStatus(status string) bool {
	for _, validStatus := range ValidStatuses {
		if validStatus == status {
			return true
		}
	}
	return false
}
