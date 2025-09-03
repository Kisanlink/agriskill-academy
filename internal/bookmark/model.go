package bookmark

import (
	"asa/internal/middleware"

	"github.com/Kisanlink/kisanlink-db/pkg/base"
	"github.com/Kisanlink/kisanlink-db/pkg/core/hash"
	"gorm.io/gorm"
)

type Bookmark struct {
	base.BaseModel
	UserID string `gorm:"type:varchar(255)" json:"user_id"`
	JobID  string `gorm:"type:varchar(255)" json:"job_id"`
}

// TableName specifies the database table name for Bookmark
func (Bookmark) TableName() string {
	return "bookmarks"
}

// NewBookmark creates a new Bookmark with proper initialization
func NewBookmark() *Bookmark {
	return &Bookmark{
		BaseModel: *base.NewBaseModel("BOOK", hash.Small),
	}
}

func InitializeCounterFromDatabase(db *gorm.DB) {
	var bookmarkIDs []string
	if err := db.Model(&Bookmark{}).Pluck("id", &bookmarkIDs).Error; err == nil {
		hash.InitializeGlobalCountersFromDatabase("BOOK", bookmarkIDs, hash.Small)
		middleware.DebugLog("Initialized BOOK counter with %d existing IDs", len(bookmarkIDs))
	}
}

// BeforeCreateGORM is called by GORM before creating a new record
func (b *Bookmark) BeforeCreateGORM(tx *gorm.DB) error {
	return b.BeforeCreate()
}

// BeforeUpdateGORM is called by GORM before updating an existing record
func (b *Bookmark) BeforeUpdateGORM(tx *gorm.DB) error {
	return b.BeforeUpdate()
}

// BeforeDeleteGORM is called by GORM before hard deleting a record
func (b *Bookmark) BeforeDeleteGORM(tx *gorm.DB) error {
	return b.BeforeDelete()
}
