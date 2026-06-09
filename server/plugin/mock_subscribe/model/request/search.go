package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

// MerchantSearch 商户配置搜索
type MerchantSearch struct {
	request.PageInfo
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	AppID          string     `json:"appId" form:"appId"`
	MchID          string     `json:"mchId" form:"mchId"`
	Active         *bool      `json:"active" form:"active"`
}

// ContractSearch 签约协议搜索
type ContractSearch struct {
	request.PageInfo
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	MerchantID     uint       `json:"merchantId" form:"merchantId"`
	OpenID         string     `json:"openId" form:"openId"`
	ContractStatus string     `json:"contractStatus" form:"contractStatus"`
	OutContractID  string     `json:"outContractId" form:"outContractId"`
}

// DeductRecordSearch 扣款记录搜索
type DeductRecordSearch struct {
	request.PageInfo
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	MerchantID     uint       `json:"merchantId" form:"merchantId"`
	ContractID     uint       `json:"contractId" form:"contractId"`
	OperationType  string     `json:"operationType" form:"operationType"`
	Status         string     `json:"status" form:"status"`
	IsFirstDeduct  *bool      `json:"isFirstDeduct" form:"isFirstDeduct"`
}

// ContractRecordSearch 签约记录搜索
type ContractRecordSearch struct {
	request.PageInfo
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	MerchantID     uint       `json:"merchantId" form:"merchantId"`
	ContractID     uint       `json:"contractId" form:"contractId"`
	OperationType  string     `json:"operationType" form:"operationType"`
	Status         string     `json:"status" form:"status"`
}
