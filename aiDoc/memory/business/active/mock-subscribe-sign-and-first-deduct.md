# 用户业务记忆

- 需求：`mockSubscribe/mockSubscribeContract` 用户协议页面需要修正详情弹窗的时间展示与详情查询参数。
- 细节：列表时间字段需要稳定格式化，避免前端展示 `NaN-aN-aN aN:aN:aN`。
- 细节：调用 `/mockSubscribeContract/findContract` 时需要传入协议 `id`，确保详情回填正常。
- 需求：微信订阅 Mock 的签约/解约回调、订单查询接口需要按微信文档对齐字段与验签规则。
- 细节：签约回调按“签约、解约结果通知”字段集验签，并区分签约回调与解约回调类型。
- 细节：`/transit/queryorder` 按“查询订单”返回订单查询字段，支持 `out_trade_no` 与 `transaction_id` 查询。
- 需求：补齐“查询签约关系”链路，并让 Python 客户端在查询到已签约后自动发起申请扣款，再继续查询订单结果。
- 细节：`/papay/querycontract` 需要稳定支持按 `contract_id` 或 `out_contract_code` 查询签约关系，并返回签约状态供客户端判定。
- 细节：Python 脚本需要新增独立的签约查询、申请扣款、订单查询命令，并在扣款成功或失败、订单成功或失败时打印明确日志。
- 需求：`/papay/preentrustweb` 签约回调 XML 的最外层标签需要与微信文档保持一致，使用 `xml` 而不是 `ContractResultNotify`。
- 细节：签约结果通知示例要求回调报文根节点为 `<xml>...</xml>`，避免因根标签错误导致下游解析失败。
- 需求：`/papay/preentrustweb` 接口需要立即返回，不能被“签约状态延时”阻塞。
- 细节：接口调用后直接构造响应 XML 并返回 HTTP，状态写入、合同号/有效期设置、回调发送等所有延后动作统一放进 `go func()` 协程中执行。
- 细节：协程内执行顺序为：`time.Sleep(SignStatusDelay)` → `UpdateContractStatus(SignTargetStatus)` → 若是 `ACTIVE` 则 `SetContractID`/`SetExpireTime` → 必要时再 `time.Sleep(SignCallbackDelay)` 后发送签约回调。
- 细节：响应中的 `contract_id`、`sign_serial_no`、`pre_entrustweb_id` 等需要提前在主流程生成并参与签名，保证立即返回的响应和协程内写库使用的值一致。
