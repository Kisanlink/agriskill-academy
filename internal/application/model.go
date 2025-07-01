// File: internal/application/model.go

package application

import (
	"time"
)

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
}
