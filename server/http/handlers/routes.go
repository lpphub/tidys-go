package handlers

import (
	"tidys-go/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 公共路由（无需认证）
	registerAuthRoutes(r.Group("/auth"))

	// API 路由（需要认证）
	api := r.Group("/", middleware.JwtAuth())
	registerTagRoutes(api)
	registerNoteRoutes(api)
	registerSpaceRoutes(api)
	registerPersonRoutes(api)
}

// registerAuthRoutes 注册认证路由
func registerAuthRoutes(r *gin.RouterGroup) {
	r.POST("/signup", AuthRegister)
	r.POST("/signin", AuthLogin)
	r.PUT("/refresh", AuthRefreshToken)
}

// registerTagRoutes 注册标签路由
func registerTagRoutes(r *gin.RouterGroup) {
	api := r.Group("/tags")
	{
		api.GET("", TagList)
		api.GET("/:id", TagGet)
		api.POST("", TagCreate)
		api.PATCH("/:id", TagUpdate)
		api.DELETE("/:id", TagDelete)
		api.POST("/reorder", TagReorder)
		api.POST("/group", TagCreateGroup)
		api.DELETE("/group/:id", TagDeleteGroup)
	}
}

// registerNoteRoutes 注册笔记路由
func registerNoteRoutes(r *gin.RouterGroup) {
	api := r.Group("/notes")
	{
		api.GET("", NoteList)
		api.POST("", NoteCreate)
		api.PATCH("/:id", NoteUpdate)
		api.DELETE("/:id", NoteDelete)
	}
}

// registerSpaceRoutes 注册空间路由
func registerSpaceRoutes(r *gin.RouterGroup) {
	api := r.Group("/spaces")
	{
		api.GET("", SpaceList)
		api.POST("", SpaceCreate)
		api.PATCH("/:id", SpaceUpdate)
		api.DELETE("/:id", SpaceDelete)

		// 成员管理
		api.GET("/:id/members", SpaceGetMembers)
		api.POST("/:id/members/invite", SpaceInviteMember)
		api.DELETE("/:id/members/:userId", SpaceRemoveMember)
	}

	// 邀请管理
	inviteApi := r.Group("/invites")
	{
		inviteApi.GET("/pending", SpaceGetPendingInvites)
		inviteApi.PATCH("/:id/respond", SpaceRespondInvite)
	}
}

// registerPersonRoutes 注册用户路由
func registerPersonRoutes(r *gin.RouterGroup) {
	api := r.Group("/person")
	{
		api.GET("/profile", UserGetProfile)
		api.PUT("/profile", UserUpdateProfile)
		api.PUT("/password", UserChangePassword)
	}
}
