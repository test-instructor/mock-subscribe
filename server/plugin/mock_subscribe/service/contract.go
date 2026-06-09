package service

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
)

type contract struct{}

func (s *contract) CreateContract(contract *model.Contract) error {
	if contract.MerchantID == 0 {
		return errors.New("商户ID不能为空")
	}
	if contract.OutContractID == "" {
		return errors.New("外部签约单号不能为空")
	}
	if contract.OpenID == "" {
		return errors.New("用户OpenID不能为空")
	}
	contract.ContractStatus = model.ContractStatusPending
	contract.IsFirstDeduct = true
	return global.GVA_DB.Create(contract).Error
}

func (s *contract) UpdateContract(contract *model.Contract) error {
	if contract.ID == 0 {
		return errors.New("签约ID不能为空")
	}
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", contract.ID).Updates(contract).Error
}

func (s *contract) GetContract(id uint) (model.Contract, error) {
	var c model.Contract
	err := global.GVA_DB.Where("id = ?", id).First(&c).Error
	return c, err
}

func (s *contract) GetContractByOutID(outContractID string) (model.Contract, error) {
	var c model.Contract
	err := global.GVA_DB.Where("out_contract_id = ?", outContractID).First(&c).Error
	return c, err
}

func (s *contract) GetContractByContractID(contractID string) (model.Contract, error) {
	var c model.Contract
	err := global.GVA_DB.Where("contract_id = ?", contractID).First(&c).Error
	return c, err
}

func (s *contract) GetContractList(info mockReq.ContractSearch) ([]model.Contract, int64, error) {
	var list []model.Contract
	var total int64

	db := global.GVA_DB.Model(&model.Contract{})
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.MerchantID != 0 {
		db = db.Where("merchant_id = ?", info.MerchantID)
	}
	if info.OpenID != "" {
		db = db.Where("open_id LIKE ?", "%"+info.OpenID+"%")
	}
	if info.ContractStatus != "" {
		db = db.Where("contract_status = ?", info.ContractStatus)
	}
	if info.OutContractID != "" {
		db = db.Where("out_contract_id LIKE ?", "%"+info.OutContractID+"%")
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

func (s *contract) UpdateContractStatus(id uint, status string, terminateType string) error {
	updates := map[string]any{
		"contract_status": status,
		"terminate_type":  terminateType,
	}
	if status == model.ContractStatusActive {
		updates["expire_time"] = time.Now().Add(24 * time.Hour)
	}
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Updates(updates).Error
}

func (s *contract) SetExpireTime(id uint, durationMinutes int) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Update("expire_time", time.Now().Add(time.Duration(durationMinutes)*time.Minute)).Error
}

func (s *contract) SetContractID(id uint, contractID string, signSerialNo string) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Updates(map[string]any{
		"contract_id":    contractID,
		"sign_serial_no": signSerialNo,
	}).Error
}

func (s *contract) HasActiveContract(outContractID string) bool {
	var c model.Contract
	err := global.GVA_DB.Where("out_contract_id = ? AND contract_status = ?", outContractID, model.ContractStatusActive).First(&c).Error
	return err == nil
}

func (s *contract) ResetFirstDeduct(id uint) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Update("is_first_deduct", true).Error
}

func (s *contract) MarkFirstDeductDone(id uint) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Updates(map[string]any{
		"is_first_deduct":   false,
		"pre_notify_called": false,
	}).Error
}

func (s *contract) MarkPreNotifyCalled(id uint) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Updates(map[string]any{
		"pre_notify_called":    true,
		"last_pre_notify_time": now,
	}).Error
}

func (s *contract) ClearPreNotify(id uint) error {
	return global.GVA_DB.Model(&model.Contract{}).Where("id = ?", id).Update("pre_notify_called", false).Error
}
