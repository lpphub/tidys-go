package handlers

import (
	"tidys-go/logic"
	"tidys-go/logic/dto"
	"tidys-go/web/rest/helper"

	"github.com/gin-gonic/gin"
)

// TagList 获取所有标签（按分组分组）
func TagList(c *gin.Context) {
	var query dto.GetTagsQuery
	if !helper.MustBindQuery(c, &query) {
		return
	}

	data, err := logic.AppSvc.Tag.GetTags(c.Request.Context(), query.SpaceID)
	helper.Respond(c, err, data)
}

// TagGet 获取单个标签
func TagGet(c *gin.Context) {
	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	tag, err := logic.AppSvc.Tag.GetOne(c, id)
	helper.Respond(c, err, tag)
}

// TagCreate 创建标签
func TagCreate(c *gin.Context) {
	var req dto.TagReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	tag, err := logic.AppSvc.Tag.CreateTag(c.Request.Context(), req)
	helper.Respond(c, err, tag)
}

// TagUpdate 更新标签
func TagUpdate(c *gin.Context) {
	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.TagReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.Respond(c, logic.AppSvc.Tag.UpdateTag(c.Request.Context(), id, req))
}

// TagDelete 删除标签
func TagDelete(c *gin.Context) {
	id, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	helper.Respond(c, logic.AppSvc.Tag.DeleteTag(c.Request.Context(), id))
}

// TagReorder 重新排序标签
func TagReorder(c *gin.Context) {
	var req dto.ReorderTagReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	helper.Respond(c, logic.AppSvc.Tag.ReorderTag(c.Request.Context(), req))
}

// TagCreateGroup 创建分组
func TagCreateGroup(c *gin.Context) {
	var req dto.GroupReq
	if !helper.MustBindJSON(c, &req) {
		return
	}

	group, err := logic.AppSvc.Tag.CreateGroup(c.Request.Context(), req.Name, req.SpaceID)
	helper.Respond(c, err, group)
}

// TagDeleteGroup 删除分组
func TagDeleteGroup(c *gin.Context) {
	groupID, ok := helper.MustParseUintParam(c, "id")
	if !ok {
		return
	}

	// 支持通过查询参数传递 spaceId
	var query dto.GetTagsQuery
	_ = c.ShouldBindQuery(&query)

	helper.Respond(c, logic.AppSvc.Tag.DeleteGroup(c.Request.Context(), groupID, query.SpaceID))
}
