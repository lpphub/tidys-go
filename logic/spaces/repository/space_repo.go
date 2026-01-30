package repository

import (
	"context"
	"tidys-go/logic/spaces/model"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type SpaceRepo struct {
	*dbx.BaseRepo[model.Space]
}

func NewSpaceRepo(db *gorm.DB) *SpaceRepo {
	return &SpaceRepo{
		BaseRepo: dbx.NewBaseRepo[model.Space](db),
	}
}

func (r *SpaceRepo) CountBySpaceIDs(ctx context.Context, spaceIDs []uint, model interface{}) (map[uint]int, error) {
	var counts []struct {
		SpaceID uint `gorm:"column:space_id"`
		Count   int  `gorm:"column:count"`
	}
	if err := r.DB().WithContext(ctx).
		Model(model).
		Select("space_id, COUNT(*) as count").
		Where("space_id IN ?", spaceIDs).
		Group("space_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	result := make(map[uint]int, len(counts))
	for _, c := range counts {
		result[c.SpaceID] = c.Count
	}
	return result, nil
}
