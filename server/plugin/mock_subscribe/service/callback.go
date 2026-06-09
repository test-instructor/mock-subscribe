package service

import (
	"bytes"
	"errors"
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
		return string(respBody), errors.New(fmt.Sprintf("回调失败，状态码: %d", resp.StatusCode))
	}
	return string(respBody), nil
}

func (s *callback) BuildContractCallbackXML(contract model.Contract, status string) string {
	return fmt.Sprintf("<xml><contract_id>%s</contract_id><out_contract_code>%s</out_contract_code><contract_status>%s</contract_status></xml>", contract.ContractID, contract.OutContractID, status)
}

func (s *callback) BuildDeductCallbackXML(record model.DeductRecord) string {
	return fmt.Sprintf("<xml><transaction_id>%s</transaction_id><status>%s</status><amount>%d</amount></xml>", record.TransactionID, record.Status, record.Amount)
}
