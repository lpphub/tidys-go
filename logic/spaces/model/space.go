package model

import (
	"time"

	"gorm.io/gorm"
)

type Space struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"user_id;not null;index" json:"userId"`
	Name        string         `gorm:"name;not null;size:100" json:"name"`
	Description string         `gorm:"description;size:500" json:"description"`
	CreatedAt   time.Time      `gorm:"created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"deleted_at" json:"-"`
}

func (*Space) TableName() string {
	return "spaces"
}

func (s *Space) IsOwnedBy(userID uint) bool {
	return s.UserID == userID
}
