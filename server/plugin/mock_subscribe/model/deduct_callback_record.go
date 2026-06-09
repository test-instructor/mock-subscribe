package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// DeductCallbackRecord 代扣回调接收记录
type DeductCallbackRecord struct {
	global.GVA_MODEL

	MerchantID        uint   `json:"merchantId" form:"merchantId" gorm:"column:merchant_id;comment:商户ID;index"`
	ContractIDRef     uint   `json:"contractIdRef" form:"contractIdRef" gorm:"column:contract_id_ref;comment:关联签约主表ID;index"`
	DeductRecordIDRef uint   `json:"deductRecordIdRef" form:"deductRecordIdRef" gorm:"column:deduct_record_id_ref;comment:关联扣款记录ID;index"`
	MchID             string `json:"mchId" form:"mchId" gorm:"column:mch_id;comment:商户号;index"`
	OutTradeNo        string `json:"outTradeNo" form:"outTradeNo" gorm:"column:out_trade_no;comment:商户订单号;index"`
	TransactionID     string `json:"transactionId" form:"transactionId" gorm:"column:transaction_id;comment:微信交易单号;index"`
	TradeType         string `json:"tradeType" form:"tradeType" gorm:"column:trade_type;comment:交易类型"`
	TradeState        string `json:"tradeState" form:"tradeState" gorm:"column:trade_state;comment:交易状态;index"`
	BankType          string `json:"bankType" form:"bankType" gorm:"column:bank_type;comment:银行类型"`
	TotalAmount       int64  `json:"totalAmount" form:"totalAmount" gorm:"column:total_amount;comment:总金额(分)"`
	CashAmount        int64  `json:"cashAmount" form:"cashAmount" gorm:"column:cash_amount;comment:现金金额(分)"`
	TimeEnd           string `json:"timeEnd" form:"timeEnd" gorm:"column:time_end;comment:支付完成时间"`
	SourceIP          string `json:"sourceIp" form:"sourceIp" gorm:"column:source_ip;comment:来源IP"`
	Headers           string `json:"headers" form:"headers" gorm:"column:headers;type:text;comment:请求头摘要"`
	RawBody           string `json:"rawBody" form:"rawBody" gorm:"column:raw_body;type:text;comment:原始请求体"`
	Sign              string `json:"sign" form:"sign" gorm:"column:sign;comment:签名"`
	SignValid         bool   `json:"signValid" form:"signValid" gorm:"column:sign_valid;comment:签名是否有效"`
	SignErrorMessage  string `json:"signErrorMessage" form:"signErrorMessage" gorm:"column:sign_error_message;comment:签名错误信息"`
	AckXML            string `json:"ackXml" form:"ackXml" gorm:"column:ack_xml;type:text;comment:响应ACK"`
}

func (DeductCallbackRecord) TableName() string {
	return "gva_mock_deduct_callback_record"
}
