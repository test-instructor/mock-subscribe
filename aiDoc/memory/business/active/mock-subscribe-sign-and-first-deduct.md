# 用户业务记忆

- 需求：`mockSubscribe/mockSubscribeContract` 用户协议页面需要修正详情弹窗的时间展示与详情查询参数。
- 细节：列表时间字段需要稳定格式化，避免前端展示 `NaN-aN-aN aN:aN:aN`。
- 细节：调用 `/mockSubscribeContract/findContract` 时需要传入协议 `id`，确保详情回填正常。
