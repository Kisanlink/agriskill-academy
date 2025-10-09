package contact

import (
	"errors"
	"fmt"
	"math"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"gorm.io/gorm"
)

type ContactRepository interface {
	Create(contact *ContactRequest) error
	GetByID(id string) (*ContactRequest, error)
	GetList(req *ContactListRequest) (*ContactListResponse, error)
	UpdateStatus(id string, status string) error
	Delete(id string) error
	GetTotalCount() (int64, error)
	GetCountByStatus(status string) (int64, error)
}

type contactRepository struct {
	*base.BaseRepository[*ContactRequest]
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{
		BaseRepository: base.NewBaseRepository[*ContactRequest](),
		db:             db,
	}
}

func (r *contactRepository) Create(contact *ContactRequest) error {
	return r.db.Create(contact).Error
}

func (r *contactRepository) GetByID(id string) (*ContactRequest, error) {
	var contact ContactRequest
	err := r.db.Where("id = ?", id).First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contact request not found")
		}
		return nil, err
	}
	return &contact, nil
}

func (r *contactRepository) GetList(req *ContactListRequest) (*ContactListResponse, error) {
	var contacts []ContactRequest
	var total int64

	// Build base query
	query := r.db.Model(&ContactRequest{})

	// Apply filters
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR subject ILIKE ?",
			search, search, search, search)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	sortBy := "created_at"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	sortOrder := "DESC"
	if req.SortOrder != "" {
		sortOrder = req.SortOrder
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&contacts).Error; err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &ContactListResponse{
		Contacts: contacts,
		Pagination: PaginationInfo{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (r *contactRepository) UpdateStatus(id string, status string) error {
	result := r.db.Model(&ContactRequest{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("contact request not found")
	}
	return nil
}

func (r *contactRepository) Delete(id string) error {
	result := r.db.Delete(&ContactRequest{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("contact request not found")
	}
	return nil
}

func (r *contactRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&ContactRequest{}).Count(&count).Error
	return count, err
}

func (r *contactRepository) GetCountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&ContactRequest{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
