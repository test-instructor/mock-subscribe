package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type merchant struct{}

func (a *merchant) CreateMerchant(c *gin.Context) {
	var info model.Merchant
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.Merchant.CreateMerchant(&info); err != nil {
		global.GVA_LOG.Error("创建商户配置失败", zap.Error(err))
		response.FailWithMessage("创建商户配置失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建商户配置成功", c)
}

func (a *merchant) UpdateMerchant(c *gin.Context) {
	var info model.Merchant
	if err := c.ShouldBindJSON(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.Merchant.UpdateMerchant(&info); err != nil {
		global.GVA_LOG.Error("更新商户配置失败", zap.Error(err))
		response.FailWithMessage("更新商户配置失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新商户配置成功", c)
}

func (a *merchant) DeleteMerchant(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceInfo.Merchant.DeleteMerchant(req.ID); err != nil {
		global.GVA_LOG.Error("删除商户配置失败", zap.Error(err))
		response.FailWithMessage("删除商户配置失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除商户配置成功", c)
}

func (a *merchant) FindMerchant(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := serviceInfo.Merchant.GetMerchant(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取商户配置失败", zap.Error(err))
		response.FailWithMessage("获取商户配置失败: "+err.Error(), c)
		return
	}
	response.OkWithData(info, c)
}

func (a *merchant) GetMerchantList(c *gin.Context) {
	var pageInfo mockReq.MerchantSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.Merchant.GetMerchantList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取商户配置列表失败", zap.Error(err))
		response.FailWithMessage("获取商户配置列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
