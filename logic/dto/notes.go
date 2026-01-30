package dto

import (
	"time"
)

// Request DTOs

type NoteReq struct {
	ID      uint   `json:"id,omitempty"`
	Content string `json:"content,omitempty" binding:"required"`
	SpaceID uint   `json:"spaceId,omitempty" binding:"required"`
}

type GetNotesQuery struct {
	Limit   int    `form:"limit"`
	Cursor  string `form:"cursor"`
	SpaceID uint   `form:"spaceId"`
}

// Response DTOs

type Note struct {
	ID        uint      `json:"id"`
	SpaceID   uint      `json:"spaceId"`
	UserID    uint      `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
