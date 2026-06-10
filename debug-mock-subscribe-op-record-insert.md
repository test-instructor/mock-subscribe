# 调试会话：mock_subscribe 操作历史无法入库

- 会话 ID：`mock-subscribe-op-record-insert`
- 状态：`[OPEN]`
- 起始时间：2026-06-10

## 现象

- 用户调用 `mockSubscribeMerchant / mockSubscribeContract / mockSubscribeDeduct / mockSubscribeDeductCallback` 等私有分组下的接口
- 期望：`system.sys_operation_records` 中新增一条 `UserID=3` 的记录
- 实际：表中没有新记录，或者记录写入异常

## 涉及改动

- `server/middleware/operation.go` 重构出 `OperationRecordWithUserID(forcedUserID int)` 与 `recordOperation(c, userId)` 公共函数
- `server/plugin/mock_subscribe/router/merchant.go / contract.go / deduct.go / deduct_callback.go` 私有分组挂 `OperationRecordWithUserID(3)`

## 静态分析阶段的候选假说

| 编号 | 假说 | 证据点 |
| ---- | ---- | ------ |
| H1 | `forcedUserID=3` 在闭包中未正确捕获，导致实际写入的 `UserID` 不是 3 | 在闭包入口与 `recordOperation` 入口各打一条 `userId` 日志，比较二者 |
| H2 | `private.Group("...").Use(middleware.OperationRecordWithUserID(3))` 写法在 Gin 中没有把中间件挂到子分组上，请求根本没走到我们的中间件 | 在 `recordOperation` 入口打日志，看 `c.Request.URL.Path` 是否在 mockSubscribe 路径下 |
| H3 | `system.SysOperationRecord` 模型上的 `User SysUser` 字段被 GORM 视为 belongs-to 关系，DB 上存在外键约束，`user_id=3` 在 `sys_users` 中不存在，导致 `Create` 直接报 FK 错误 | 通过 SQLite/MySQL 看 `sys_operation_records` 的 DDL；并在 `global.GVA_DB.Create(&record).Error` 失败分支里把 `err` 与表 DDL 一并日志化 |
| H4 | GVA 启用了多库 / 租户隔离，`global.GVA_DB` 不是写入 `sys_operation_records` 那个库 | 在 Create 前后分别 `SELECT DATABASE()` / 当前库名 |
| H5 | `c.Next()` 之后 `c.Writer` 已经被框架层替换或写入流被 close，`responseBodyWriter` 的 `body` 缓冲区始终为空，且 `Status` / `ErrorMessage` 字段取不到正常值，但这本身不阻塞 Insert；只解释字段异常 | 在 Create 之前打 `record` 全部字段 + `Create` 错误日志 |
| H6 | 中间件执行顺序问题：`OperationRecordWithUserID` 跑在 JWT/Casbin 之前，被 Casbin 403 提前 `Abort`，导致 `c.Next()` 后面的写库逻辑没执行 | 在 `recordOperation` 入口与 `c.Next()` 之后分别打日志 |

## 静态分析结论

- **H3 部分排除**：`server/service/system/sys_initdb_mysql.go:67` 与各 DB 初始化都设了 `DisableForeignKeyConstraintWhenMigrating: true`，GORM AutoMigrate 不会建外键约束，DB 层不存在 `user_id REFERENCES sys_users(id)` 的硬约束。
- **H3 真正成因**：GORM v2 在 `Create` 阶段，对 `SysOperationRecord` 上的 `User SysUser`（belongs-to，默认外键 `UserID`）会主动去装载 / 关联 User。UserID 强制写为 3 时，sys_users 中可能并不存在这条记录，GORM 会先做额外 SELECT / 尝试插入关联 User，整个 Create 阶段静默失败，sys_operation_records 里就看不到新行。

## 修复

- `server/middleware/operation.go`：把 `Create(&record)` 改成 `Omit("User").Create(&record)`，切断 GORM 对 User 关系的关联写入路径，让它只写 `user_id` 外键列。
- 错误日志增强：`create operation record error` 现在会带 `err` / `path` / `user_id`，未来如果还出问题可以直接看到 DB 真实错误。

## 运行记录

待用户调用后补充入库是否成功。


