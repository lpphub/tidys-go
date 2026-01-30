package model

import (
	"tidys-go/pkg/consts"
	"time"

	"gorm.io/gorm"
)

type Invite struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SpaceID   uint           `gorm:"space_id;not null;index" json:"spaceId"`
	InviteeID uint           `gorm:"invitee_id;not null;index" json:"inviteeId"`
	InviterID uint           `gorm:"inviter_id;not null;index" json:"inviterId"`
	Status    string         `gorm:"status;not null;size:20;default:pending;index" json:"status"`
	CreatedAt time.Time      `gorm:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"deleted_at" json:"-"`
}

func (*Invite) TableName() string {
	return "space_invites"
}

func (i *Invite) IsPending() bool {
	return i.Status == consts.InvitePending
}
