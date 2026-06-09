package api

import (
	"encoding/json"
	"io"
	"net"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type callback struct{}

func (a *callback) ReceiveContract(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "读取请求失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var req model.ContractCallbackRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "XML解析失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	headers, _ := json.Marshal(c.Request.Header)
	callbackType := model.CallbackTypeContractSign
	if req.ChangeType == "DELETE" {
		callbackType = model.CallbackTypeTerminate
	}
	record := model.CallbackRecord{
		MchID:           req.MchID,
		OutContractCode: req.ContractCode,
		ContractCode:    req.ContractID,
		CallbackType:    callbackType,
		SourceIP:        clientIP(c),
		Headers:         string(headers),
		RawBody:         string(body),
		ContractStatus:  req.ChangeType,
		TimeStamp:       req.OperateTime,
		Sign:            req.Sign,
	}

	if err = serviceInfo.CallbackRecord.ValidateContractCallback(req); err != nil {
		record.SignValid = false
		record.SignErrorMessage = err.Error()
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.CallbackRecord.Create(&record)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	merchant, contract, locateErr := serviceInfo.CallbackRecord.LocateMerchantAndContract(req)
	if locateErr != nil {
		record.SignValid = false
		record.SignErrorMessage = locateErr.Error()
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "未找到签约或商户"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.CallbackRecord.Create(&record)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	record.MerchantID = merchant.ID
	record.ContractIDRef = contract.ID
	verifyErr := serviceInfo.CallbackRecord.VerifyContractCallback(req, merchant.SignKey)
	record.SignValid = verifyErr == nil
	if verifyErr != nil {
		record.SignErrorMessage = verifyErr.Error()
	}

	ack := model.GenericACK{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK"}
	if verifyErr != nil {
		ack = model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: verifyErr.Error()}
	}
	xml, _ := serviceInfo.XMLCodec.Marshal(ack)
	record.AckXML = xml
	if err = serviceInfo.CallbackRecord.Create(&record); err != nil {
		global.GVA_LOG.Error("保存回调记录失败", zap.Error(err))
	}
	c.Data(200, "application/xml; charset=utf-8", []byte(xml))
}

func (a *callback) GetCallbackRecordList(c *gin.Context) {
	var pageInfo mockReq.CallbackRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.CallbackRecord.GetList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取回调记录列表失败", zap.Error(err))
		response.FailWithMessage("获取回调记录列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *callback) FindCallbackRecord(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := serviceInfo.CallbackRecord.GetByID(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取回调记录详情失败", zap.Error(err))
		response.FailWithMessage("获取回调记录详情失败: "+err.Error(), c)
		return
	}
	response.OkWithData(info, c)
}

func clientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if strings.TrimSpace(ip) != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err == nil {
		return host
	}
	return c.Request.RemoteAddr
}
