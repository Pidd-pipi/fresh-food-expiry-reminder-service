# Bug 是什么

家庭组详情查询不存在的家庭组、以及用无效邀请码加入家庭组时，接口都返回 500，而不是 404。根因是仓储层 `FindByID`/`FindByInviteCode` 用 `%v` 丢掉了 `ErrNotFound` 错误链，服务层 `errors.Is` 误判为系统错误。

# 如何触发

进入 backend 目录，执行：

```
go test ./internal/service -run '^TestFamilyGroupGetMissing$' -count=1
go test ./internal/service -run '^TestFamilyGroupJoinBadCode$' -count=1
go test ./internal/handler -run '^TestFamilyGroupDetail404$' -count=1
```

# 错误信息

```
--- FAIL: TestFamilyGroupGetMissing
    family_group_service_error_test.go:24: expected AppError, got *errors.errorString: family group 999 not found: record not found
--- FAIL: TestFamilyGroupJoinBadCode
    family_group_service_error_test.go:42: expected AppError, got *errors.errorString: family group invite code NO_SUCH_CODE not found: record not found
--- FAIL: TestFamilyGroupDetail404
    family_group_handler_test.go:53: status = 500, want 404
```
