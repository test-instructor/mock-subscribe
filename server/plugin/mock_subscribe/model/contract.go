package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Contract 订阅主信息（仅保存订阅信息）
type Contract struct {
	global.GVA_MODEL

	MerchantID uint `json:"merchantId" form:"merchantId" gorm:"column:merchant_id;comment:商户ID"`

	OpenID    string `json:"openId" form:"openId" gorm:"column:open_id;comment:用户OpenID"`
	OutUserID string `json:"outUserId" form:"outUserId" gorm:"column:out_user_id;comment:外部用户标识"`

	OutContractID string `json:"outContractId" form:"outContractId" gorm:"column:out_contract_id;comment:外部签约单号;uniqueIndex"`
	ContractID    string `json:"contractId" form:"contractId" gorm:"column:contract_id;comment:协议号;index"`
	SignSerialNo  string `json:"signSerialNo" form:"signSerialNo" gorm:"column:sign_serial_no;comment:签约号"`
	PlanID        string `json:"planId" form:"planId" gorm:"column:plan_id;comment:签约模板ID"`

	NotifyURL string `json:"notifyUrl" form:"notifyUrl" gorm:"column:notify_url;comment:签约回调地址"`

	ExpireTime time.Time `json:"expireTime" form:"expireTime" gorm:"column:expire_time;comment:到期时间"`

	RequestData string `json:"requestData" form:"requestData" gorm:"column:request_data;type:text;comment:原始请求数据"`
}

func (Contract) TableName() string {
	return "gva_mock_contract"
}
