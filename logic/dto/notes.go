package dto

import (
	"tidys-go/pkg/pagination"
	"time"
)

// Request DTOs

type NoteReq struct {
	ID      uint   `json:"id,omitempty"`
	Content string `json:"content,omitempty" binding:"required"`
	SpaceID uint   `json:"spaceId,omitempty" binding:"required"`
}

type GetNotesQuery struct {
	pagination.Cursor
	SpaceID uint   `form:"spaceId"`
	Day     string `form:"day"` // YYYY-MM-DD format for filtering by day
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
