package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	mockReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model/request"
)

type deductCallback struct{}

func (s *deductCallback) Create(record *model.DeductCallbackRecord) error {
	return global.GVA_DB.Create(record).Error
}

func (s *deductCallback) GetByID(id uint) (model.DeductCallbackRecord, error) {
	var record model.DeductCallbackRecord
	err := global.GVA_DB.Where("id = ?", id).First(&record).Error
	return record, err
}

func (s *deductCallback) GetList(info mockReq.DeductCallbackRecordSearch) ([]model.DeductCallbackRecord, int64, error) {
	var list []model.DeductCallbackRecord
	var total int64

	db := global.GVA_DB.Model(&model.DeductCallbackRecord{})
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if strings.TrimSpace(info.MchID) != "" {
		db = db.Where("mch_id LIKE ?", "%"+info.MchID+"%")
	}
	if strings.TrimSpace(info.OutTradeNo) != "" {
		db = db.Where("out_trade_no LIKE ?", "%"+info.OutTradeNo+"%")
	}
	if strings.TrimSpace(info.TransactionID) != "" {
		db = db.Where("transaction_id LIKE ?", "%"+info.TransactionID+"%")
	}
	if strings.TrimSpace(info.TradeState) != "" {
		db = db.Where("trade_state = ?", info.TradeState)
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

func (s *deductCallback) BuildDeductSignParams(req model.DeductNotifyRequest) map[string]string {
	params := map[string]string{
		"return_code":    req.ReturnCode,
		"appid":          req.AppID,
		"mch_id":         req.MchID,
		"out_trade_no":   req.OutTradeNo,
		"transaction_id": req.TransactionID,
		"trade_type":     req.TradeType,
		"trade_state":    req.TradeState,
		"bank_type":      req.BankType,
		"total_amount":   strconv.FormatInt(req.TotalAmount, 10),
		"cash_amount":    strconv.FormatInt(req.CashAmount, 10),
		"time_end":       req.TimeEnd,
		"timestamp":      req.TimeStamp,
		"nonce":          req.Nonce,
		"sign":           req.Sign,
	}
	if strings.TrimSpace(req.SignType) != "" {
		params["sign_type"] = req.SignType
	}
	return params
}

func (s *deductCallback) ValidateDeductCallback(req model.DeductNotifyRequest, verifySign bool) error {
	if strings.TrimSpace(req.MchID) == "" {
		return errors.New("mch_id不能为空")
	}
	if strings.TrimSpace(req.OutTradeNo) == "" && strings.TrimSpace(req.TransactionID) == "" {
		return errors.New("out_trade_no和transaction_id不能同时为空")
	}
	if !Service.Signature.RequireSignIfNeeded(verifySign, req.Sign) {
		return errors.New("sign不能为空")
	}
	return nil
}

func (s *deductCallback) VerifyDeductCallback(req model.DeductNotifyRequest, verifySign bool, key string) error {
	params := s.BuildDeductSignParams(req)
	return Service.Signature.VerifyIfNeeded(verifySign, params, key)
}

func (s *deductCallback) LocateMerchantAndDeduct(req model.DeductNotifyRequest) (model.Merchant, model.DeductRecord, model.Contract, error) {
	merchant, err := Service.Deduct.GetMerchantFromDB(req.MchID)
	if err != nil {
		return model.Merchant{}, model.DeductRecord{}, model.Contract{}, err
	}

	var record model.DeductRecord
	if strings.TrimSpace(req.OutTradeNo) != "" {
		record, err = Service.Deduct.GetDeductRecordByOutTradeNo(req.OutTradeNo)
	}
	if err != nil && strings.TrimSpace(req.TransactionID) != "" {
		record, err = Service.Deduct.GetDeductRecordByTransactionID(req.TransactionID)
	}
	if err != nil {
		return model.Merchant{}, model.DeductRecord{}, model.Contract{}, err
	}

	contract, err := Service.Deduct.GetContractByID(record.ContractID)
	if err != nil {
		return model.Merchant{}, model.DeductRecord{}, model.Contract{}, err
	}
	return merchant, record, contract, nil
}
