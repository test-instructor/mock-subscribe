package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	toolsModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
	"github.com/gin-gonic/gin"
)

type environment struct{}

func (a *environment) CreateEnvironment(c *gin.Context) {
	var info toolsModel.Environment
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.Environment.CreateEnvironment(&info); err != nil {
		response.FailWithMessage("创建环境配置失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建环境配置成功", c)
}

func (a *environment) UpdateEnvironment(c *gin.Context) {
	var info toolsModel.Environment
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.Environment.UpdateEnvironment(&info); err != nil {
		response.FailWithMessage("更新环境配置失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新环境配置成功", c)
}

func (a *environment) DeleteEnvironment(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.Environment.DeleteEnvironment(req.ID); err != nil {
		response.FailWithMessage("删除环境配置失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除环境配置成功", c)
}

func (a *environment) FindEnvironment(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := serviceInfo.Environment.GetEnvironment(req.ID)
	if err != nil {
		response.FailWithMessage("获取环境配置失败: "+err.Error(), c)
		return
	}
	response.OkWithData(info, c)
}

func (a *environment) GetEnvironmentList(c *gin.Context) {
	var pageInfo toolsReq.EnvironmentSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.Environment.GetEnvironmentList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取环境配置列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
