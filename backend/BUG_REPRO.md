# Bug 复现说明

## Bug 是什么
CSV 批量导入食品时，如果先读到合法行、后面又出现格式错误的行，已解析成功的食品已经被写入数据库，返回失败后库存里会残留部分导入数据。

## 如何触发
进入 backend 后运行：
```bash
go test -count 1 ./internal/service -run TestFoodItemService_ImportCSV_AllOrNothingOnMalformedRow
```
也可以调用 `FoodItemService.ImportCSV`，传入包含“合法行 + 未闭合引号行”的 CSV 文本。

## 错误信息
测试失败，期望失败导入后食品数为 0，实际为 1：
```text
--- FAIL: TestFoodItemService_ImportCSV_AllOrNothingOnMalformedRow
    food_import_csv_atomic_test.go:30: food count = 1, want 0 after failed import
```
