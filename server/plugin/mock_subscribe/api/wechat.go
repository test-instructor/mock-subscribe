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
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
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
		"nonce":             resp.Nonce,
	}, merchant.SignKey)
	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
	_ = serviceInfo.Deduct.UpdateContractRecordResponse(record.ID, xmlResp, merchant.SignTargetStatus)
	if merchant.SignCallbackEnabled {
		if merchant.SignCallbackDelay > 0 {
			time.Sleep(time.Duration(merchant.SignCallbackDelay) * time.Second)
		}
		contract.ContractID = contractID
		callbackXML := serviceInfo.Callback.BuildContractCallbackXML(contract, merchant.MchID, merchant.SignTargetStatus, merchant.SignKey)
		result, callbackErr := serviceInfo.Callback.DoXMLCallback(req.NotifyURL, callbackXML)
		if callbackErr != nil {
			result = callbackErr.Error() + "; " + result
		}
		_ = serviceInfo.Deduct.UpdateContractRecordStatus(record.ID, merchant.SignTargetStatus, "", result)
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
	params := map[string]string{
		"appid":             req.AppID,
		"mch_id":            req.MchID,
		"contract_id":       req.ContractID,
		"out_contract_code": req.OutContractCode,
		"sign_type":         req.SignType,
		"timestamp":         req.TimeStamp,
		"nonce":             req.Nonce,
		"sign":              req.Sign,
	}
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var contract model.Contract
	if strings.TrimSpace(req.ContractID) != "" {
		contract, err = serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	}
	if err != nil && strings.TrimSpace(req.OutContractCode) != "" {
		contract, err = serviceInfo.Deduct.GetContractFromDB(req.OutContractCode)
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

	resp := model.QueryContractResponse{
		ReturnCode:     model.ErrCodeSuccess,
		ReturnMsg:      "OK",
		ResultCode:     model.ErrCodeSuccess,
		ContractID:     contract.ContractID,
		ContractStatus: statusRecord.ContractStatus,
		ContractExt:    contract.SignSerialNo,
		PlanID:         contract.PlanID,
		SignStatus:     statusRecord.ContractStatus,
		SignType:       req.SignType,
		TimeStamp:      strconv.FormatInt(time.Now().Unix(), 10),
		Nonce:          req.Nonce,
	}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{
		"return_code":     resp.ReturnCode,
		"result_code":     resp.ResultCode,
		"contract_id":     resp.ContractID,
		"contract_status": resp.ContractStatus,
		"contract_ext":    resp.ContractExt,
		"plan_id":         resp.PlanID,
		"sign_status":     resp.SignStatus,
		"sign_type":       resp.SignType,
		"timestamp":       resp.TimeStamp,
		"nonce":           resp.Nonce,
	}, merchant.SignKey)
	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
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
	params := map[string]string{
		"appid":                req.AppID,
		"mch_id":               req.MchID,
		"contract_id":          req.ContractID,
		"out_contract_code":    req.OutContractCode,
		"contract_status":      req.ContractStatus,
		"contract_ending_type": req.ContractEndingType,
		"sign_type":            req.SignType,
		"timestamp":            req.TimeStamp,
		"nonce":                req.Nonce,
		"sign":                 req.Sign,
	}
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var contract model.Contract
	if strings.TrimSpace(req.ContractID) != "" {
		contract, err = serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	}
	if err != nil && strings.TrimSpace(req.OutContractCode) != "" {
		contract, err = serviceInfo.Deduct.GetContractFromDB(req.OutContractCode)
	}
	if err != nil {
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	terminateType := strings.TrimSpace(req.ContractEndingType)
	if terminateType == "" {
		terminateType = model.TerminateTypeMerchantRequest
	}
	_ = serviceInfo.Contract.UpdateContractStatus(contract.ID, model.ContractStatusTerminated, terminateType)
	_ = serviceInfo.Deduct.SetContractStatus(contract.ID, model.ContractStatusTerminated, terminateType)

	resp := model.TerminateContractResponse{
		ReturnCode:     model.ErrCodeSuccess,
		ReturnMsg:      "OK",
		ResultCode:     model.ErrCodeSuccess,
		ContractID:     contract.ContractID,
		ContractStatus: model.ContractStatusTerminated,
		SignType:       req.SignType,
		TimeStamp:      strconv.FormatInt(time.Now().Unix(), 10),
		Nonce:          req.Nonce,
	}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{
		"return_code":     resp.ReturnCode,
		"result_code":     resp.ResultCode,
		"contract_id":     resp.ContractID,
		"contract_status": resp.ContractStatus,
		"sign_type":       resp.SignType,
		"timestamp":       resp.TimeStamp,
		"nonce":           resp.Nonce,
	}, merchant.SignKey)
	xmlResp, _ := serviceInfo.XMLCodec.Marshal(resp)
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

	buildResp := func(returnMsg, resultCode, errCode, errCodeDes, mchID, outTradeNo, transactionID string, amount int64, signType, nonce string) model.DeductApplyResponse {
		resp := model.DeductApplyResponse{
			ReturnCode:    model.ErrCodeFail,
			ReturnMsg:     returnMsg,
			ResultCode:    resultCode,
			ErrCode:       errCode,
			ErrCodeDes:    errCodeDes,
			MchID:         mchID,
			OutTradeNo:    outTradeNo,
			TransactionID: transactionID,
			TotalAmount:   amount,
			SignType:      signType,
			TimeStamp:     strconv.FormatInt(time.Now().Unix(), 10),
			Nonce:         nonce,
		}
		if resultCode == model.ErrCodeSuccess {
			resp.ReturnCode = model.ErrCodeSuccess
			resp.ReturnMsg = "OK"
		}
		return resp
	}

	writeXML := func(resp model.DeductApplyResponse) string {
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return xml
	}

	effectiveAmount := req.TotalFee
	if effectiveAmount <= 0 {
		effectiveAmount = req.TotalAmount
	}
	effectiveNonce := strings.TrimSpace(req.NonceStr)
	if effectiveNonce == "" {
		effectiveNonce = strings.TrimSpace(req.Nonce)
	}

	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		resp := buildResp("商户不存在", model.ErrCodeFail, model.ErrCodeInvalidParams, "商户配置不存在", req.MchID, req.OutTradeNo, req.TransactionID, effectiveAmount, req.SignType, effectiveNonce)
		writeXML(resp)
		return
	}

	params := map[string]string{
		"appid":          req.AppID,
		"mch_id":         req.MchID,
		"body":           req.Body,
		"detail":         req.Detail,
		"attach":         req.Attach,
		"out_trade_no":   req.OutTradeNo,
		"contract_id":    req.ContractID,
		"transaction_id": req.TransactionID,
		"total_fee":      strconv.FormatInt(effectiveAmount, 10),
		"fee_type":       req.FeeType,
		"notify_url":     req.NotifyURL,
		"trade_type":     req.TradeType,
		"device_info":    req.DeviceInfo,
		"nonce_str":      effectiveNonce,
		"sign_type":      req.SignType,
		"timestamp":      req.TimeStamp,
		"sign":           req.Sign,
	}
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		resp := buildResp("签名校验失败", model.ErrCodeFail, model.ErrCodeInvalidSign, err.Error(), merchant.MchID, req.OutTradeNo, req.TransactionID, effectiveAmount, req.SignType, effectiveNonce)
		writeXML(resp)
		return
	}

	contract, err := serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	if err != nil {
		resp := buildResp("签约不存在", model.ErrCodeFail, model.ErrCodeSignNotFound, "未找到签约关系", merchant.MchID, req.OutTradeNo, req.TransactionID, effectiveAmount, req.SignType, effectiveNonce)
		writeXML(resp)
		return
	}

	transactionID := strings.TrimSpace(req.TransactionID)
	if transactionID == "" {
		transactionID = fmt.Sprintf("MOCK-T-%d", time.Now().UnixNano())
	}
	record := model.DeductRecord{
		ContractID:      contract.ID,
		MerchantID:      merchant.ID,
		OperationType:   "deduct",
		RequestData:     string(body),
		CallbackURL:     req.NotifyURL,
		TransactionID:   transactionID,
		OutTradeNo:      req.OutTradeNo,
		Amount:          effectiveAmount,
		Status:          model.DeductStatusPending,
		IsFirstDeduct:   false,
		PreNotifyCalled: false,
	}

	statusRecord, statusErr := serviceInfo.Contract.GetContractStatusByContractID(contract.ID)
	if statusErr != nil {
		record.Status = model.DeductStatusFailed
		record.ErrorCode = model.ErrCodeSignNotFound
		record.ErrorMessage = "订阅状态不存在"
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		resp := buildResp("订阅信息不存在", model.ErrCodeFail, model.ErrCodeSignNotFound, "未找到订阅状态", merchant.MchID, req.OutTradeNo, transactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xml, record.Status)
		return
	}

	record.IsFirstDeduct = statusRecord.IsFirstDeduct
	record.PreNotifyCalled = statusRecord.PreNotifyCalled

	if statusRecord.ContractStatus != model.ContractStatusActive {
		record.Status = model.DeductStatusFailed
		record.ErrorCode = model.ErrCodeDeductNotAllowed
		record.ErrorMessage = "签约未生效或已解约"
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		resp := buildResp("签约状态不可扣款", model.ErrCodeFail, model.ErrCodeDeductNotAllowed, "签约未生效或已解约", merchant.MchID, req.OutTradeNo, transactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xml, record.Status)
		return
	}

	if !record.IsFirstDeduct && merchant.StrictDeductRule && !statusRecord.PreNotifyCalled {
		record.Status = model.DeductStatusFailed
		record.ErrorCode = model.ErrCodePreNotifyRequired
		record.ErrorMessage = "非首次扣款前必须先调用预扣费通知API"
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		resp := buildResp("未先调用预扣费通知", model.ErrCodeFail, model.ErrCodePreNotifyRequired, "非首次扣款前必须先调用预扣费通知API", merchant.MchID, req.OutTradeNo, transactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xml, record.Status)
		return
	}

	_ = serviceInfo.Deduct.SaveDeductRecord(&record)
	if merchant.DeductStatusDelay > 0 {
		time.Sleep(time.Duration(merchant.DeductStatusDelay) * time.Second)
	}

	finalStatus := strings.TrimSpace(merchant.DeductTargetStatus)
	if finalStatus == "" {
		finalStatus = model.DeductStatusSuccess
	}
	_ = serviceInfo.Deduct.UpdateDeductRecordStatus(record.ID, finalStatus, "", "")
	record.Status = finalStatus

	if record.IsFirstDeduct {
		_ = serviceInfo.Contract.MarkFirstDeductDone(contract.ID)
	} else {
		_ = serviceInfo.Contract.ClearPreNotify(contract.ID)
	}

	resp := buildResp("OK", model.ErrCodeSuccess, "", "", merchant.MchID, req.OutTradeNo, transactionID, effectiveAmount, req.SignType, effectiveNonce)
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{
		"return_code":    resp.ReturnCode,
		"result_code":    resp.ResultCode,
		"mch_id":         resp.MchID,
		"out_trade_no":   resp.OutTradeNo,
		"transaction_id": resp.TransactionID,
		"total_amount":   strconv.FormatInt(resp.TotalAmount, 10),
		"sign_type":      resp.SignType,
		"timestamp":      resp.TimeStamp,
		"nonce":          resp.Nonce,
	}, merchant.SignKey)
	xmlResp := writeXML(resp)
	_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xmlResp, finalStatus)

	if merchant.DeductCallbackEnabled {
		if merchant.DeductCallbackDelay > 0 {
			time.Sleep(time.Duration(merchant.DeductCallbackDelay) * time.Second)
		}
		callbackXML := serviceInfo.Callback.BuildDeductCallbackXML(merchant, contract, record, req.SignType)
		callbackTarget := strings.TrimSpace(req.NotifyURL)
		if callbackTarget == "" {
			callbackTarget = strings.TrimSpace(contract.NotifyURL)
		}
		result, callbackErr := serviceInfo.Callback.DoXMLCallback(callbackTarget, callbackXML)
		if callbackErr != nil {
			result = callbackErr.Error() + "; " + result
		}
		_ = serviceInfo.Deduct.SetCallbackResult(record.ID, result, time.Now().Unix())
	}
}

func (a *wechat) QueryDeduct(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.QueryDeductRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	effectiveNonce := strings.TrimSpace(req.NonceStr)
	if effectiveNonce == "" {
		effectiveNonce = strings.TrimSpace(req.Nonce)
	}
	params := map[string]string{
		"appid":          req.AppID,
		"mch_id":         req.MchID,
		"out_trade_no":   req.OutTradeNo,
		"transaction_id": req.TransactionID,
		"sign_type":      req.SignType,
		"timestamp":      req.TimeStamp,
		"nonce_str":      effectiveNonce,
		"sign":           req.Sign,
	}
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var record model.DeductRecord
	if strings.TrimSpace(req.OutTradeNo) != "" {
		record, err = serviceInfo.Deduct.GetDeductRecordByOutTradeNo(req.OutTradeNo)
	}
	if err != nil && strings.TrimSpace(req.TransactionID) != "" {
		record, err = serviceInfo.Deduct.GetDeductRecordByTransactionID(req.TransactionID)
	}
	if err != nil {
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "订单不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeFail, ErrCodeDes: "未找到扣款记录"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	contract, contractErr := serviceInfo.Deduct.GetContractByID(record.ContractID)
	if contractErr != nil {
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	tradeState := record.Status
	if strings.TrimSpace(tradeState) == "" {
		tradeState = model.DeductStatusPending
	}
	timeEnd := ""
	if !record.UpdatedAt.IsZero() {
		timeEnd = record.UpdatedAt.Format("20060102150405")
	}
	resp := model.QueryDeductResponse{
		ReturnCode:    model.ErrCodeSuccess,
		ReturnMsg:     "OK",
		ResultCode:    model.ErrCodeSuccess,
		AppID:         req.AppID,
		MchID:         merchant.MchID,
		OpenID:        contract.OpenID,
		TradeType:     "PAP",
		TradeState:    tradeState,
		BankType:      "MOCK",
		TotalAmount:   record.Amount,
		CashAmount:    record.Amount,
		TransactionID: record.TransactionID,
		OutTradeNo:    record.OutTradeNo,
		TimeEnd:       timeEnd,
		SignType:      req.SignType,
		TimeStamp:     strconv.FormatInt(time.Now().Unix(), 10),
		Nonce:         req.Nonce,
	}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{
		"return_code":    resp.ReturnCode,
		"result_code":    resp.ResultCode,
		"appid":          resp.AppID,
		"mch_id":         resp.MchID,
		"openid":         resp.OpenID,
		"trade_type":     resp.TradeType,
		"trade_state":    resp.TradeState,
		"bank_type":      resp.BankType,
		"total_amount":   strconv.FormatInt(resp.TotalAmount, 10),
		"cash_amount":    strconv.FormatInt(resp.CashAmount, 10),
		"transaction_id": resp.TransactionID,
		"out_trade_no":   resp.OutTradeNo,
		"time_end":       resp.TimeEnd,
		"sign_type":      resp.SignType,
		"timestamp":      resp.TimeStamp,
		"nonce":          resp.Nonce,
	}, merchant.SignKey)
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
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
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
