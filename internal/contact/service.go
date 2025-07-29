package contact

import (
	"fmt"
	"log"
)

type ContactService interface {
	SubmitContactForm(req *ContactSubmissionRequest) (*ContactSubmissionResponse, error)
	GetContactRequests(req *ContactListRequest) (*ContactListResponse, error)
	GetContactByID(id string) (*ContactRequest, error)
	UpdateContactStatus(id string, status string) error
	DeleteContact(id string) error
	GetContactAnalytics() (*ContactAnalytics, error)
}

type contactService struct {
	repo ContactRepository
}

func NewContactService(repo ContactRepository) ContactService {
	return &contactService{repo: repo}
}

func (s *contactService) SubmitContactForm(req *ContactSubmissionRequest) (*ContactSubmissionResponse, error) {
	// Create contact request entity
	contact := &ContactRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Subject:   req.Subject,
		Message:   req.Message,
		Status:    "new",
	}

	// Save to database
	if err := s.repo.Create(contact); err != nil {
		log.Printf("Error creating contact request: %v", err)
		return nil, fmt.Errorf("failed to submit contact form")
	}

	return &ContactSubmissionResponse{
		Success: true,
		Message: "Your message has been submitted successfully. We will get back to you soon.",
		ID:      contact.ID,
	}, nil
}

func (s *contactService) GetContactRequests(req *ContactListRequest) (*ContactListResponse, error) {
	// Set default values if not provided
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	return s.repo.GetList(req)
}

func (s *contactService) GetContactByID(id string) (*ContactRequest, error) {
	return s.repo.GetByID(id)
}

func (s *contactService) UpdateContactStatus(id string, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"new":       true,
		"read":      true,
		"responded": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	return s.repo.UpdateStatus(id, status)
}

func (s *contactService) DeleteContact(id string) error {
	return s.repo.Delete(id)
}

func (s *contactService) GetContactAnalytics() (*ContactAnalytics, error) {
	// Get total count
	total, err := s.repo.GetTotalCount()
	if err != nil {
		return nil, err
	}

	// Get counts by status
	newCount, err := s.repo.GetCountByStatus("new")
	if err != nil {
		return nil, err
	}

	readCount, err := s.repo.GetCountByStatus("read")
	if err != nil {
		return nil, err
	}

	respondedCount, err := s.repo.GetCountByStatus("responded")
	if err != nil {
		return nil, err
	}

	return &ContactAnalytics{
		TotalContacts:     int(total),
		NewContacts:       int(newCount),
		ReadContacts:      int(readCount),
		RespondedContacts: int(respondedCount),
	}, nil
}

// Contact analytics model
type ContactAnalytics struct {
	TotalContacts     int `json:"total_contacts"`
	NewContacts       int `json:"new_contacts"`
	ReadContacts      int `json:"read_contacts"`
	RespondedContacts int `json:"responded_contacts"`
}
