# BUG_REPRO

## Bug 是什么
食品列表先用数据库中的旧 status 过滤，再在内存中刷新新鲜度，导致按 fresh 查询时仍会返回刷新后已经 expired 的食品。

## 如何触发
在 backend 目录运行：

```bash
go test ./internal/service -run TestFoodItemService_ListRefreshesStatusBeforeFiltering -count=1
```

测试会创建一条 status=fresh、到期日已过的食品，再按 fresh 查询。

## 错误信息
```
fresh items = 1, want 0 after refreshing status
```
