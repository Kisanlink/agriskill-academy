package employerapplication

import (
	"gorm.io/gorm"
)

type EmployerApplicationRepository interface {
	GetApplicationsForJob(jobID, status string) ([]JobApplicationWithApplicant, error)
	GetApplicationsByStudent(studentID string) ([]JobApplicationWithApplicant, error)
	UpdateStatus(applicationID, status string) error
	GetApplicantProfile(studentID string) (*ApplicantProfile, error)
	AddMessage(msg *Message) error
	GetMessages(applicationID string) ([]Message, error)
}

type employerApplicationRepository struct {
	db *gorm.DB
}

func NewEmployerApplicationRepository(db *gorm.DB) EmployerApplicationRepository {
	return &employerApplicationRepository{db}
}

func (r *employerApplicationRepository) GetApplicationsForJob(jobID, status string) ([]JobApplicationWithApplicant, error) {
	var results []JobApplicationWithApplicant

	rows, err := r.db.Raw(`
		SELECT 
			a.id AS application_id, a.job_id, a.student_id, a.applied_at, a.status AS application_status,
			a.cover_letter, a.resume_file AS student_resume_file,
			a.job_title, a.company, a.location AS job_location, a.job_type AS job_type,

			u.id AS user_id, u.name AS user_name, u.email AS user_email,

			up.profile_photo AS avatar, up.resume AS resume_url, 
			up.skills, up.location AS user_location, 
			up.experience AS user_experience, up.education, 
			up.portfolio, up.linkedin, up.github, up.name AS profile_name

		FROM applications a
		JOIN users u ON u.id = a.student_id
		JOIN user_profiles up ON up.user_id = a.student_id
		WHERE a.job_id = ? AND a.status = ?

	`, jobID, status).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var app JobApplicationWithApplicant
		var skills []string
		err = r.db.ScanRows(rows, &app)
		if err == nil {
			app.Applicant.Skills = skills
			results = append(results, app)
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
			up.phone, up.location, up.skills, up.experience, up.education,
			up.profile_photo as avatar, up.resume as resume_url, up.portfolio, up.linkedin, up.github, up.name as name
		FROM users u
		JOIN user_profiles up ON up.user_id = u.id
		WHERE u.id = ?
	`, studentID).Scan(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *employerApplicationRepository) AddMessage(msg *Message) error {
	return r.db.Create(msg).Error
}

func (r *employerApplicationRepository) GetApplicationsByStudent(studentID string) ([]JobApplicationWithApplicant, error) {
	var results []JobApplicationWithApplicant

	rows, err := r.db.Raw(`
		SELECT 
			a.id, a.job_id, a.student_id, a.applied_at, a.status, a.cover_letter, a.resume_file,
			a.job_title, a.company, a.location, a.job_type, a.experience,
			u.id as id, u.name, u.email,
			up.profile_photo as avatar, up.resume as resume_url, up.skills, up.location, 
			up.experience as experience, up.education, up.portfolio, up.linkedin, up.github, up.name as name
		FROM applications a
		JOIN users u ON u.id = a.student_id
		JOIN user_profiles up ON up.user_id = a.student_id
		WHERE a.student_id = ?
	`, studentID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var app JobApplicationWithApplicant
		var skills []string
		err = r.db.ScanRows(rows, &app)
		if err == nil {
			app.Applicant.Skills = skills
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
