package api

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	"github.com/gin-gonic/gin"
)

type wechat struct{}

func (a *wechat) ContractSign(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}

	var req model.SignContractRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}

	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	params := map[string]string{
		"appid": req.AppID, "mch_id": req.MchID, "plan_id": req.PlanID, "out_contract_code": req.OutContractCode,
		"outer_openid": req.OutUserID, "contract_display_account": req.ContractDisplayAccount, "notify_url": req.NotifyURL,
		"sign_type": req.SignType, "version": req.Version, "timestamp": req.TimeStamp, "nonce": req.Nonce, "sign": req.Sign,
	}
	if err = serviceInfo.Signature.Verify(params, merchant.SignKey); err != nil {
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	if serviceInfo.Contract.HasActiveContract(req.OutContractCode) || serviceInfo.Contract.HasActiveContractByUser(merchant.ID, req.OutUserID, req.OutUserID) {
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "重复签约", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignExists, ErrCodeDes: "已有有效签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	contract := model.Contract{
		MerchantID:    merchant.ID,
		OpenID:        req.OutUserID,
		OutUserID:     req.OutUserID,
		OutContractID: req.OutContractCode,
		PlanID:        req.PlanID,
		NotifyURL:     req.NotifyURL,
		RequestData:   string(body),
	}
	statusRecord := model.ContractStatusRecord{
		MerchantID:     merchant.ID,
		OutContractID:  req.OutContractCode,
		ContractStatus: model.ContractStatusPending,
		IsFirstDeduct:  true,
	}
	if err = serviceInfo.Contract.CreateContractWithStatus(&contract, &statusRecord); err != nil {
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "创建签约失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeFail, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	record := model.ContractRecord{
		ContractID:    contract.ID,
		MerchantID:    merchant.ID,
		OperationType: "sign",
		RequestXML:    string(body),
		CallbackURL:   req.NotifyURL,
		Status:        model.ContractStatusPending,
	}
	_ = serviceInfo.Deduct.CreateContractRecord(&record)

	if merchant.SignStatusDelay > 0 {
		time.Sleep(time.Duration(merchant.SignStatusDelay) * time.Second)
	}
	_ = serviceInfo.Contract.UpdateContractStatus(contract.ID, merchant.SignTargetStatus, "")

	contractID := fmt.Sprintf("MOCK-C-%d", contract.ID)
	signSerialNo := fmt.Sprintf("MOCK-S-%d", time.Now().UnixNano())
	if merchant.SignTargetStatus == model.ContractStatusActive {
		_ = serviceInfo.Contract.SetContractID(contract.ID, contractID, signSerialNo)
		_ = serviceInfo.Contract.SetExpireTime(contract.ID, merchant.SignDurationMinutes)
	}

	resp := model.SignContractResponse{
		ReturnCode:      model.ErrCodeSuccess,
		ReturnMsg:       "OK",
		ResultCode:      model.ErrCodeSuccess,
		ContractID:      contractID,
		ContractExtID:   signSerialNo,
		OperationType:   "sign",
		MchID:           merchant.MchID,
		OutContractCode: req.OutContractCode,
		SignType:        req.SignType,
		TimeStamp:       strconv.FormatInt(time.Now().Unix(), 10),
		Nonce:           req.Nonce,
	}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{
		"return_code":       model.ErrCodeSuccess,
		"result_code":       model.ErrCodeSuccess,
		"contract_id":       contractID,
		"contract_ext_id":   signSerialNo,
		"mch_id":            merchant.MchID,
		"out_contract_code": req.OutContractCode,
		"sign_type":         req.SignType,
		"timestamp":         resp.TimeStamp,
		"nonce":             req.Nonce,
	}, merchant.SignKey)

	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
	_ = serviceInfo.Deduct.UpdateContractRecordResponse(record.ID, xmlResp, merchant.SignTargetStatus)

	if merchant.SignCallbackEnabled {
		if merchant.SignCallbackDelay > 0 {
			time.Sleep(time.Duration(merchant.SignCallbackDelay) * time.Second)
		}
		contract.ContractID = contractID
		contract.SignSerialNo = signSerialNo
		callbackXML := serviceInfo.Callback.BuildContractCallbackXML(contract, merchant.SignTargetStatus)
		result, callbackErr := serviceInfo.Callback.DoXMLCallback(req.NotifyURL, callbackXML)
		callbackTime := time.Now().Unix()
		if callbackErr != nil {
			result = callbackErr.Error() + "; " + result
		}
		_ = serviceInfo.Deduct.SetContractRecordCallbackResult(record.ID, result, callbackTime)
	}

	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) QueryContract(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.QueryContractRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	params := map[string]string{"appid": req.AppID, "mch_id": req.MchID, "contract_id": req.ContractID, "out_contract_code": req.OutContractCode, "sign_type": req.SignType, "timestamp": req.TimeStamp, "nonce": req.Nonce, "sign": req.Sign}
	if err = serviceInfo.Signature.Verify(params, merchant.SignKey); err != nil {
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	contract, err := serviceInfo.Deduct.GetContractFromDB(req.OutContractCode)
	if err != nil && req.ContractID != "" {
		contract, err = serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	}
	if err != nil {
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	statusRecord, statusErr := serviceInfo.Contract.GetContractStatusByContractID(contract.ID)
	if statusErr != nil {
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约状态不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约状态"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	resp := model.QueryContractResponse{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK", ResultCode: model.ErrCodeSuccess, ContractID: contract.ContractID, ContractStatus: statusRecord.ContractStatus, ContractExt: contract.SignSerialNo, PlanID: contract.PlanID, SignStatus: statusRecord.ContractStatus, SignType: req.SignType, TimeStamp: strconv.FormatInt(time.Now().Unix(), 10), Nonce: req.Nonce}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{"return_code": resp.ReturnCode, "result_code": resp.ResultCode, "contract_id": resp.ContractID, "contract_status": resp.ContractStatus, "contract_ext": resp.ContractExt, "plan_id": resp.PlanID, "sign_status": resp.SignStatus, "sign_type": resp.SignType, "timestamp": resp.TimeStamp, "nonce": resp.Nonce}, merchant.SignKey)
	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
	record := model.ContractRecord{ContractID: contract.ID, MerchantID: merchant.ID, OperationType: "query", RequestXML: string(body), ResponseXML: xmlResp, Status: statusRecord.ContractStatus}
	_ = serviceInfo.Deduct.CreateContractRecord(&record)
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) TerminateContract(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.TerminateContractRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	params := map[string]string{"appid": req.AppID, "mch_id": req.MchID, "contract_id": req.ContractID, "out_contract_code": req.OutContractCode, "contract_status": req.ContractStatus, "contract_ending_type": req.ContractEndingType, "sign_type": req.SignType, "timestamp": req.TimeStamp, "nonce": req.Nonce, "sign": req.Sign}
	if err = serviceInfo.Signature.Verify(params, merchant.SignKey); err != nil {
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	contract, err := serviceInfo.Deduct.GetContractFromDB(req.OutContractCode)
	if err != nil && req.ContractID != "" {
		contract, err = serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	}
	if err != nil {
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	_ = serviceInfo.Contract.UpdateContractStatus(contract.ID, model.ContractStatusTerminated, req.ContractEndingType)
	_ = serviceInfo.Contract.ResetFirstDeduct(contract.ID)
	record := model.ContractRecord{ContractID: contract.ID, MerchantID: merchant.ID, OperationType: "terminate", RequestXML: string(body), CallbackURL: contract.NotifyURL, Status: model.ContractStatusTerminated}
	_ = serviceInfo.Deduct.CreateContractRecord(&record)
	resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK", ResultCode: model.ErrCodeSuccess, ContractID: contract.ContractID, ContractStatus: model.ContractStatusTerminated, SignType: req.SignType, TimeStamp: strconv.FormatInt(time.Now().Unix(), 10), Nonce: req.Nonce}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{"return_code": resp.ReturnCode, "result_code": resp.ResultCode, "contract_id": resp.ContractID, "contract_status": resp.ContractStatus, "sign_type": resp.SignType, "timestamp": resp.TimeStamp, "nonce": resp.Nonce}, merchant.SignKey)
	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
	_ = serviceInfo.Deduct.UpdateContractRecordResponse(record.ID, xmlResp, model.ContractStatusTerminated)
	if merchant.TerminateNotifyEnabled && strings.TrimSpace(contract.NotifyURL) != "" {
		callbackXML := serviceInfo.Callback.BuildContractCallbackXML(contract, model.ContractStatusTerminated)
		result, callbackErr := serviceInfo.Callback.DoXMLCallback(contract.NotifyURL, callbackXML)
		if callbackErr != nil {
			result = callbackErr.Error() + "; " + result
		}
		_ = serviceInfo.Deduct.SetContractRecordCallbackResult(record.ID, result, time.Now().Unix())
	}
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) ApplyDeduct(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.DeductApplyRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	params := map[string]string{"appid": req.AppID, "mch_id": req.MchID, "out_trade_no": req.OutTradeNo, "contract_id": req.ContractID, "transaction_id": req.TransactionID, "total_amount": strconv.FormatInt(req.TotalAmount, 10), "fee_type": req.Currency, "notify_url": req.NotifyURL, "sign_type": req.SignType, "timestamp": req.TimeStamp, "nonce": req.Nonce, "sign": req.Sign}
	if err = serviceInfo.Signature.Verify(params, merchant.SignKey); err != nil {
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	contract, err := serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	if err != nil {
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	statusRecord, statusErr := serviceInfo.Contract.GetContractStatusByContractID(contract.ID)
	if statusErr != nil {
		record := model.DeductRecord{ContractID: contract.ID, MerchantID: merchant.ID, OperationType: "deduct", RequestData: string(body), CallbackURL: req.NotifyURL, OutTradeNo: req.OutTradeNo, TransactionID: req.TransactionID, Amount: req.TotalAmount, Status: model.DeductStatusFailed, IsFirstDeduct: false, PreNotifyCalled: false, ErrorCode: model.ErrCodeSignNotFound, ErrorMessage: "订阅状态不存在"}
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "订阅信息不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到订阅状态"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	if statusRecord.ContractStatus != model.ContractStatusActive {
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约状态不可扣款", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeDeductNotAllowed, ErrCodeDes: "签约未生效或已解约"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	isFirst := statusRecord.IsFirstDeduct
	if !isFirst && merchant.StrictDeductRule && !statusRecord.PreNotifyCalled {
		record := model.DeductRecord{ContractID: contract.ID, MerchantID: merchant.ID, OperationType: "deduct", RequestData: string(body), CallbackURL: req.NotifyURL, OutTradeNo: req.OutTradeNo, TransactionID: req.TransactionID, Amount: req.TotalAmount, Status: model.DeductStatusFailed, IsFirstDeduct: false, PreNotifyCalled: false, ErrorCode: model.ErrCodePreNotifyRequired, ErrorMessage: "非首次扣款前必须先调用预扣费通知API"}
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "未先调用预扣费通知", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodePreNotifyRequired, ErrCodeDes: "非首次扣款前必须先调用预扣费通知API"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	record := model.DeductRecord{ContractID: contract.ID, MerchantID: merchant.ID, OperationType: "deduct", RequestData: string(body), CallbackURL: req.NotifyURL, OutTradeNo: req.OutTradeNo, TransactionID: req.TransactionID, Amount: req.TotalAmount, Status: model.DeductStatusPending, IsFirstDeduct: isFirst, PreNotifyCalled: statusRecord.PreNotifyCalled}
	_ = serviceInfo.Deduct.SaveDeductRecord(&record)
	if merchant.DeductStatusDelay > 0 {
		time.Sleep(time.Duration(merchant.DeductStatusDelay) * time.Second)
	}
	_ = serviceInfo.Deduct.UpdateDeductRecordStatus(record.ID, merchant.DeductTargetStatus, "", "")
	if isFirst {
		_ = serviceInfo.Contract.MarkFirstDeductDone(contract.ID)
	} else {
		_ = serviceInfo.Contract.ClearPreNotify(contract.ID)
	}
	transactionID := req.TransactionID
	if strings.TrimSpace(transactionID) == "" {
		transactionID = fmt.Sprintf("MOCK-T-%d", time.Now().UnixNano())
		_ = serviceInfo.Deduct.SetDeductRecordTransactionID(record.ID, transactionID)
	}
	resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK", ResultCode: model.ErrCodeSuccess, MchID: merchant.MchID, OutTradeNo: req.OutTradeNo, TransactionID: transactionID, Amount: req.TotalAmount, SignType: req.SignType, TimeStamp: strconv.FormatInt(time.Now().Unix(), 10), Nonce: req.Nonce}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{"return_code": resp.ReturnCode, "result_code": resp.ResultCode, "mch_id": resp.MchID, "out_trade_no": resp.OutTradeNo, "transaction_id": resp.TransactionID, "amount": strconv.FormatInt(resp.Amount, 10), "sign_type": resp.SignType, "timestamp": resp.TimeStamp, "nonce": resp.Nonce}, merchant.SignKey)
	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
	_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xmlResp, merchant.DeductTargetStatus)
	if merchant.DeductCallbackEnabled {
		if merchant.DeductCallbackDelay > 0 {
			time.Sleep(time.Duration(merchant.DeductCallbackDelay) * time.Second)
		}
		record.TransactionID = transactionID
		record.Status = merchant.DeductTargetStatus
		callbackXML := serviceInfo.Callback.BuildDeductCallbackXML(record)
		result, callbackErr := serviceInfo.Callback.DoXMLCallback(req.NotifyURL, callbackXML)
		if callbackErr != nil {
			result = callbackErr.Error() + "; " + result
		}
		_ = serviceInfo.Deduct.SetCallbackResult(record.ID, result, time.Now().Unix())
	}
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) QueryDeduct(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.DeductApplyRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	record, err := serviceInfo.Deduct.GetDeductRecordByOutTradeNo(req.OutTradeNo)
	if err != nil {
		resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "扣款记录不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeFail, ErrCodeDes: "未找到扣款记录"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	resp := model.DeductApplyResponse{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK", ResultCode: model.ErrCodeSuccess, MchID: merchant.MchID, OutTradeNo: record.OutTradeNo, TransactionID: record.TransactionID, Amount: record.Amount, SignType: req.SignType, TimeStamp: strconv.FormatInt(time.Now().Unix(), 10), Nonce: req.Nonce}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{"return_code": resp.ReturnCode, "result_code": resp.ResultCode, "mch_id": resp.MchID, "out_trade_no": resp.OutTradeNo, "transaction_id": resp.TransactionID, "amount": strconv.FormatInt(resp.Amount, 10), "sign_type": resp.SignType, "timestamp": resp.TimeStamp, "nonce": resp.Nonce}, merchant.SignKey)
	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) PreDeductNotify(c *gin.Context) {
	var req model.PreDeductNotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: err.Error(), ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: err.Error()})
		return
	}
	if req.ContractID == "" {
		req.ContractID = c.Param("contract_id")
	}
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"})
		return
	}
	params := map[string]string{"appid": req.AppID, "mch_id": req.MchID, "contract_id": req.ContractID, "out_trade_no": req.OutTradeNo, "trade_no": req.TradeNo, "action_type": strconv.Itoa(req.ActionType), "account_id": req.AccountID, "notify_url": req.NotifyURL, "request_serial": strconv.FormatInt(req.RequestSerial, 10), "sign_type": req.SignType, "timestamp": req.TimeStamp, "nonce": req.Nonce, "sign": req.Sign}
	if err = serviceInfo.Signature.Verify(params, merchant.SignKey); err != nil {
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()})
		return
	}
	contract, err := serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	if err != nil {
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"})
		return
	}
	if _, statusErr := serviceInfo.Contract.GetContractStatusByContractID(contract.ID); statusErr != nil {
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "订阅状态不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到订阅状态"})
		return
	}
	bodyBytes := new(bytes.Buffer)
	_ = c.Request.Body.Close()
	bodyBytes.WriteString(fmt.Sprintf("appid=%s&mch_id=%s&contract_id=%s&out_trade_no=%s&trade_no=%s", req.AppID, req.MchID, req.ContractID, req.OutTradeNo, req.TradeNo))
	record := model.DeductRecord{ContractID: contract.ID, MerchantID: merchant.ID, OperationType: "pre_notify", RequestData: bodyBytes.String(), CallbackURL: req.NotifyURL, OutTradeNo: req.OutTradeNo, TransactionID: req.TradeNo, Status: model.DeductStatusSuccess, IsFirstDeduct: false, PreNotifyCalled: true}
	_ = serviceInfo.Deduct.SaveDeductRecord(&record)
	_ = serviceInfo.Contract.MarkPreNotifyCalled(contract.ID)
	resp := model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK", ResultCode: model.ErrCodeSuccess, AppID: req.AppID, MchID: req.MchID, SignType: req.SignType, TimeStamp: strconv.FormatInt(time.Now().Unix(), 10), Nonce: req.Nonce}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{"return_code": resp.ReturnCode, "result_code": resp.ResultCode, "appid": resp.AppID, "mch_id": resp.MchID, "sign_type": resp.SignType, "timestamp": resp.TimeStamp, "nonce": resp.Nonce}, merchant.SignKey)
	c.JSON(200, resp)
}
