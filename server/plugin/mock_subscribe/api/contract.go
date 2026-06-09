package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type contract struct{}

func (a *contract) GetContractList(c *gin.Context) {
	var pageInfo mockReq.ContractSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.Contract.GetContractList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取用户协议列表失败", zap.Error(err))
		response.FailWithMessage("获取用户协议列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *contract) FindContract(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := serviceInfo.Contract.GetContract(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取用户协议详情失败", zap.Error(err))
		response.FailWithMessage("获取用户协议详情失败: "+err.Error(), c)
		return
	}
	response.OkWithData(info, c)
}

func (a *contract) UpdateContractStatus(c *gin.Context) {
	var req struct {
		ID             uint   `json:"id"`
		ContractStatus string `json:"contractStatus"`
		TerminateType  string `json:"terminateType"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.ID == 0 {
		response.FailWithMessage("签约ID不能为空", c)
		return
	}
	if err := serviceInfo.Contract.UpdateContractStatus(req.ID, req.ContractStatus, req.TerminateType); err != nil {
		global.GVA_LOG.Error("更新用户协议状态失败", zap.Error(err))
		response.FailWithMessage("更新用户协议状态失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新用户协议状态成功", c)
}

func (a *contract) GetContractRecordList(c *gin.Context) {
	var pageInfo mockReq.ContractRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.Deduct.GetContractRecordList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取协议流水列表失败", zap.Error(err))
		response.FailWithMessage("获取协议流水列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

type updateContractReq struct {
	ID             uint   `json:"id"`
	ContractStatus string `json:"contractStatus"`
	TerminateType  string `json:"terminateType"`
}

func (a *contract) UpdateContract(c *gin.Context) {
	var req updateContractReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.ID == 0 {
		response.FailWithMessage("签约ID不能为空", c)
		return
	}
	if err := serviceInfo.Contract.UpdateContractStatus(req.ID, req.ContractStatus, req.TerminateType); err != nil {
		global.GVA_LOG.Error("更新用户协议失败", zap.Error(err))
		response.FailWithMessage("更新用户协议失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新用户协议成功", c)
}

func (a *contract) UpdateContractStatusV2(c *gin.Context) {
	var req struct {
		ID             uint   `json:"id" form:"id"`
		ContractStatus string `json:"contractStatus" form:"contractStatus"`
		TerminateType  string `json:"terminateType" form:"terminateType"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.ID == 0 {
		response.FailWithMessage("签约ID不能为空", c)
		return
	}
	if err := serviceInfo.Contract.UpdateContractStatus(req.ID, req.ContractStatus, req.TerminateType); err != nil {
		global.GVA_LOG.Error("更新用户协议状态失败", zap.Error(err))
		response.FailWithMessage("更新用户协议状态失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新用户协议状态成功", c)
}
