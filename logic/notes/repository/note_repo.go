package repository

import (
	"context"
	"tidys-go/logic/notes/model"
	"tidys-go/pkg/pagination"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type NoteRepo struct {
	*dbx.BaseRepo[model.Note]
}

func NewNoteRepo(db *gorm.DB) *NoteRepo {
	return &NoteRepo{
		BaseRepo: dbx.NewBaseRepo[model.Note](db),
	}
}

// CursorListBySpaceID lists notes with cursor pagination for a space
func (r *NoteRepo) CursorListBySpaceID(ctx context.Context, spaceID uint, page pagination.CursorQuery) (*pagination.CursorPageData[model.Note], error) {
	query := r.DB().WithContext(ctx).Model(&model.Note{}).
		Where("space_id = ?", spaceID)

	return pagination.QueryCursor[model.Note](query, page, pagination.OrderBy("id", true, func(n model.Note) uint { return n.ID }))
}
