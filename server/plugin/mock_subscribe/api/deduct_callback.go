package api

import (
	"encoding/json"
	"io"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type deductCallback struct{}

func (a *deductCallback) ReceiveDeduct(c *gin.Context) {
	start := time.Now()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		LogError(c, "ReceiveDeduct:读取请求体", err)
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "读取请求失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		LogResponse(c, "ReceiveDeduct", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var req model.DeductNotifyRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		LogError(c, "ReceiveDeduct:XML解析", err, zap.String("raw_body", string(body)))
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "XML解析失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		LogResponse(c, "ReceiveDeduct", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	LogRequest(c, "ReceiveDeduct", gin.H{
		"raw_body":  string(body),
		"parsed":    req,
		"client_ip": clientIP(c),
	})

	headers, _ := json.Marshal(c.Request.Header)
	record := model.DeductCallbackRecord{
		MchID:         req.MchID,
		OutTradeNo:    req.OutTradeNo,
		TransactionID: req.TransactionID,
		TradeType:     req.TradeType,
		TradeState:    req.TradeState,
		BankType:      req.BankType,
		TotalAmount:   req.TotalAmount,
		CashAmount:    req.CashAmount,
		TimeEnd:       req.TimeEnd,
		SourceIP:      clientIP(c),
		Headers:       string(headers),
		RawBody:       string(body),
		Sign:          req.Sign,
	}

	LogServiceCall(c, "DeductCallback", "LocateMerchantAndDeduct", zap.Any("req", req))
	merchant, deductRecord, contract, locateErr := serviceInfo.DeductCallback.LocateMerchantAndDeduct(req)
	if locateErr == nil {
		record.MerchantID = merchant.ID
		record.ContractIDRef = contract.ID
		record.DeductRecordIDRef = deductRecord.ID
	}

	verifySign := true
	if locateErr == nil {
		verifySign = merchant.VerifySign
	}

	LogServiceCall(c, "DeductCallback", "ValidateDeductCallback", zap.Any("req", req), zap.Bool("verify_sign", verifySign))
	if err = serviceInfo.DeductCallback.ValidateDeductCallback(req, verifySign); err != nil {
		record.SignValid = false
		record.SignErrorMessage = err.Error()
		LogError(c, "ReceiveDeduct:扣款参数校验", err)
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.DeductCallback.Create(&record)
		LogResponse(c, "ReceiveDeduct", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	if locateErr != nil {
		record.SignValid = false
		record.SignErrorMessage = locateErr.Error()
		LogError(c, "ReceiveDeduct:定位商户或扣款记录", locateErr)
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "未找到商户或扣款记录"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.DeductCallback.Create(&record)
		LogResponse(c, "ReceiveDeduct", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	LogServiceCall(c, "DeductCallback", "VerifyDeductCallback", zap.Any("req", req))
	verifyErr := serviceInfo.DeductCallback.VerifyDeductCallback(req, merchant.VerifySign, merchant.SignKey)
	record.SignValid = verifyErr == nil
	if verifyErr != nil {
		record.SignErrorMessage = verifyErr.Error()
		LogError(c, "ReceiveDeduct:验签", verifyErr)
	}

	ack := model.GenericACK{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK"}
	if verifyErr != nil {
		ack = model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: verifyErr.Error()}
	}
	xml, _ := serviceInfo.XMLCodec.Marshal(ack)
	record.AckXML = xml
	LogServiceCall(c, "DeductCallback", "Create", zap.Any("record", record))
	if err = serviceInfo.DeductCallback.Create(&record); err != nil {
		LogError(c, "ReceiveDeduct:保存代扣回调记录", err)
		global.GVA_LOG.Error("保存代扣回调记录失败", zap.Error(err))
	}
	if verifyErr == nil {
		callbackResult := string(body)
		LogServiceCall(c, "Deduct", "UpdateDeductRecordByCallback", zap.Any("id", deductRecord.ID), zap.String("trade_state", req.TradeState))
		if err = serviceInfo.Deduct.UpdateDeductRecordByCallback(deductRecord.ID, req.TradeState, req.TransactionID, callbackResult, time.Now().Unix(), "", ""); err != nil {
			LogError(c, "ReceiveDeduct:回写扣款记录状态", err)
			global.GVA_LOG.Error("回写扣款记录状态失败", zap.Error(err))
		}
	}
	LogResponse(c, "ReceiveDeduct", string(xml), start)
	c.Data(200, "application/xml; charset=utf-8", []byte(xml))
}

func (a *deductCallback) GetDeductCallbackRecordList(c *gin.Context) {
	start := time.Now()
	var pageInfo mockReq.DeductCallbackRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		LogError(c, "GetDeductCallbackRecordList:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "GetDeductCallbackRecordList", pageInfo)
	LogServiceCall(c, "DeductCallback", "GetList", zap.Any("page", pageInfo))

	list, total, err := serviceInfo.DeductCallback.GetList(pageInfo)
	if err != nil {
		LogError(c, "GetDeductCallbackRecordList:获取代扣回调记录列表", err)
		global.GVA_LOG.Error("获取代扣回调记录列表失败", zap.Error(err))
		response.FailWithMessage("获取代扣回调记录列表失败: "+err.Error(), c)
		return
	}
	pageResult := response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}
	LogResponse(c, "GetDeductCallbackRecordList", pageResult, start)
	response.OkWithDetailed(pageResult, "获取成功", c)
}

func (a *deductCallback) FindDeductCallbackRecord(c *gin.Context) {
	start := time.Now()
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		LogError(c, "FindDeductCallbackRecord:参数绑定", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	LogRequest(c, "FindDeductCallbackRecord", req)
	LogServiceCall(c, "DeductCallback", "GetByID", zap.Any("id", req.ID))

	info, err := serviceInfo.DeductCallback.GetByID(req.ID)
	if err != nil {
		LogError(c, "FindDeductCallbackRecord:获取代扣回调记录详情", err)
		global.GVA_LOG.Error("获取代扣回调记录详情失败", zap.Error(err))
		response.FailWithMessage("获取代扣回调记录详情失败: "+err.Error(), c)
		return
	}
	LogResponse(c, "FindDeductCallbackRecord", info, start)
	response.OkWithData(info, c)
}
