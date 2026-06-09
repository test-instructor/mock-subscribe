package service

import (
	"errors"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
)

type deduct struct{}

func (s *deduct) CreateDeductRecord(record *model.DeductRecord) error {
	if record.ContractID == 0 {
		return errors.New("签约ID不能为空")
	}
	record.Status = model.DeductStatusPending
	return global.GVA_DB.Create(record).Error
}

func (s *deduct) UpdateDeductRecord(record *model.DeductRecord) error {
	if record.ID == 0 {
		return errors.New("扣款记录ID不能为空")
	}
	return global.GVA_DB.Model(&model.DeductRecord{}).Where("id = ?", record.ID).Updates(record).Error
}

func (s *deduct) GetDeductRecord(id uint) (model.DeductRecord, error) {
	var r model.DeductRecord
	err := global.GVA_DB.Where("id = ?", id).First(&r).Error
	return r, err
}

func (s *deduct) GetDeductRecordByContractAndTradeNo(contractID uint, outTradeNo string) (model.DeductRecord, error) {
	var r model.DeductRecord
	err := global.GVA_DB.Where("contract_id = ? AND out_trade_no = ?", contractID, outTradeNo).First(&r).Error
	return r, err
}

func (s *deduct) GetDeductRecordList(info mockReq.DeductRecordSearch) ([]model.DeductRecord, int64, error) {
	var list []model.DeductRecord
	var total int64

	db := global.GVA_DB.Model(&model.DeductRecord{})
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.MerchantID != 0 {
		db = db.Where("merchant_id = ?", info.MerchantID)
	}
	if info.ContractID != 0 {
		db = db.Where("contract_id = ?", info.ContractID)
	}
	if info.OperationType != "" {
		db = db.Where("operation_type = ?", info.OperationType)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.IsFirstDeduct != nil {
		db = db.Where("is_first_deduct = ?", *info.IsFirstDeduct)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := info.Page
	if page <= 0 {
		page = 1
	}
	pageSize := info.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := db.Order("id desc").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *deduct) UpdateStatus(id uint, status string) error {
	return global.GVA_DB.Model(&model.DeductRecord{}).Where("id = ?", id).Update("status", status).Error
}

func (s *deduct) SetCallbackResult(id uint, result string, callbackTime int64) error {
	return global.GVA_DB.Model(&model.DeductRecord{}).Where("id = ?", id).Updates(map[string]any{
		"callback_result": result,
		"callback_time":   callbackTime,
	}).Error
}

func (s *deduct) CreateContractRecord(record *model.ContractRecord) error {
	return global.GVA_DB.Create(record).Error
}

func (s *deduct) UpdateContractRecord(record *model.ContractRecord) error {
	if record.ID == 0 {
		return errors.New("签约记录ID不能为空")
	}
	return global.GVA_DB.Model(&model.ContractRecord{}).Where("id = ?", record.ID).Updates(record).Error
}

func (s *deduct) GetContractRecordList(info mockReq.ContractRecordSearch) ([]model.ContractRecord, int64, error) {
	var list []model.ContractRecord
	var total int64

	db := global.GVA_DB.Model(&model.ContractRecord{})
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.MerchantID != 0 {
		db = db.Where("merchant_id = ?", info.MerchantID)
	}
	if info.ContractID != 0 {
		db = db.Where("contract_id = ?", info.ContractID)
	}
	if info.OperationType != "" {
		db = db.Where("operation_type = ?", info.OperationType)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := info.Page
	if page <= 0 {
		page = 1
	}
	pageSize := info.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := db.Order("id desc").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *deduct) GetContractRecord(id uint) (model.ContractRecord, error) {
	var r model.ContractRecord
	err := global.GVA_DB.Where("id = ?", id).First(&r).Error
	return r, err
}

func (s *deduct) GetContractRecordByContract(contractID uint, operationType string) (model.ContractRecord, error) {
	var r model.ContractRecord
	query := global.GVA_DB.Where("contract_id = ?", contractID)
	if operationType != "" {
		query = query.Where("operation_type = ?", operationType)
	}
	err := query.Order("id desc").First(&r).Error
	return r, err
}

func (s *deduct) SetContractRecordCallbackResult(id uint, result string, callbackTime int64) error {
	return global.GVA_DB.Model(&model.ContractRecord{}).Where("id = ?", id).Updates(map[string]any{
		"callback_result": result,
		"callback_time":   callbackTime,
	}).Error
}

func (s *deduct) GetMerchantFromDB(mchID string) (model.Merchant, error) {
	var m model.Merchant
	err := global.GVA_DB.Where("mch_id = ? AND active = ?", mchID, true).First(&m).Error
	return m, err
}

func (s *deduct) GetContractFromDB(outContractCode string) (model.Contract, error) {
	var c model.Contract
	err := global.GVA_DB.Where("out_contract_id = ?", outContractCode).First(&c).Error
	return c, err
}

func (s *deduct) GetContractByContractIDFromDB(contractID string) (model.Contract, error) {
	var c model.Contract
	err := global.GVA_DB.Where("contract_id = ?", contractID).First(&c).Error
	return c, err
}

func (s *deduct) GetContractByID(id uint) (model.Contract, error) {
	var c model.Contract
	err := global.GVA_DB.Where("id = ?", id).First(&c).Error
	return c, err
}

func (s *deduct) SaveContractRecord(record *model.ContractRecord) error {
	return global.GVA_DB.Create(record).Error
}

func (s *deduct) SaveDeductRecord(record *model.DeductRecord) error {
	return global.GVA_DB.Create(record).Error
}

func (s *deduct) UpdateContractRecordStatus(id uint, status string, errCode string, errMsg string) error {
	updates := map[string]any{
		"status":        status,
		"error_code":    errCode,
		"error_message": errMsg,
	}
	return global.GVA_DB.Model(&model.ContractRecord{}).Where("id = ?", id).Updates(updates).Error
}

func (s *deduct) UpdateContractRecordResponse(id uint, responseXML string, status string) error {
	return global.GVA_DB.Model(&model.ContractRecord{}).Where("id = ?", id).Updates(map[string]any{
		"response_xml": responseXML,
		"status":       status,
	}).Error
}

func (s *deduct) UpdateDeductRecordStatus(id uint, status string, errCode string, errMsg string) error {
	updates := map[string]any{
		"status":        status,
		"error_code":    errCode,
		"error_message": errMsg,
	}
	return global.GVA_DB.Model(&model.DeductRecord{}).Where("id = ?", id).Updates(updates).Error
}

func (s *deduct) UpdateDeductRecordResponse(id uint, responseData string, status string) error {
	return global.GVA_DB.Model(&model.DeductRecord{}).Where("id = ?", id).Updates(map[string]any{
		"response_data": responseData,
		"status":        status,
	}).Error
}

func (s *deduct) SetContractExpireTime(id uint, durationMinutes int) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Update("expire_time", time.Now().Add(time.Duration(durationMinutes)*time.Minute)).Error
}

func (s *deduct) SetContractContractID(id uint, contractID string, signSerialNo string) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Updates(map[string]any{
		"contract_id":    contractID,
		"sign_serial_no": signSerialNo,
	}).Error
}

func (s *deduct) SetContractStatus(id uint, status string, terminateType string) error {
	updates := map[string]any{
		"contract_status": status,
		"terminate_type":  terminateType,
	}
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Updates(updates).Error
}

func (s *deduct) GetDeductRecordByOutTradeNo(outTradeNo string) (model.DeductRecord, error) {
	var r model.DeductRecord
	err := global.GVA_DB.Where("out_trade_no = ?", outTradeNo).First(&r).Error
	return r, err
}

func (s *deduct) GetDeductRecordByTransactionID(transactionID string) (model.DeductRecord, error) {
	var r model.DeductRecord
	err := global.GVA_DB.Where("transaction_id = ?", transactionID).First(&r).Error
	return r, err
}

func (s *deduct) SetDeductRecordTransactionID(id uint, transactionID string) error {
	return global.GVA_DB.Model(&model.DeductRecord{}).Where("id = ?", id).Update("transaction_id", transactionID).Error
}

func (s *deduct) UpdateDeductRecordByCallback(id uint, status string, transactionID string, callbackResult string, callbackTime int64, errCode string, errMsg string) error {
	updates := map[string]any{
		"status":          status,
		"callback_result": callbackResult,
		"callback_time":   callbackTime,
		"error_code":      errCode,
		"error_message":   errMsg,
	}
	if strings.TrimSpace(transactionID) != "" {
		updates["transaction_id"] = transactionID
	}
	return global.GVA_DB.Model(&model.DeductRecord{}).Where("id = ?", id).Updates(updates).Error
}

func (s *deduct) ClearContractPreNotify(contractID uint) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", contractID).Update("pre_notify_called", false).Error
}

func (s *deduct) SetContractPreNotify(contractID uint) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", contractID).Updates(map[string]any{
		"pre_notify_called":    true,
		"last_pre_notify_time": now,
	}).Error
}
