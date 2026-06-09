package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
)

type callback struct{}

func (s *callback) DoXMLCallback(url string, body string) (string, error) {
	if strings.TrimSpace(url) == "" {
		return "", nil
	}
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
		ReturnCode:   model.ErrCodeSuccess,
		ResultCode:   model.ErrCodeSuccess,
		MchID:        mchID,
		ContractCode: contract.OutContractID,
		OpenID:       contract.OpenID,
		PlanID:       contract.PlanID,
		ChangeType:   changeType,
		OperateTime:  operateTime,
		ContractID:   contract.ContractID,
	}
	notify.Sign = Service.Signature.Sign(map[string]string{
		"return_code":   notify.ReturnCode,
		"result_code":   notify.ResultCode,
		"mch_id":        notify.MchID,
		"contract_code": notify.ContractCode,
		"openid":        notify.OpenID,
		"plan_id":       notify.PlanID,
		"change_type":   notify.ChangeType,
		"operate_time":  notify.OperateTime,
		"contract_id":   notify.ContractID,
	}, signKey)
	xml, _ := Service.XMLCodec.Marshal(notify)
	return xml
}

func (s *callback) BuildDeductCallbackXML(record model.DeductRecord) string {
	return fmt.Sprintf("<xml><transaction_id>%s</transaction_id><status>%s</status><amount>%d</amount></xml>", record.TransactionID, record.Status, record.Amount)
}
