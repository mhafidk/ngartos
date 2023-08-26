package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid"`
	Name string
	Content string
	ParentID *uuid.UUID
	Parent *Topic
}

type Topics struct {
	Topics []Topic `json:"topics"`
}

func (topic *Topic) BeforeCreate(tx *gorm.DB) (err error) {
	topic.ID = uuid.New()
	return
}