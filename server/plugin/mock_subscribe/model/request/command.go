package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Merchant 商户配置(响应用,继承gorm.Model用于ID传递)
type MerchantInfo struct {
	global.GVA_MODEL

	AppID            string `json:"appId" gorm:"column:app_id"`
	MchID            string `json:"mchId" gorm:"column:mch_id"`
	ContractMchID    string `json:"contractMchId" gorm:"column:contract_mch_id"`
	ContractAppID    string `json:"contractAppId" gorm:"column:contract_app_id"`
	DisplayName      string `json:"displayName" gorm:"column:display_name"`
	SignKey          string `json:"signKey" gorm:"column:sign_key"`
	ContractTemplate string `json:"contractTemplate" gorm:"column:contract_template"`

	SignCallbackEnabled bool `json:"signCallbackEnabled" gorm:"column:sign_callback_enabled"`
	SignCallbackDelay   int  `json:"signCallbackDelay" gorm:"column:sign_callback_delay"`

	DeductCallbackEnabled bool `json:"deductCallbackEnabled" gorm:"column:deduct_callback_enabled"`
	DeductCallbackDelay   int  `json:"deductCallbackDelay" gorm:"column:deduct_callback_delay"`

	SignTargetStatus string `json:"signTargetStatus" gorm:"column:sign_target_status"`
	SignStatusDelay  int    `json:"signStatusDelay" gorm:"column:sign_status_delay"`

	DeductTargetStatus string `json:"deductTargetStatus" gorm:"column:deduct_target_status"`
	DeductStatusDelay  int    `json:"deductStatusDelay" gorm:"column:deduct_status_delay"`

	TerminateNotifyEnabled bool `json:"terminateNotifyEnabled" gorm:"column:terminate_notify_enabled"`
	SignDurationMinutes    int  `json:"signDurationMinutes" gorm:"column:sign_duration_minutes"`
	StrictDeductRule       bool `json:"strictDeductRule" gorm:"column:strict_deduct_rule"`

	Active bool `json:"active" gorm:"column:active"`
}

// ContractInfo 签约信息详情
type ContractInfo struct {
	global.GVA_MODEL

	MerchantID uint   `json:"merchantId" gorm:"column:merchant_id"`
	OpenID     string `json:"openId" gorm:"column:open_id"`
	OutUserID  string `json:"outUserId" gorm:"column:out_user_id"`

	OutContractID string `json:"outContractId" gorm:"column:out_contract_id"`
	ContractID    string `json:"contractId" gorm:"column:contract_id"`
	SignSerialNo  string `json:"signSerialNo" gorm:"column:sign_serial_no"`
	PlanID        string `json:"planId" gorm:"column:plan_id"`

	NotifyURL string `json:"notifyUrl" gorm:"column:notify_url"`

	ContractStatus string `json:"contractStatus" gorm:"column:contract_status"`
	TerminateType  string `json:"terminateType" gorm:"column:terminate_type"`

	ExpireTime time.Time `json:"expireTime" gorm:"column:expire_time"`

	IsFirstDeduct     bool       `json:"isFirstDeduct" gorm:"column:is_first_deduct"`
	LastPreNotifyTime *time.Time `json:"lastPreNotifyTime" gorm:"column:last_pre_notify_time"`
	PreNotifyCalled   bool       `json:"preNotifyCalled" gorm:"column:pre_notify_called"`

	RequestData string `json:"requestData" gorm:"column:request_data"`
}

// DeductRecordInfo 扣款记录详情
type DeductRecordInfo struct {
	global.GVA_MODEL

	ContractID      uint   `json:"contractId" gorm:"column:contract_id"`
	MerchantID      uint   `json:"merchantId" gorm:"column:merchant_id"`
	OperationType   string `json:"operationType" gorm:"column:operation_type"`
	RequestData     string `json:"requestData" gorm:"column:request_data"`
	CallbackURL     string `json:"callbackUrl" gorm:"column:callback_url"`
	TransactionID   string `json:"transactionId" gorm:"column:transaction_id"`
	Amount          int64  `json:"amount" gorm:"column:amount"`
	ResponseData    string `json:"responseData" gorm:"column:response_data"`
	Status          string `json:"status" gorm:"column:status"`
	IsFirstDeduct   bool   `json:"isFirstDeduct" gorm:"column:is_first_deduct"`
	PreNotifyCalled bool   `json:"preNotifyCalled" gorm:"column:pre_notify_called"`
	ErrorCode       string `json:"errorCode" gorm:"column:error_code"`
	ErrorMessage    string `json:"errorMessage" gorm:"column:error_message"`
	CallbackResult  string `json:"callbackResult" gorm:"column:callback_result"`
	CallbackTime    *int64 `json:"callbackTime" gorm:"column:callback_time"`
}

// ContractRecordInfo 签约记录详情
type ContractRecordInfo struct {
	global.GVA_MODEL

	ContractID     uint   `json:"contractId" gorm:"column:contract_id"`
	MerchantID     uint   `json:"merchantId" gorm:"column:merchant_id"`
	OperationType  string `json:"operationType" gorm:"column:operation_type"`
	RequestXML     string `json:"requestXml" gorm:"column:request_xml"`
	CallbackURL    string `json:"callbackUrl" gorm:"column:callback_url"`
	ResponseXML    string `json:"responseXml" gorm:"column:response_xml"`
	Status         string `json:"status" gorm:"column:status"`
	ErrorCode      string `json:"errorCode" gorm:"column:error_code"`
	ErrorMessage   string `json:"errorMessage" gorm:"column:error_message"`
	CallbackResult string `json:"callbackResult" gorm:"column:callback_result"`
	CallbackTime   *int64 `json:"callbackTime" gorm:"column:callback_time"`
}
