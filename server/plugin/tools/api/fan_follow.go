package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
	"github.com/gin-gonic/gin"
)

type fanFollow struct{}

func (a *fanFollow) CreateFanFollow(c *gin.Context) {
	var req toolsReq.FanFollowCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	recordID, count, err := serviceInfo.FanFollow.CreateFanFollow(req)
	if err != nil {
		response.FailWithMessage("执行粉丝/关注/好友操作失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{
		"recordId":     recordID,
		"successCount": count,
	}, "操作完成", c)
}

func (a *fanFollow) GetFanFollowList(c *gin.Context) {
	var pageInfo toolsReq.FanFollowSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.FanFollow.GetFanFollowList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取执行记录失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
