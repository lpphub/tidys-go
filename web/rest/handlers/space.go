package handlers

import (
	"tidys-go/logic"
	"tidys-go/logic/dto"
	"tidys-go/web/rest/helper"

	"github.com/gin-gonic/gin"
)

// SpaceList 获取用户的所有空间
func SpaceList(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	data, err := logic.AppSvc.Space.GetSpaces(c.Request.Context(), userID)
	helper.Respond(c, err, data)
}

// SpaceCreate 创建空间
func SpaceCreate(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	var req dto.SpaceReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	var spaceID uint
	space, err := logic.AppSvc.Space.CreateSpace(c.Request.Context(), userID, req)
	if space != nil {
		spaceID = space.ID
	}
	helper.Respond(c, err, spaceID)
}

// SpaceUpdate 更新空间
func SpaceUpdate(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	spaceID, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.SpaceReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.Respond(c, logic.AppSvc.Space.UpdateSpace(c.Request.Context(), spaceID, userID, req))
}

// SpaceDelete 删除空间
func SpaceDelete(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	spaceID, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	helper.Respond(c, logic.AppSvc.Space.DeleteSpace(c.Request.Context(), spaceID, userID))
}

// SpaceGetMembers 获取空间成员列表
func SpaceGetMembers(c *gin.Context) {
	spaceID, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	data, err := logic.AppSvc.Space.GetMembers(c.Request.Context(), spaceID)
	helper.Respond(c, err, data)
}

// SpaceInviteMember 邀请成员
func SpaceInviteMember(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	spaceID, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.SpaceInviteMemberReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.Respond(c, logic.AppSvc.Space.InviteMember(c.Request.Context(), spaceID, userID, req.Emails))
}

// SpaceRemoveMember 移除成员
func SpaceRemoveMember(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	spaceID, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	userIDToRemove, ok := helper.MustParseUintParam(c, "userId")
	if !ok {
		return
	}

	helper.Respond(c, logic.AppSvc.Space.RemoveMember(c.Request.Context(), spaceID, userID, userIDToRemove))
}

// SpaceGetPendingInvites 获取待处理邀请
func SpaceGetPendingInvites(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	data, err := logic.AppSvc.Space.GetPendingInvites(c.Request.Context(), userID)
	helper.Respond(c, err, data)
}

// SpaceRespondInvite 响应邀请
func SpaceRespondInvite(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.Respond(c, logic.AppSvc.Space.RespondInvite(c.Request.Context(), id, userID, req.Action))
}
