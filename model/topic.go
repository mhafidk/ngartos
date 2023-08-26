package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid"`
	Name string `json:"name"`
	Content string `json:"content"`
	Slug string `gorm:"uniqueIndex"`
	ParentID *uuid.UUID `json:"parent_id"`
	Parent *Topic
}

type Topics struct {
	Topics []Topic `json:"topics"`
}

func (topic *Topic) BeforeCreate(tx *gorm.DB) (err error) {
	topic.ID = uuid.New()
	return
}