package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type deduct struct{}

func (a *deduct) GetDeductRecordList(c *gin.Context) {
	var pageInfo mockReq.DeductRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.Deduct.GetDeductRecordList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取扣款记录列表失败", zap.Error(err))
		response.FailWithMessage("获取扣款记录列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *deduct) FindDeductRecord(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := serviceInfo.Deduct.GetDeductRecord(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取扣款记录详情失败", zap.Error(err))
		response.FailWithMessage("获取扣款记录详情失败: "+err.Error(), c)
		return
	}
	response.OkWithData(info, c)
}
