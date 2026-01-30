package notes

import (
	"context"
	"tidys-go/logic/dto"
	"tidys-go/logic/notes/model"
	"tidys-go/logic/notes/repository"
	"tidys-go/pkg/pagination"
	"tidys-go/pkg/slices"
	"time"
)

type Service struct {
	noteRepo *repository.NoteRepo
}

func NewService(noteRepo *repository.NoteRepo) *Service {
	return &Service{
		noteRepo: noteRepo,
	}
}

// GetNotesList returns paginated notes list
func (s *Service) GetNotesList(ctx context.Context, query dto.GetNotesQuery) (*pagination.CursorPageData[dto.Note], error) {
	page := pagination.CursorQuery{
		Limit:  query.Limit,
		Cursor: query.Cursor,
	}

	result, err := s.noteRepo.CursorListBySpaceID(ctx, query.SpaceID, page)
	if err != nil {
		return nil, err
	}

	dtoNotes := slices.Map(result.List, func(n model.Note) dto.Note {
		return dto.Note{
			ID:        n.ID,
			SpaceID:   n.SpaceID,
			UserID:    n.UserID,
			Content:   n.Content,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		}
	})

	return &pagination.CursorPageData[dto.Note]{
		List:       dtoNotes,
		HasMore:    result.HasMore,
		NextCursor: result.NextCursor,
	}, nil
}

func (s *Service) CreateNote(ctx context.Context, userID uint, req dto.NoteReq) (*model.Note, error) {
	note := model.Note{
		SpaceID:   req.SpaceID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := s.noteRepo.Create(ctx, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

func (s *Service) UpdateNote(ctx context.Context, id uint, req dto.NoteReq) error {
	_, err := s.noteRepo.First(ctx, id)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"content":    req.Content,
		"updated_at": time.Now(),
	}

	return s.noteRepo.Update(ctx, id, updates)
}

// DeleteNote soft deletes a note
func (s *Service) DeleteNote(ctx context.Context, userID uint, id uint) error {
	return s.noteRepo.DB().WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Note{}).Error
}
