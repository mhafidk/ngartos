package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Exercise struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Slug      string    `gorm:"uniqueIndex"`
	TopicID   uuid.UUID `json:"topic_id"`
	Bookmarks []Bookmark
}

type Exercises struct {
	Exercises []Topic `json:"exercises"`
}

func (exercise *Exercise) BeforeCreate(tx *gorm.DB) (err error) {
	exercise.ID = uuid.New()
	return
}
