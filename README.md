# BowlUtils

一个轻量级、模块化的 Go 工具库集合，提供常用的业务功能和基础设施组件。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.25-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 特性

- 🎯 **模块化设计** - 每个包独立干净，可按需引入
- 🚀 **高性能** - 基于 Go 泛型，零反射开销
- 🔧 **开箱即用** - 提供常用业务场景的最佳实践
- 📦 **依赖清晰** - 最小化依赖，避免包膨胀
- ✅ **测试完善** - 全面的单元测试覆盖

## 安装

```bash
go get github.com/lazyfury/bowlutils
```

## 包列表

### 核心工具

#### `crud` - CRUD 操作封装
基于 GORM 的通用 CRUD 仓储模式实现，支持分页、条件查询等。

```go
import "github.com/lazyfury/bowlutils/crud"

type User struct {
    crud.BaseModel
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (u User) TableName() string {
    return "users"
}

// 创建仓储
repo := crud.NewRepository(User{}, db)

// 查询
user, err := repo.FindByID(1)

// 分页查询
page, err := repo.FindPage(1, 10, map[string]interface{}{
    "name": "John",
})
```

**特性：**
- 泛型支持，类型安全
- 内置软删除支持
- 灵活的条件查询（支持 eq, ne, gt, lt, like, in 等）
- 分页查询
- 排序支持

#### `ioc` - IOC 容器
轻量级的依赖注入容器，支持单例和工厂模式。

```go
import "github.com/lazyfury/bowlutils/ioc"

// 注册服务
ioc.Provide("db", func() (any, error) {
    return db.NewDB("postgres", dsn), nil
}, true) // true = 单例

// 获取服务
db, ok := ioc.Get("db")
db := ioc.MustGet[*gorm.DB]("db") // 泛型方式
```

**特性：**
- 单例模式支持
- 懒加载
- 线程安全
- 泛型类型断言

#### `eventbus` - 事件总线
线程安全的发布订阅模式实现。

```go
import "github.com/lazyfury/bowlutils/eventbus"

bus := eventbus.New()

// 订阅
id, ch := bus.Subscribe("user.created", 10)
go func() {
    for event := range ch {
        // 处理事件
    }
}()

// 发布
bus.Publish("user.created", userData)

// 取消订阅
bus.Unsubscribe("user.created", id)
```

**特性：**
- 非阻塞发布
- 支持多订阅者
- 自动丢弃满缓冲区消息
- 线程安全

#### `isvlid` - 数据验证
基于 validator/v10 的验证增强，支持自定义条件。

```go
import "github.com/lazyfury/bowlutils/isvlid"

type UserInput struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age"`
}

validator := isvlid.NewValidator(&input,
    isvlid.WithCondition("Name", isvlid.Required()),
    isvlid.WithCondition("Email", isvlid.IsValidEmail("", false)),
    isvlid.WithCondition("Age", isvlid.Min(18), isvlid.Max(100)),
)

if err := validator.Validate(); err != nil {
    // 处理验证错误
}
```

**内置验证器：**
- `Required()` - 必填
- `IsEnum()` / `IsOneOf()` - 枚举值
- `IsValidPhone()` - 手机号
- `IsValidEmail()` - 邮箱
- `Min()` / `Max()` - 数值范围
- `Default()` - 默认值

### 基础设施

#### `db` - 数据库连接
简化的数据库连接管理，支持 MySQL 和 PostgreSQL。

```go
import "github.com/lazyfury/bowlutils/db"

db := db.NewDB("postgres", dsn)
// 或
db := db.NewDB("mysql", dsn)
```

#### `email` - 邮件发送
SMTP 邮件发送器，支持同步和异步发送。

```go
import "github.com/lazyfury/bowlutils/email"

// 配置
config := &email.Config{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-password",
    From:     "your-email@gmail.com",
    FromName: "Your App",
    TLS:      false,
}

sender := email.NewSMTPSender(config)

// 发送邮件
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Hello",
    Body:    "Plain text body",
    HTML:    "<h1>HTML body</h1>",
}

err := sender.Send(context.Background(), msg)
```

#### `resp` - HTTP 响应
基于 Gin 的统一响应格式。

```go
import "github.com/lazyfury/bowlutils/resp"

// 成功响应
resp.Ok(c, data)

// 失败响应
resp.Fail[any](c, "操作失败")

// 错误响应
resp.Error(c, 500, "服务器错误", nil)

// 其他响应
resp.NotFound[any](c, "资源不存在")
resp.Unauthorized[any](c, "未授权")
resp.Forbidden[any](c, "无权限")
```

#### `openapi` - OpenAPI 文档
OpenAPI 3.0 文档生成工具。

```go
import "github.com/lazyfury/bowlutils/openapi"

doc := openapi.NewDocument("3.0.0").
    WithInfo("My API", "1.0.0").
    AddServer(openapi.Server{URL: "http://localhost:8080"})

// 添加路径和操作
doc.AddOperation("/users", "get", openapi.Operation{
    Summary: "获取用户列表",
    Responses: openapi.NewResponses(
        openapi.NewResponseFrom(200, "成功", []User{}),
    ),
})
```

#### `viperinit` - 配置管理
简化的 Viper 配置初始化。

```go
import "github.com/lazyfury/bowlutils/viperinit"

v := viperinit.NewViper("config", "yaml", ".")
port := v.GetInt("server.port")
```

#### `utils` - 通用工具
常用的辅助函数。

```go
import "github.com/lazyfury/bowlutils/utils"

// 类型转换
str := utils.ToString(123)           // "123"
m, _ := utils.ToMap(struct{Name: "John"})

// 零值检查
utils.IsZero(0)    // true
utils.IsEmpty("")  // true

// 默认值
value := utils.Def("", "default") // "default"
```

## 依赖关系

```
独立包（无内部依赖）：
├── utils       - 通用工具
├── viperinit   - 配置管理
├── db          - 数据库连接
├── eventbus    - 事件总线
├── ioc         - IOC容器
├── isvlid      - 数据验证
└── openapi     - OpenAPI文档

有限依赖包：
├── crud        → gorm (外部)
├── resp        → gin (外部)
└── email       → ioc (可选)
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...

# 运行特定包的测试
go test ./crud/...
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 致谢

基于以下优秀的开源项目：
- [GORM](https://gorm.io/) - ORM 库
- [Gin](https://gin-gonic.com/) - Web 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Validator](https://github.com/go-playground/validator) - 数据验证

---

**注意**: 本项目处于积极开发中，API 可能会有变动。建议在生产环境使用前固定版本。
