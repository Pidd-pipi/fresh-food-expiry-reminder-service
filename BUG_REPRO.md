# Bug 是什么

食谱详情接口在查询不存在的食谱 ID 时返回 500，而不是 404。根因是仓储层用 `%v` 丢掉了 `ErrNotFound` 错误链，服务层 `errors.Is` 误判为系统错误，处理器又把错误硬编码成 500。

# 如何触发

进入 backend 目录，执行：

```
go test ./internal/service -run '^TestRecipeGetMissing$' -count=1
go test ./internal/handler -run '^TestRecipeDetail404$' -count=1
```

# 错误信息

```
--- FAIL: TestRecipeGetMissing (0.00s)
    recipe_service_test.go:28: expected AppError, got *errors.errorString: recipe 999 not found: record not found
--- FAIL: TestRecipeDetail404 (0.00s)
    recipe_handler_test.go:56: status = 500, want 404; body={"code":1007,"data":null,"message":"服务内部错误"}
```
