package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Merchant 商户配置
type Merchant struct {
	global.GVA_MODEL

	// 商户标识配置
	AppID            string `json:"appId" form:"appId" gorm:"column:app_id;comment:应用ID"`
	MchID            string `json:"mchId" form:"mchId" gorm:"column:mch_id;comment:商户号"`
	ContractMchID    string `json:"contractMchId" form:"contractMchId" gorm:"column:contract_mch_id;comment:签约商户号"`
	ContractAppID    string `json:"contractAppId" form:"contractAppId" gorm:"column:contract_app_id;comment:签约AppID"`
	DisplayName      string `json:"displayName" form:"displayName" gorm:"column:display_name;comment:用户账户展示名称"`
	SignKey          string `json:"signKey" form:"signKey" gorm:"column:sign_key;comment:签名key"`
	ContractTemplate string `json:"contractTemplate" form:"contractTemplate" gorm:"column:contract_template;comment:签约模板"`

	// 签约回调配置
	SignCallbackEnabled bool `json:"signCallbackEnabled" form:"signCallbackEnabled" gorm:"column:sign_callback_enabled;comment:签约回调开关"`
	SignCallbackDelay   int  `json:"signCallbackDelay" form:"signCallbackDelay" gorm:"column:sign_callback_delay;comment:签约回调延时(秒)"`

	// 扣款回调配置
	DeductCallbackEnabled bool `json:"deductCallbackEnabled" form:"deductCallbackEnabled" gorm:"column:deduct_callback_enabled;comment:扣款回调开关"`
	DeductCallbackDelay   int  `json:"deductCallbackDelay" form:"deductCallbackDelay" gorm:"column:deduct_callback_delay;comment:扣款回调延时(秒)"`

	// 签约状态配置
	SignTargetStatus string `json:"signTargetStatus" form:"signTargetStatus" gorm:"column:sign_target_status;comment:签约目标状态"`
	SignStatusDelay  int    `json:"signStatusDelay" form:"signStatusDelay" gorm:"column:sign_status_delay;comment:签约状态延时(秒)"`

	// 扣款状态配置
	DeductTargetStatus string `json:"deductTargetStatus" form:"deductTargetStatus" gorm:"column:deduct_target_status;comment:扣款目标状态"`
	DeductStatusDelay  int    `json:"deductStatusDelay" form:"deductStatusDelay" gorm:"column:deduct_status_delay;comment:扣款状态延时(秒)"`

	// 行为配置
	TerminateNotifyEnabled bool `json:"terminateNotifyEnabled" form:"terminateNotifyEnabled" gorm:"column:terminate_notify_enabled;comment:解约是否通知商户"`
	SignDurationMinutes    int  `json:"signDurationMinutes" form:"signDurationMinutes" gorm:"column:sign_duration_minutes;comment:签约时长(分钟)"`
	StrictDeductRule       bool `json:"strictDeductRule" form:"strictDeductRule" gorm:"column:strict_deduct_rule;comment:是否严格按扣费规则"`

	// 状态
	Active bool `json:"active" form:"active" gorm:"column:active;default:true;comment:是否启用"`
}

func (Merchant) TableName() string {
	return "gva_mock_merchant"
}
