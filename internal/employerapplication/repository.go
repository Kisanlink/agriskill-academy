package employerapplication

import (
	"fmt"

	"gorm.io/gorm"
)

type EmployerApplicationRepository interface {
	GetApplicationsForJob(jobID, status string) ([]JobApplicationWithApplicant, error)
	GetApplicationsByStudent(studentID string) ([]JobApplicationWithApplicant, error)
	UpdateStatus(applicationID, status string) error
	GetApplicantProfile(studentID string) (*ApplicantProfile, error)
	AddMessage(msg *Message) error
	GetMessages(applicationID string) ([]Message, error)
	GetMessagesWithSenderInfo(applicationID string) ([]MessageWithSender, error)
	IsUserAuthorizedForApplication(applicationID, userID string) (bool, error)
	GetJobEmployerID(jobID string) (string, error)
}

type employerApplicationRepository struct {
	db *gorm.DB
}

func NewEmployerApplicationRepository(db *gorm.DB) EmployerApplicationRepository {
	return &employerApplicationRepository{db}
}

func (r *employerApplicationRepository) GetApplicationsForJob(jobID, status string) ([]JobApplicationWithApplicant, error) {
	fmt.Printf("DEBUG: Repository GetApplicationsForJob - JobID: %s, Status: '%s'\n", jobID, status)

	var results []JobApplicationWithApplicant

	// Build query based on whether status is provided
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT 
				a.id AS application_id, a.job_id, a.student_id, a.applied_at, a.status AS application_status,
				a.cover_letter, a.resume_file AS student_resume_file,
				a.job_title, a.company, a.location AS job_location, a.job_type AS job_type,
				u.id AS user_id, u.name AS user_name, u.email AS user_email,
				up.profile_photo AS avatar, 
				up.skills::text AS skills, 
				COALESCE(up.location, '') AS user_location, 
				COALESCE(CAST(up.experience AS TEXT), '') AS user_experience, 
				COALESCE(up.education, '') AS education, 
				COALESCE(up.portfolio, '') AS portfolio, 
				COALESCE(up.linkedin, '') AS linkedin, 
				COALESCE(up.github, '') AS github, 
				COALESCE(up.name, u.name) AS profile_name,
				COALESCE(up.phone_number, '') AS phone
			FROM applications a
			JOIN users u ON u.id = a.student_id
			LEFT JOIN student_profiles up ON up.user_id = a.student_id
			WHERE a.job_id = ? AND a.status = ?
		`
		args = []interface{}{jobID, status}
	} else {
		query = `
			SELECT 
				a.id AS application_id, a.job_id, a.student_id, a.applied_at, a.status AS application_status,
				a.cover_letter, a.resume_file AS student_resume_file,
				a.job_title, a.company, a.location AS job_location, a.job_type AS job_type,
				u.id AS user_id, u.name AS user_name, u.email AS user_email,
				up.profile_photo AS avatar, 
				up.skills::text AS skills, 
				COALESCE(up.location, '') AS user_location, 
				COALESCE(CAST(up.experience AS TEXT), '') AS user_experience, 
				COALESCE(up.education, '') AS education, 
				COALESCE(up.portfolio, '') AS portfolio, 
				COALESCE(up.linkedin, '') AS linkedin, 
				COALESCE(up.github, '') AS github, 
				COALESCE(up.name, u.name) AS profile_name,
				COALESCE(up.phone_number, '') AS phone
			FROM applications a
			JOIN users u ON u.id = a.student_id
			LEFT JOIN student_profiles up ON up.user_id = a.student_id
			WHERE a.job_id = ?
		`
		args = []interface{}{jobID}
	}

	fmt.Printf("DEBUG: Repository executing query with args: %+v\n", args)
	fmt.Printf("DEBUG: Repository executing query: %s\n", query)

	// First, let's check if the applications exist at all
	var appCount int64
	err := r.db.Model(&struct{}{}).Table("applications").Where("job_id = ?", jobID).Count(&appCount).Error
	if err != nil {
		fmt.Printf("DEBUG: Error counting applications: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Total applications for job %s: %d\n", jobID, appCount)
	}

	// Use Scan instead of ScanRows for better struct mapping
	err = r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		fmt.Printf("DEBUG: Repository query error: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Repository found %d applications after JOIN\n", len(results))

	// Debug: Let's see what the raw query returns
	var rawResults []map[string]interface{}
	err = r.db.Raw(query, args...).Scan(&rawResults).Error
	if err != nil {
		fmt.Printf("DEBUG: Raw query error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Raw query returned %d rows\n", len(rawResults))
		if len(rawResults) > 0 {
			fmt.Printf("DEBUG: First raw result: %+v\n", rawResults[0])
		}
	}

	return results, err
}

func (r *employerApplicationRepository) UpdateStatus(applicationID, status string) error {
	return r.db.Table("applications").Where("id = ?", applicationID).Update("status", status).Error
}

func (r *employerApplicationRepository) GetApplicantProfile(studentID string) (*ApplicantProfile, error) {
	var profile ApplicantProfile
	err := r.db.Raw(`
		SELECT 
			u.id, u.name, u.email, 
			up.phone_number as phone, up.location, up.skills::text as skills, CAST(up.experience AS TEXT) as experience, up.education,
			up.profile_photo as avatar, up.resume as resume_url, up.portfolio, up.linkedin, up.github, up.name as name
		FROM users u
		JOIN student_profiles up ON up.user_id = u.id
		WHERE u.id = ?
	`, studentID).Scan(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *employerApplicationRepository) AddMessage(msg *Message) error {
	fmt.Printf("DEBUG: Repository AddMessage - Message timestamp: %v\n", msg.SentAt)

	// Create message without the sent_at field to let database set it
	messageToSave := map[string]interface{}{
		"application_id": msg.ApplicationID,
		"sender_id":      msg.SenderID,
		"message":        msg.Message,
		// Explicitly omit sent_at to let database use DEFAULT CURRENT_TIMESTAMP
	}

	err := r.db.Table("messages").Create(messageToSave).Error
	if err != nil {
		fmt.Printf("DEBUG: Repository AddMessage error: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Repository AddMessage success - Message saved with ID: %s\n", msg.ID)
	}
	return err
}

func (r *employerApplicationRepository) GetApplicationsByStudent(studentID string) ([]JobApplicationWithApplicant, error) {
	var results []JobApplicationWithApplicant

	rows, err := r.db.Raw(`
		SELECT 
			a.id, a.job_id, a.student_id, a.applied_at, a.status, a.cover_letter, a.resume_file,
			a.job_title, a.company, a.location, a.job_type, a.experience,
			u.id as id, u.name, u.email,
			up.profile_photo as avatar, up.resume as resume_url, up.skills::text as skills, up.location, 
			up.experience as experience, up.education, up.portfolio, up.linkedin, up.github, up.name as name,
			up.phone_number as phone
		FROM applications a
		JOIN users u ON u.id = a.student_id
		JOIN student_profiles up ON up.user_id = a.student_id
		WHERE a.student_id = ?
	`, studentID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var app JobApplicationWithApplicant
		err = r.db.ScanRows(rows, &app)
		if err == nil {
			results = append(results, app)
		}
	}
	return results, err
}

func (r *employerApplicationRepository) GetMessages(applicationID string) ([]Message, error) {
	var messages []Message
	err := r.db.Where("application_id = ?", applicationID).
		Order("sent_at asc").
		Find(&messages).Error
	return messages, err
}

func (r *employerApplicationRepository) GetMessagesWithSenderInfo(applicationID string) ([]MessageWithSender, error) {
	var messages []MessageWithSender

	// Query messages with sender information
	err := r.db.Raw(`
		SELECT 
			m.id, m.application_id, m.sender_id, m.message, m.sent_at,
			u.name as sender_name,
			CASE 
				WHEN a.student_id = m.sender_id THEN 'student'
				WHEN jp.employer_id = m.sender_id THEN 'employer'
				ELSE 'unknown'
			END as sender_type
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		JOIN applications a ON a.id = m.application_id
		JOIN job_posts jp ON jp.id = a.job_id
		WHERE m.application_id = ?
		ORDER BY m.sent_at ASC
	`, applicationID).Scan(&messages).Error

	return messages, err
}

func (r *employerApplicationRepository) IsUserAuthorizedForApplication(applicationID, userID string) (bool, error) {
	// Check if user is the student who applied or the employer who posted the job
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) 
		FROM applications a
		JOIN job_posts jp ON jp.id = a.job_id
		WHERE a.id = ? AND (a.student_id = ? OR jp.employer_id = ?)
	`, applicationID, userID, userID).Count(&count).Error

	return count > 0, err
}

func (r *employerApplicationRepository) GetJobEmployerID(jobID string) (string, error) {
	fmt.Printf("DEBUG: Repository GetJobEmployerID - JobID: %s\n", jobID)

	var employerID string
	err := r.db.Raw("SELECT employer_id FROM job_posts WHERE id = ?", jobID).Scan(&employerID).Error

	fmt.Printf("DEBUG: Repository GetJobEmployerID result - EmployerID: %s, Error: %v\n", employerID, err)
	return employerID, err
}
