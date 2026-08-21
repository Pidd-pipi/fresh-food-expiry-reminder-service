# Bug 是什么

通知已读状态机错位：`MarkRead` 用 `food_item_id` 当主键导致标错对象，`MarkAllRead` 漏写 `read_at` 且筛选错集合，`UnreadCount` 传错家庭参数，最终表现为标已读标错、全部已读后未读数字仍不为 0。

# 如何触发

进入 backend 目录，执行：

```
go test ./internal/service -run '^TestNotifyMarkReadTarget$' -count=1
go test ./internal/service -run '^TestNotifyMarkAllReadState$' -count=1
go test ./internal/service -run '^TestNotifyUnreadCount$' -count=1
```

# 错误信息

```
--- FAIL: TestNotifyMarkReadTarget
    notification_service_test.go:59: n2 should remain unread
--- FAIL: TestNotifyMarkAllReadState
    notification_service_test.go:90: notification 1 read_at should be set
--- FAIL: TestNotifyUnreadCount
    notification_service_test.go:108: unread count = 0, want 2
```
