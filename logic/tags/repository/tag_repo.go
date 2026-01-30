package repository

import (
	"context"
	"tidys-go/logic/tags/model"
	"time"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type TagRepo struct {
	*dbx.BaseRepo[model.Tag]
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{
		BaseRepo: dbx.NewBaseRepo[model.Tag](db),
	}
}

func (r *TagRepo) ListBySpaceID(ctx context.Context, spaceID uint) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.DB().WithContext(ctx).
		Where("space_id = ?", spaceID).
		Order("group_id, order_no, id ASC").
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepo) GetNextOrder(ctx context.Context, spaceID uint, groupID uint) (int, error) {
	var maxOrder int
	err := r.DB().WithContext(ctx).Model(&model.Tag{}).
		Where("space_id = ? AND group_id = ?", spaceID, groupID).
		Select("COALESCE(MAX(`order_no`), 0)").
		Scan(&maxOrder).Error
	return maxOrder + 1, err
}

func (r *TagRepo) CountByGroupID(ctx context.Context, spaceID, groupID uint) (int64, error) {
	var count int64
	err := r.DB().WithContext(ctx).Model(&model.Tag{}).
		Where("space_id = ? and group_id = ?", spaceID, groupID).
		Count(&count).Error
	return count, err
}

// ListByGroupIDOrdered 获取指定分组内按 order_no 排序的标签列表
func (r *TagRepo) ListByGroupIDOrdered(ctx context.Context, spaceID, groupID uint) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.DB().WithContext(ctx).
		Where("space_id = ? AND group_id = ?", spaceID, groupID).
		Order("order_no ASC").
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// UpdateOrderOrGroup 更新单个 tag 的 order_no，可选更新 group_id
func (r *TagRepo) UpdateOrderOrGroup(tx *gorm.DB, id, groupID uint, orderNo int) error {
	updates := map[string]interface{}{
		"order_no":   orderNo,
		"updated_at": time.Now(),
	}
	if groupID > 0 {
		updates["group_id"] = groupID
	}
	return tx.Model(&model.Tag{}).Where("id = ?", id).Updates(updates).Error
}

// IncrementOrderRange 批量更新 order_no + delta, to <= 0 表示到末尾
func (r *TagRepo) IncrementOrderRange(tx *gorm.DB, groupID uint, from, to, delta int) error {
	db := tx.Model(&model.Tag{}).Where("group_id = ? AND order_no >= ?", groupID, from)
	if to > 0 {
		db = db.Where("order_no <= ?", to)
	}
	return db.UpdateColumn("order_no", gorm.Expr("order_no + ?", delta)).Error
}
