package model

import (
	"bytes"
	"encoding/xml"
	"strings"
)

// ====================== APP纯签约请求/响应 ======================

// SignContractRequest APP纯签约请求(XML)
type SignContractRequest struct {
	XMLName                xml.Name `xml:"xml"`
	AppID                  string   `xml:"appid"`                    // 应用ID
	MchID                  string   `xml:"mch_id"`                   // 商户号
	PlanID                 string   `xml:"plan_id"`                  // 签约模板ID
	ContractCode           string   `xml:"contract_code"`            // 用户侧签约协议号（官方字段）
	OutContractCode        string   `xml:"out_contract_code"`        // 用户侧签约协议号（兼容旧字段）
	OpenID                 string   `xml:"openid"`                   // 用户OpenID（官方字段）
	OutUserID              string   `xml:"outer_openid"`             // 用户OpenID（兼容旧字段）
	ContractDisplayAccount string   `xml:"contract_display_account"` // 用户账户展示名称
	NotifyURL              string   `xml:"notify_url"`               // 回调地址
	SignType               string   `xml:"sign_type"`                // 签名类型
	Version                string   `xml:"version"`                  // 版本号
	TimeStamp              string   `xml:"timestamp"`                // 时间戳
	Nonce                  string   `xml:"nonce"`                    // 随机字符串
	Sign                   string   `xml:"sign"`                     // 签名
}

// SignContractResponse APP纯签约响应(XML)
type SignContractResponse struct {
	XMLName         xml.Name `xml:"xml"`
	ReturnCode      string   `xml:"return_code"`     // 返回状态码
	ReturnMsg       string   `xml:"return_msg"`      // 返回信息
	ResultCode      string   `xml:"result_code"`     // 业务结果码
	ErrCode         string   `xml:"err_code"`        // 错误码
	ErrCodeDes      string   `xml:"err_code_des"`    // 错误码描述
	ContractID      string   `xml:"contract_id"`     // 协议号
	ContractExtID   string   `xml:"contract_ext_id"` // 签约扩展ID
	OperationType   string   `xml:"operation_type"`  // 操作类型
	MchID           string   `xml:"mch_id"`          // 商户号
	OutContractCode string   `xml:"contract_code"`   // 用户侧签约协议号
	SignType        string   `xml:"sign_type"`       // 签名类型
	TimeStamp       string   `xml:"timestamp"`       // 时间戳
	Nonce           string   `xml:"nonce"`           // 随机字符串
	Sign            string   `xml:"sign"`            // 签名
}

// ====================== 申请扣款请求/响应 ======================

// DeductApplyRequest 申请扣款请求(XML)
type DeductApplyRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid"`          // 应用ID
	MchID         string   `xml:"mch_id"`         // 商户号
	Body          string   `xml:"body"`           // 商品描述
	Detail        string   `xml:"detail"`         // 商品详情
	Attach        string   `xml:"attach"`         // 附加数据
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	ContractID    string   `xml:"contract_id"`    // 协议号
	TransactionID string   `xml:"transaction_id"` // 微信交易单号
	TotalFee      int64    `xml:"total_fee"`      // 扣款金额（官方V2字段）
	TotalAmount   int64    `xml:"total_amount"`   // 扣款金额（兼容字段）
	FeeType       string   `xml:"fee_type"`       // 货币类型
	NotifyURL     string   `xml:"notify_url"`     // 回调地址
	TradeType     string   `xml:"trade_type"`     // 交易类型
	DeviceInfo    string   `xml:"device_info"`    // 设备号
	NonceStr      string   `xml:"nonce_str"`      // 随机字符串（官方V2字段）
	SignType      string   `xml:"sign_type"`      // 签名类型
	TimeStamp     string   `xml:"timestamp"`      // 时间戳
	Nonce         string   `xml:"nonce"`          // 随机字符串（兼容字段）
	Sign          string   `xml:"sign"`           // 签名
}

// DeductApplyResponse 申请扣款响应(XML)
type DeductApplyResponse struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`    // 返回状态码
	ReturnMsg     string   `xml:"return_msg"`     // 返回信息
	ResultCode    string   `xml:"result_code"`    // 业务结果码
	ErrCode       string   `xml:"err_code"`       // 错误码
	ErrCodeDes    string   `xml:"err_code_des"`   // 错误码描述
	MchID         string   `xml:"mch_id"`         // 商户号
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	TransactionID string   `xml:"transaction_id"` // 微信交易单号
	TotalAmount   int64    `xml:"total_amount"`   // 扣款金额
	SignType      string   `xml:"sign_type"`      // 签名类型
	TimeStamp     string   `xml:"timestamp"`      // 时间戳
	Nonce         string   `xml:"nonce"`          // 随机字符串
	Sign          string   `xml:"sign"`           // 签名
}

// QueryDeductRequest 查询订单请求(XML)
type QueryDeductRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppID         string   `xml:"appid"`
	MchID         string   `xml:"mch_id"`
	OutTradeNo    string   `xml:"out_trade_no"`
	TransactionID string   `xml:"transaction_id"`
	SignType      string   `xml:"sign_type"`
	TimeStamp     string   `xml:"timestamp"`
	NonceStr      string   `xml:"nonce_str"`
	Nonce         string   `xml:"nonce"`
	Sign          string   `xml:"sign"`
}

// QueryDeductResponse 查询订单响应(XML)
type QueryDeductResponse struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`
	ReturnMsg     string   `xml:"return_msg"`
	ResultCode    string   `xml:"result_code"`
	ErrCode       string   `xml:"err_code"`
	ErrCodeDes    string   `xml:"err_code_des"`
	AppID         string   `xml:"appid"`
	MchID         string   `xml:"mch_id"`
	OpenID        string   `xml:"openid"`
	TradeType     string   `xml:"trade_type"`
	TradeState    string   `xml:"trade_state"`
	BankType      string   `xml:"bank_type"`
	TotalAmount   int64    `xml:"total_amount"`
	CashAmount    int64    `xml:"cash_amount"`
	TransactionID string   `xml:"transaction_id"`
	OutTradeNo    string   `xml:"out_trade_no"`
	TimeEnd       string   `xml:"time_end"`
	SignType      string   `xml:"sign_type"`
	TimeStamp     string   `xml:"timestamp"`
	Nonce         string   `xml:"nonce"`
	Sign          string   `xml:"sign"`
}

// ====================== 扣款结果通知 ======================

// DeductNotifyRequest 扣款结果通知请求(XML)
type DeductNotifyRequest struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`    // 返回状态码
	ReturnMsg     string   `xml:"return_msg"`     // 返回信息
	AppID         string   `xml:"appid"`          // 应用ID
	MchID         string   `xml:"mch_id"`         // 商户号
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	TransactionID string   `xml:"transaction_id"` // 微信交易单号
	TradeType     string   `xml:"trade_type"`     // 交易类型
	TradeState    string   `xml:"trade_state"`    // 交易状态
	BankType      string   `xml:"bank_type"`      // 银行类型
	TotalAmount   int64    `xml:"total_amount"`   // 总金额
	CashAmount    int64    `xml:"cash_amount"`    // 现金金额
	TimeEnd       string   `xml:"time_end"`       // 支付完成时间
	SignType      string   `xml:"sign_type"`      // 签名类型
	TimeStamp     string   `xml:"timestamp"`      // 时间戳
	Nonce         string   `xml:"nonce"`          // 随机字符串
	Sign          string   `xml:"sign"`           // 签名
}

// DeductNotifyResponse 扣款结果通知响应(XML)
type DeductNotifyResponse struct {
	XMLName       xml.Name `xml:"xml"`
	ReturnCode    string   `xml:"return_code"`    // 返回状态码
	ReturnMsg     string   `xml:"return_msg"`     // 返回信息
	AppID         string   `xml:"appid"`          // 应用ID
	MchID         string   `xml:"mch_id"`         // 商户号
	OutTradeNo    string   `xml:"out_trade_no"`   // 商户订单号
	TransactionID string   `xml:"transaction_id"` // 微信交易单号
	TradeType     string   `xml:"trade_type"`     // 交易类型
	TradeState    string   `xml:"trade_state"`    // 交易状态
	BankType      string   `xml:"bank_type"`      // 银行类型
	TotalAmount   int64    `xml:"total_amount"`   // 总金额
	CashAmount    int64    `xml:"cash_amount"`    // 现金金额
	TimeEnd       string   `xml:"time_end"`       // 支付完成时间
	SignType      string   `xml:"sign_type"`      // 签名类型
	TimeStamp     string   `xml:"timestamp"`      // 时间戳
	Nonce         string   `xml:"nonce"`          // 随机字符串
	Sign          string   `xml:"sign"`           // 签名
}

// ====================== 查询签约关系请求/响应 ======================

// QueryContractRequest 查询签约关系请求(XML)
type QueryContractRequest struct {
	XMLName         xml.Name `xml:"xml"`
	AppID           string   `xml:"appid"`
	MchID           string   `xml:"mch_id"`
	ContractID      string   `xml:"contract_id"`
	PlanID          string   `xml:"plan_id"`
	ContractCode    string   `xml:"contract_code"`
	OutContractCode string   `xml:"out_contract_code"`
	SignType        string   `xml:"sign_type"`
	TimeStamp       string   `xml:"timestamp"`
	Nonce           string   `xml:"nonce"`
	Sign            string   `xml:"sign"`
}

// QueryContractResponse 查询签约关系响应(XML)
type QueryContractResponse struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
	ResultCode     string   `xml:"result_code"`
	ErrCode        string   `xml:"err_code"`
	ErrCodeDes     string   `xml:"err_code_des"`
	ContractID     string   `xml:"contract_id"`
	ContractStatus string   `xml:"contract_status"`
	ContractExt    string   `xml:"contract_ext"`
	PlanID         string   `xml:"plan_id"`
	SignStatus     string   `xml:"sign_status"`
	SignType       string   `xml:"sign_type"`
	TimeStamp      string   `xml:"timestamp"`
	Nonce          string   `xml:"nonce"`
	Sign           string   `xml:"sign"`
}

// ====================== 申请解约请求/响应 ======================

// TerminateContractRequest 申请解约请求(XML)
type TerminateContractRequest struct {
	XMLName                   xml.Name `xml:"xml"`
	AppID                     string   `xml:"appid"`
	MchID                     string   `xml:"mch_id"`
	ContractID                string   `xml:"contract_id"`
	PlanID                    string   `xml:"plan_id"`
	ContractCode              string   `xml:"contract_code"`
	OutContractCode           string   `xml:"out_contract_code"`
	ContractTerminationRemark string   `xml:"contract_termination_remark"`
	ContractStatus            string   `xml:"contract_status"`
	ContractEndingType        string   `xml:"contract_ending_type"`
	Version                   string   `xml:"version"`
	SignType                  string   `xml:"sign_type"`
	TimeStamp                 string   `xml:"timestamp"`
	Nonce                     string   `xml:"nonce"`
	Sign                      string   `xml:"sign"`
}

// TerminateContractResponse 申请解约响应(XML)
type TerminateContractResponse struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
	ResultCode     string   `xml:"result_code"`
	ErrCode        string   `xml:"err_code"`
	ErrCodeDes     string   `xml:"err_code_des"`
	ContractID     string   `xml:"contract_id"`
	ContractStatus string   `xml:"contract_status"`
	SignType       string   `xml:"sign_type"`
	TimeStamp      string   `xml:"timestamp"`
	Nonce          string   `xml:"nonce"`
	Sign           string   `xml:"sign"`
}

// ====================== 通用ACK ======================

// ContractCallbackRequest 签约/解约回调请求(XML)
type ContractCallbackRequest struct {
	XMLName         xml.Name `xml:"xml"`
	ReturnCode      string   `xml:"return_code"`
	ResultCode      string   `xml:"result_code"`
	Sign            string   `xml:"sign"`
	MchID           string   `xml:"mch_id"`
	OutContractCode string   `xml:"contract_code"`
	OpenID          string   `xml:"openid"`
	PlanID          string   `xml:"plan_id"`
	ChangeType      string   `xml:"change_type"`
	OperateTime     string   `xml:"operate_time"`
	ContractID      string   `xml:"contract_id"`
}

// ContractResultNotify 签约结果通知(XML)
type ContractResultNotify struct {
	XMLName         xml.Name `xml:"xml"`
	ReturnCode      string   `xml:"return_code"`
	ResultCode      string   `xml:"result_code"`
	Sign            string   `xml:"sign"`
	MchID           string   `xml:"mch_id"`
	OutContractCode string   `xml:"contract_code"`
	OpenID          string   `xml:"openid"`
	PlanID          string   `xml:"plan_id"`
	ChangeType      string   `xml:"change_type"`
	OperateTime     string   `xml:"operate_time"`
	ContractID      string   `xml:"contract_id"`
}

// GenericACK 通用成功/失败ACK
type GenericACK struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

// SignContractResponseV2 APP纯签约响应(XML) - 严格按微信字段顺序输出
// 用于 /papay/preentrustweb 成功响应，必须与微信官方字段顺序完全一致
type SignContractResponseV2 struct {
	ReturnCode          string `xml:"return_code"`
	ReturnMsg           string `xml:"return_msg"`
	ResultCode          string `xml:"result_code"`
	AppID               string `xml:"appid"`
	MchID               string `xml:"mch_id"`
	MiniprogramUsername string `xml:"miniprogram_username"`
	MiniprogramPath     string `xml:"miniprogram_path"`
	NonceStr            string `xml:"nonce_str"`
	Sign                string `xml:"sign"`
	PreEntrustwebID     string `xml:"pre_entrustweb_id"`
}

// ToXMLBytes 将响应结构体序列化为严格按字段顺序的 XML 字符串，使用 CDATA 包裹值
func (r SignContractResponseV2) ToXMLBytes() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(`<xml>`)
	b.WriteString(`<return_code><![CDATA[`)
	b.WriteString(escape(r.ReturnCode))
	b.WriteString(`]]></return_code>`)
	b.WriteString(`<return_msg><![CDATA[`)
	b.WriteString(escape(r.ReturnMsg))
	b.WriteString(`]]></return_msg>`)
	b.WriteString(`<result_code><![CDATA[`)
	b.WriteString(escape(r.ResultCode))
	b.WriteString(`]]></result_code>`)
	b.WriteString(`<appid><![CDATA[`)
	b.WriteString(escape(r.AppID))
	b.WriteString(`]]></appid>`)
	b.WriteString(`<mch_id><![CDATA[`)
	b.WriteString(escape(r.MchID))
	b.WriteString(`]]></mch_id>`)
	b.WriteString(`<miniprogram_username><![CDATA[`)
	b.WriteString(escape(r.MiniprogramUsername))
	b.WriteString(`]]></miniprogram_username>`)
	b.WriteString(`<miniprogram_path><![CDATA[`)
	b.WriteString(escape(r.MiniprogramPath))
	b.WriteString(`]]></miniprogram_path>`)
	b.WriteString(`<nonce_str><![CDATA[`)
	b.WriteString(escape(r.NonceStr))
	b.WriteString(`]]></nonce_str>`)
	b.WriteString(`<sign><![CDATA[`)
	b.WriteString(escape(r.Sign))
	b.WriteString(`]]></sign>`)
	b.WriteString(`<pre_entrustweb_id><![CDATA[`)
	b.WriteString(escape(r.PreEntrustwebID))
	b.WriteString(`]]></pre_entrustweb_id>`)
	b.WriteString(`</xml>`)
	return b.Bytes(), nil
}

// escape 对 XML 特殊字符进行转义，防止注入
func escape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, `'`, "&apos;")
	return s
}

// ====================== 预扣费通知API(JSON) ======================

// PreDeductNotifyRequest 预扣费通知请求(JSON)
type PreDeductNotifyRequest struct {
	MchID string `json:"mchid"`
	AppID string `json:"appid"`

	DeductDuration struct {
		Count int    `json:"count"`
		Unit  string `json:"unit"`
	} `json:"deduct_duration"`

	EstimatedAmount struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	} `json:"estimated_amount"`
}

// PreDeductNotifyResponse 预扣费通知响应(JSON)
type PreDeductNotifyResponse struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	ResultCode string `json:"result_code"`
	ErrCode    string `json:"err_code,omitempty"`
	ErrCodeDes string `json:"err_code_des,omitempty"`
}
