package model

import (
	"time"

	"gorm.io/gorm"
)

type TagGroup struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SpaceID   uint           `gorm:"space_id;not null;default:0;index" json:"spaceId"`
	Name      string         `gorm:"name;not null;size:100" json:"name"`
	CreatedAt time.Time      `gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"deleted_at" json:"-"`
}

func (TagGroup) TableName() string {
	return "tag_groups"
}
