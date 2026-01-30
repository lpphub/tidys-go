package tags

import (
	"context"
	"tidys-go/logic/dto"
	"tidys-go/logic/tags/model"
	"tidys-go/logic/tags/repository"
	"tidys-go/pkg/errs"
	"tidys-go/pkg/slices"

	"gorm.io/gorm"
)

type Service struct {
	tagRepo   *repository.TagRepo
	groupRepo *repository.TagGroupRepo
}

func NewService(tagRepo *repository.TagRepo, groupRepo *repository.TagGroupRepo) *Service {
	return &Service{
		tagRepo:   tagRepo,
		groupRepo: groupRepo,
	}
}

func (s *Service) GetTags(ctx context.Context, spaceID uint) ([]dto.TagGroupDetail, error) {
	groups, err := s.groupRepo.ListBySpaceID(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	allTags, err := s.tagRepo.ListBySpaceID(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	tagsByGroup := slices.GroupBy(allTags, func(t model.Tag) uint { return t.GroupID })

	result := make([]dto.TagGroupDetail, 0, len(groups))
	for _, g := range groups {
		dtoTags := slices.Map(tagsByGroup[g.ID], func(t model.Tag) dto.Tag {
			return dto.Tag{
				ID:        t.ID,
				SpaceID:   t.SpaceID,
				Content:   t.Content,
				GroupID:   t.GroupID,
				Color:     t.Color,
				OrderNo:   t.OrderNo,
				CreatedAt: t.CreatedAt,
				UpdatedAt: &t.UpdatedAt,
			}
		})

		result = append(result, dto.TagGroupDetail{
			ID:      g.ID,
			Name:    g.Name,
			SpaceID: g.SpaceID,
			Tags:    dtoTags,
		})
	}

	return result, nil
}

func (s *Service) GetOne(ctx context.Context, id uint) (*model.Tag, error) {
	return s.tagRepo.First(ctx, id)
}

func (s *Service) validateGroupExists(ctx context.Context, spaceID, groupID uint) error {
	exists, err := s.groupRepo.ExistsByID(ctx, spaceID, groupID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.ErrRecordNotFound
	}
	return nil
}

func (s *Service) CreateTag(ctx context.Context, req dto.TagReq) (*model.Tag, error) {
	if err := s.validateGroupExists(ctx, req.SpaceID, req.GroupID); err != nil {
		return nil, err
	}

	order, err := s.tagRepo.GetNextOrder(ctx, req.SpaceID, req.GroupID)
	if err != nil {
		return nil, err
	}

	tag := model.Tag{
		SpaceID: req.SpaceID,
		Content: req.Content,
		GroupID: req.GroupID,
		Color:   req.Color,
		OrderNo: order,
	}

	if err = s.tagRepo.Create(ctx, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *Service) UpdateTag(ctx context.Context, id uint, req dto.TagReq) error {
	_, err := s.tagRepo.First(ctx, id)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"content": req.Content,
		"color":   req.Color,
	}
	return s.tagRepo.Update(ctx, id, updates)
}

func (s *Service) DeleteTag(ctx context.Context, id uint) error {
	return s.tagRepo.Delete(ctx, id)
}

func (s *Service) ReorderTag(ctx context.Context, req dto.ReorderTagReq) error {
	if req.ToIndex < 0 {
		return errs.ErrInvalidParam
	}

	tag, err := s.tagRepo.First(ctx, req.FromID)
	if err != nil {
		return err
	}

	// 获取目标分组的标签列表（按 order_no 排序）
	targetTags, err := s.tagRepo.ListByGroupIDOrdered(ctx, tag.SpaceID, req.ToGroupID)
	if err != nil {
		return err
	}

	// 将数组索引转换为实际 order_no
	targetOrderNo := s.indexToOrderNo(targetTags, req.ToIndex)

	if tag.GroupID == req.ToGroupID {
		return s.reorderWithinGroup(ctx, tag, targetOrderNo)
	}
	return s.reorderAcrossGroups(ctx, tag, req.ToGroupID, targetOrderNo)
}

func (s *Service) indexToOrderNo(tags []model.Tag, index int) int {
	if len(tags) == 0 {
		return 1
	}
	if index >= len(tags) {
		return tags[len(tags)-1].OrderNo + 1
	}
	return tags[index].OrderNo
}

func (s *Service) reorderWithinGroup(ctx context.Context, tag *model.Tag, toOrderNo int) error {
	oldOrderNo := tag.OrderNo
	if oldOrderNo == toOrderNo {
		return nil
	}

	return s.tagRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if oldOrderNo > toOrderNo {
			if err := s.tagRepo.IncrementOrderRange(tx, tag.GroupID, toOrderNo, oldOrderNo-1, 1); err != nil {
				return err
			}
		} else {
			if err := s.tagRepo.IncrementOrderRange(tx, tag.GroupID, oldOrderNo+1, toOrderNo, -1); err != nil {
				return err
			}
		}
		return s.tagRepo.UpdateOrderOrGroup(tx, tag.ID, 0, toOrderNo)
	})
}

func (s *Service) reorderAcrossGroups(ctx context.Context, tag *model.Tag, newGroupID uint, toOrderNo int) error {
	oldOrderNo := tag.OrderNo

	return s.tagRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 原分组 oldOrderNo+1 到末尾 -1
		if err := s.tagRepo.IncrementOrderRange(tx, tag.GroupID, oldOrderNo+1, 0, -1); err != nil {
			return err
		}
		// 2. 新分组 toOrderNo 到末尾 +1
		if err := s.tagRepo.IncrementOrderRange(tx, newGroupID, toOrderNo, 0, 1); err != nil {
			return err
		}
		// 3. 更新 tag group_id 和 order_no
		return s.tagRepo.UpdateOrderOrGroup(tx, tag.ID, newGroupID, toOrderNo)
	})
}

func (s *Service) CreateGroup(ctx context.Context, name string, spaceID uint) (*model.TagGroup, error) {
	group := model.TagGroup{
		SpaceID: spaceID,
		Name:    name,
	}

	if err := s.groupRepo.Create(ctx, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *Service) DeleteGroup(ctx context.Context, groupID, spaceID uint) error {
	count, _ := s.tagRepo.CountByGroupID(ctx, spaceID, groupID)
	if count > 0 {
		return errs.ErrInvalidParam
	}

	return s.groupRepo.Delete(ctx, spaceID, groupID)
}
