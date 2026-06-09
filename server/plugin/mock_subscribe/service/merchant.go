package service

import (
	"errors"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
)

type merchant struct{}

func (s *merchant) CreateMerchant(info *model.Merchant) error {
	if err := s.validateMerchant(info); err != nil {
		return err
	}
	return global.GVA_DB.Create(info).Error
}

func (s *merchant) UpdateMerchant(info *model.Merchant) error {
	if info.ID == 0 {
		return errors.New("商户ID不能为空")
	}
	if err := s.validateMerchant(info); err != nil {
		return err
	}
	updates := map[string]any{
		"app_id":                   info.AppID,
		"mch_id":                   info.MchID,
		"contract_mch_id":          info.ContractMchID,
		"contract_app_id":          info.ContractAppID,
		"display_name":             info.DisplayName,
		"sign_key":                 info.SignKey,
		"verify_sign":              info.VerifySign,
		"contract_template":        info.ContractTemplate,
		"sign_callback_enabled":    info.SignCallbackEnabled,
		"sign_callback_delay":      info.SignCallbackDelay,
		"deduct_callback_enabled":  info.DeductCallbackEnabled,
		"deduct_callback_delay":    info.DeductCallbackDelay,
		"sign_target_status":       info.SignTargetStatus,
		"sign_status_delay":        info.SignStatusDelay,
		"deduct_target_status":     info.DeductTargetStatus,
		"deduct_status_delay":      info.DeductStatusDelay,
		"terminate_notify_enabled": info.TerminateNotifyEnabled,
		"sign_duration_minutes":    info.SignDurationMinutes,
		"strict_deduct_rule":       info.StrictDeductRule,
		"active":                   info.Active,
		"updated_at":               time.Now(),
	}
	return global.GVA_DB.Model(&model.Merchant{}).Where("id = ?", info.ID).Updates(updates).Error
}

func (s *merchant) DeleteMerchant(id uint) error {
	if id == 0 {
		return errors.New("商户ID不能为空")
	}
	return global.GVA_DB.Delete(&model.Merchant{}, id).Error
}

func (s *merchant) GetMerchant(id uint) (model.Merchant, error) {
	var merchant model.Merchant
	err := global.GVA_DB.Where("id = ?", id).First(&merchant).Error
	return merchant, err
}

func (s *merchant) GetMerchantByMchID(mchID string) (model.Merchant, error) {
	var merchant model.Merchant
	err := global.GVA_DB.Where("mch_id = ? AND active = ?", mchID, true).First(&merchant).Error
	return merchant, err
}

func (s *merchant) GetMerchantList(info mockReq.MerchantSearch) ([]model.Merchant, int64, error) {
	var list []model.Merchant
	var total int64

	db := global.GVA_DB.Model(&model.Merchant{})
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.AppID != "" {
		db = db.Where("app_id LIKE ?", "%"+info.AppID+"%")
	}
	if info.MchID != "" {
		db = db.Where("mch_id LIKE ?", "%"+info.MchID+"%")
	}
	if info.Active != nil {
		db = db.Where("active = ?", *info.Active)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes((&commonReq.PageInfo{Page: info.Page, PageSize: info.PageSize}).Paginate()).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *merchant) validateMerchant(info *model.Merchant) error {
	if strings.TrimSpace(info.AppID) == "" {
		return errors.New("应用ID不能为空")
	}
	if strings.TrimSpace(info.MchID) == "" {
		return errors.New("商户号不能为空")
	}
	if info.VerifySign && strings.TrimSpace(info.SignKey) == "" {
		return errors.New("验签开启时签名key不能为空")
	}
	if info.SignStatusDelay < 0 || info.SignCallbackDelay < 0 || info.DeductStatusDelay < 0 || info.DeductCallbackDelay < 0 {
		return errors.New("延时配置不能小于0")
	}
	if info.SignDurationMinutes <= 0 {
		return errors.New("签约时长必须大于0")
	}
	if info.SignTargetStatus == "" {
		info.SignTargetStatus = model.ContractStatusActive
	}
	if info.DeductTargetStatus == "" {
		info.DeductTargetStatus = model.DeductStatusSuccess
	}
	if !s.validContractStatus(info.SignTargetStatus) {
		return errors.New("签约目标状态不合法")
	}
	if !s.validDeductStatus(info.DeductTargetStatus) {
		return errors.New("扣款目标状态不合法")
	}
	if info.CreatedAt.IsZero() {
		info.CreatedAt = time.Now()
	}
	return nil
}

func (s *merchant) validContractStatus(status string) bool {
	switch status {
	case model.ContractStatusPending, model.ContractStatusActive, model.ContractStatusFailed, model.ContractStatusTerminated, model.ContractStatusExpired, model.ContractStatusPause:
		return true
	default:
		return false
	}
}

func (s *merchant) validDeductStatus(status string) bool {
	switch status {
	case model.DeductStatusPending, model.DeductStatusSuccess, model.DeductStatusFailed, model.DeductStatusRefunding, model.DeductStatusRefunded:
		return true
	default:
		return false
	}
}
