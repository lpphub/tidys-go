package handlers

import (
	"tidys-go/logic"
	"tidys-go/logic/dto"
	"tidys-go/web/rest/helper"

	"github.com/gin-gonic/gin"
)

func UserGetProfile(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	user, err := logic.AppSvc.User.Get(c.Request.Context(), userID)
	helper.ResponseResult(c, err, user)
}

func UserUpdateProfile(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	var req dto.UpdateProfileReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.ResponseResult(c, logic.AppSvc.User.UpdateProfile(c.Request.Context(), userID, req))
}

func UserChangePassword(c *gin.Context) {
	userID, ok := helper.MustGetUserID(c)
	if !ok {
		return
	}

	var req dto.ChangePasswordReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.ResponseResult(c, logic.AppSvc.User.ChangePassword(c.Request.Context(), userID, req))
}
