package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                  uuid.UUID `gorm:"type:uuid"`
	Username            string    `gorm:"uniqueIndex" json:"username"`
	Email               string    `gorm:"uniqueIndex" json:"email"`
	Password            string    `json:"password"`
	VerificationToken   string
	Verified            bool
	VerifyAt            *time.Time
	Bookmarks           []Bookmark
	ForgotPasswordToken string
}

type Users struct {
	Users []User `json:"users"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}
