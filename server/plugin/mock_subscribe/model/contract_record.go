package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ContractRecord 签约/查询/解约状态流水
type ContractRecord struct {
	global.GVA_MODEL

	// 关联
	ContractID uint `json:"contractId" form:"contractId" gorm:"column:contract_id;comment:签约ID;index"`
	MerchantID uint `json:"merchantId" form:"merchantId" gorm:"column:merchant_id;comment:商户ID;index"`

	// 操作类型
	OperationType string `json:"operationType" form:"operationType" gorm:"column:operation_type;comment:操作类型(sign/query/terminate)"`

	// 请求信息
	RequestXML  string `json:"requestXml" form:"requestXml" gorm:"column:request_xml;type:text;comment:请求XML"`
	CallbackURL string `json:"callbackUrl" form:"callbackUrl" gorm:"column:callback_url;comment:回调地址"`

	// 响应信息
	ResponseXML string `json:"responseXml" form:"responseXml" gorm:"column:response_xml;type:text;comment:响应XML"`

	// 状态
	Status string `json:"status" form:"status" gorm:"column:status;comment:当前状态"`

	// 错误信息
	ErrorCode    string `json:"errorCode" form:"errorCode" gorm:"column:error_code;comment:错误码"`
	ErrorMessage string `json:"errorMessage" form:"errorMessage" gorm:"column:error_message;comment:错误信息"`

	// 回调结果
	CallbackResult string `json:"callbackResult" form:"callbackResult" gorm:"column:callback_result;type:text;comment:回调结果"`
	CallbackTime   *int64 `json:"callbackTime" form:"callbackTime" gorm:"column:callback_time;comment:回调时间戳"`
}

func (ContractRecord) TableName() string {
	return "gva_mock_contract_record"
}
