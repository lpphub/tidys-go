package handlers

import (
	"tidys-go/logic"
	"tidys-go/logic/dto"
	"tidys-go/web/rest/helper"

	"github.com/gin-gonic/gin"
)

// AuthRegister 用户注册
func AuthRegister(c *gin.Context) {
	var req dto.AuthReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	authData, err := logic.AppSvc.Auth.Register(c, req)
	helper.ResponseResult(c, err, authData)
}

// AuthLogin 用户登录
func AuthLogin(c *gin.Context) {
	var req dto.AuthReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	authData, err := logic.AppSvc.Auth.Login(c, req)
	helper.ResponseResult(c, err, authData)
}

// AuthRefreshToken 刷新 token
func AuthRefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if !helper.MustBindJSON(c, &req) {
		return
	}

	tokenPair, err := logic.AppSvc.Auth.RefreshToken(c, req.RefreshToken)
	helper.ResponseResult(c, err, tokenPair)
}
