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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "读取请求失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var req model.DeductNotifyRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "XML解析失败"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

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

	if err = serviceInfo.DeductCallback.ValidateDeductCallback(req, verifySign); err != nil {
		record.SignValid = false
		record.SignErrorMessage = err.Error()
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.DeductCallback.Create(&record)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	if locateErr != nil {
		record.SignValid = false
		record.SignErrorMessage = locateErr.Error()
		ack := model.GenericACK{ReturnCode: model.ErrCodeFail, ReturnMsg: "未找到商户或扣款记录"}
		xml, _ := serviceInfo.XMLCodec.Marshal(ack)
		record.AckXML = xml
		_ = serviceInfo.DeductCallback.Create(&record)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	verifyErr := serviceInfo.DeductCallback.VerifyDeductCallback(req, merchant.VerifySign, merchant.SignKey)
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
	if err = serviceInfo.DeductCallback.Create(&record); err != nil {
		global.GVA_LOG.Error("保存代扣回调记录失败", zap.Error(err))
	}
	if verifyErr == nil {
		callbackResult := string(body)
		if err = serviceInfo.Deduct.UpdateDeductRecordByCallback(deductRecord.ID, req.TradeState, req.TransactionID, callbackResult, time.Now().Unix(), "", ""); err != nil {
			global.GVA_LOG.Error("回写扣款记录状态失败", zap.Error(err))
		}
	}
	c.Data(200, "application/xml; charset=utf-8", []byte(xml))
}

func (a *deductCallback) GetDeductCallbackRecordList(c *gin.Context) {
	var pageInfo mockReq.DeductCallbackRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceInfo.DeductCallback.GetList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取代扣回调记录列表失败", zap.Error(err))
		response.FailWithMessage("获取代扣回调记录列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

func (a *deductCallback) FindDeductCallbackRecord(c *gin.Context) {
	var req struct {
		ID uint `form:"ID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info, err := serviceInfo.DeductCallback.GetByID(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取代扣回调记录详情失败", zap.Error(err))
		response.FailWithMessage("获取代扣回调记录详情失败: "+err.Error(), c)
		return
	}
	response.OkWithData(info, c)
}
