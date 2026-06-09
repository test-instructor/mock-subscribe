# 用户协议前端适配嵌套 contract/status 返回

## 基本信息

- 提出日期：2026-06-09
- 当前状态：`active`
- 需求类型：前端展示调整
- 优先级：中
- 需求文件：`aiDoc/memory/business/active/mock-subscribe-contract-frontend-display.md`

## 用户原始意图摘要

根据 `mockSubscribe/mockSubscribeContract/getContractList` 与详情接口返回的 `contract`、`status` 嵌套结构，调整“用户协议”页面前端显示，使页面正确展示协议主体信息与状态信息。

## 影响范围

- 后端：无接口结构改动，仅以前端适配现有返回契约为准
- 前端：`web/src/plugin/mock_subscribe/view/contract.vue`
- 文档：业务记忆索引与本需求记录
- 插件 / 模块：`mock_subscribe`

## 涉及对象

- 模块：`mock_subscribe`
- 接口：`/mockSubscribeContract/getContractList`、`/mockSubscribeContract/findContract`
- 页面：`web/src/plugin/mock_subscribe/view/contract.vue`
- 配置：无

## 已确认约束

- 列表接口每条记录返回 `{ contract, status }`
- 详情接口返回 `{ contract, status }`
- 前端需保持统一响应结构 `{ code, data, msg }` 的使用方式
- 仅做与当前页面展示相关的最小改动，不扩展无关功能

## 当前进展

- 已定位到用户协议页面仍按扁平字段读取列表数据
- 已确认接口真实返回为嵌套结构，属于前端展示层适配问题

## 后续待办

- 将列表列渲染改为从 `row.contract` / `row.status` 读取
- 适配弹窗编辑入口的 ID 与展示字段读取
- 完成后做页面级 lint 检查
