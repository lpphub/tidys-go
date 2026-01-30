package dto

import "time"

type SpaceReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SpaceDetail struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TagCount    int       `json:"tagCount"`
	MemberCount int       `json:"memberCount"`
	Owner       uint      `json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type SpaceMemberDetail struct {
	ID       uint       `json:"id"`
	SpaceID  uint       `json:"spaceId"`
	UserID   uint       `json:"userId"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Avatar   string     `json:"avatar"`
	IsOwner  bool       `json:"isOwner"`
	JoinedAt *time.Time `json:"joinedAt"`
}

type InviteDetail struct {
	ID            uint      `json:"id"`
	SpaceID       uint      `json:"spaceId"`
	SpaceName     string    `json:"spaceName"`
	InviterID     uint      `json:"inviterId"`
	InviterName   string    `json:"inviterName"`
	InviterEmail  string    `json:"inviterEmail"`
	InviterAvatar string    `json:"inviterAvatar"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}

type SpaceInviteMemberReq struct {
	Emails []string `json:"emails" binding:"required,min=1,dive,email"`
}
