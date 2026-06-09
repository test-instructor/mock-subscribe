package model

const (
	// ContractStatus 签约状态
	ContractStatusPending    = "PENDING"    // 签约中
	ContractStatusActive     = "ACTIVE"     // 签约成功
	ContractStatusFailed     = "FAILED"     // 签约失败
	ContractStatusTerminated = "TERMINATED" // 已解约
	ContractStatusExpired    = "EXPIRED"    // 已到期
	ContractStatusPause      = "PAUSE"      // 已暂停

	// DeductStatus 扣款状态
	DeductStatusPending   = "PENDING"   // 扣款中
	DeductStatusSuccess   = "SUCCESS"   // 扣款成功
	DeductStatusFailed    = "FAILED"    // 扣款失败
	DeductStatusRefunding = "REFUNDING" // 退款中
	DeductStatusRefunded  = "REFUNDED"  // 已退款

	// TerminateType 解约方式
	TerminateTypeUserRequest     = "USER_REQUEST"     // 用户申请解约
	TerminateTypeMerchantRequest = "MERCHANT_REQUEST" // 商户申请解约
	TerminateTypeExpired         = "EXPIRED"          // 到期自动解约
	TerminateTypeSystemRevoke    = "SYSTEM_REVOKE"    // 系统撤销

	// CallbackType 回调类型
	CallbackTypeContractSign = "CONTRACT_SIGN" // 签约回调
	CallbackTypeDeductFirst  = "DEDUCT_FIRST"  // 首次扣款回调
	CallbackTypeDeductRepeat = "DEDUCT_REPEAT" // 非首次扣款回调
	CallbackTypeTerminate    = "TERMINATE"     // 解约回调

	// ErrorCode 错误码
	ErrCodeSuccess           = "SUCCESS"
	ErrCodeFail              = "FAIL"
	ErrCodeSignExists        = "SIGN_EXISTS"
	ErrCodeSignNotFound      = "SIGN_NOT_FOUND"
	ErrCodeInvalidParams     = "INVALID_PARAMS"
	ErrCodeInvalidSign       = "INVALID_SIGN"
	ErrCodeDeductNotAllowed  = "DEDUCT_NOT_ALLOWED"
	ErrCodePreNotifyRequired = "PRE_NOTIFY_REQUIRED"

	// XMLRoot XML根节点
	XMLRootSign              = "xml"
	XMLRootSignResponse      = "xml"
	XMLRootDeduct            = "xml"
	XMLRootDeductResponse    = "xml"
	XMLRootQuery             = "xml"
	XMLRootQueryResponse     = "xml"
	XMLRootTerminate         = "xml"
	XMLRootTerminateResponse = "xml"
	XMLRootNotify            = "xml"
	XMLRootPreNotify         = "paymsg"
)
