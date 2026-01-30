package repository

import (
	"context"
	"tidys-go/logic/tags/model"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type TagGroupRepo struct {
	*dbx.BaseRepo[model.TagGroup]
}

func NewTagGroupRepo(db *gorm.DB) *TagGroupRepo {
	return &TagGroupRepo{
		BaseRepo: dbx.NewBaseRepo[model.TagGroup](db),
	}
}

func (r *TagGroupRepo) ExistsByID(ctx context.Context, spaceID, groupID uint) (bool, error) {
	var count int64
	if err := r.DB().WithContext(ctx).Model(&model.TagGroup{}).
		Where("id = ? AND space_id = ?", groupID, spaceID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TagGroupRepo) ListBySpaceID(ctx context.Context, spaceID uint) ([]model.TagGroup, error) {
	var groups []model.TagGroup
	if err := r.DB().WithContext(ctx).Where("space_id = ?", spaceID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *TagGroupRepo) Delete(ctx context.Context, spaceID, groupID uint) error {
	return r.DB().WithContext(ctx).Delete(&model.TagGroup{}, "id = ? AND space_id = ?", groupID, spaceID).Error
}
