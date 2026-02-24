package handlers

import (
	"tidys-go/logic"
	"tidys-go/logic/dto"
	"tidys-go/server/http/helper"

	"github.com/gin-gonic/gin"
)

// NoteList returns paginated notes list
func NoteList(c *gin.Context) {
	var query dto.GetNotesQuery
	if !helper.MustBindQuery(c, &query) {
		return
	}

	data, err := logic.AppSvc.Note.GetNotesList(c.Request.Context(), query)
	helper.Respond(c, err, data)
}

// NoteCreate creates a new note
func NoteCreate(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	var req dto.NoteReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	note, err := logic.AppSvc.Note.CreateNote(c.Request.Context(), userID, req)
	helper.Respond(c, err, note)
}

// NoteUpdate updates an existing note
func NoteUpdate(c *gin.Context) {
	var req dto.NoteReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.Respond(c, logic.AppSvc.Note.UpdateNote(c.Request.Context(), req.ID, req))
}

// NoteDelete deletes a note
func NoteDelete(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	helper.Respond(c, logic.AppSvc.Note.DeleteNote(c.Request.Context(), userID, id))
}
