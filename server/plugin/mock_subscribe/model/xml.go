package model

// ====================== APP纯签约请求/响应 ======================

// SignContractRequest APP纯签约请求(XML)
type SignContractRequest struct {
	AppID                  string `xml:"appid"`                    // 应用ID
	MchID                  string `xml:"mch_id"`                   // 商户号
	ContractAppID          string `xml:"plan_id"`                  // 签约模板ID
	OutContractCode        string `xml:"out_contract_code"`        // 用户侧签约协议号
	OutUserID              string `xml:"outer_openid"`             // 用户OpenID
	ContractDisplayAccount string `xml:"contract_display_account"` // 用户账户展示名称
	NotifyURL              string `xml:"notify_url"`               // 回调地址
	SignType               string `xml:"sign_type"`                // 签名类型
	Version                string `xml:"version"`                  // 版本号
	TimeStamp              string `xml:"timestamp"`                // 时间戳
	Nonce                  string `xml:"nonce"`                    // 随机字符串
	Sign                   string `xml:"sign"`                     // 签名
}

// SignContractResponse APP纯签约响应(XML)
type SignContractResponse struct {
	ReturnCode      string `xml:"return_code"`       // 返回状态码
	ReturnMsg       string `xml:"return_msg"`        // 返回信息
	ResultCode      string `xml:"result_code"`       // 业务结果码
	ErrCode         string `xml:"err_code"`          // 错误码
	ErrCodeDes      string `xml:"err_code_des"`      // 错误码描述
	ContractID      string `xml:"contract_id"`       // 协议号
	ContractExtID   string `xml:"contract_ext_id"`   // 签约扩展ID
	OperationType   string `xml:"operation_type"`    // 操作类型
	MchID           string `xml:"mch_id"`            // 商户号
	OutContractCode string `xml:"out_contract_code"` // 用户侧签约协议号
	SignType        string `xml:"sign_type"`         // 签名类型
	TimeStamp       string `xml:"timestamp"`         // 时间戳
	Nonce           string `xml:"nonce"`             // 随机字符串
	Sign            string `xml:"sign"`              // 签名
}

// ====================== 申请扣款请求/响应 ======================

// DeductApplyRequest 申请扣款请求(XML)
type DeductApplyRequest struct {
	AppID         string `xml:"appid"`          // 应用ID
	MchID         string `xml:"mch_id"`         // 商户号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	ContractID    string `xml:"contract_id"`    // 协议号
	TransactionID string `xml:"transaction_id"` // 微信交易单号
	TotalAmount   int64  `xml:"total_amount"`   // 扣款金额
	Currency      string `xml:"fee_type"`       // 货币类型
	NotifyURL     string `xml:"notify_url"`     // 回调地址
	SignType      string `xml:"sign_type"`      // 签名类型
	TimeStamp     string `xml:"timestamp"`      // 时间戳
	Nonce         string `xml:"nonce"`          // 随机字符串
	Sign          string `xml:"sign"`           // 签名
}

// DeductApplyResponse 申请扣款响应(XML)
type DeductApplyResponse struct {
	ReturnCode    string `xml:"return_code"`    // 返回状态码
	ReturnMsg     string `xml:"return_msg"`     // 返回信息
	ResultCode    string `xml:"result_code"`    // 业务结果码
	ErrCode       string `xml:"err_code"`       // 错误码
	ErrCodeDes    string `xml:"err_code_des"`   // 错误码描述
	MchID         string `xml:"mch_id"`         // 商户号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	TransactionID string `xml:"transaction_id"` // 微信交易单号
	Amount        int64  `xml:"amount"`         // 扣款金额
	SignType      string `xml:"sign_type"`      // 签名类型
	TimeStamp     string `xml:"timestamp"`      // 时间戳
	Nonce         string `xml:"nonce"`          // 随机字符串
	Sign          string `xml:"sign"`           // 签名
}

// ====================== 扣款结果通知 ======================

// DeductNotify 扣款结果通知(XML)
type DeductNotify struct {
	ReturnCode    string `xml:"return_code"`    // 返回状态码
	ReturnMsg     string `xml:"return_msg"`     // 返回信息
	AppID         string `xml:"appid"`          // 应用ID
	MchID         string `xml:"mch_id"`         // 商户号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	TransactionID string `xml:"transaction_id"` // 微信交易单号
	TradeType     string `xml:"trade_type"`     // 交易类型
	TradeState    string `xml:"trade_state"`    // 交易状态
	BankType      string `xml:"bank_type"`      // 银行类型
	TotalAmount   int64  `xml:"total_amount"`   // 总金额
	CashAmount    int64  `xml:"cash_amount"`    // 现金金额
	TimeStamp     string `xml:"timestamp"`      // 时间戳
	Nonce         string `xml:"nonce"`          // 随机字符串
	Sign          string `xml:"sign"`           // 签名
}

// ====================== 查询签约关系请求/响应 ======================

// QueryContractRequest 查询签约关系请求(XML)
type QueryContractRequest struct {
	AppID           string `xml:"appid"`
	MchID           string `xml:"mch_id"`
	ContractID      string `xml:"contract_id"`
	OutContractCode string `xml:"out_contract_code"`
	SignType        string `xml:"sign_type"`
	TimeStamp       string `xml:"timestamp"`
	Nonce           string `xml:"nonce"`
	Sign            string `xml:"sign"`
}

// QueryContractResponse 查询签约关系响应(XML)
type QueryContractResponse struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	ResultCode     string `xml:"result_code"`
	ErrCode        string `xml:"err_code"`
	ErrCodeDes     string `xml:"err_code_des"`
	ContractID     string `xml:"contract_id"`
	ContractStatus string `xml:"contract_status"`
	ContractExt    string `xml:"contract_ext"`
	PlanID         string `xml:"plan_id"`
	SignStatus     string `xml:"sign_status"`
	SignType       string `xml:"sign_type"`
	TimeStamp      string `xml:"timestamp"`
	Nonce          string `xml:"nonce"`
	Sign           string `xml:"sign"`
}

// ====================== 申请解约请求/响应 ======================

// TerminateContractRequest 申请解约请求(XML)
type TerminateContractRequest struct {
	AppID              string `xml:"appid"`
	MchID              string `xml:"mch_id"`
	ContractID         string `xml:"contract_id"`
	OutContractCode    string `xml:"out_contract_code"`
	ContractStatus     string `xml:"contract_status"`
	ContractEndingType string `xml:"contract_ending_type"`
	SignType           string `xml:"sign_type"`
	TimeStamp          string `xml:"timestamp"`
	Nonce              string `xml:"nonce"`
	Sign               string `xml:"sign"`
}

// TerminateContractResponse 申请解约响应(XML)
type TerminateContractResponse struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	ResultCode     string `xml:"result_code"`
	ErrCode        string `xml:"err_code"`
	ErrCodeDes     string `xml:"err_code_des"`
	ContractID     string `xml:"contract_id"`
	ContractStatus string `xml:"contract_status"`
	SignType       string `xml:"sign_type"`
	TimeStamp      string `xml:"timestamp"`
	Nonce          string `xml:"nonce"`
	Sign           string `xml:"sign"`
}

// ====================== 通用ACK ======================

// ContractCallbackRequest 签约回调请求(XML)
type ContractCallbackRequest struct {
	AppID              string `xml:"appid"`
	MchID              string `xml:"mch_id"`
	ContractID         string `xml:"contract_id"`
	OutContractCode    string `xml:"out_contract_code"`
	ContractStatus     string `xml:"contract_status"`
	ContractEndingType string `xml:"contract_ending_type"`
	ContractExtID      string `xml:"contract_ext_id"`
	PlanID             string `xml:"plan_id"`
	OpenID             string `xml:"openid"`
	SignType           string `xml:"sign_type"`
	TimeStamp          string `xml:"timestamp"`
	Nonce              string `xml:"nonce"`
	Sign               string `xml:"sign"`
}

// GenericACK 通用成功/失败ACK
type GenericACK struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// ====================== 预扣费通知API(JSON) ======================

// PreDeductNotifyRequest 预扣费通知请求(JSON)
type PreDeductNotifyRequest struct {
	AppID         string `json:"appid"`
	MchID         string `json:"mch_id"`
	ContractID    string `json:"contract_id"`
	OutTradeNo    string `json:"out_trade_no"`
	TradeNo       string `json:"trade_no"`
	ActionType    int    `json:"action_type"`
	AccountID     string `json:"account_id"`
	NotifyURL     string `json:"notify_url"`
	RequestSerial int64  `json:"request_serial"`
	SignType      string `json:"sign_type"`
	TimeStamp     string `json:"timestamp"`
	Nonce         string `json:"nonce"`
	Sign          string `json:"sign"`
}

// PreDeductNotifyResponse 预扣费通知响应(JSON)
type PreDeductNotifyResponse struct {
	ReturnCode string `json:"return_code"`
	ReturnMsg  string `json:"return_msg"`
	ResultCode string `json:"result_code"`
	ErrCode    string `json:"err_code"`
	ErrCodeDes string `json:"err_code_des"`
	AppID      string `json:"appid"`
	MchID      string `json:"mch_id"`
	SignType   string `json:"sign_type"`
	TimeStamp  string `json:"timestamp"`
	Nonce      string `json:"nonce"`
	Sign       string `json:"sign"`
}
