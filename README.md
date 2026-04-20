# Go Echo Template v0.1.0

基于 Echo + GORM 的 Go 后端模板，默认提供完整的用户认证链路、统一响应与错误处理、CLI 启动方式、可配置限流、健康检查和可复用模块结构。

## 核心特性

- **Web 框架**: Echo v4
- **ORM**: GORM v2，支持 SQLite / MySQL
- **认证体系**: JWT Access Token + Refresh Token
- **工程骨架**: `handler/service/repository` 分层
- **统一协议**: 标准化 API 响应、统一错误结构、请求 ID
- **生产基础能力**: 优雅关闭、读写超时、可配置限流、日志切割
- **模板复用性**: 模块边界清晰，新增业务模块成本低

## 快速开始

```bash
# 1. 克隆项目
git clone <repo-url> myapp
cd myapp

# 2. 配置环境变量
cp .env.example .env
# 至少修改 JWT_SECRET

# 3. 安装依赖
go mod tidy

# 4. 执行迁移与种子数据
go run . migrate
go run . seed

# 5. 启动服务
go run . serve
```

服务默认启动于：`http://localhost:8080`

## 文档索引

- **接口文档**: `./api.md`
- **环境变量模板**: `./.env.example`

## 当前 API 概览

### 健康检查

```text
GET /health/live
GET /health/ready
```

### 认证相关

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
GET  /api/v1/auth/profile
PUT  /api/v1/auth/profile
POST /api/v1/auth/password
```

### 管理员用户管理

```text
GET    /api/v1/admin/users
GET    /api/v1/admin/users/:id
POST   /api/v1/admin/users
PUT    /api/v1/admin/users/:id
PATCH  /api/v1/admin/users/:id/status
DELETE /api/v1/admin/users/:id
```

## 启动流程

服务入口为 `main.go`，通过 CLI 命令分发到 `cmd/commands/serve.go`。

`serve` 命令的核心启动流程如下：

1. 加载 `.env` 与运行配置
2. 初始化日志
3. 初始化数据库并做健康检查
4. 创建 Echo 实例并挂载：
   - Request ID
   - Recover
   - CORS
   - 可选全局限流
   - Debug Logger
5. 注册路由与业务模块
6. 使用 HTTP Server 配置超时参数
7. 监听系统信号并优雅关闭

## 默认测试账号

执行 `go run . seed` 后：

| 用户名 | 密码 | 角色 |
|------|------|------|
| admin | admin123 | admin |
| user01 | user123 | user |

## CLI 命令

```bash
go run . serve                  # 启动 HTTP 服务（默认 localhost:8080）
go run . serve --port 9000      # 指定端口
go run . serve --host 0.0.0.0   # 指定监听地址
go run . migrate                # 数据库迁移（建表/更新表结构）
go run . seed                   # 填充测试数据
go run . --version              # 查看版本
go run . help                   # 查看所有命令
```

## 测试账号

执行 `go run . seed` 后：

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |
| user01 | user123 | 普通用户 |

## 项目结构

```
.
├── main.go                          # 入口
├── .env.example                     # 环境变量模板
│
├── cmd/commands/                    # CLI 命令
│   ├── command.go                   #   Command 接口 + Registry + ParseFlags
│   ├── commands.go                  #   全局变量 + RegisterCommands
│   ├── migrate.go                   #   数据库迁移命令
│   ├── seed.go                      #   数据填充命令
│   └── serve.go                     #   启动服务命令
│
├── config/                          # 配置层
│   ├── env.go                       #   Config 主结构 + Load()
│   ├── app.go                       #   AppConfig
│   ├── database.go                  #   DatabaseConfig (MySQL/SQLite)
│   ├── jwt.go                       #   JWTConfig
│   └── logger.go                    #   LoggerConfig
│
├── common/                          # 公共工具
│   └── response.go                  #   统一 API 响应
│
├── models/                          # GORM 模型
│   ├── base.go                      #   BaseModel (ID + 时间戳 + 软删除)
│   ├── user.go                      #   User 模型
│   └── setting.go                   #   Setting 模型
│
├── server/
│   ├── database/
│   │   ├── db.go                    #   InitDB + CloseDB + HealthCheck
│   │   └── migration.go            #   AutoMigrate
│   ├── middlewares/
│   │   ├── auth.go                  #   JWT 认证 + AdminOnly
│   │   └── rate_limit.go           #   限流中间件
│   └── router/
│       └── router.go               #   路由注册
│
├── modules/                         # 业务模块（按领域划分）
│   └── user/                        #   用户模块（完整示例）
│       ├── dto.go                   #     请求/响应 DTO
│       ├── handler.go               #     HTTP Handler
│       ├── repository.go            #     数据仓库
│       └── service.go               #     业务逻辑
│
├── services/                        # 公共服务层
│   ├── interfaces.go                #   Cache/Lock/RateLimiter/Notifier/SMS 接口
│   ├── redis_service.go             #   Redis 连接 + 全套操作
│   ├── cache_service.go             #   缓存（内存 / Redis 双实现）
│   ├── memory_lock_service.go       #   内存锁
│   ├── redis_lock_service.go        #   Redis 分布式锁
│   ├── rate_limiter_service.go      #   限流器（内存 / Redis 双实现）
│   ├── jwt_service.go               #   JWT 生成/解析/刷新
│   ├── mail_service.go              #   SMTP 邮件
│   ├── sms_service.go               #   短信（接口 + 阿里云示例 + Mock）
│   ├── pagination.go                #   分页查询辅助
│   ├── setting_service.go           #   系统设置 CRUD
│   ├── notify_service.go            #   统一通知（聚合多渠道）
│   ├── notify_telegram.go           #   Telegram Bot
│   ├── notify_dingtalk.go           #   钉钉 Webhook
│   ├── notify_feishu.go             #   飞书 Webhook
│   ├── queue_service.go             #   队列接口 + Task/Handler 定义
│   ├── sqlite_queue_service.go      #   SQLite 队列实现
│   ├── job_executor.go              #   Worker Pool 执行器
│   ├── job_scheduler.go             #   定时任务调度器
│   └── job_dispatcher.go            #   任务分发器
│
└── utils/                           # 工具函数
    ├── logger.go                    #   Logrus 日志
    ├── gorm_logger.go               #   GORM → Logrus 桥接
    └── validation.go                #   Validator 服务
```

## API 接口

完整接口说明请查看 `api.md`。

### 快速示例

```bash
# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 刷新 token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh-token>"}'

# 访问个人资料
curl http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <access-token>"

# 健康检查
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DB_TYPE` | `sqlite` | 数据库类型：`sqlite` / `mysql` |
| `SQLITE_PATH` | `./data/app.db` | SQLite 文件路径 |
| `MYSQL_DSN` | - | MySQL 连接串（DB_TYPE=mysql 时必填） |
| `DB_HOST` | `127.0.0.1` | MySQL 主机 |
| `DB_PORT` | `3306` | MySQL 端口 |
| `DB_USER` | `root` | MySQL 用户名 |
| `DB_PASSWORD` | - | MySQL 密码 |
| `DB_NAME` | `myapp` | MySQL 数据库名 |
| `JWT_SECRET` | (内置默认值) | JWT 签名密钥（生产环境必须修改） |
| `JWT_EXPIRATION` | `168h` | JWT 过期时间 |
| `JWT_REFRESH_EXPIRATION` | `336h` | Refresh Token 过期时间 |
| `APP_NAME` | `MyApp` | 应用名 |
| `APP_VERSION` | `0.1.0` | 应用版本 |
| `APP_ENV` | `development` | 运行环境 |
| `APP_HOST` | `localhost` | 监听地址 |
| `APP_PORT` | `8080` | 监听端口 |
| `APP_DEBUG` | `true` | 调试模式 |
| `APP_READ_TIMEOUT` | `15s` | 读取超时 |
| `APP_WRITE_TIMEOUT` | `15s` | 写入超时 |
| `APP_IDLE_TIMEOUT` | `60s` | 空闲连接超时 |
| `APP_SHUTDOWN_TIMEOUT` | `30s` | 优雅关闭超时 |
| `RATE_LIMIT_ENABLED` | `false` | 是否启用全局限流 |
| `RATE_LIMIT_LIMIT` | `120` | 窗口内允许请求数 |
| `RATE_LIMIT_WINDOW` | `1m` | 限流时间窗口 |
| `LOG_LEVEL` | `info` | 日志级别：trace/debug/info/warn/error |
| `LOG_FORMAT` | `text` | 日志格式：text/json |
| `LOG_OUTPUT` | `both` | 输出目标：stdout/file/both |
| `LOG_FILE` | `./logs/app.log` | 日志文件路径 |

## 使用指南

### 1. 新增业务模块

参照 `modules/user/` 的四件套模式：

```
modules/product/
├── dto.go          # 请求/响应结构体
├── handler.go      # HTTP Handler（路由处理）
├── repository.go   # 数据仓库（数据库操作）
└── service.go      # 业务逻辑
```

```go
// modules/product/handler.go
package product

type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
    g.GET("/products", h.List)
    g.POST("/products", h.Create)
    g.GET("/products/:id", h.Get)
    g.PUT("/products/:id", h.Update)
    g.DELETE("/products/:id", h.Delete)
}
```

然后在 `server/router/router.go` 中注册：

```go
// router.go
productRepo := product.NewRepository(db)
productService := product.NewService(productRepo)
productHandler := product.NewHandler(productService)
productHandler.RegisterRoutes(authGroup.Group("/products"))
```

### 2. 使用缓存

```go
// 内存缓存（单机）
cache := services.NewMemoryCacheService()

// Redis 缓存（分布式）
redisSvc, _ := services.NewRedisService(services.RedisConfig{
    Host: "127.0.0.1", Port: "6379",
})
cache := services.NewRedisCacheService(redisSvc)

// 两种实现接口相同
val, _ := cache.Get("key")
cache.Set("key", "value", 5*time.Minute)

// GetOrSet：缓存不存在时自动执行回调
result, _ := cache.GetOrSet("user:1", 10*time.Minute, func() (string, error) {
    user, err := repo.FindByID(1)
    data, _ := json.Marshal(user)
    return string(data), err
})

// JSON 便捷方法
services.CacheSetJSON(cache, "user:1", user, 10*time.Minute)
services.CacheGetJSON(cache, "user:1", &user)
```

### 3. 使用锁

```go
// 内存锁（单机）
locker := services.NewMemoryLockService()

// Redis 锁（分布式）
locker := services.NewRedisLockService(redisSvc)

// 方式一：手动获取/释放
lock, err := locker.Obtain("order:123", 30*time.Second, 3) // 重试 3 次
if err == nil {
    defer lock.Release()
    // do work...
}

// 方式二：WithLock 自动管理
err := locker.WithLock("order:123", 30*time.Second, func() error {
    // do work...
    return nil
})
```

### 4. 使用队列 & 定时任务

```go
// === 启动时初始化 (serve.go) ===
queue := services.NewSQLiteQueueService(db)
lock := services.NewMemoryLockService()

// 注册任务处理器
executor := services.NewJobExecutor(queue, lock, 2)
executor.RegisterHandler(&SendEmailHandler{})
executor.RegisterHandler(&GenerateReportHandler{})

// 启动执行器 + 调度器
go executor.Start()
go services.NewJobScheduler(queue, 10*time.Second).Start()

// === 业务代码中投递任务 ===
dispatcher := services.NewJobDispatcher(queue)

// 立即执行
dispatcher.Dispatch("send_email", map[string]string{
    "to":      "user@example.com",
    "subject": "Welcome",
})

// 5 分钟后执行
dispatcher.DispatchLater("generate_report", payload, 5*time.Minute)

// 指定时间执行
dispatcher.DispatchAt("send_email", payload, time.Date(2026, 3, 20, 10, 0, 0, 0, time.Local))

// 带优先级（值越大越优先）
dispatcher.DispatchWithPriority("urgent_task", payload, 10)
```

**自定义 Handler：**

```go
package jobs

import (
    "context"
    "encoding/json"
    "myapp/services"
)

type SendEmailHandler struct {
    mailer *services.MailService
}

func (h *SendEmailHandler) GetTaskType() string            { return "send_email" }
func (h *SendEmailHandler) CanHandle(t string) bool        { return t == "send_email" }

func (h *SendEmailHandler) Handle(ctx context.Context, task services.Task) error {
    var payload struct {
        To      string `json:"to"`
        Subject string `json:"subject"`
        Body    string `json:"body"`
    }
    json.Unmarshal(task.GetPayload(), &payload)
    return h.mailer.Send([]string{payload.To}, payload.Subject, payload.Body)
}
```

### 5. 使用限流

```go
// 内存限流（单机）
limiter := services.NewMemoryRateLimiter()

// Redis 限流（分布式）
limiter := services.NewRedisRateLimiter(redisSvc)

// 作为中间件使用
e.Use(middlewares.RateLimit(middlewares.RateLimitConfig{
    Limiter: limiter,
    Limit:   100,              // 每窗口 100 次
    Window:  time.Minute,      // 窗口 1 分钟
}))

// 按用户限流
apiGroup.Use(middlewares.RateLimit(middlewares.RateLimitConfig{
    Limiter: limiter,
    Limit:   30,
    Window:  time.Minute,
    KeyFunc: func(c echo.Context) string {
        return fmt.Sprintf("user:%v", c.Get("user_id"))
    },
}))
```

### 6. 使用通知

```go
// 创建通知渠道
tg := services.NewTelegramNotifier("bot_token", "chat_id")
dd := services.NewDingtalkNotifier("https://oapi.dingtalk.com/robot/send?access_token=xxx")
fs := services.NewFeishuNotifier("https://open.feishu.cn/open-apis/bot/v2/hook/xxx")

// 聚合通知服务
notifier := services.NewNotifyService(tg, dd, fs)

// 广播到所有渠道
notifier.Send("部署通知", "v0.1.0 已发布到生产环境")

// 发送到指定渠道
notifier.SendTo("telegram", "告警", "CPU > 90%")

// 查看已注册渠道
notifier.Channels() // ["telegram", "dingtalk", "feishu"]
```

### 7. 使用邮件 & 短信

```go
// 邮件
mailer := services.NewMailService(services.MailConfig{
    Host: "smtp.example.com", Port: 465,
    Username: "noreply@example.com", Password: "xxx",
    From: "noreply@example.com", UseTLS: true,
})
mailer.Send([]string{"user@example.com"}, "Subject", "Body")
mailer.SendHTML([]string{"user@example.com"}, "Subject", "<h1>Hello</h1>", true)

// 短信（阿里云）
sms := services.NewSMSService(services.NewAliyunSMS("keyID", "keySecret", "签名"))
sms.Send("13800138000", "SMS_123456", map[string]string{"code": "1234"})

// 短信（开发环境 Mock，仅打印到控制台）
sms := services.NewSMSService(services.NewMockSMS())
```

### 8. 使用分页

```go
// 方式一：GORM Scope
p := services.NewPagination(page, pageSize)
var users []models.User
db.Model(&models.User{}).Scopes(p.Paginate()).Find(&users)

// 方式二：一步到位（count + find）
var products []models.Product
p, err := services.PaginateQuery(
    db.Model(&models.Product{}).Where("status = ?", "active"),
    page, pageSize, &products,
)
// p.Total, p.Pages, p.Page, p.PageSize 全部自动填充
```

### 9. 切换数据库

`.env` 中修改即可，代码无需改动：

```bash
# SQLite（默认，零配置）
DB_TYPE=sqlite
SQLITE_PATH=./data/app.db

# MySQL
DB_TYPE=mysql
MYSQL_DSN=root:password@tcp(127.0.0.1:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local
# 或者分开配置
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=myapp
```

## 构建

```bash
# 开发
go run . serve

# 生产构建
go build -ldflags "-X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o myapp .

# 运行
./myapp serve
./myapp migrate
./myapp --version
# MyApp v0.1.0
```

## License

MIT
