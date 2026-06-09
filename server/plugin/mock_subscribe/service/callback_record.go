package service

import (
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
)

type callbackRecord struct{}

func (s *callbackRecord) Create(record *model.CallbackRecord) error {
	return global.GVA_DB.Create(record).Error
}

func (s *callbackRecord) GetByID(id uint) (model.CallbackRecord, error) {
	var record model.CallbackRecord
	err := global.GVA_DB.Where("id = ?", id).First(&record).Error
	return record, err
}

func (s *callbackRecord) GetList(info mockReq.CallbackRecordSearch) ([]model.CallbackRecord, int64, error) {
	var list []model.CallbackRecord
	var total int64

	db := global.GVA_DB.Model(&model.CallbackRecord{})
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if strings.TrimSpace(info.MchID) != "" {
		db = db.Where("mch_id LIKE ?", "%"+info.MchID+"%")
	}
	if strings.TrimSpace(info.OutContractCode) != "" {
		db = db.Where("out_contract_code LIKE ?", "%"+info.OutContractCode+"%")
	}
	if strings.TrimSpace(info.ContractCode) != "" {
		db = db.Where("contract_code LIKE ?", "%"+info.ContractCode+"%")
	}
	if strings.TrimSpace(info.CallbackType) != "" {
		db = db.Where("callback_type = ?", info.CallbackType)
	}
	if info.SignValid != nil {
		db = db.Where("sign_valid = ?", *info.SignValid)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := commonReq.PageInfo{Page: info.Page, PageSize: info.PageSize}
	if err := db.Scopes(page.Paginate()).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *callbackRecord) BuildContractSignParams(req model.ContractCallbackRequest) map[string]string {
	return map[string]string{
		"appid":                req.AppID,
		"mch_id":               req.MchID,
		"contract_id":          req.ContractID,
		"contract_status":      req.ContractStatus,
		"contract_ending_type": req.ContractEndingType,
		"out_contract_code":    req.OutContractCode,
		"contract_ext_id":      req.ContractExtID,
		"sign_type":            req.SignType,
		"timestamp":            req.TimeStamp,
		"nonce":                req.Nonce,
		"sign":                 req.Sign,
		"plan_id":              req.PlanID,
		"openid":               req.OpenID,
	}
}

func (s *callbackRecord) LocateMerchantAndContract(req model.ContractCallbackRequest) (model.Merchant, model.Contract, error) {
	contract, err := Service.Contract.GetContractByContractID(req.ContractID)
	if err != nil && strings.TrimSpace(req.OutContractCode) != "" {
		contract, err = Service.Contract.GetContractByOutID(req.OutContractCode)
	}
	if err != nil {
		return model.Merchant{}, model.Contract{}, err
	}
	merchant, err := Service.Merchant.GetMerchant(contract.MerchantID)
	if err != nil {
		return model.Merchant{}, model.Contract{}, err
	}
	return merchant, contract, nil
}

func (s *callbackRecord) VerifyContractCallback(req model.ContractCallbackRequest, key string) error {
	params := s.BuildContractSignParams(req)
	return Service.Signature.Verify(params, key)
}

func (s *callbackRecord) ValidateContractCallback(req model.ContractCallbackRequest) error {
	if strings.TrimSpace(req.MchID) == "" {
		return errors.New("mch_id不能为空")
	}
	if strings.TrimSpace(req.ContractID) == "" && strings.TrimSpace(req.OutContractCode) == "" {
		return errors.New("contract_id和out_contract_code不能同时为空")
	}
	if strings.TrimSpace(req.Sign) == "" {
		return errors.New("sign不能为空")
	}
	return nil
}
