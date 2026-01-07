# HTTP Client

一个简单、通用且易用的 Go HTTP 客户端库，提供链式调用、拦截器、重试机制等特性。

## 特性

- 🔗 **链式调用** - 流畅的 API 设计
- 🔄 **自动重试** - 可配置的重试策略
- 🎯 **拦截器** - 请求/响应拦截器支持
- ⏱️ **超时控制** - 灵活的超时配置
- 🔐 **认证支持** - Basic Auth、Bearer Token
- 📝 **多种请求体** - JSON、表单、自定义
- 🎨 **简洁易用** - 简单明了的 API

## 安装

```bash
go get github.com/lazyfury/bowlutils/httpclient
```

## 快速开始

### 基础用法

```go
import "github.com/lazyfury/bowlutils/httpclient"

// 创建客户端
client := httpclient.New()

// GET 请求
resp, err := client.Get("https://api.example.com/users").Do()
if err != nil {
    log.Fatal(err)
}
defer resp.Close()

// 获取响应内容
body, _ := resp.String()
fmt.Println(body)
```

### 使用配置选项

```go
client := httpclient.New(
    httpclient.WithBaseURL("https://api.example.com"),
    httpclient.WithTimeout(10*time.Second),
    httpclient.WithHeader("X-API-Key", "your-api-key"),
)
```

### GET 请求（带查询参数）

```go
var result map[string]interface{}
err := client.Get("/users").
    Query("page", "1").
    Query("size", "10").
    DoJSON(&result)
```

### POST JSON 数据

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{
    Name:  "John Doe",
    Email: "john@example.com",
}

var response map[string]interface{}
err := client.Post("/users").
    JSONBody(user).
    DoJSON(&response)
```

### 表单提交

```go
resp, err := client.Post("/login").
    FormBody(map[string]string{
        "username": "user",
        "password": "pass",
    }).
    Do()
```

## 配置选项

### 基础配置

```go
// 设置基础 URL
httpclient.WithBaseURL("https://api.example.com")

// 设置超时时间
httpclient.WithTimeout(10 * time.Second)

// 设置默认请求头
httpclient.WithHeader("X-Custom-Header", "value")
httpclient.WithHeaders(map[string]string{
    "X-Header-1": "value1",
    "X-Header-2": "value2",
})

// 设置 User-Agent
httpclient.WithUserAgent("MyApp/1.0")
```

### 认证配置

```go
// Basic Auth
httpclient.WithBasicAuth("username", "password")

// Bearer Token
httpclient.WithBearerToken("your-token-here")
```

### 重试配置

```go
// 最大重试3次，每次延迟1秒，对500、502、503、504状态码重试
httpclient.WithRetry(3, time.Second, 500, 502, 503, 504)
```

### 拦截器

```go
// 日志拦截器
logInterceptor := &httpclient.LogInterceptor{
    Logger: func(format string, args ...interface{}) {
        log.Printf(format, args...)
    },
}

client := httpclient.New(
    httpclient.WithInterceptor(logInterceptor),
)
```

## 自定义拦截器

```go
type CustomInterceptor struct{}

func (c *CustomInterceptor) Before(req *http.Request) error {
    // 请求前处理
    req.Header.Set("X-Request-ID", generateRequestID())
    return nil
}

func (c *CustomInterceptor) After(resp *http.Response) error {
    // 响应后处理
    log.Printf("Status: %d", resp.StatusCode)
    return nil
}

// 使用自定义拦截器
client := httpclient.New(
    httpclient.WithInterceptor(&CustomInterceptor{}),
)
```

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/lazyfury/bowlutils/httpclient"
)

func main() {
    // 创建配置完整的客户端
    client := httpclient.New(
        httpclient.WithBaseURL("https://api.github.com"),
        httpclient.WithTimeout(15*time.Second),
        httpclient.WithUserAgent("MyApp/1.0"),
        httpclient.WithHeader("Accept", "application/json"),
        httpclient.WithRetry(3, time.Second, 500, 502, 503),
        httpclient.WithInterceptor(&httpclient.LogInterceptor{
            Logger: log.Printf,
        }),
    )
    
    // 发送请求
    var result map[string]interface{}
    err := client.Get("/users/octocat").
        Header("X-Custom-Header", "value").
        DoJSON(&result)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User: %v\n", result)
}
```

## API 文档

### Client 方法

- `New(options ...Option) *Client` - 创建新客户端
- `Get(url string) *Request` - 创建 GET 请求
- `Post(url string) *Request` - 创建 POST 请求
- `Put(url string) *Request` - 创建 PUT 请求
- `Delete(url string) *Request` - 创建 DELETE 请求
- `Patch(url string) *Request` - 创建 PATCH 请求

### Request 方法

- `Header(key, value string) *Request` - 设置请求头
- `Headers(headers map[string]string) *Request` - 批量设置请求头
- `Query(key, value string) *Request` - 设置查询参数
- `QueryParams(params map[string]string) *Request` - 批量设置查询参数
- `Body(body io.Reader) *Request` - 设置请求体
- `JSONBody(v interface{}) *Request` - 设置 JSON 请求体
- `FormBody(data map[string]string) *Request` - 设置表单请求体
- `Context(ctx context.Context) *Request` - 设置上下文
- `Do() (*Response, error)` - 执行请求
- `DoJSON(v interface{}) error` - 执行请求并解析 JSON
- `DoString() (string, error)` - 执行请求并返回字符串
- `DoBytes() ([]byte, error)` - 执行请求并返回字节数组

### Response 方法

- `Bytes() ([]byte, error)` - 获取响应字节
- `String() (string, error)` - 获取响应字符串
- `JSON(v interface{}) error` - 解析 JSON 响应
- `IsSuccess() bool` - 判断请求是否成功
- `Error() error` - 获取错误信息
- `Close() error` - 关闭响应体

## 许可证

MIT License
