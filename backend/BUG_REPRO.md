# Bug 复现说明

## Bug 是什么
月度消耗统计使用了数据库专用的 `to_char` 函数，在 SQLite 测试环境执行时会报 SQL 错误，导致统计分析无法完成。

## 如何触发
进入 backend 后运行：
```bash
go test -count 1 ./internal/service -run TestConsumptionRecordService_Analysis_OnSQLite
```

## 错误信息
```text
SQL logic error: no such function: to_char (1)
--- FAIL: TestConsumptionRecordService_Analysis_OnSQLite
    consumption_analysis_sqlite_test.go:22: Analysis() error = monthly stats: SQL logic error: no such function: to_char (1)
```
