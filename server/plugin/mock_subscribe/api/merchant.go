package api

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type merchant struct{}

func (a *merchant) CreateMerchant(c *gin.Context) {
	start := time.Now()
	var info model.Merchant
	if err := c.ShouldBindJSON(&info); err != nil {
		LogError(c, "CreateMerchant:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "CreateMerchant", info)
	LogServiceCall(c, "Merchant", "CreateMerchant", zap.Any("merchant", info))

	if err := serviceInfo.Merchant.CreateMerchant(&info); err != nil {
		LogError(c, "CreateMerchant:创建商户配置", err)
		global.GVA_LOG.Error("创建商户配置失败", zap.Error(err))
		response.FailWithMessage("创建商户配置失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "CreateMerchant", info, start)
	response.OkWithMessage("创建商户配置成功", c)
}

func (a *merchant) UpdateMerchant(c *gin.Context) {
	start := time.Now()
	var info model.Merchant
	if err := c.ShouldBindJSON(&info); err != nil {
		LogError(c, "UpdateMerchant:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "UpdateMerchant", info)
	LogServiceCall(c, "Merchant", "UpdateMerchant", zap.Any("merchant", info))

	if err := serviceInfo.Merchant.UpdateMerchant(&info); err != nil {
		LogError(c, "UpdateMerchant:更新商户配置", err)
		global.GVA_LOG.Error("更新商户配置失败", zap.Error(err))
		response.FailWithMessage("更新商户配置失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "UpdateMerchant", info, start)
	response.OkWithMessage("更新商户配置成功", c)
}

func (a *merchant) DeleteMerchant(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		LogError(c, "DeleteMerchant:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "DeleteMerchant", req)
	LogServiceCall(c, "Merchant", "DeleteMerchant", zap.Any("id", req.ID))

	if err := serviceInfo.Merchant.DeleteMerchant(req.ID); err != nil {
		LogError(c, "DeleteMerchant:删除商户配置", err)
		global.GVA_LOG.Error("删除商户配置失败", zap.Error(err))
		response.FailWithMessage("删除商户配置失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "DeleteMerchant", req, start)
	response.OkWithMessage("删除商户配置成功", c)
}

func (a *merchant) FindMerchant(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		LogError(c, "FindMerchant:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "FindMerchant", req)
	LogServiceCall(c, "Merchant", "GetMerchant", zap.Any("id", req.ID))

	info, err := serviceInfo.Merchant.GetMerchant(req.ID)
	if err != nil {
		LogError(c, "FindMerchant:获取商户配置", err)
		global.GVA_LOG.Error("获取商户配置失败", zap.Error(err))
		response.FailWithMessage("获取商户配置失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "FindMerchant", info, start)
	response.OkWithData(info, c)
}

func (a *merchant) GetMerchantList(c *gin.Context) {
	start := time.Now()
	var pageInfo mockReq.MerchantSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		LogError(c, "GetMerchantList:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "GetMerchantList", pageInfo)
	LogServiceCall(c, "Merchant", "GetMerchantList", zap.Any("page", pageInfo))

	list, total, err := serviceInfo.Merchant.GetMerchantList(pageInfo)
	if err != nil {
		LogError(c, "GetMerchantList:获取商户配置列表", err)
		global.GVA_LOG.Error("获取商户配置列表失败", zap.Error(err))
		response.FailWithMessage("获取商户配置列表失败: "+err.Error(), c)
		return
	}
	pageResult := response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}
	LogResponse(c, "GetMerchantList", pageResult, start)
	response.OkWithDetailed(pageResult, "获取成功", c)
}
