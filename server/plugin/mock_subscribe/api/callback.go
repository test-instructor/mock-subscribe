package api

import (
	"encoding/json"
	"io"
	"net"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type callback struct{}

func (a *callback) ReceiveContract(c *gin.Context) {
	start := time.Now()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		LogError(c, "ReceiveContract:读取请求体", err)
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "读取请求失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		LogResponse(c, "ReceiveContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var req model.ContractCallbackRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		LogError(c, "ReceiveContract:XML解析", err, zap.String("raw_body", string(body)))
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "XML解析失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		LogResponse(c, "ReceiveContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	LogRequest(c, "ReceiveContract", gin.H{
		"raw_body":  string(body),
		"parsed":    req,
		"client_ip": clientIP(c),
	})

	headers, _ := json.Marshal(c.Request.Header)
	callbackType := model.CallbackTypeContractSign
	if req.ChangeType == "DELETE" {
		callbackType = model.CallbackTypeTerminate
	}
	record := model.CallbackRecord{
		MchID:           req.MchID,
		OutContractCode: req.OutContractCode,
		ContractCode:    req.ContractID,
		CallbackType:    callbackType,
		SourceIP:        clientIP(c),
		Headers:         string(headers),
		RawBody:         string(body),
		ContractStatus:  req.ChangeType,
		TimeStamp:       req.OperateTime,
		Sign:            req.Sign,
	}

	LogServiceCall(c, "CallbackRecord", "LocateMerchantAndContract", zap.Any("req", req))
	merchant, contract, locateErr := serviceInfo.CallbackRecord.LocateMerchantAndContract(req)
	if locateErr == nil {
		record.MerchantID = merchant.ID
		record.ContractIDRef = contract.ID
	}

	verifySign := true
	if locateErr == nil {
		verifySign = merchant.VerifySign
	}

	LogServiceCall(c, "CallbackRecord", "ValidateContractCallback", zap.Any("req", req), zap.Bool("verify_sign", verifySign))
	if err = serviceInfo.CallbackRecord.ValidateContractCallback(req, verifySign); err != nil {
		record.SignValid = false
		record.SignErrorMessage = err.Error()
		LogError(c, "ReceiveContract:签约参数校验", err)
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.CallbackRecord.Create(&record)
		LogResponse(c, "ReceiveContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	if locateErr != nil {
		record.SignValid = false
		record.SignErrorMessage = locateErr.Error()
		LogError(c, "ReceiveContract:定位商户或签约", locateErr)
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "未找到签约或商户"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.CallbackRecord.Create(&record)
		LogResponse(c, "ReceiveContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	LogServiceCall(c, "CallbackRecord", "VerifyContractCallback", zap.Any("req", req))
	verifyErr := serviceInfo.CallbackRecord.VerifyContractCallback(req, merchant.VerifySign, merchant.SignKey)
	record.SignValid = verifyErr == nil
	if verifyErr != nil {
		record.SignErrorMessage = verifyErr.Error()
		LogError(c, "ReceiveContract:验签", verifyErr)
	}

	ack := model.GenericACK{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK"}
	if verifyErr != nil {
		ack = model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: verifyErr.Error()}
	}
	xml, _ := serviceInfo.XMLCodec.Marshal(ack)
	record.AckXML = xml
	LogServiceCall(c, "CallbackRecord", "Create", zap.Any("record", record))
	if err = serviceInfo.CallbackRecord.Create(&record); err != nil {
		LogError(c, "ReceiveContract:保存回调记录", err)
		global.GVA_LOG.Error("保存回调记录失败", zap.Error(err))
	}
	LogResponse(c, "ReceiveContract", string(xml), start)
	c.Data(200, "application/xml; charset=utf-8", []byte(xml))
}

func (a *callback) GetCallbackRecordList(c *gin.Context) {
	start := time.Now()
	var pageInfo mockReq.CallbackRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		LogError(c, "GetCallbackRecordList:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "GetCallbackRecordList", pageInfo)
	LogServiceCall(c, "CallbackRecord", "GetList", zap.Any("page", pageInfo))

	list, total, err := serviceInfo.CallbackRecord.GetList(pageInfo)
	if err != nil {
		LogError(c, "GetCallbackRecordList:获取回调记录列表", err)
		global.GVA_LOG.Error("获取回调记录列表失败", zap.Error(err))
		response.FailWithMessage("获取回调记录列表失败: "+err.Error(), c)
		return
	}
	pageResult := response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}
	LogResponse(c, "GetCallbackRecordList", pageResult, start)
	response.OkWithDetailed(pageResult, "获取成功", c)
}

func (a *callback) FindCallbackRecord(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		LogError(c, "FindCallbackRecord:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "FindCallbackRecord", req)
	LogServiceCall(c, "CallbackRecord", "GetByID", zap.Any("id", req.ID))

	info, err := serviceInfo.CallbackRecord.GetByID(req.ID)
	if err != nil {
		LogError(c, "FindCallbackRecord:获取回调记录详情", err)
		global.GVA_LOG.Error("获取回调记录详情失败", zap.Error(err))
		response.FailWithMessage("获取回调记录详情失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "FindCallbackRecord", info, start)
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
