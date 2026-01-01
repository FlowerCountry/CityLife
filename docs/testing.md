# 测试接口说明

本项目提供测试专用接口用于提升余额，方便覆盖高金额场景。该接口仅在 Gin 的测试模式下可用，生产路由不会注册。

## 测试专用余额接口

- 路径：`POST /api/v2/sessions/:session_id/test/money`
- 请求体：
  - `wallet`：钱包总额（非负整数）
  - `bank`：银行余额（非负整数）
- 说明：
  - 至少提供一个字段
  - 仅测试模式可用，非测试模式返回 403

示例：

```bash
curl -X POST "http://localhost:8080/api/v2/sessions/<session_id>/test/money" \
  -H "Content-Type: application/json" \
  -d '{"wallet": 1200, "bank": 300}'
```
