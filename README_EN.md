# BowlUtils

A lightweight, modular Go utility library collection providing common business functionality and infrastructure components.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.25-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Features

- 🎯 **Modular Design** - Each package is independent and clean, import only what you need
- 🚀 **High Performance** - Built with Go generics, zero reflection overhead
- 🔧 **Ready to Use** - Provides best practices for common business scenarios
- 📦 **Clear Dependencies** - Minimized dependencies, avoid package bloat
- ✅ **Well Tested** - Comprehensive unit test coverage

## Installation

```bash
go get github.com/lazyfury/bowlutils
```

## Packages

### Core Tools

#### `crud` - CRUD Operations Wrapper
Generic CRUD repository pattern implementation based on GORM, supporting pagination and conditional queries.

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

// Create repository
repo := crud.NewRepository(User{}, db)

// Query
user, err := repo.FindByID(1)

// Paginated query
page, err := repo.FindPage(1, 10, map[string]interface{}{
    "name": "John",
})
```

**Features:**
- Generic support with type safety
- Built-in soft delete support
- Flexible conditional queries (eq, ne, gt, lt, like, in, etc.)
- Pagination
- Sorting support

#### `ioc` - IOC Container
Lightweight dependency injection container supporting singleton and factory patterns.

```go
import "github.com/lazyfury/bowlutils/ioc"

// Register service
ioc.Provide("db", func() (any, error) {
    return db.NewDB("postgres", dsn), nil
}, true) // true = singleton

// Get service
db, ok := ioc.Get("db")
db := ioc.MustGet[*gorm.DB]("db") // Generic way
```

**Features:**
- Singleton pattern support
- Lazy loading
- Thread-safe
- Generic type assertions

#### `eventbus` - Event Bus
Thread-safe publish-subscribe pattern implementation.

```go
import "github.com/lazyfury/bowlutils/eventbus"

bus := eventbus.New()

// Subscribe
id, ch := bus.Subscribe("user.created", 10)
go func() {
    for event := range ch {
        // Handle event
    }
}()

// Publish
bus.Publish("user.created", userData)

// Unsubscribe
bus.Unsubscribe("user.created", id)
```

**Features:**
- Non-blocking publish
- Multiple subscribers support
- Auto-drop messages on full buffer
- Thread-safe

#### `isvlid` - Data Validation
Enhanced validation based on validator/v10 with custom conditions support.

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
    // Handle validation error
}
```

**Built-in Validators:**
- `Required()` - Required field
- `IsEnum()` / `IsOneOf()` - Enum values
- `IsValidPhone()` - Phone number
- `IsValidEmail()` - Email address
- `Min()` / `Max()` - Numeric range
- `Default()` - Default value

### Infrastructure

#### `db` - Database Connection
Simplified database connection management, supporting MySQL and PostgreSQL.

```go
import "github.com/lazyfury/bowlutils/db"

db := db.NewDB("postgres", dsn)
// or
db := db.NewDB("mysql", dsn)
```

#### `email` - Email Sending
SMTP email sender supporting synchronous and asynchronous sending.

```go
import "github.com/lazyfury/bowlutils/email"

// Configuration
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

// Send email
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Hello",
    Body:    "Plain text body",
    HTML:    "<h1>HTML body</h1>",
}

err := sender.Send(context.Background(), msg)
```

#### `resp` - HTTP Response
Unified response format based on Gin.

```go
import "github.com/lazyfury/bowlutils/resp"

// Success response
resp.Ok(c, data)

// Failure response
resp.Fail[any](c, "Operation failed")

// Error response
resp.Error(c, 500, "Server error", nil)

// Other responses
resp.NotFound[any](c, "Resource not found")
resp.Unauthorized[any](c, "Unauthorized")
resp.Forbidden[any](c, "Forbidden")
```

#### `openapi` - OpenAPI Documentation
OpenAPI 3.0 document generation tool.

```go
import "github.com/lazyfury/bowlutils/openapi"

doc := openapi.NewDocument("3.0.0").
    WithInfo("My API", "1.0.0").
    AddServer(openapi.Server{URL: "http://localhost:8080"})

// Add paths and operations
doc.AddOperation("/users", "get", openapi.Operation{
    Summary: "Get user list",
    Responses: openapi.NewResponses(
        openapi.NewResponseFrom(200, "Success", []User{}),
    ),
})
```

#### `viperinit` - Configuration Management
Simplified Viper configuration initialization.

```go
import "github.com/lazyfury/bowlutils/viperinit"

v := viperinit.NewViper("config", "yaml", ".")
port := v.GetInt("server.port")
```

#### `utils` - General Utilities
Common helper functions.

```go
import "github.com/lazyfury/bowlutils/utils"

// Type conversion
str := utils.ToString(123)           // "123"
m, _ := utils.ToMap(struct{Name: "John"})

// Zero value check
utils.IsZero(0)    // true
utils.IsEmpty("")  // true

// Default value
value := utils.Def("", "default") // "default"
```

## Dependency Graph

```
Independent packages (no internal dependencies):
├── utils       - General utilities
├── viperinit   - Configuration management
├── db          - Database connection
├── eventbus    - Event bus
├── ioc         - IOC container
├── isvlid      - Data validation
└── openapi     - OpenAPI documentation

Limited dependency packages:
├── crud        → gorm (external)
├── resp        → gin (external)
└── email       → ioc (optional)
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./crud/...
```

## Contributing

Issues and Pull Requests are welcome!

## License

MIT License

## Acknowledgments

Built on top of these excellent open source projects:
- [GORM](https://gorm.io/) - ORM library
- [Gin](https://gin-gonic.com/) - Web framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Validator](https://github.com/go-playground/validator) - Data validation

---

**Note**: This project is under active development, APIs may change. Consider pinning to a specific version for production use.
