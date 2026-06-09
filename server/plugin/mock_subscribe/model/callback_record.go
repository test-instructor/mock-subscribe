package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// CallbackRecord 回调接收记录
type CallbackRecord struct {
	global.GVA_MODEL

	MerchantID       uint   `json:"merchantId" form:"merchantId" gorm:"column:merchant_id;comment:商户ID;index"`
	ContractIDRef    uint   `json:"contractIdRef" form:"contractIdRef" gorm:"column:contract_id_ref;comment:关联签约主表ID;index"`
	OutContractCode  string `json:"outContractCode" form:"outContractCode" gorm:"column:out_contract_code;comment:外部签约单号;index"`
	ContractCode     string `json:"contractCode" form:"contractCode" gorm:"column:contract_code;comment:协议号;index"`
	SignSerialNo     string `json:"signSerialNo" form:"signSerialNo" gorm:"column:sign_serial_no;comment:签约号"`
	MchID            string `json:"mchId" form:"mchId" gorm:"column:mch_id;comment:商户号;index"`
	CallbackType     string `json:"callbackType" form:"callbackType" gorm:"column:callback_type;comment:回调类型"`
	SourceIP         string `json:"sourceIp" form:"sourceIp" gorm:"column:source_ip;comment:来源IP"`
	Headers          string `json:"headers" form:"headers" gorm:"column:headers;type:text;comment:请求头摘要"`
	RawBody          string `json:"rawBody" form:"rawBody" gorm:"column:raw_body;type:text;comment:原始请求体"`
	ContractStatus   string `json:"contractStatus" form:"contractStatus" gorm:"column:contract_status;comment:协议状态"`
	TimeStamp        string `json:"timeStamp" form:"timeStamp" gorm:"column:time_stamp;comment:时间戳"`
	Nonce            string `json:"nonce" form:"nonce" gorm:"column:nonce;comment:随机串"`
	Sign             string `json:"sign" form:"sign" gorm:"column:sign;comment:签名"`
	SignValid        bool   `json:"signValid" form:"signValid" gorm:"column:sign_valid;comment:签名是否有效"`
	SignErrorMessage string `json:"signErrorMessage" form:"signErrorMessage" gorm:"column:sign_error_message;comment:签名错误信息"`
	AckXML           string `json:"ackXml" form:"ackXml" gorm:"column:ack_xml;type:text;comment:响应ACK"`
}

func (CallbackRecord) TableName() string {
	return "gva_mock_callback_record"
}
