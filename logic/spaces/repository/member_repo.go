package repository

import (
	"context"
	"tidys-go/logic/spaces/model"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type MemberRepo struct {
	*dbx.BaseRepo[model.SpaceMember]
}

func NewMemberRepo(db *gorm.DB) *MemberRepo {
	return &MemberRepo{
		BaseRepo: dbx.NewBaseRepo[model.SpaceMember](db),
	}
}

func (r *MemberRepo) GetBySpaceID(ctx context.Context, spaceID uint) ([]model.SpaceMember, error) {
	var members []model.SpaceMember
	if err := r.DB().WithContext(ctx).Where("space_id = ?", spaceID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *MemberRepo) Delete(ctx context.Context, spaceID, userID uint) error {
	return r.DB().WithContext(ctx).Where("space_id = ? AND user_id = ?", spaceID, userID).Delete(&model.SpaceMember{}).Error
}
