package contact

import (
	"github.com/Kisanlink/agriskill-academy/internal/middleware"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"gorm.io/gorm"
)

// @Description Contact request model for admin management
type ContactRequest struct {
	base.BaseModel
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
	Email     string `json:"email" binding:"required,email" example:"john@example.com"`
	Phone     string `json:"phone" binding:"required" example:"9876543210"`
	Subject   string `json:"subject" binding:"required" example:"General Inquiry"`
	Message   string `json:"message" binding:"required" example:"I have a question about your services"`
	Status    string `json:"status" gorm:"default:new" example:"new"` // new, read, responded
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

func InitializeCounterFromDatabase(db *gorm.DB) {
	var contactIDs []string
	if err := db.Model(&ContactRequest{}).Pluck("id", &contactIDs).Error; err == nil {
		hash.InitializeGlobalCountersFromDatabase("CONT", contactIDs, hash.Small)
		middleware.DebugLog("Initialized CONT counter with %d existing IDs", len(contactIDs))
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (c *ContactRequest) BeforeCreateGORM(tx *gorm.DB) error {
	return c.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (c *ContactRequest) BeforeUpdateGORM(tx *gorm.DB) error {
	return c.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (c *ContactRequest) BeforeDeleteGORM(tx *gorm.DB) error {
	return c.BeforeDelete()
}

// @Description Contact form submission request
type ContactSubmissionRequest struct {
	FirstName string `json:"first_name" binding:"required" example:"John"`
	LastName  string `json:"last_name" binding:"required" example:"Doe"`
	Email     string `json:"email" binding:"required,email" example:"john@example.com"`
	Phone     string `json:"phone" binding:"required" example:"9876543210"`
	Subject   string `json:"subject" binding:"required" example:"General Inquiry"`
	Message   string `json:"message" binding:"required" example:"I have a question about your services"`
}

// @Description Contact form submission response
type ContactSubmissionResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Contact form submitted successfully"`
	ID      string `json:"id,omitempty" example:"contact_123"`
}

// @Description Admin list request for contact requests
type ContactListRequest struct {
	Page      int    `form:"page" binding:"min=1" example:"1"`
	Limit     int    `form:"limit" binding:"min=1,max=100" example:"10"`
	Status    string `form:"status" example:"new"`
	Search    string `form:"search" example:"john"`
	SortBy    string `form:"sort_by" example:"created_at"`
	SortOrder string `form:"sort_order" example:"DESC"`
}

// @Description Admin list response for contact requests
type ContactListResponse struct {
	Contacts   []ContactRequest `json:"contacts"`
	Pagination PaginationInfo   `json:"pagination"`
}

// @Description Update status request
type UpdateContactStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=new read responded" example:"read"`
}

// @Description Pagination information
type PaginationInfo struct {
	Page       int `json:"page" example:"1"`
	Limit      int `json:"limit" example:"10"`
	Total      int `json:"total" example:"25"`
	TotalPages int `json:"total_pages" example:"3"`
}

// @Description Contact response wrapper
type ContactResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}
