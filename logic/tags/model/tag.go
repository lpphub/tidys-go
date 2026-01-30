package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SpaceID   uint           `gorm:"space_id;not null;default:0;index" json:"spaceId"`
	GroupID   uint           `gorm:"group_id;index" json:"groupId"`
	Content   string         `gorm:"content;type:text" json:"content"`
	Color     string         `gorm:"color;not null;size:20" json:"color"`
	OrderNo   int            `gorm:"order_no;default:0" json:"order"`
	CreatedAt time.Time      `gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"deleted_at" json:"-"`
}

func (Tag) TableName() string {
	return "tags"
}
