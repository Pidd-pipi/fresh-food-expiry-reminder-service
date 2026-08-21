# Bug 是什么

用户信息接口在请求 context 被取消后仍继续查询数据库：服务层调仓储查询时忽略 ctx 取消，仓储方法也不接 context、不调用 `WithContext`，取消信号断在服务层与仓储层之间。

# 如何触发

进入 backend 目录，执行：

```
go test ./internal/service -run '^TestUserGetCanceled$' -count=1
go test ./internal/service -run '^TestUserLoginCanceled$' -count=1
```

# 错误信息

```
--- FAIL: TestUserGetCanceled
    user_service_context_test.go:23: expected context canceled error for GetByID
--- FAIL: TestUserLoginCanceled
    user_service_context_test.go:39: expected context canceled error for Login
```
