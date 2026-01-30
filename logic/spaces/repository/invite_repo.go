package repository

import (
	"context"
	"time"

	"tidys-go/logic/spaces/model"
	"tidys-go/pkg/consts"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type InviteRepo struct {
	*dbx.BaseRepo[model.Invite]
}

func NewInviteRepo(db *gorm.DB) *InviteRepo {
	return &InviteRepo{
		BaseRepo: dbx.NewBaseRepo[model.Invite](db),
	}
}

func (r *InviteRepo) GetPendingByUserID(ctx context.Context, userID uint) ([]model.Invite, error) {
	var invites []model.Invite
	if err := r.DB().WithContext(ctx).
		Where("invitee_id = ? AND status = ?", userID, consts.InvitePending).
		Order("id DESC").
		Find(&invites).Error; err != nil {
		return nil, err
	}
	return invites, nil
}

func (r *InviteRepo) UpdateStatus(ctx context.Context, id uint, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	return dbx.TxAwareDB(ctx, r.DB()).Model(&model.Invite{}).Where("id = ?", id).Updates(updates).Error
}

func (r *InviteRepo) GetPendingBySpaceAndUserIDs(ctx context.Context, spaceID uint, userIDs []uint) ([]uint, error) {
	var inviteeIDs []uint
	if err := r.DB().WithContext(ctx).
		Model(&model.Invite{}).
		Where("space_id = ? AND invitee_id IN ? AND status = ?", spaceID, userIDs, consts.InvitePending).
		Pluck("invitee_id", &inviteeIDs).Error; err != nil {
		return nil, err
	}
	return inviteeIDs, nil
}
