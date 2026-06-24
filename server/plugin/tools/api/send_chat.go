package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
	"github.com/gin-gonic/gin"
)

type sendChat struct{}

func (a *sendChat) CreateSendChatTask(c *gin.Context) {
	var req toolsReq.SendChatCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	taskID, err := serviceInfo.SendChat.CreateSendChatTask(req)
	if err != nil {
		response.FailWithMessage("创建发送任务失败: "+err.Error(), c)
		return
	}
	response.OkWithData(map[string]uint{"taskId": taskID}, c)
}

func (a *sendChat) GetSendChatTaskList(c *gin.Context) {
	var pageInfo toolsReq.SendChatSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.SendChat.GetSendChatTaskList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取任务列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *sendChat) StopSendChatTask(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.SendChat.StopSendChatTask(req.ID); err != nil {
		response.FailWithMessage("停止任务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("任务已停止", c)
}
