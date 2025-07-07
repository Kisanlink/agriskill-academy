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
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	JobID       string    `json:"jobId"`
	StudentID   string    `json:"studentId"`
	AppliedAt   time.Time `json:"appliedAt"`
	Status      string    `json:"status"`
	CoverLetter string    `json:"coverLetter"`
	ResumeFile  string    `json:"resumeFile"`
	JobTitle    string    `json:"jobTitle"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	JobType     string    `json:"jobType"`
	Experience  string    `json:"experience"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
