package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Contract 订阅主信息（订阅信息和订阅状态分表存放）
type Contract struct {
	global.GVA_MODEL

	// 商户关联
	MerchantID uint `json:"merchantId" form:"merchantId" gorm:"column:merchant_id;comment:商户ID"`

	// 用户维度
	OpenID    string `json:"openId" form:"openId" gorm:"column:open_id;comment:用户OpenID"`
	OutUserID string `json:"outUserId" form:"outUserId" gorm:"column:out_user_id;comment:外部用户标识"`

	// 协议维度
	OutContractID string `json:"outContractId" form:"outContractId" gorm:"column:out_contract_id;comment:外部签约单号;uniqueIndex"`
	ContractID    string `json:"contractId" form:"contractId" gorm:"column:contract_id;comment:协议号"`
	SignSerialNo  string `json:"signSerialNo" form:"signSerialNo" gorm:"column:sign_serial_no;comment:签约号"`
	PlanID        string `json:"planId" form:"planId" gorm:"column:plan_id;comment:签约模板ID"`

	// 回调地址
	NotifyURL string `json:"notifyUrl" form:"notifyUrl" gorm:"column:notify_url;comment:签约回调地址"`

	// 当前协议状态
	ContractStatus string `json:"contractStatus" form:"contractStatus" gorm:"column:contract_status;comment:协议状态"`
	TerminateType  string `json:"terminateType" form:"terminateType" gorm:"column:terminate_type;comment:解约方式"`

	// 有效期
	ExpireTime time.Time `json:"expireTime" form:"expireTime" gorm:"column:expire_time;comment:到期时间"`

	// 扣款相关
	IsFirstDeduct     bool       `json:"isFirstDeduct" form:"isFirstDeduct" gorm:"column:is_first_deduct;default:true;comment:是否首次扣款"`
	LastPreNotifyTime *time.Time `json:"lastPreNotifyTime" form:"lastPreNotifyTime" gorm:"column:last_pre_notify_time;comment:最近一次预扣费通知时间"`
	PreNotifyCalled   bool       `json:"preNotifyCalled" form:"preNotifyCalled" gorm:"column:pre_notify_called;default:false;comment:是否已调用预扣费通知"`

	// 原始请求数据（JSON存储）
	RequestData string `json:"requestData" form:"requestData" gorm:"column:request_data;type:text;comment:原始请求数据"`
}

func (Contract) TableName() string {
	return "gva_mock_contract"
}
