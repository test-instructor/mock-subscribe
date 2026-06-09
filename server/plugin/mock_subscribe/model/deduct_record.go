package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// DeductRecord 扣款/预扣费通知流水
type DeductRecord struct {
	global.GVA_MODEL

	// 关联
	ContractID uint `json:"contractId" form:"contractId" gorm:"column:contract_id;comment:签约ID;index"`
	MerchantID uint `json:"merchantId" form:"merchantId" gorm:"column:merchant_id;comment:商户ID;index"`

	// 操作类型
	OperationType string `json:"operationType" form:"operationType" gorm:"column:operation_type;comment:操作类型(pre_notify/deduct)"`

	// 请求信息
	RequestData   string `json:"requestData" form:"requestData" gorm:"column:request_data;type:text;comment:请求数据(XML或JSON)"`
	CallbackURL   string `json:"callbackUrl" form:"callbackUrl" gorm:"column:callback_url;comment:回调地址"`
	TransactionID string `json:"transactionId" form:"transactionId" gorm:"column:transaction_id;comment:微信交易单号"`
	OutTradeNo    string `json:"outTradeNo" form:"outTradeNo" gorm:"column:out_trade_no;comment:商户订单号"`

	// 扣款金额
	Amount int64 `json:"amount" form:"amount" gorm:"column:amount;comment:扣款金额(分)"`

	// 响应信息
	ResponseData string `json:"responseData" form:"responseData" gorm:"column:response_data;type:text;comment:响应数据"`

	// 业务状态
	Status string `json:"status" form:"status" gorm:"column:status;comment:业务状态"`

	// 是否首次扣款
	IsFirstDeduct bool `json:"isFirstDeduct" form:"isFirstDeduct" gorm:"column:is_first_deduct;default:false;comment:是否首次扣款"`

	// 是否已预扣费通知
	PreNotifyCalled bool `json:"preNotifyCalled" form:"preNotifyCalled" gorm:"column:pre_notify_called;default:false;comment:是否已调用预扣费通知"`

	// 错误信息
	ErrorCode    string `json:"errorCode" form:"errorCode" gorm:"column:error_code;comment:错误码"`
	ErrorMessage string `json:"errorMessage" form:"errorMessage" gorm:"column:error_message;comment:错误信息"`

	// 回调结果
	CallbackResult string `json:"callbackResult" form:"callbackResult" gorm:"column:callback_result;type:text;comment:回调结果"`
	CallbackTime   *int64 `json:"callbackTime" form:"callbackTime" gorm:"column:callback_time;comment:回调时间戳"`
}

func (DeductRecord) TableName() string {
	return "gva_mock_deduct_record"
}
