package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	"go.uber.org/zap"
)

type callback struct{}

func (s *callback) DoXMLCallback(url string, body string) (string, error) {
	if strings.TrimSpace(url) == "" {
		return "", nil
	}
	global.GVA_LOG.Info("DoXMLCallback", zap.String("url", url))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/xml")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		global.GVA_LOG.Error("DoXMLCallback", zap.String("url", url))
		return string(respBody), fmt.Errorf("回调失败，状态码: %d", resp.StatusCode)
	}
	return string(respBody), nil
}

func (s *callback) BuildContractCallbackXML(contract model.Contract, mchID string, status string, signKey string) string {
	changeType := "MODIFY"
	if status == model.ContractStatusActive {
		changeType = "ADD"
	}
	if status == model.ContractStatusTerminated {
		changeType = "DELETE"
	}
	operateTime := time.Now().Format("2006-01-02 15:04:05")
	notify := model.ContractResultNotify{
		ReturnCode:      model.ErrCodeSuccess,
		ResultCode:      model.ErrCodeSuccess,
		MchID:           mchID,
		OutContractCode: contract.OutContractID,
		OpenID:          contract.OpenID,
		PlanID:          contract.PlanID,
		ChangeType:      changeType,
		OperateTime:     operateTime,
		ContractID:      contract.ContractID,
	}
	notify.Sign = Service.Signature.Sign(map[string]string{
		"return_code":   notify.ReturnCode,
		"result_code":   notify.ResultCode,
		"mch_id":        notify.MchID,
		"contract_code": notify.OutContractCode,
		"openid":        notify.OpenID,
		"plan_id":       notify.PlanID,
		"change_type":   notify.ChangeType,
		"operate_time":  notify.OperateTime,
		"contract_id":   notify.ContractID,
	}, signKey)
	xml, _ := Service.XMLCodec.Marshal(notify)
	return xml
}

func (s *callback) BuildDeductCallbackXML(merchant model.Merchant, contract model.Contract, record model.DeductRecord, signType string) string {
	tradeState := strings.TrimSpace(record.Status)
	if tradeState == "" {
		tradeState = model.DeductStatusPending
	}
	timeEnd := time.Now().Format("20060102150405")
	if tradeState == model.DeductStatusPending {
		timeEnd = ""
	}
	nonce := fmt.Sprintf("mock-deduct-%d", time.Now().UnixNano())
	notify := model.DeductNotifyResponse{
		ReturnCode:    model.ErrCodeSuccess,
		ReturnMsg:     "OK",
		AppID:         merchant.AppID,
		MchID:         merchant.MchID,
		OutTradeNo:    record.OutTradeNo,
		TransactionID: record.TransactionID,
		TradeType:     "PAP",
		TradeState:    tradeState,
		BankType:      "MOCK",
		TotalAmount:   record.Amount,
		CashAmount:    record.Amount,
		TimeEnd:       timeEnd,
		SignType:      signType,
		TimeStamp:     fmt.Sprintf("%d", time.Now().Unix()),
		Nonce:         nonce,
	}
	params := map[string]string{
		"return_code":    notify.ReturnCode,
		"appid":          notify.AppID,
		"mch_id":         notify.MchID,
		"out_trade_no":   notify.OutTradeNo,
		"transaction_id": notify.TransactionID,
		"trade_type":     notify.TradeType,
		"trade_state":    notify.TradeState,
		"bank_type":      notify.BankType,
		"total_amount":   fmt.Sprintf("%d", notify.TotalAmount),
		"cash_amount":    fmt.Sprintf("%d", notify.CashAmount),
		"time_end":       notify.TimeEnd,
		"timestamp":      notify.TimeStamp,
		"nonce":          notify.Nonce,
	}
	if strings.TrimSpace(signType) != "" {
		params["sign_type"] = signType
	}
	notify.Sign = Service.Signature.Sign(params, merchant.SignKey)
	xml, _ := Service.XMLCodec.Marshal(notify)
	return xml
}
