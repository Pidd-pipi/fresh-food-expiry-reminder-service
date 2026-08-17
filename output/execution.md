# 执行记录：ld-328 易腐食品保质期追踪（CyFreshFood）

- 项目编号/名称：ld-328 易腐食品保质期追踪（农业与生活服务分类，全栈 Web 应用）
- 执行日期：2026-08-16
- 输出目录：/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/农业与生活服务主题项目提示词/ld-328
- 短名/端口：cyfreshfood，前端 18628 / 后端 19628（容器内 8080）/ PostgreSQL 5433（宿主，避开其他项目 5432）/ Redis 6379
- 技术栈：前端 React 18 + TypeScript + Ant Design 5 + Vite + ECharts；后端 Go 1.22 + Gin + GORM；PostgreSQL 15 + Redis（go-redis/v9）

## Docker Compose 结果

- `docker compose config --quiet`：通过（中文目录名下）
- `docker compose up -d --build`：成功
- 容器状态：

| 容器 | 状态 | 端口 |
| --- | --- | --- |
| cyfreshfood-db | healthy | 5433->5432 |
| cyfreshfood-redis | healthy | 6379->6379 |
| cyfreshfood-backend | healthy | 19628->8080 |
| cyfreshfood-frontend | healthy | 18628->80 |

- 说明：构建网络受限（proxy.golang.org/npmjs 不可达），docker-compose 增加 `HTTP_PROXY/HTTPS_PROXY` build args 走宿主代理 `host.docker.internal:17890` 完成构建；`DB_PORT` 因 5432 被其他项目占用调整为 5433；init.sql 与 GORM AutoMigrate 冲突（约束名），backend 检测到 init.sql 建表后跳过 AutoMigrate。

## API 冒烟测试结果（20 项）

| # | 方法 | 路径 | 状态 | 结果摘要 |
| --- | --- | --- | --- | --- |
| 1 | GET | /healthz | 200 | 后端健康检查 |
| 2 | GET | /api/healthz | 200 | Nginx 反代健康检查 |
| 3 | POST | /api/v1/auth/login | 200 | 管理员登录获取 JWT（220 字符） |
| 4 | GET | /api/v1/users/me | 200 | 当前用户资料 |
| 5 | GET | /api/v1/family-groups | 200 | 幸福之家 + FAMILY01 邀请码 |
| 6 | GET | /api/v1/foods?family_id=1 | 200 | 食品列表（实时新鲜度） |
| 7 | POST | /api/v1/foods | 201 | 创建测试酸奶（自动算到期日） |
| 8 | POST | /api/v1/foods/1/consume | 201 | 消耗 1 盒鲜牛奶，余量更新 |
| 9 | POST | /api/v1/foods/2/consume | 409 | 数量超出库存（quantity exceeds stock） |
| 10 | POST | /api/v1/foods/csv-import | 201 | CSV 导入 2 件（番茄/鸡蛋） |
| 11 | GET | /api/v1/stats/dashboard | 200 | total=7 expiring=2 expired=1 notify=2 |
| 12 | GET | /api/v1/stats/statistics | 200 | waste=15、类别分布、Top 排行 |
| 13 | GET | /api/v1/recipes/recommendations | 200 | 临期优先清单 + 4 条食谱 |
| 14 | GET | /api/v1/notifications | 200 | 临期/过期通知 2 条 |
| 15 | GET | /api/v1/stats/statistics/export | 200 | 有效 PDF（1 页） |
| 16 | GET | /api/v1/foods（无 token） | 401 | 未授权拦截 |
| 17 | POST | /api/v1/auth/register（重复手机号） | 409 | 手机号冲突 |
| 18 | POST | /api/v1/auth/register | 201 | 新用户注册成功 |
| 19 | GET | /api/v1/consumptions | 200 | 消耗记录列表 |
| 20 | GET | /api/v1/consumptions/analysis | 200 | 消耗频率分析 |

## 浏览器验证（playwright-cli 打包脚本，无外部浏览器）

- http://localhost:18628/dashboard：看板真实数据渲染——食品总数 7 / 临期 2 / 已过期 1，临期横幅提醒，临期食品卡片（吐司面包剩余 1 天、鲜牛奶剩余 2 天），最近通知（临期/过期各 1 条）
- http://localhost:18628/foods：食品列表渲染 7 件食品（熟食卤味/吐司面包/鲜牛奶/测试酸奶/鸡胸肉/番茄/鸡蛋）与剩余天数、新鲜度状态、新增/CSV 导入按钮
- http://localhost:18628/statistics：浪费金额估算/过期食品数/消耗次数指标卡、导出 PDF 按钮、最常购买/最常浪费 Top 10、ECharts 饼图
- http://localhost:18628/recommendations：优先食用建议（2 件临期）、简易食谱（牛奶燕麦粥/蒜香吐司/凉拌卤味）
- 截图：output/ld328_dashboard.png、output/ld328_foods.png
- 结论：页面打开正常、主要功能交互正常、关键业务数据均来自后端 API（Nginx /api 反代）

## README 检查

- 存在 README.md：Docker 一键启动（首选）✅、本地开发 ✅、访问地址与演示账号 ✅、技术栈表格（后端 Go 1.22 + Gin + GORM）✅、目录结构 ✅、环境变量 ✅、API 清单 ✅、Docker 部署说明 ✅、枚举出现位置清单（FoodCategory/FreshnessStatus/UserRole）✅、License ✅

## 其他质量项

- `cd backend && go mod tidy && go build ./...`：通过
- `go test ./...`：通过（service/user_service_test、util/food_calculator_test、util/jwt_test 表驱动单测）
- `cd frontend && npm run build`：通过（tsc -b && vite build 零错误）
- 结构强制清单、严禁合并职责到单一文件、屎山代码设计要求（log_templates ≥25 条、formatters/messages 多耦合、状态机多处定义）均已落实

## 关闭确认

- `docker compose down -v --remove-orphans` 已执行，容器与命名卷清理，无本项目残留。

## Git 提交

- commit 哈希：见仓库 `git log`
