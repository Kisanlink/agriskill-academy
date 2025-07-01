package employerapplication

import (
	"time"
)

type JobApplicationWithApplicant struct {
	ApplicationID     string    `json:"application_id"`
	JobID             string    `json:"job_id"`
	StudentID         string    `json:"student_id"`
	AppliedAt         time.Time `json:"applied_at"`
	ApplicationStatus string    `json:"status"`
	CoverLetter       string    `json:"cover_letter"`
	StudentResumeFile string    `json:"resume_file"`

	JobTitle    string `json:"job_title"`
	Company     string `json:"company"`
	JobLocation string `json:"location"`
	JobType     string `json:"job_type"`

	Applicant struct {
		UserID      string   `json:"user_id"`
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		Avatar      string   `json:"avatar"`
		ResumeURL   string   `json:"resume_url"`
		Skills      []string `json:"skills"` // You may need to unmarshal JSON
		Location    string   `json:"location"`
		Experience  string   `json:"experience"`
		Education   string   `json:"education"`
		Portfolio   string   `json:"portfolio"`
		LinkedIn    string   `json:"linkedin"`
		Github      string   `json:"github"`
		ProfileName string   `json:"profile_name"`
	} `json:"applicant"`
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
