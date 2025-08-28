package contact

import (
	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"gorm.io/gorm"
)

type ContactRequest struct {
	base.BaseModel
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Status    string `json:"status" gorm:"default:new"` // new, read, responded
}

// TableName specifies the database table name for ContactRequest
func (ContactRequest) TableName() string {
	return "contact_requests"
}

// NewContactRequest creates a new ContactRequest with proper initialization
func NewContactRequest() *ContactRequest {
	return &ContactRequest{
		BaseModel: *base.NewBaseModel("CONT", hash.Small),
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (c *ContactRequest) BeforeCreateGORM(tx *gorm.DB) error {
	return c.BaseModel.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (c *ContactRequest) BeforeUpdateGORM(tx *gorm.DB) error {
	return c.BaseModel.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (c *ContactRequest) BeforeDeleteGORM(tx *gorm.DB) error {
	return c.BaseModel.BeforeDelete()
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
