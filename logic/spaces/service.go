package spaces

import (
	"context"
	"fmt"
	"tidys-go/logic/dto"
	"tidys-go/pkg/consts"
	"time"

	"tidys-go/logic/spaces/model"
	"tidys-go/logic/spaces/repository"
	tagsModel "tidys-go/logic/tags/model"
	"tidys-go/logic/user"
	"tidys-go/pkg/errs"
	"tidys-go/pkg/slices"

	"github.com/lpphub/goweb/ext/dbx"
	"github.com/lpphub/goweb/pkg/logging"
	"gorm.io/gorm"
)

type Service struct {
	spaceRepo  *repository.SpaceRepo
	memberRepo *repository.MemberRepo
	inviteRepo *repository.InviteRepo
	userSvc    *user.Service
}

func NewService(spaceRepo *repository.SpaceRepo, memberRepo *repository.MemberRepo,
	inviteRepo *repository.InviteRepo, userSvc *user.Service) *Service {
	return &Service{
		spaceRepo:  spaceRepo,
		memberRepo: memberRepo,
		inviteRepo: inviteRepo,
		userSvc:    userSvc,
	}
}

func (s *Service) GetSpaces(ctx context.Context, userID uint) ([]dto.SpaceDetail, error) {
	// 获取用户参与的所有 space IDs（包括自己创建的）
	spaceIDs, err := s.memberRepo.GetSpaceIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(spaceIDs) == 0 {
		return []dto.SpaceDetail{}, nil
	}

	spaces, err := s.spaceRepo.FindByIDs(ctx, spaceIDs)
	if err != nil {
		return nil, err
	}

	tagCountMap, _ := s.spaceRepo.CountBySpaceIDs(ctx, spaceIDs, &tagsModel.Tag{})
	memberCountMap, _ := s.spaceRepo.CountBySpaceIDs(ctx, spaceIDs, &model.SpaceMember{})

	result := make([]dto.SpaceDetail, len(spaces))
	for i, sp := range spaces {
		result[i] = dto.SpaceDetail{
			ID:          sp.ID,
			Name:        sp.Name,
			Description: sp.Description,
			TagCount:    tagCountMap[sp.ID],
			MemberCount: memberCountMap[sp.ID],
			Owner:       sp.UserID,
			CreatedAt:   sp.CreatedAt,
			UpdatedAt:   sp.UpdatedAt,
		}
	}

	return result, nil
}

func (s *Service) CreateSpace(ctx context.Context, userID uint, req dto.SpaceReq) (*model.Space, error) {
	var space *model.Space
	err := s.spaceRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		space = &model.Space{
			UserID:      userID,
			Name:        req.Name,
			Description: req.Description,
			CreatedAt:   now,
		}
		if err := tx.Create(space).Error; err != nil {
			return err
		}

		member := model.SpaceMember{
			SpaceID:  space.ID,
			UserID:   userID,
			IsOwner:  true,
			JoinedAt: &now,
		}
		return tx.Create(&member).Error
	})
	if err != nil {
		return nil, err
	}
	return space, nil
}

func (s *Service) UpdateSpace(ctx context.Context, id, userID uint, req dto.SpaceReq) error {
	space, err := s.spaceRepo.First(ctx, id)
	if err != nil {
		return err
	}

	if !space.IsOwnedBy(userID) {
		logging.Warn(ctx, fmt.Sprintf("spaceId=[%d] update by not owner=[%d]", id, userID))
		//return errs.ErrSpaceNotOwned
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"updated_at":  time.Now(),
	}
	return s.spaceRepo.Update(ctx, id, updates)
}

func (s *Service) DeleteSpace(ctx context.Context, id, userID uint) error {
	space, err := s.spaceRepo.First(ctx, id)
	if err != nil {
		return err
	}

	if space.IsOwnedBy(userID) {
		return s.spaceRepo.Delete(ctx, id)
	}

	// Non-owner: leave the space by deleting their membership
	return s.memberRepo.Delete(ctx, id, userID)
}

func (s *Service) GetMembers(ctx context.Context, spaceID uint) ([]dto.SpaceMemberDetail, error) {
	members, err := s.memberRepo.GetBySpaceID(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []dto.SpaceMemberDetail{}, nil
	}

	userIDs := slices.Map(members, func(m model.SpaceMember) uint { return m.UserID })

	users, err := s.userSvc.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	userMap := slices.IndexBy(users, func(u user.User) uint { return u.ID })

	result := slices.Map(members, func(m model.SpaceMember) dto.SpaceMemberDetail {
		vo := dto.SpaceMemberDetail{
			ID:       m.ID,
			SpaceID:  m.SpaceID,
			UserID:   m.UserID,
			IsOwner:  m.IsOwner,
			JoinedAt: m.JoinedAt,
		}

		if u, exists := userMap[m.UserID]; exists {
			vo.Name = u.Name
			vo.Email = u.Email
			vo.Avatar = u.Avatar
		}
		return vo
	})
	return result, nil
}

func (s *Service) InviteMember(ctx context.Context, spaceID, userID uint, emails []string) error {
	if len(emails) > 100 {
		return errs.ErrMaxInviteCount
	}

	space, err := s.spaceRepo.First(ctx, spaceID)
	if err != nil {
		return err
	}

	if !space.IsOwnedBy(userID) {
		return errs.ErrSpaceNotOwned
	}

	// 批量查询用户信息
	users, err := s.userSvc.GetByEmails(ctx, emails)
	if err != nil {
		return err
	}
	userMap := slices.IndexBy(users, func(u *user.User) string { return u.Email })
	userIDs := slices.Map(users, func(u *user.User) uint { return u.ID })

	// 批量查询待处理邀请
	existingInviteeIDs, _ := s.inviteRepo.GetPendingBySpaceAndUserIDs(ctx, spaceID, userIDs)
	existingInviteIDs := make(map[uint]bool)
	for _, id := range existingInviteeIDs {
		existingInviteIDs[id] = true
	}

	// 批量查询现有成员
	existingMembers, _ := s.memberRepo.GetBySpaceID(ctx, spaceID)
	existingMemberIDs := make(map[uint]bool)
	for _, m := range existingMembers {
		existingMemberIDs[m.UserID] = true
	}

	// 处理每个邮箱 - 在程序中判断
	for _, email := range emails {
		invitee, exists := userMap[email]
		if !exists {
			return errs.ErrUserNotExists
		}

		if existingInviteIDs[invitee.ID] {
			return errs.ErrDuplicateInvite
		}

		if existingMemberIDs[invitee.ID] {
			return errs.ErrDuplicateInvite
		}

		invite := model.Invite{
			SpaceID:   spaceID,
			InviteeID: invitee.ID,
			InviterID: userID,
			Status:    consts.InvitePending,
		}

		if err := s.inviteRepo.Create(ctx, &invite); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) RemoveMember(ctx context.Context, spaceID, userID, memberUserID uint) error {
	space, err := s.spaceRepo.First(ctx, spaceID)
	if err != nil {
		return err
	}

	// Only space owner can remove members
	if !space.IsOwnedBy(userID) {
		return errs.ErrSpaceNotOwned
	}

	// Owner cannot remove themselves
	if memberUserID == userID {
		return errs.ErrInvalidParam
	}

	return s.memberRepo.Delete(ctx, spaceID, memberUserID)
}

func (s *Service) GetPendingInvites(ctx context.Context, userID uint) ([]dto.InviteDetail, error) {
	invites, err := s.inviteRepo.GetPendingByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(invites) == 0 {
		return []dto.InviteDetail{}, nil
	}

	inviterIDs := slices.Map(invites, func(inv model.Invite) uint { return inv.InviterID })
	spaceIDs := slices.Map(invites, func(inv model.Invite) uint { return inv.SpaceID })

	users, _ := s.userSvc.GetByIDs(ctx, inviterIDs)
	userMap := slices.IndexBy(users, func(u user.User) uint { return u.ID })

	spaces, _ := s.spaceRepo.FindByIDs(ctx, spaceIDs)
	spaceMap := slices.IndexBy(spaces, func(sp model.Space) uint { return sp.ID })

	result := make([]dto.InviteDetail, len(invites))
	for i, invite := range invites {
		result[i] = dto.InviteDetail{
			ID:        invite.ID,
			SpaceID:   invite.SpaceID,
			InviterID: invite.InviterID,
			Status:    invite.Status,
			CreatedAt: invite.CreatedAt,
		}

		if inviter, exists := userMap[invite.InviterID]; exists {
			result[i].InviterName = inviter.Name
			result[i].InviterEmail = inviter.Email
			result[i].InviterAvatar = inviter.Avatar
		}

		if sp, exists := spaceMap[invite.SpaceID]; exists {
			result[i].SpaceName = sp.Name
		}
	}

	return result, nil
}

func (s *Service) RespondInvite(ctx context.Context, inviteID, userID uint, action string) error {
	invite, err := s.inviteRepo.First(ctx, inviteID)
	if err != nil {
		return err
	}

	if invite.InviteeID != userID {
		return errs.ErrInvalidParam
	}

	if !invite.IsPending() {
		return errs.ErrInvalidParam
	}

	if action == consts.InviteReject {
		return s.inviteRepo.UpdateStatus(ctx, inviteID, consts.InviteReject)
	}
	if action == consts.InviteAccept {
		return s.inviteRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err = s.inviteRepo.UpdateStatus(dbx.WithTx(ctx, tx), inviteID, consts.InviteAccept); err != nil {
				return err
			}

			acceptedAt := time.Now()
			member := model.SpaceMember{
				SpaceID:  invite.SpaceID,
				UserID:   userID,
				IsOwner:  false,
				JoinedAt: &acceptedAt,
			}
			return tx.Create(&member).Error
		})
	}

	return errs.ErrInvalidParam
}
