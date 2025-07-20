// File: internal/application/model.go

package application

import (
	"time"
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
	ID                 string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	JobID              string    `gorm:"type:uuid" json:"job_id" binding:"required"`
	StudentID          string    `gorm:"type:uuid" json:"student_id" binding:"required"`
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
	UpdatedAt          time.Time `json:"updated_at"`
	StudentPhoneNumber string    `json:"student_phone_number,omitempty"`
}

// TableName specifies the database table name for Application
func (Application) TableName() string {
	return "applications"
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
