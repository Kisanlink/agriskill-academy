package employerapplication

import (
	"time"
)

type JobApplicationWithApplicant struct {
	ApplicationID     string    `json:"application_id" gorm:"column:application_id;type:uuid"`
	JobID             string    `json:"job_id" gorm:"column:job_id;type:uuid"`
	StudentID         string    `json:"student_id" gorm:"column:student_id;type:uuid"`
	AppliedAt         time.Time `json:"applied_at" gorm:"column:applied_at"`
	ApplicationStatus string    `json:"status" gorm:"column:application_status"`
	CoverLetter       string    `json:"cover_letter" gorm:"column:cover_letter"`
	// S3 key for resume
	ResumeKey string `json:"resume_key" gorm:"column:resume_key"`

	JobTitle    string `json:"job_title" gorm:"column:job_title"`
	Company     string `json:"company" gorm:"column:company"`
	JobLocation string `json:"job_location" gorm:"column:job_location"`
	JobType     string `json:"job_type" gorm:"column:job_type"`

	// Applicant fields with proper column mapping
	UserID string `json:"user_id" gorm:"column:user_id"`
	Name   string `json:"name" gorm:"column:user_name"`
	Email  string `json:"email" gorm:"column:user_email"`
	// S3 key for avatar
	AvatarKey   string `json:"avatar_key" gorm:"column:avatar_key"`
	Skills      string `json:"skills" gorm:"column:skills"`
	Location    string `json:"user_location" gorm:"column:user_location"`
	Experience  string `json:"experience" gorm:"column:user_experience"`
	Education   string `json:"education" gorm:"column:education"`
	Portfolio   string `json:"portfolio" gorm:"column:portfolio"`
	LinkedIn    string `json:"linkedin" gorm:"column:linkedin"`
	Github      string `json:"github" gorm:"column:github"`
	ProfileName string `json:"profile_name" gorm:"column:profile_name"`
	Phone       string `json:"phone" gorm:"column:phone"`
}

// New response structure for frontend
// Use ResumeKey for S3 storage
// Remove ProfilePhoto, use ProfilePhotoKey
// Remove Avatar, use AvatarKey
type JobApplicationResponse struct {
	ApplicationID string    `json:"application_id"`
	JobID         string    `json:"job_id"`
	StudentID     string    `json:"student_id"`
	AppliedAt     time.Time `json:"applied_at"`
	Status        string    `json:"status"`
	CoverLetter   string    `json:"cover_letter"`
	ResumeKey     string    `json:"resume_key"`
	JobTitle      string    `json:"job_title"`
	Company       string    `json:"company"`
	JobType       string    `json:"job_type"`
	UserID        string    `json:"user_id"`
	ID            string    `json:"id" gorm:"column:application_id"` // Map to application_id from query

	Applicant ApplicantInfo `json:"applicant"`
}

type ApplicantInfo struct {
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	ProfilePhotoKey string   `json:"profile_photo_key"`
	Skills          []string `json:"skills"`
	Experience      string   `json:"experience"`
	Education       string   `json:"education"`
	Portfolio       string   `json:"portfolio"`
	LinkedIn        string   `json:"linkedin"`
	Github          string   `json:"github"`
	ProfileName     string   `json:"profile_name"`
	Location        string   `json:"location"`
	Summary         string   `json:"summary"`
	Phone           string   `json:"phone"`
}

type ApplicantProfile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Location   string `json:"location"`
	Skills     string `json:"skills"`
	Experience string `json:"experience"`
	Education  string `json:"education"`
	AvatarKey  string `json:"avatar_key"`
	ResumeKey  string `json:"resume_key"`
	Portfolio  string `json:"portfolio"`
	LinkedIn   string `json:"linkedin"`
	Github     string `json:"github"`
	Summary    string `json:"summary"`
}

type Message struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ApplicationID string     `json:"application_id"`
	SenderID      string     `gorm:"type:uuid;column:sender_id" json:"sender_id"`
	Message       string     `json:"message" gorm:"column:message"`
	SentAt        *time.Time `json:"sent_at" gorm:"column:sent_at;autoCreateTime"`
}

type MessageWithSender struct {
	ID            string     `json:"id"`
	ApplicationID string     `json:"application_id" gorm:"column:application_id"`
	SenderID      string     `json:"sender_id" gorm:"column:sender_id"`
	SenderName    string     `json:"sender_name"`
	SenderType    string     `json:"sender_type"` // "student" or "employer"
	Message       string     `json:"message" gorm:"column:message"`
	SentAt        *time.Time `json:"sent_at" gorm:"column:sent_at"`
}
