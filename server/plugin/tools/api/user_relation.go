package api

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
	"github.com/gin-gonic/gin"
)

type userRelation struct{}

func (a *userRelation) CreateUserRelation(c *gin.Context) {
	var info toolsReq.UserRelationCreate
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	created, err := serviceInfo.UserRelation.CreateUserRelation(info)
	if err != nil {
		response.FailWithMessage("批量创建用户关联失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("批量创建成功，已新增"+strconv.Itoa(created)+"条记录", c)
}

func (a *userRelation) FindUserRelation(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := serviceInfo.UserRelation.GetUserRelation(req.ID)
	if err != nil {
		response.FailWithMessage("获取用户关联失败: "+err.Error(), c)
		return
	}
	response.OkWithData(info, c)
}

func (a *userRelation) GetUserRelationList(c *gin.Context) {
	var pageInfo toolsReq.UserRelationSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.UserRelation.GetUserRelationList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取用户关联列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *userRelation) DeleteUserRelation(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.UserRelation.DeleteUserRelation(req.ID); err != nil {
		response.FailWithMessage("删除用户关联失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除用户关联成功", c)
}

func (a *userRelation) GetUserIdsByEnvironment(c *gin.Context) {
	var req struct {
		EnvironmentKey string `form:"environmentKey"`
		Limit          int    `form:"limit"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.EnvironmentKey == "" {
		response.FailWithMessage("环境Key不能为空", c)
		return
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}
	ids, err := serviceInfo.UserRelation.GetUserIdsByEnvironmentKey(req.EnvironmentKey, req.Limit)
	if err != nil {
		response.FailWithMessage("获取用户ID列表失败: "+err.Error(), c)
		return
	}
	response.OkWithData(ids, c)
}
