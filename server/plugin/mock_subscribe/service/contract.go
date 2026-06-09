package service

import (
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
	"gorm.io/gorm"
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
	return global.GVA_DB.Create(contract).Error
}

func (s *contract) CreateContractStatus(status *model.ContractStatusRecord) error {
	if status.ContractID == 0 {
		return errors.New("签约ID不能为空")
	}
	if status.OutContractID == "" {
		return errors.New("外部签约单号不能为空")
	}
	if status.ContractStatus == "" {
		status.ContractStatus = model.ContractStatusPending
	}
	if !status.IsFirstDeduct {
		status.IsFirstDeduct = true
	}
	return global.GVA_DB.Create(status).Error
}

func (s *contract) CreateContractWithStatus(contract *model.Contract, status *model.ContractStatusRecord) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(contract).Error; err != nil {
			return err
		}
		status.ContractID = contract.ID
		if status.MerchantID == 0 {
			status.MerchantID = contract.MerchantID
		}
		if status.OutContractID == "" {
			status.OutContractID = contract.OutContractID
		}
		if status.ContractStatus == "" {
			status.ContractStatus = model.ContractStatusPending
		}
		status.IsFirstDeduct = true
		return tx.Create(status).Error
	})
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

func (s *contract) GetContractStatusByContractID(contractID uint) (model.ContractStatusRecord, error) {
	var status model.ContractStatusRecord
	err := global.GVA_DB.Where("contract_id = ?", contractID).Order("id desc").First(&status).Error
	return status, err
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

func (s *contract) GetContractList(info mockReq.ContractSearch) ([]map[string]any, int64, error) {
	var contracts []model.Contract
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
	if err := db.Order("id desc").Offset(offset).Limit(pageSize).Find(&contracts).Error; err != nil {
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(contracts))
	for _, item := range contracts {
		status, err := s.GetContractStatusByContractID(item.ID)
		if err != nil {
			if info.ContractStatus != "" {
				continue
			}
			list = append(list, map[string]any{
				"contract": item,
			})
			continue
		}
		if info.ContractStatus != "" && status.ContractStatus != info.ContractStatus {
			continue
		}
		list = append(list, map[string]any{
			"contract": item,
			"status":   status,
		})
	}
	if info.ContractStatus != "" {
		total = int64(len(list))
	}
	return list, total, nil
}

func (s *contract) UpdateContractStatus(id uint, status string, terminateType string) error {
	updates := map[string]any{
		"contract_status": status,
		"terminate_type":  terminateType,
	}
	if status == model.ContractStatusActive {
		expireTime := time.Now().Add(24 * time.Hour)
		updates["expire_time"] = &expireTime
	}
	return global.GVA_DB.Model(&model.ContractStatusRecord{}).Where("contract_id = ?", id).Updates(updates).Error
}

func (s *contract) SetExpireTime(id uint, durationMinutes int) error {
	expireTime := time.Now().Add(time.Duration(durationMinutes) * time.Minute)
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Contract{}).Where("id = ?", id).Update("expire_time", expireTime).Error; err != nil {
			return err
		}
		return tx.Model(&model.ContractStatusRecord{}).Where("contract_id = ?", id).Update("expire_time", expireTime).Error
	})
}

func (s *contract) SetContractID(id uint, contractID string, signSerialNo string) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Contract{}).Where("id = ?", id).Updates(map[string]any{
			"contract_id":    contractID,
			"sign_serial_no": signSerialNo,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&model.ContractStatusRecord{}).Where("contract_id = ?", id).Updates(map[string]any{
			"contract_no":    contractID,
			"sign_serial_no": signSerialNo,
		}).Error
	})
}

func (s *contract) HasActiveContract(outContractID string) bool {
	var status model.ContractStatusRecord
	err := global.GVA_DB.Where("out_contract_id = ? AND contract_status IN ?", outContractID, []string{model.ContractStatusPending, model.ContractStatusActive}).First(&status).Error
	return err == nil
}

func (s *contract) HasActiveContractByUser(merchantID uint, outUserID string, openID string) bool {
	if merchantID == 0 || (outUserID == "" && openID == "") {
		return false
	}
	query := global.GVA_DB.Model(&model.Contract{}).Where("merchant_id = ?", merchantID)
	if outUserID != "" {
		query = query.Where("out_user_id = ?", outUserID)
	} else {
		query = query.Where("open_id = ?", openID)
	}
	var contracts []model.Contract
	if err := query.Find(&contracts).Error; err != nil {
		return false
	}
	for _, item := range contracts {
		var status model.ContractStatusRecord
		if err := global.GVA_DB.Where("contract_id = ? AND contract_status IN ?", item.ID, []string{model.ContractStatusPending, model.ContractStatusActive}).Order("id desc").First(&status).Error; err == nil {
			return true
		}
	}
	return false
}

func (s *contract) ResetFirstDeduct(id uint) error {
	return global.GVA_DB.Model(&model.ContractStatusRecord{}).Where("contract_id = ?", id).Update("is_first_deduct", true).Error
}

func (s *contract) MarkFirstDeductDone(id uint) error {
	return global.GVA_DB.Model(&model.ContractStatusRecord{}).Where("contract_id = ?", id).Updates(map[string]any{
		"is_first_deduct":   false,
		"pre_notify_called": false,
	}).Error
}

func (s *contract) MarkPreNotifyCalled(id uint) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.ContractStatusRecord{}).Where("contract_id = ?", id).Updates(map[string]any{
		"pre_notify_called":    true,
		"last_pre_notify_time": now,
	}).Error
}

func (s *contract) ClearPreNotify(id uint) error {
	return global.GVA_DB.Model(&model.ContractStatusRecord{}).Where("contract_id = ?", id).Update("pre_notify_called", false).Error
}
