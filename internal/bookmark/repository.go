// File: internal/bookmark/repository.go

package bookmark

import (
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	Save(bookmark *Bookmark) error
	Remove(userID, jobID string) error
	GetByUser(userID string) ([]Bookmark, error)
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepository{db}
}

func (r *bookmarkRepository) Save(bookmark *Bookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *bookmarkRepository) Remove(userID, jobID string) error {
	return r.db.Where("user_id = ? AND job_id = ?", userID, jobID).Delete(&Bookmark{}).Error
}

func (r *bookmarkRepository) GetByUser(userID string) ([]Bookmark, error) {
	var bookmarks []Bookmark
	err := r.db.Where("user_id = ?", userID).Find(&bookmarks).Error
	return bookmarks, err
}
