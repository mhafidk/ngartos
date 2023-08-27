package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid"`
	UserID uuid.UUID `json:"user_id"`
	TopicID uuid.UUID `json:"topic_id"`
}

type Bookmarks struct {
	Bookmarks []Topic `json:"bookmarks"`
}

func (bookmark *Bookmark) BeforeCreate(tx *gorm.DB) (err error) {
	bookmark.ID = uuid.New()
	return
}