# Bug 是什么

通知列表接口按类型聚合时，`util.CountStrings`/`util.GroupStrings` 用 nil map 直接写入，首次命中即 panic；同时 `RateLimiter` 的 `buckets` 也漏初始化为 nil map，首个请求触发 `assignment to entry in nil map`。

# 如何触发

进入 backend 目录，执行：

```
go test ./internal/util -run '^TestCountStringsSafe$' -count=1
go test ./internal/util -run '^TestGroupStringsSafe$' -count=1
go test ./internal/middleware -run '^TestRateLimiterSafe$' -count=1
```

# 错误信息

```
panic: assignment to entry in nil map [recovered]
	panic: assignment to entry in nil map

goroutine 35 [running]:
github.com/blueship581/cyfreshfood/internal/util.CountStrings(...)
	backend/internal/util/formatters.go:111
```
