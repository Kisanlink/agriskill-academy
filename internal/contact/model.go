package contact

import (
	"time"
)

type ContactRequest struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Phone     string    `json:"phone" binding:"required"`
	Subject   string    `json:"subject" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Status    string    `json:"status" gorm:"default:new"` // new, read, responded
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the database table name for ContactRequest
func (ContactRequest) TableName() string {
	return "contact_requests"
}

// Contact form submission request
type ContactSubmissionRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

// Contact form submission response
type ContactSubmissionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

// Admin list request for contact requests
type ContactListRequest struct {
	Page      int    `form:"page" binding:"min=1"`
	Limit     int    `form:"limit" binding:"min=1,max=100"`
	Status    string `form:"status"`
	Search    string `form:"search"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}

// Admin list response for contact requests
type ContactListResponse struct {
	Contacts   []ContactRequest `json:"contacts"`
	Pagination PaginationInfo   `json:"pagination"`
}

// Update status request
type UpdateContactStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=new read responded"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ContactResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
