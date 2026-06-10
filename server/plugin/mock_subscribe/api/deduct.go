package api

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type deduct struct{}

func (a *deduct) GetDeductRecordList(c *gin.Context) {
	start := time.Now()
	var pageInfo mockReq.DeductRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		LogError(c, "GetDeductRecordList:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "GetDeductRecordList", pageInfo)
	LogServiceCall(c, "Deduct", "GetDeductRecordList", zap.Any("page", pageInfo))

	list, total, err := serviceInfo.Deduct.GetDeductRecordList(pageInfo)
	if err != nil {
		LogError(c, "GetDeductRecordList:获取扣款记录列表", err)
		global.GVA_LOG.Error("获取扣款记录列表失败", zap.Error(err))
		response.FailWithMessage("获取扣款记录列表失败: "+err.Error(), c)
		return
	}
	pageResult := response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}
	LogResponse(c, "GetDeductRecordList", pageResult, start)
	response.OkWithDetailed(pageResult, "获取成功", c)
}

func (a *deduct) FindDeductRecord(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		LogError(c, "FindDeductRecord:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "FindDeductRecord", req)
	LogServiceCall(c, "Deduct", "GetDeductRecord", zap.Any("id", req.ID))

	info, err := serviceInfo.Deduct.GetDeductRecord(req.ID)
	if err != nil {
		LogError(c, "FindDeductRecord:获取扣款记录详情", err)
		global.GVA_LOG.Error("获取扣款记录详情失败", zap.Error(err))
		response.FailWithMessage("获取扣款记录详情失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "FindDeductRecord", info, start)
	response.OkWithData(info, c)
}
