package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid"`
	UserID     uuid.UUID `json:"user_id"`
	ExerciseID uuid.UUID `json:"exercise_id"`
}

type Bookmarks struct {
	Bookmarks []Bookmark `json:"bookmarks"`
}

func (bookmark *Bookmark) BeforeCreate(tx *gorm.DB) (err error) {
	bookmark.ID = uuid.New()
	return
}
