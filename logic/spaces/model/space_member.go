package model

import (
	"time"

	"gorm.io/gorm"
)

type SpaceMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SpaceID   uint           `gorm:"space_id;not null;index:idx_space_members_space_user,priority:1" json:"spaceId"`
	UserID    uint           `gorm:"user_id;not null;index:idx_space_members_space_user,priority:2" json:"userId"`
	IsOwner   bool           `gorm:"is_owner;default:false" json:"isOwner"`
	JoinedAt  *time.Time     `gorm:"joined_at" json:"joinedAt"`
	CreatedAt time.Time      `gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"deleted_at" json:"-"`
}

func (SpaceMember) TableName() string {
	return "space_members"
}
