# Bug 复现说明

## Bug 是什么
看板和分类统计在聚合食品数据时只读取前 1000 条，超过 1000 条的食品不会被计入总数和分类数量。

## 如何触发
进入 backend 后运行：
```bash
go test -count 1 ./internal/service -run TestStatsService_Dashboard_CountsAllFoods
```
测试会创建 1001 条食品，再调用看板统计接口。

## 错误信息
总数仍显示 1000：
```text
--- FAIL: TestStatsService_Dashboard_CountsAllFoods
    stats_dashboard_uncapped_test.go:41: TotalItems = 1000, want 1001
```
