# 用户业务记忆

- 需求：`mockSubscribe/mockSubscribeContract` 用户协议页面需要修正详情弹窗的时间展示与详情查询参数。
- 细节：列表时间字段需要稳定格式化，避免前端展示 `NaN-aN-aN aN:aN:aN`。
- 细节：调用 `/mockSubscribeContract/findContract` 时需要传入协议 `id`，确保详情回填正常。
- 需求：微信订阅 Mock 的签约/解约回调、订单查询接口需要按微信文档对齐字段与验签规则。
- 细节：签约回调按“签约、解约结果通知”字段集验签，并区分签约回调与解约回调类型。
- 细节：`/transit/queryorder` 按“查询订单”返回订单查询字段，支持 `out_trade_no` 与 `transaction_id` 查询。
