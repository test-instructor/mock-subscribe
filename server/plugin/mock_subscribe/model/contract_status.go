package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ContractStatusRecord 订阅状态信息
type ContractStatusRecord struct {
	global.GVA_MODEL

	ContractID uint `json:"contractId" form:"contractId" gorm:"column:contract_id;comment:签约ID;index"`
	MerchantID uint `json:"merchantId" form:"merchantId" gorm:"column:merchant_id;comment:商户ID;index"`

	OutContractID string `json:"outContractId" form:"outContractId" gorm:"column:out_contract_id;comment:外部签约单号;index"`
	ContractNo    string `json:"contractNo" form:"contractNo" gorm:"column:contract_no;comment:协议号;index"`
	SignSerialNo  string `json:"signSerialNo" form:"signSerialNo" gorm:"column:sign_serial_no;comment:签约号"`

	ContractStatus string `json:"contractStatus" form:"contractStatus" gorm:"column:contract_status;comment:协议状态;index"`
	TerminateType  string `json:"terminateType" form:"terminateType" gorm:"column:terminate_type;comment:解约方式"`

	IsFirstDeduct     bool       `json:"isFirstDeduct" form:"isFirstDeduct" gorm:"column:is_first_deduct;default:true;comment:是否首次扣款"`
	LastPreNotifyTime *time.Time `json:"lastPreNotifyTime" form:"lastPreNotifyTime" gorm:"column:last_pre_notify_time;comment:最近一次预扣费通知时间"`
	PreNotifyCalled   bool       `json:"preNotifyCalled" form:"preNotifyCalled" gorm:"column:pre_notify_called;default:false;comment:是否已调用预扣费通知"`

	ExpireTime *time.Time `json:"expireTime" form:"expireTime" gorm:"column:expire_time;comment:到期时间"`
}

func (ContractStatusRecord) TableName() string {
	return "gva_mock_contract_status"
}
