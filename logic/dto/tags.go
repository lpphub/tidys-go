package dto

import (
	"time"
)

type TagReq struct {
	Content string `json:"content"`
	GroupID uint   `json:"groupId"`
	Color   string `json:"color"`
	SpaceID uint   `json:"spaceId"`
}

type ReorderTagReq struct {
	FromID    uint `json:"fromId"`
	ToGroupID uint `json:"toGroupId"`
	ToIndex   int  `json:"toIndex"`
}

type GroupReq struct {
	Name    string `json:"name"`
	SpaceID uint   `json:"spaceId"`
}

type GetTagsQuery struct {
	SpaceID uint `form:"spaceId"`
}

type TagGroupDetail struct {
	ID      uint   `json:"id"`
	SpaceID uint   `json:"spaceId"`
	Name    string `json:"name"`
	Tags    []Tag  `json:"tags"`
}

type Tag struct {
	ID        uint       `json:"id"`
	SpaceID   uint       `json:"spaceId"`
	Content   string     `json:"content"`
	GroupID   uint       `json:"groupId"`
	Color     string     `json:"color"`
	OrderNo   int        `json:"order"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
