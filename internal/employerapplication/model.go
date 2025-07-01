// File: internal/employerapplication/model.go

package employerapplication

import (
	"time"
)

type JobApplicationWithApplicant struct {
	ID          string           `json:"id"`
	JobID       string           `json:"jobId"`
	StudentID   string           `json:"studentId"`
	AppliedAt   time.Time        `json:"appliedAt"`
	Status      string           `json:"status"`
	CoverLetter string           `json:"coverLetter"`
	ResumeFile  string           `json:"resumeFile"`
	JobTitle    string           `json:"jobTitle"`
	Company     string           `json:"company"`
	Location    string           `json:"location"`
	JobType     string           `json:"jobType"`
	Experience  string           `json:"experience"`
	Applicant   ApplicantProfile `json:"applicant"`
}

type ApplicantProfile struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Location   string   `json:"location"`
	Skills     []string `json:"skills"`
	Experience string   `json:"experience"`
	Education  string   `json:"education"`
	Avatar     string   `json:"avatar"`
	ResumeUrl  string   `json:"resumeUrl"`
	Portfolio  string   `json:"portfolio"`
	LinkedIn   string   `json:"linkedIn"`
	Github     string   `json:"github"`
	Summary    string   `json:"summary"`
}

type Message struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ApplicationID string    `json:"applicationId"`
	SenderID      string    `json:"senderId"`
	Message       string    `json:"message"`
	SentAt        time.Time `json:"sentAt"`
}
