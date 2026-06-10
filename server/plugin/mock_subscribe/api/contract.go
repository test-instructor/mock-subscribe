package api

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type contract struct{}

func (a *contract) GetContractList(c *gin.Context) {
	start := time.Now()
	var pageInfo mockReq.ContractSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		LogError(c, "GetContractList:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "GetContractList", pageInfo)
	LogServiceCall(c, "Contract", "GetContractList", zap.Any("page", pageInfo))

	list, total, err := serviceInfo.Contract.GetContractList(pageInfo)
	if err != nil {
		LogError(c, "GetContractList:获取用户协议列表", err)
		global.GVA_LOG.Error("获取用户协议列表失败", zap.Error(err))
		response.FailWithMessage("获取用户协议列表失败: "+err.Error(), c)
		return
	}
	pageResult := response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}
	LogResponse(c, "GetContractList", pageResult, start)
	response.OkWithDetailed(pageResult, "获取成功", c)
}

func (a *contract) FindContract(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		LogError(c, "FindContract:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "FindContract", req)
	LogServiceCall(c, "Contract", "GetContract", zap.Any("id", req.ID))

	info, err := serviceInfo.Contract.GetContract(req.ID)
	if err != nil {
		LogError(c, "FindContract:获取用户协议详情", err)
		global.GVA_LOG.Error("获取用户协议详情失败", zap.Error(err))
		response.FailWithMessage("获取用户协议详情失败: "+err.Error(), c)
		return
	}
	LogServiceCall(c, "Contract", "GetContractStatusByContractID", zap.Any("contract_id", info.ID))
	status, statusErr := serviceInfo.Contract.GetContractStatusByContractID(info.ID)
	if statusErr != nil {
		LogError(c, "FindContract:获取用户协议状态", statusErr)
		respData := gin.H{"contract": info}
		LogResponse(c, "FindContract", respData, start)
		response.OkWithData(respData, c)
		return
	}
	respData := gin.H{"contract": info, "status": status}
	LogResponse(c, "FindContract", respData, start)
	response.OkWithData(respData, c)
}

func (a *contract) UpdateContractStatus(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID             uint   `json:"id"`
		ContractStatus string `json:"contractStatus"`
		TerminateType  string `json:"terminateType"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		LogError(c, "UpdateContractStatus:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.ID == 0 {
		LogError(c, "UpdateContractStatus:参数校验", nil, zap.String("reason", "签约ID不能为空"))
		response.FailWithMessage("签约ID不能为空", c)
		return
	}
	LogRequest(c, "UpdateContractStatus", req)
	LogServiceCall(c, "Contract", "UpdateContractStatus", zap.Any("req", req))

	if err := serviceInfo.Contract.UpdateContractStatus(req.ID, req.ContractStatus, req.TerminateType); err != nil {
		LogError(c, "UpdateContractStatus:更新用户协议状态", err)
		global.GVA_LOG.Error("更新用户协议状态失败", zap.Error(err))
		response.FailWithMessage("更新用户协议状态失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "UpdateContractStatus", req, start)
	response.OkWithMessage("更新用户协议状态成功", c)
}

func (a *contract) GetContractRecordList(c *gin.Context) {
	start := time.Now()
	var pageInfo mockReq.ContractRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		LogError(c, "GetContractRecordList:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "GetContractRecordList", pageInfo)
	LogServiceCall(c, "Deduct", "GetContractRecordList", zap.Any("page", pageInfo))

	list, total, err := serviceInfo.Deduct.GetContractRecordList(pageInfo)
	if err != nil {
		LogError(c, "GetContractRecordList:获取协议流水列表", err)
		global.GVA_LOG.Error("获取协议流水列表失败", zap.Error(err))
		response.FailWithMessage("获取协议流水列表失败: "+err.Error(), c)
		return
	}
	pageResult := response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}
	LogResponse(c, "GetContractRecordList", pageResult, start)
	response.OkWithDetailed(pageResult, "获取成功", c)
}

type updateContractReq struct {
	ID             uint   `json:"id"`
	ContractStatus string `json:"contractStatus"`
	TerminateType  string `json:"terminateType"`
}

func (a *contract) UpdateContract(c *gin.Context) {
	start := time.Now()
	var req updateContractReq
	if err := c.ShouldBindJSON(&req); err != nil {
		LogError(c, "UpdateContract:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.ID == 0 {
		LogError(c, "UpdateContract:参数校验", nil, zap.String("reason", "签约ID不能为空"))
		response.FailWithMessage("签约ID不能为空", c)
		return
	}
	LogRequest(c, "UpdateContract", req)
	LogServiceCall(c, "Contract", "UpdateContractStatus", zap.Any("req", req))

	if err := serviceInfo.Contract.UpdateContractStatus(req.ID, req.ContractStatus, req.TerminateType); err != nil {
		LogError(c, "UpdateContract:更新用户协议", err)
		global.GVA_LOG.Error("更新用户协议失败", zap.Error(err))
		response.FailWithMessage("更新用户协议失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "UpdateContract", req, start)
	response.OkWithMessage("更新用户协议成功", c)
}

func (a *contract) UpdateContractStatusV2(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID             uint   `json:"id" form:"id"`
		ContractStatus string `json:"contractStatus" form:"contractStatus"`
		TerminateType  string `json:"terminateType" form:"terminateType"`
	}
	if err := c.ShouldBind(&req); err != nil {
		LogError(c, "UpdateContractStatusV2:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.ID == 0 {
		LogError(c, "UpdateContractStatusV2:参数校验", nil, zap.String("reason", "签约ID不能为空"))
		response.FailWithMessage("签约ID不能为空", c)
		return
	}
	LogRequest(c, "UpdateContractStatusV2", req)
	LogServiceCall(c, "Contract", "UpdateContractStatus", zap.Any("req", req))

	if err := serviceInfo.Contract.UpdateContractStatus(req.ID, req.ContractStatus, req.TerminateType); err != nil {
		LogError(c, "UpdateContractStatusV2:更新用户协议状态", err)
		global.GVA_LOG.Error("更新用户协议状态失败", zap.Error(err))
		response.FailWithMessage("更新用户协议状态失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "UpdateContractStatusV2", req, start)
	response.OkWithMessage("更新用户协议状态成功", c)
}
