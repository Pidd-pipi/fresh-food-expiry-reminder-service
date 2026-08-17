# BUG_REPRO

## Bug 是什么
临期扫描逻辑里，食品应在剩余 3 天内变为 expiring 并产生临期通知，但当前阈值和扫描候选状态配合错误，导致临期食品仍显示 fresh，也不会生成通知。

## 如何触发
在 backend 目录运行：

```bash
go test ./internal/service -run TestReminderService_Scan_CreatesExpiringNotification -count=1
```

测试会创建一个 2 天后到期的食品，再执行一轮扫描。

## 错误信息
```
created notifications = 0, want 1
```
