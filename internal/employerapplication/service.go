package employerapplication

import (
	"encoding/json"
	"fmt"
	"strings"
)

type EmployerApplicationService interface {
	GetApplicationsForJob(jobID, status string) ([]JobApplicationResponse, error)
	GetApplicationsByStudent(studentID string) ([]JobApplicationResponse, error)
	UpdateStatus(applicationID, status string) error
	GetApplicantProfile(studentID string) (*ApplicantProfile, error)
	SendMessage(msg *Message) error
	GetMessages(applicationID string) ([]Message, error)
	GetMessagesWithSenderInfo(applicationID string) ([]MessageWithSender, error)
	IsUserAuthorizedForApplication(applicationID, userID string) (bool, error)
	GetJobEmployerID(jobID string) (string, error)
}

type employerApplicationService struct {
	repo EmployerApplicationRepository
}

func NewEmployerApplicationService(repo EmployerApplicationRepository) EmployerApplicationService {
	return &employerApplicationService{repo}
}

func (s *employerApplicationService) GetApplicationsForJob(jobID, status string) ([]JobApplicationResponse, error) {
	fmt.Printf("DEBUG: Service GetApplicationsForJob - JobID: %s, Status: '%s'\n", jobID, status)

	apps, err := s.repo.GetApplicationsForJob(jobID, status)
	if err != nil {
		fmt.Printf("DEBUG: Service GetApplicationsForJob error: %v\n", err)
		return nil, err
	}

	// Transform the data to match frontend requirements
	var responses []JobApplicationResponse
	for _, app := range apps {
		// Parse skills from string to array
		var skills []string
		if app.Skills != "" {
			// Try to parse as JSON array, if it fails, treat as comma-separated string
			if err := json.Unmarshal([]byte(app.Skills), &skills); err != nil {
				// If JSON parsing fails, split by comma
				skills = strings.Split(app.Skills, ",")
				// Trim spaces from each skill
				for i, skill := range skills {
					skills[i] = strings.TrimSpace(skill)
				}
			}
		}

		response := JobApplicationResponse{
			ApplicationID: app.ApplicationID,
			JobID:         app.JobID,
			StudentID:     app.StudentID,
			AppliedAt:     app.AppliedAt,
			Status:        app.ApplicationStatus,
			CoverLetter:   app.CoverLetter,
			JobTitle:      app.JobTitle,
			Company:       app.Company,
			JobType:       app.JobType,
			UserID:        app.UserID,
			ID:            app.ApplicationID, // For consistency
			Applicant: ApplicantInfo{
				Name:        app.Name,
				Email:       app.Email,
				Skills:      skills,
				Experience:  app.Experience,
				Education:   app.Education,
				Portfolio:   app.Portfolio,
				LinkedIn:    app.LinkedIn,
				Github:      app.Github,
				ProfileName: app.ProfileName,
				Location:    app.Location,
				Summary:     "",        // Not available in current data
				Phone:       app.Phone, // Use phone number from database
			},
		}
		responses = append(responses, response)
	}

	fmt.Printf("DEBUG: Service GetApplicationsForJob success - Found %d applications\n", len(responses))
	fmt.Printf("DEBUG: Service returning applications: %+v\n", responses)
	return responses, err
}

func (s *employerApplicationService) GetApplicationsByStudent(studentID string) ([]JobApplicationResponse, error) {
	fmt.Printf("DEBUG: Service GetApplicationsByStudent - StudentID: %s\n", studentID)

	apps, err := s.repo.GetApplicationsByStudent(studentID)
	if err != nil {
		fmt.Printf("DEBUG: Service GetApplicationsByStudent error: %v\n", err)
		return nil, err
	}

	// Transform the data to match frontend requirements
	var responses []JobApplicationResponse
	for _, app := range apps {
		// Parse skills from string to array
		var skills []string
		if app.Skills != "" {
			// Try to parse as JSON array, if it fails, treat as comma-separated string
			if err := json.Unmarshal([]byte(app.Skills), &skills); err != nil {
				// If JSON parsing fails, split by comma
				skills = strings.Split(app.Skills, ",")
				// Trim spaces from each skill
				for i, skill := range skills {
					skills[i] = strings.TrimSpace(skill)
				}
			}
		}

		response := JobApplicationResponse{
			ApplicationID: app.ApplicationID,
			JobID:         app.JobID,
			StudentID:     app.StudentID,
			AppliedAt:     app.AppliedAt,
			Status:        app.ApplicationStatus,
			CoverLetter:   app.CoverLetter,
			JobTitle:      app.JobTitle,
			Company:       app.Company,
			JobType:       app.JobType,
			UserID:        app.UserID,
			ID:            app.ApplicationID, // For consistency
			Applicant: ApplicantInfo{
				Name:        app.Name,
				Email:       app.Email,
				Skills:      skills,
				Experience:  app.Experience,
				Education:   app.Education,
				Portfolio:   app.Portfolio,
				LinkedIn:    app.LinkedIn,
				Github:      app.Github,
				ProfileName: app.ProfileName,
				Location:    app.Location,
				Summary:     "",        // Not available in current data
				Phone:       app.Phone, // Use phone number from database
			},
		}
		responses = append(responses, response)
	}

	fmt.Printf("DEBUG: Service GetApplicationsByStudent success - Found %d applications\n", len(responses))
	fmt.Printf("DEBUG: Service returning applications: %+v\n", responses)
	return responses, err
}

func (s *employerApplicationService) UpdateStatus(applicationID, status string) error {
	return s.repo.UpdateStatus(applicationID, status)
}

func (s *employerApplicationService) GetApplicantProfile(studentID string) (*ApplicantProfile, error) {
	return s.repo.GetApplicantProfile(studentID)
}

func (s *employerApplicationService) SendMessage(msg *Message) error {
	return s.repo.AddMessage(msg)
}

func (s *employerApplicationService) GetMessages(applicationID string) ([]Message, error) {
	return s.repo.GetMessages(applicationID)
}

func (s *employerApplicationService) GetMessagesWithSenderInfo(applicationID string) ([]MessageWithSender, error) {
	return s.repo.GetMessagesWithSenderInfo(applicationID)
}

func (s *employerApplicationService) IsUserAuthorizedForApplication(applicationID, userID string) (bool, error) {
	return s.repo.IsUserAuthorizedForApplication(applicationID, userID)
}

func (s *employerApplicationService) GetJobEmployerID(jobID string) (string, error) {
	fmt.Printf("DEBUG: Service GetJobEmployerID - JobID: %s\n", jobID)

	employerID, err := s.repo.GetJobEmployerID(jobID)

	fmt.Printf("DEBUG: Service GetJobEmployerID result - EmployerID: %s, Error: %v\n", employerID, err)
	return employerID, err
}
