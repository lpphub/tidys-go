package errs

import (
	"net/http"

	"github.com/lpphub/goweb/base"
)

var (
	// 系统错误
	ErrServerError = base.NewErrorWithStatus(500, "server internal error", http.StatusInternalServerError)

	// 通用错误
	ErrNoToken        = base.NewErrorWithStatus(1000, "no token", http.StatusUnauthorized)
	ErrInvalidToken   = base.NewErrorWithStatus(1001, "invalid token", http.StatusUnauthorized)
	ErrInvalidParam   = base.NewError(1100, "参数错误")
	ErrRecordNotFound = base.NewError(1101, "数据不存在")

	// 业务错误
	ErrUserExists      = base.NewError(2101, "用户已存在")
	ErrUserNotExists   = base.NewError(2102, "用户不存在")
	ErrInvalidPassword = base.NewError(2103, "密码格式错误")
	ErrLoginFailed     = base.NewError(2104, "用户名或密码错误")
	ErrUserDisabled    = base.NewError(2105, "用户已禁用")

	ErrSpaceNotOwned   = base.NewError(2201, "空间无权限")
	ErrDuplicateInvite = base.NewError(2202, "已邀请，请勿重复邀请")
	ErrMaxInviteCount  = base.NewError(2203, "单次最多邀请100人")
)
