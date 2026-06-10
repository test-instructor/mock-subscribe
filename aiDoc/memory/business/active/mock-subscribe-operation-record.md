# mock_subscribe 接口接入操作历史（UserID=3）

## 基本信息

- 提出日期：2026-06-10
- 当前状态：`active`
- 需求类型：后端中间件 / 操作历史
- 优先级：中
- 需求文件：`aiDoc/memory/business/active/mock-subscribe-operation-record.md`

## 用户原始意图摘要

把 `mock_subscribe` 插件对外的几个管理类接口的请求 / 响应数据，按中间件方式写入到 `system.sys_operation_records` 操作历史表里，并把 `UserID` 统一记为 `3`。

## 影响范围

- 后端：
  - `server/middleware/operation.go`
  - `server/plugin/mock_subscribe/router/merchant.go`
  - `server/plugin/mock_subscribe/router/contract.go`
  - `server/plugin/mock_subscribe/router/deduct.go`
  - `server/plugin/mock_subscribe/router/deduct_callback.go`
- 前端：无
- 文档：`aiDoc/memory/business/demand-index.md`
- 插件 / 模块：`server/plugin/mock_subscribe`

## 涉及对象

- 模块：`mock_subscribe` 插件
- 接口：
  - `mockSubscribeMerchant/*`（5 个）
  - `mockSubscribeContract/*`（4 个）
  - `mockSubscribeDeduct/*`（2 个）
  - `mockSubscribeDeductCallback/*`（2 个，仅私有组）
  - 公开的 `papay/notify`、`pay/pappaynotify`、`papay/preentrustweb` 等微信侧 / 回调接口，本次不在范围内
- 页面：无
- 配置：无

## 已确认约束

- 写入方式必须是中间件，不在 API 层手动 `Create(&SysOperationRecord{...})`。
- `UserID` 固定为 `3`，不走 JWT / `x-user-id` 解析。
- 原有的 `OperationRecord()` 中间件行为保持不变，避免影响其它模块。
- 公共路由（`callback` / `wechat`）由微信侧调用，不挂此中间件，避免污染操作历史。
- 不新增数据库迁移，操作历史表 `sys_operation_records` 已存在。

## 当前进展

- [x] 在 `server/middleware/operation.go` 中拆出 `recordOperation` / `resolveOperationUserID` 公共逻辑
- [x] 新增 `OperationRecordWithUserID(forcedUserID int) gin.HandlerFunc`
- [x] 保留原 `OperationRecord()`，行为不变
- [x] `merchant.go`：`OperationRecord()` → `OperationRecordWithUserID(3)`
- [x] `contract.go` 私有分组挂 `OperationRecordWithUserID(3)`
- [x] `deduct.go` 私有分组挂 `OperationRecordWithUserID(3)`
- [x] `deduct_callback.go` 私有分组挂 `OperationRecordWithUserID(3)`
- [x] 编译与 `go vet` 检查通过

## 后续待办

- 如需把 `papay/notify` 等微信侧回调也纳入操作历史，需要单独确认是否同样使用 `UserID=3`。
- 若未来希望按调用方区分操作历史归属，可以把 `UserID=3` 抽成插件级配置项。

## 更新规则

- 同一需求始终维护在同一个文件中
- 新信息优先补充到对应段落，不要另起一份重复记录
- 只有需求状态变化时，才在 `active/` 与 `done/` 之间移动文件
