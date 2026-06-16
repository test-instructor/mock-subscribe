package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type wechat struct{}

func normalizeContractCode(contractCode string, legacyCode string) string {
	if strings.TrimSpace(contractCode) != "" {
		return strings.TrimSpace(contractCode)
	}
	return strings.TrimSpace(legacyCode)
}

func normalizeOpenID(openID string, legacyOpenID string) string {
	if strings.TrimSpace(openID) != "" {
		return strings.TrimSpace(openID)
	}
	return strings.TrimSpace(legacyOpenID)
}

func (a *wechat) ContractSign(c *gin.Context) {
	start := time.Now()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		LogError(c, "ContractSign:读取请求体", err)
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}

	var req model.SignContractRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		LogError(c, "ContractSign:XML解析", err, zap.String("raw_body", string(body)))
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	LogRequest(c, "ContractSign", gin.H{
		"raw_body":  string(body),
		"parsed":    req,
		"client_ip": clientIP(c),
	})

	LogServiceCall(c, "Merchant", "GetMerchantByMchID", zap.String("mch_id", req.MchID))
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		LogError(c, "ContractSign:获取商户", err)
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "ContractSign", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	contractCode := normalizeContractCode(req.ContractCode, req.OutContractCode)
	openID := normalizeOpenID(req.OpenID, req.OutUserID)

	params := map[string]string{
		"appid": req.AppID, "mch_id": req.MchID, "plan_id": req.PlanID, "contract_code": contractCode,
		"openid": openID, "contract_display_account": req.ContractDisplayAccount, "notify_url": req.NotifyURL,
		"sign_type": req.SignType, "version": req.Version, "timestamp": req.TimeStamp, "nonce": req.Nonce, "sign": req.Sign,
	}
	LogServiceCall(c, "Signature", "VerifyIfNeeded", zap.Bool("verify_sign", merchant.VerifySign))
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		LogError(c, "ContractSign:签名校验", err)
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "ContractSign", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	LogServiceCall(c, "Contract", "HasActiveContract", zap.String("contract_code", contractCode))
	LogServiceCall(c, "Contract", "HasActiveContractByUser", zap.Any("merchant_id", merchant.ID), zap.String("open_id", openID))
	if serviceInfo.Contract.HasActiveContract(contractCode) || serviceInfo.Contract.HasActiveContractByUser(merchant.ID, openID, openID) {
		LogError(c, "ContractSign:重复签约", nil, zap.String("contract_code", contractCode))
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "重复签约", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignExists, ErrCodeDes: "已有有效签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "ContractSign", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	contract := model.Contract{
		MerchantID:    merchant.ID,
		OpenID:        openID,
		OutUserID:     openID,
		OutContractID: contractCode,
		PlanID:        req.PlanID,
		NotifyURL:     req.NotifyURL,
		RequestData:   string(body),
	}
	statusRecord := model.ContractStatusRecord{
		MerchantID:     merchant.ID,
		OutContractID:  contractCode,
		ContractStatus: model.ContractStatusPending,
		IsFirstDeduct:  true,
	}
	LogServiceCall(c, "Contract", "CreateContractWithStatus", zap.Any("contract", contract))
	if err = serviceInfo.Contract.CreateContractWithStatus(&contract, &statusRecord); err != nil {
		LogError(c, "ContractSign:创建签约", err)
		resp := model.SignContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "创建签约失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeFail, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "ContractSign", string(xml), start)
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
	LogServiceCall(c, "Deduct", "CreateContractRecord", zap.Any("record", record))
	_ = serviceInfo.Deduct.CreateContractRecord(&record)

	if merchant.SignStatusDelay > 0 {
		time.Sleep(time.Duration(merchant.SignStatusDelay) * time.Second)
	}
	LogServiceCall(c, "Contract", "UpdateContractStatus", zap.Any("id", contract.ID), zap.String("status", merchant.SignTargetStatus))
	_ = serviceInfo.Contract.UpdateContractStatus(contract.ID, merchant.SignTargetStatus, "")

	contractID := fmt.Sprintf("MOCK-C-%d", contract.ID)
	signSerialNo := fmt.Sprintf("MOCK-S-%d", time.Now().UnixNano())
	if merchant.SignTargetStatus == model.ContractStatusActive {
		LogServiceCall(c, "Contract", "SetContractID", zap.Any("id", contract.ID), zap.String("contract_id", contractID))
		_ = serviceInfo.Contract.SetContractID(contract.ID, contractID, signSerialNo)
		LogServiceCall(c, "Contract", "SetExpireTime", zap.Any("id", contract.ID))
		_ = serviceInfo.Contract.SetExpireTime(contract.ID, merchant.SignDurationMinutes)
	}

	preEntrustwebID := model.RandomMixed(27)
	miniprogramPath := fmt.Sprintf("pages/index/index?sign_scene=app&domain_type=cn&pre_entrustweb_id=%s", preEntrustwebID)

	resp := model.SignContractResponseV2{
		ReturnCode:          model.ErrCodeSuccess,
		ReturnMsg:           "OK",
		ResultCode:          model.ErrCodeSuccess,
		AppID:               req.AppID,
		MchID:               merchant.MchID,
		MiniprogramUsername: merchant.MiniprogramUsername,
		MiniprogramPath:     miniprogramPath,
		NonceStr:            req.Nonce,
		PreEntrustwebID:     preEntrustwebID,
	}
	resp.Sign = serviceInfo.Signature.Sign(map[string]string{
		"return_code":          model.ErrCodeSuccess,
		"result_code":          model.ErrCodeSuccess,
		"appid":                req.AppID,
		"mch_id":               merchant.MchID,
		"miniprogram_username": merchant.MiniprogramUsername,
		"miniprogram_path":     miniprogramPath,
		"nonce_str":            req.Nonce,
		"pre_entrustweb_id":    preEntrustwebID,
	}, merchant.SignKey)
	xmlRespBytes, _ := resp.ToXMLBytes()
	xmlResp := string(xmlRespBytes)
	LogServiceCall(c, "Deduct", "UpdateContractRecordResponse", zap.Any("id", record.ID))
	_ = serviceInfo.Deduct.UpdateContractRecordResponse(record.ID, xmlResp, merchant.SignTargetStatus)
	go func() {
		if merchant.SignCallbackEnabled {
			if merchant.SignCallbackDelay > 0 {
				time.Sleep(time.Duration(merchant.SignCallbackDelay) * time.Second)
			}
			contract.ContractID = contractID
			callbackXML := serviceInfo.Callback.BuildContractCallbackXML(contract, merchant.MchID, merchant.SignTargetStatus, merchant.SignKey)
			LogServiceCall(c, "Callback", "DoXMLCallback", zap.String("url", req.NotifyURL))
			result, callbackErr := serviceInfo.Callback.DoXMLCallback(req.NotifyURL, callbackXML)
			if callbackErr != nil {
				LogError(c, "ContractSign:异步回调", callbackErr, zap.String("result", result))
				result = callbackErr.Error() + "; " + result
			}
			LogServiceCall(c, "Deduct", "UpdateContractRecordStatus", zap.Any("id", record.ID), zap.String("status", merchant.SignTargetStatus))
			_ = serviceInfo.Deduct.UpdateContractRecordStatus(record.ID, merchant.SignTargetStatus, "", result)
		}
	}()
	LogResponse(c, "ContractSign", string(xmlResp), start)
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) QueryContract(c *gin.Context) {
	start := time.Now()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		LogError(c, "QueryContract:读取请求体", err)
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.QueryContractRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		LogError(c, "QueryContract:XML解析", err, zap.String("raw_body", string(body)))
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	LogRequest(c, "QueryContract", gin.H{
		"raw_body": string(body),
		"parsed":   req,
	})

	LogServiceCall(c, "Merchant", "GetMerchantByMchID", zap.String("mch_id", req.MchID))
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		LogError(c, "QueryContract:获取商户", err)
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	contractCode := normalizeContractCode(req.ContractCode, req.OutContractCode)
	params := map[string]string{
		"appid":         req.AppID,
		"mch_id":        req.MchID,
		"contract_id":   req.ContractID,
		"plan_id":       req.PlanID,
		"contract_code": contractCode,
		"sign_type":     req.SignType,
		"timestamp":     req.TimeStamp,
		"nonce":         req.Nonce,
		"sign":          req.Sign,
	}
	LogServiceCall(c, "Signature", "VerifyIfNeeded", zap.Bool("verify_sign", merchant.VerifySign))
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		LogError(c, "QueryContract:签名校验", err)
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var contract model.Contract
	if strings.TrimSpace(req.ContractID) != "" {
		LogServiceCall(c, "Deduct", "GetContractByContractIDFromDB", zap.String("contract_id", req.ContractID))
		contract, err = serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	}
	if err != nil && strings.TrimSpace(contractCode) != "" {
		LogServiceCall(c, "Deduct", "GetContractFromDB", zap.String("contract_code", contractCode))
		contract, err = serviceInfo.Deduct.GetContractFromDB(contractCode)
	}
	if err != nil {
		LogError(c, "QueryContract:获取签约", err)
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	LogServiceCall(c, "Contract", "GetContractStatusByContractID", zap.Any("contract_id", contract.ID))
	statusRecord, statusErr := serviceInfo.Contract.GetContractStatusByContractID(contract.ID)
	if statusErr != nil {
		LogError(c, "QueryContract:获取签约状态", statusErr)
		resp := model.QueryContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约状态不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约状态"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryContract", string(xml), start)
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
	LogResponse(c, "QueryContract", string(xmlResp), start)
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) TerminateContract(c *gin.Context) {
	start := time.Now()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		LogError(c, "TerminateContract:读取请求体", err)
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.TerminateContractRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		LogError(c, "TerminateContract:XML解析", err, zap.String("raw_body", string(body)))
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	LogRequest(c, "TerminateContract", gin.H{
		"raw_body": string(body),
		"parsed":   req,
	})

	LogServiceCall(c, "Merchant", "GetMerchantByMchID", zap.String("mch_id", req.MchID))
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		LogError(c, "TerminateContract:获取商户", err)
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "TerminateContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}
	contractCode := normalizeContractCode(req.ContractCode, req.OutContractCode)
	params := map[string]string{
		"appid":                       req.AppID,
		"mch_id":                      req.MchID,
		"contract_id":                 req.ContractID,
		"plan_id":                     req.PlanID,
		"contract_code":               contractCode,
		"contract_termination_remark": req.ContractTerminationRemark,
		"version":                     req.Version,
		"sign_type":                   req.SignType,
		"timestamp":                   req.TimeStamp,
		"nonce":                       req.Nonce,
		"sign":                        req.Sign,
	}
	LogServiceCall(c, "Signature", "VerifyIfNeeded", zap.Bool("verify_sign", merchant.VerifySign))
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		LogError(c, "TerminateContract:签名校验", err)
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "TerminateContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var contract model.Contract
	if strings.TrimSpace(req.ContractID) != "" {
		LogServiceCall(c, "Deduct", "GetContractByContractIDFromDB", zap.String("contract_id", req.ContractID))
		contract, err = serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	}
	if err != nil && strings.TrimSpace(contractCode) != "" {
		LogServiceCall(c, "Deduct", "GetContractFromDB", zap.String("contract_code", contractCode))
		contract, err = serviceInfo.Deduct.GetContractFromDB(contractCode)
	}
	if err != nil {
		LogError(c, "TerminateContract:获取签约", err)
		resp := model.TerminateContractResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "TerminateContract", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	terminateType := strings.TrimSpace(req.ContractEndingType)
	if terminateType == "" {
		terminateType = model.TerminateTypeMerchantRequest
	}
	terminateStatus := strings.TrimSpace(merchant.TerminateTargetStatus)
	if terminateStatus == "" {
		terminateStatus = model.ContractStatusTerminated
	}
	terminateCallbackEnabled := merchant.TerminateCallbackEnabled || merchant.TerminateNotifyEnabled

	record := model.ContractRecord{
		ContractID:    contract.ID,
		MerchantID:    merchant.ID,
		OperationType: "terminate",
		RequestXML:    string(body),
		CallbackURL:   contract.NotifyURL,
		Status:        model.ContractStatusPending,
	}
	LogServiceCall(c, "Deduct", "CreateContractRecord", zap.Any("record", record))
	_ = serviceInfo.Deduct.CreateContractRecord(&record)

	if merchant.TerminateStatusDelay > 0 {
		time.Sleep(time.Duration(merchant.TerminateStatusDelay) * time.Second)
	}
	LogServiceCall(c, "Contract", "UpdateContractStatus", zap.Any("id", contract.ID), zap.String("status", terminateStatus))
	_ = serviceInfo.Contract.UpdateContractStatus(contract.ID, terminateStatus, terminateType)
	LogServiceCall(c, "Deduct", "SetContractStatus", zap.Any("id", contract.ID), zap.String("status", terminateStatus))
	_ = serviceInfo.Deduct.SetContractStatus(contract.ID, terminateStatus, terminateType)

	resp := model.TerminateContractResponse{
		ReturnCode:     model.ErrCodeSuccess,
		ReturnMsg:      "OK",
		ResultCode:     model.ErrCodeSuccess,
		ContractID:     contract.ContractID,
		ContractStatus: terminateStatus,
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
	LogServiceCall(c, "Deduct", "UpdateContractRecordResponse", zap.Any("id", record.ID))
	_ = serviceInfo.Deduct.UpdateContractRecordResponse(record.ID, xmlResp, terminateStatus)

	go func() {
		if terminateCallbackEnabled {
			if merchant.TerminateCallbackDelay > 0 {
				time.Sleep(time.Duration(merchant.TerminateCallbackDelay) * time.Second)
			}
			callbackXML := serviceInfo.Callback.BuildContractCallbackXML(contract, merchant.MchID, terminateStatus, merchant.SignKey)
			LogServiceCall(c, "Callback", "DoXMLCallback", zap.String("url", contract.NotifyURL))
			result, callbackErr := serviceInfo.Callback.DoXMLCallback(contract.NotifyURL, callbackXML)
			if callbackErr != nil {
				LogError(c, "TerminateContract:异步回调", callbackErr, zap.String("result", result))
				result = callbackErr.Error() + "; " + result
			}
			LogServiceCall(c, "Deduct", "UpdateContractRecordStatus", zap.Any("id", record.ID))
			_ = serviceInfo.Deduct.UpdateContractRecordStatus(record.ID, terminateStatus, "", result)
			LogServiceCall(c, "Deduct", "SetContractRecordCallbackResult", zap.Any("id", record.ID))
			_ = serviceInfo.Deduct.SetContractRecordCallbackResult(record.ID, result, time.Now().Unix())
		}
	}()
	LogResponse(c, "TerminateContract", string(xmlResp), start)
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) ApplyDeduct(c *gin.Context) {
	start := time.Now()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		LogError(c, "ApplyDeduct:读取请求体", err)
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}

	var req model.DeductApplyRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		LogError(c, "ApplyDeduct:XML解析", err, zap.String("raw_body", string(body)))
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	LogRequest(c, "ApplyDeduct", gin.H{
		"raw_body": string(body),
		"parsed":   req,
	})

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

	LogServiceCall(c, "Merchant", "GetMerchantByMchID", zap.String("mch_id", req.MchID))
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		LogError(c, "ApplyDeduct:获取商户", err)
		resp := buildResp("商户不存在", model.ErrCodeFail, model.ErrCodeInvalidParams, "商户配置不存在", req.MchID, req.OutTradeNo, req.TransactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		LogResponse(c, "ApplyDeduct", xml, start)
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
	LogServiceCall(c, "Signature", "VerifyIfNeeded", zap.Bool("verify_sign", merchant.VerifySign))
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		LogError(c, "ApplyDeduct:签名校验", err)
		resp := buildResp("签名校验失败", model.ErrCodeFail, model.ErrCodeInvalidSign, err.Error(), merchant.MchID, req.OutTradeNo, req.TransactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		LogResponse(c, "ApplyDeduct", xml, start)
		return
	}

	LogServiceCall(c, "Deduct", "GetContractByContractIDFromDB", zap.String("contract_id", req.ContractID))
	contract, err := serviceInfo.Deduct.GetContractByContractIDFromDB(req.ContractID)
	if err != nil {
		LogError(c, "ApplyDeduct:获取签约", err)
		resp := buildResp("签约不存在", model.ErrCodeFail, model.ErrCodeSignNotFound, "未找到签约关系", merchant.MchID, req.OutTradeNo, req.TransactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		LogResponse(c, "ApplyDeduct", xml, start)
		return
	}

	transactionID := strings.TrimSpace(req.TransactionID)
	if transactionID == "" {
		transactionID = fmt.Sprintf("MOCK-T-%d", time.Now().UnixNano())
	}

	LogServiceCall(c, "Contract", "GetContractStatusByContractID", zap.Any("contract_id", contract.ID))
	statusRecord, statusErr := serviceInfo.Contract.GetContractStatusByContractID(contract.ID)
	if statusErr != nil {
		LogError(c, "ApplyDeduct:获取签约状态", statusErr)
		record := model.DeductRecord{
			ContractID:      contract.ID,
			MerchantID:      merchant.ID,
			OperationType:   model.DeductOperationTypeDeduct,
			RequestData:     string(body),
			CallbackURL:     req.NotifyURL,
			TransactionID:   transactionID,
			OutTradeNo:      req.OutTradeNo,
			Amount:          effectiveAmount,
			Status:          model.DeductStatusFailed,
			IsFirstDeduct:   false,
			PreNotifyCalled: false,
			ErrorCode:       model.ErrCodeSignNotFound,
			ErrorMessage:    "订阅状态不存在",
		}
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		resp := buildResp("订阅信息不存在", model.ErrCodeFail, model.ErrCodeSignNotFound, "未找到订阅状态", merchant.MchID, req.OutTradeNo, transactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xml, record.Status)
		LogResponse(c, "ApplyDeduct", xml, start)
		return
	}

	callbackTarget := strings.TrimSpace(req.NotifyURL)
	if callbackTarget == "" {
		callbackTarget = strings.TrimSpace(contract.NotifyURL)
	}

	var record model.DeductRecord
	usingExistingPreNotify := false
	if !statusRecord.IsFirstDeduct {
		LogServiceCall(c, "Deduct", "GetLatestPendingPreNotifyRecord", zap.Any("contract_id", contract.ID))
		record, err = serviceInfo.Deduct.GetLatestPendingPreNotifyRecord(contract.ID)
		if err == nil {
			usingExistingPreNotify = true
		}
	}

	if !usingExistingPreNotify {
		record = model.DeductRecord{
			ContractID:      contract.ID,
			MerchantID:      merchant.ID,
			OperationType:   model.DeductOperationTypeDeduct,
			RequestData:     string(body),
			CallbackURL:     callbackTarget,
			TransactionID:   transactionID,
			OutTradeNo:      req.OutTradeNo,
			Amount:          effectiveAmount,
			Status:          model.DeductStatusPending,
			IsFirstDeduct:   statusRecord.IsFirstDeduct,
			PreNotifyCalled: statusRecord.PreNotifyCalled,
		}
	}

	if statusRecord.ContractStatus != model.ContractStatusActive {
		record.Status = model.DeductStatusFailed
		record.ErrorCode = model.ErrCodeDeductNotAllowed
		record.ErrorMessage = "签约未生效或已解约"
		LogError(c, "ApplyDeduct:签约状态不可扣款", nil, zap.Any("contract_id", contract.ID), zap.String("status", statusRecord.ContractStatus))
		if usingExistingPreNotify {
			_ = serviceInfo.Deduct.UpdatePreNotifyRecordResponse(record.ID, "", record.Status, record.ErrorCode, record.ErrorMessage)
		} else {
			_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		}
		resp := buildResp("签约状态不可扣款", model.ErrCodeFail, model.ErrCodeDeductNotAllowed, "签约未生效或已解约", merchant.MchID, req.OutTradeNo, transactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xml, record.Status)
		LogResponse(c, "ApplyDeduct", xml, start)
		return
	}

	if !statusRecord.IsFirstDeduct && merchant.StrictDeductRule && !usingExistingPreNotify {
		record.Status = model.DeductStatusFailed
		record.ErrorCode = model.ErrCodePreNotifyRequired
		record.ErrorMessage = "非首次扣款前必须先调用预扣费通知API"
		LogError(c, "ApplyDeduct:非首次扣款前必须先调用预扣费通知", nil, zap.Any("contract_id", contract.ID))
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
		resp := buildResp("未先调用预扣费通知", model.ErrCodeFail, model.ErrCodePreNotifyRequired, "非首次扣款前必须先调用预扣费通知API", merchant.MchID, req.OutTradeNo, transactionID, effectiveAmount, req.SignType, effectiveNonce)
		xml := writeXML(resp)
		_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xml, record.Status)
		LogResponse(c, "ApplyDeduct", xml, start)
		return
	}

	if usingExistingPreNotify {
		LogServiceCall(c, "Deduct", "ConsumePreNotifyRecord", zap.Any("id", record.ID))
		_ = serviceInfo.Deduct.ConsumePreNotifyRecord(record.ID, req.OutTradeNo, transactionID, string(body), callbackTarget, "", effectiveAmount, statusRecord.IsFirstDeduct)
		record.OperationType = model.DeductOperationTypeDeduct
		record.OutTradeNo = req.OutTradeNo
		record.TransactionID = transactionID
		record.RequestData = string(body)
		record.CallbackURL = callbackTarget
		record.Amount = effectiveAmount
		record.Status = model.DeductStatusPending
		record.IsFirstDeduct = statusRecord.IsFirstDeduct
		record.PreNotifyCalled = true
	} else {
		_ = serviceInfo.Deduct.SaveDeductRecord(&record)
	}

	if merchant.DeductStatusDelay > 0 {
		time.Sleep(time.Duration(merchant.DeductStatusDelay) * time.Second)
	}

	finalStatus := strings.TrimSpace(merchant.DeductTargetStatus)
	if finalStatus == "" {
		finalStatus = model.DeductStatusSuccess
	}
	LogServiceCall(c, "Deduct", "UpdateDeductRecordStatus", zap.Any("id", record.ID), zap.String("status", finalStatus))
	_ = serviceInfo.Deduct.UpdateDeductRecordStatus(record.ID, finalStatus, "", "")
	record.Status = finalStatus

	if record.IsFirstDeduct {
		LogServiceCall(c, "Contract", "MarkFirstDeductDone", zap.Any("contract_id", contract.ID))
		_ = serviceInfo.Contract.MarkFirstDeductDone(contract.ID)
	} else {
		LogServiceCall(c, "Contract", "ClearPreNotify", zap.Any("contract_id", contract.ID))
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
	LogServiceCall(c, "Deduct", "UpdateDeductRecordResponse", zap.Any("id", record.ID))
	_ = serviceInfo.Deduct.UpdateDeductRecordResponse(record.ID, xmlResp, finalStatus)

	go func() {
		if merchant.DeductCallbackEnabled {
			if merchant.DeductCallbackDelay > 0 {
				time.Sleep(time.Duration(merchant.DeductCallbackDelay) * time.Second)
			}
			callbackXML := serviceInfo.Callback.BuildDeductCallbackXML(merchant, contract, record, req.SignType)
			LogServiceCall(c, "Callback", "DoXMLCallback", zap.String("url", callbackTarget))
			result, callbackErr := serviceInfo.Callback.DoXMLCallback(callbackTarget, callbackXML)
			if callbackErr != nil {
				LogError(c, "ApplyDeduct:异步回调", callbackErr, zap.String("result", result))
				result = callbackErr.Error() + "; " + result
			}
			LogServiceCall(c, "Deduct", "SetCallbackResult", zap.Any("id", record.ID))
			_ = serviceInfo.Deduct.SetCallbackResult(record.ID, result, time.Now().Unix())
		}
	}()
	LogResponse(c, "ApplyDeduct", string(xmlResp), start)
}

func (a *wechat) QueryDeduct(c *gin.Context) {
	start := time.Now()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		LogError(c, "QueryDeduct:读取请求体", err)
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>读取请求失败</return_msg></xml>")
		return
	}
	var req model.QueryDeductRequest
	if err = serviceInfo.XMLCodec.Unmarshal(body, &req); err != nil {
		LogError(c, "QueryDeduct:XML解析", err, zap.String("raw_body", string(body)))
		c.String(200, "<xml><return_code>FAIL</return_code><return_msg>XML解析失败</return_msg></xml>")
		return
	}
	LogRequest(c, "QueryDeduct", gin.H{
		"raw_body": string(body),
		"parsed":   req,
	})

	LogServiceCall(c, "Merchant", "GetMerchantByMchID", zap.String("mch_id", req.MchID))
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		LogError(c, "QueryDeduct:获取商户", err)
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryDeduct", string(xml), start)
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
	LogServiceCall(c, "Signature", "VerifyIfNeeded", zap.Bool("verify_sign", merchant.VerifySign))
	if err = serviceInfo.Signature.VerifyIfNeeded(merchant.VerifySign, params, merchant.SignKey); err != nil {
		LogError(c, "QueryDeduct:签名校验", err)
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签名校验失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidSign, ErrCodeDes: err.Error()}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryDeduct", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	var record model.DeductRecord
	if strings.TrimSpace(req.OutTradeNo) != "" {
		LogServiceCall(c, "Deduct", "GetDeductRecordByOutTradeNo", zap.String("out_trade_no", req.OutTradeNo))
		record, err = serviceInfo.Deduct.GetDeductRecordByOutTradeNo(req.OutTradeNo)
	}
	if err != nil && strings.TrimSpace(req.TransactionID) != "" {
		LogServiceCall(c, "Deduct", "GetDeductRecordByTransactionID", zap.String("transaction_id", req.TransactionID))
		record, err = serviceInfo.Deduct.GetDeductRecordByTransactionID(req.TransactionID)
	}
	if err != nil {
		LogError(c, "QueryDeduct:获取扣款记录", err)
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "订单不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeFail, ErrCodeDes: "未找到扣款记录"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryDeduct", string(xml), start)
		c.Data(200, "application/xml; charset=utf-8", []byte(xml))
		return
	}

	LogServiceCall(c, "Deduct", "GetContractByID", zap.Any("id", record.ContractID))
	contract, contractErr := serviceInfo.Deduct.GetContractByID(record.ContractID)
	if contractErr != nil {
		LogError(c, "QueryDeduct:获取签约", contractErr)
		resp := model.QueryDeductResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"}
		xml, _ := serviceInfo.XMLCodec.Marshal(resp)
		LogResponse(c, "QueryDeduct", string(xml), start)
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
	LogResponse(c, "QueryDeduct", string(xmlResp), start)
	c.Data(200, "application/xml; charset=utf-8", []byte(xmlResp))
}

func (a *wechat) PreDeductNotify(c *gin.Context) {
	start := time.Now()
	var req model.PreDeductNotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		LogError(c, "PreDeductNotify:参数绑定", err)
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: err.Error(), ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: err.Error()})
		return
	}

	contractID := strings.TrimSpace(c.Param("contract_id"))
	if contractID == "" {
		LogError(c, "PreDeductNotify:参数校验", nil, zap.String("reason", "contract_id不能为空"))
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "contract_id不能为空", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "contract_id不能为空"})
		return
	}
	if strings.TrimSpace(req.MchID) == "" || strings.TrimSpace(req.AppID) == "" {
		LogError(c, "PreDeductNotify:参数校验", nil, zap.String("reason", "mchid或appid不能为空"))
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "mchid或appid不能为空", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "mchid或appid不能为空"})
		return
	}
	LogRequest(c, "PreDeductNotify", gin.H{
		"contract_id": contractID,
		"request":     req,
	})

	if req.EstimatedAmount.Amount <= 0 {
		LogError(c, "PreDeductNotify:参数校验", nil, zap.String("reason", "预计扣费金额必须大于0"))
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "预计扣费金额必须大于0", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "estimated_amount.amount必须大于0"})
		return
	}
	if req.DeductDuration.Count <= 0 {
		req.DeductDuration.Count = 1
	}
	if strings.TrimSpace(req.DeductDuration.Unit) == "" {
		req.DeductDuration.Unit = "DAY"
	}
	if strings.TrimSpace(req.EstimatedAmount.Currency) == "" {
		req.EstimatedAmount.Currency = "CNY"
	}

	LogServiceCall(c, "Merchant", "GetMerchantByMchID", zap.String("mch_id", req.MchID))
	merchant, err := serviceInfo.Merchant.GetMerchantByMchID(req.MchID)
	if err != nil {
		LogError(c, "PreDeductNotify:获取商户", err)
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "商户不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeInvalidParams, ErrCodeDes: "商户配置不存在"})
		return
	}

	LogServiceCall(c, "Deduct", "GetContractByContractIDFromDB", zap.String("contract_id", contractID))
	contract, err := serviceInfo.Deduct.GetContractByContractIDFromDB(contractID)
	if err != nil {
		LogError(c, "PreDeductNotify:获取签约", err)
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到签约关系"})
		return
	}

	LogServiceCall(c, "Contract", "GetContractStatusByContractID", zap.Any("contract_id", contract.ID))
	statusRecord, statusErr := serviceInfo.Contract.GetContractStatusByContractID(contract.ID)
	if statusErr != nil {
		LogError(c, "PreDeductNotify:获取签约状态", statusErr)
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "订阅状态不存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeSignNotFound, ErrCodeDes: "未找到订阅状态"})
		return
	}
	if statusRecord.ContractStatus != model.ContractStatusActive {
		LogError(c, "PreDeductNotify:签约状态不可预扣费", nil, zap.Any("contract_id", contract.ID), zap.String("status", statusRecord.ContractStatus))
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "签约状态不可预扣费", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeDeductNotAllowed, ErrCodeDes: "签约未生效或已解约"})
		return
	}

	if existing, existingErr := serviceInfo.Deduct.GetLatestPendingPreNotifyRecord(contract.ID); existingErr == nil && existing.ID != 0 {
		LogError(c, "PreDeductNotify:已存在预扣费通知", nil, zap.Any("contract_id", contract.ID), zap.Any("existing_id", existing.ID))
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "预扣费通知已存在", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeAlreadyExists, ErrCodeDes: "已经成功发送通知，无需重复调用"})
		return
	}

	now := time.Now()
	if merchant.StrictDeductRule {
		loc := time.FixedZone("CST", 8*3600)
		beijingHour := now.In(loc).Hour()
		if beijingHour < 7 || beijingHour >= 22 {
			LogError(c, "PreDeductNotify:不在可通知时间段", nil, zap.Int("beijing_hour", beijingHour))
			c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "不在可通知时间段", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeDeductNotAllowed, ErrCodeDes: "预扣费通知只允许在北京时间7:00-22:00调用"})
			return
		}
	}

	requestData, _ := json.Marshal(map[string]any{
		"contract_id":      contractID,
		"mchid":            req.MchID,
		"appid":            req.AppID,
		"deduct_duration":  req.DeductDuration,
		"estimated_amount": req.EstimatedAmount,
	})
	record := model.DeductRecord{
		ContractID:      contract.ID,
		MerchantID:      merchant.ID,
		OperationType:   model.DeductOperationTypePreNotify,
		RequestData:     string(requestData),
		CallbackURL:     strings.TrimSpace(contract.NotifyURL),
		Amount:          req.EstimatedAmount.Amount,
		Status:          model.DeductStatusWaitDeduct,
		IsFirstDeduct:   false,
		PreNotifyCalled: true,
	}
	LogServiceCall(c, "Deduct", "SaveDeductRecord", zap.Any("record", record))
	if err = serviceInfo.Deduct.SaveDeductRecord(&record); err != nil {
		LogError(c, "PreDeductNotify:保存预扣费记录", err)
		c.JSON(200, model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeFail, ReturnMsg: "预扣费记录保存失败", ResultCode: model.ErrCodeFail, ErrCode: model.ErrCodeFail, ErrCodeDes: err.Error()})
		return
	}
	LogServiceCall(c, "Contract", "MarkPreNotifyCalled", zap.Any("contract_id", contract.ID))
	_ = serviceInfo.Contract.MarkPreNotifyCalled(contract.ID)
	responseData := `{"http_status":204}`
	_ = serviceInfo.Deduct.UpdatePreNotifyRecordResponse(record.ID, responseData, model.DeductStatusWaitDeduct, "", "")
	resp := model.PreDeductNotifyResponse{ReturnCode: model.ErrCodeSuccess, ReturnMsg: "OK", ResultCode: model.ErrCodeSuccess}
	LogResponse(c, "PreDeductNotify", resp, start)
	c.JSON(200, resp)
}
