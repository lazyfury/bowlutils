package httpclient_test

import (
	"fmt"
	"log"
	"time"

	"github.com/lazyfury/bowlutils/httpclient"
)

// Example_basic 基础用法示例
func Example_basic() {
	client := httpclient.New()

	// GET请求
	resp, err := client.Get("https://api.example.com/users").Do()
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()

	body, _ := resp.String()
	fmt.Println(body)
}

// Example_withOptions 使用配置选项
func Example_withOptions() {
	client := httpclient.New(
		httpclient.WithBaseURL("https://api.example.com"),
		httpclient.WithTimeout(10*time.Second),
		httpclient.WithHeader("X-API-Key", "your-api-key"),
	)

	// 发送请求
	var result map[string]interface{}
	err := client.Get("/users").
		Query("page", "1").
		Query("size", "10").
		DoJSON(&result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

// Example_postJSON POST JSON数据
func Example_postJSON() {
	client := httpclient.New(
		httpclient.WithBaseURL("https://api.example.com"),
	)

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

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

// Example_withRetry 使用重试机制
func Example_withRetry() {
	client := httpclient.New(
		httpclient.WithBaseURL("https://api.example.com"),
		httpclient.WithRetry(3, time.Second, 500, 502, 503, 504),
	)

	resp, err := client.Get("/unstable-endpoint").Do()
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()

	fmt.Println(resp.StatusCode)
}

// Example_withInterceptor 使用拦截器
func Example_withInterceptor() {
	logInterceptor := &httpclient.LogInterceptor{
		Logger: func(format string, args ...interface{}) {
			log.Printf(format, args...)
		},
	}

	client := httpclient.New(
		httpclient.WithBaseURL("https://api.example.com"),
		httpclient.WithInterceptor(logInterceptor),
	)

	resp, err := client.Get("/users").Do()
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()
}

// Example_withAuth 使用认证
func Example_withAuth() {
	// Basic Auth
	client1 := httpclient.New(
		httpclient.WithBaseURL("https://api.example.com"),
		httpclient.WithBasicAuth("username", "password"),
	)

	resp1, _ := client1.Get("/protected").Do()
	defer resp1.Close()

	// Bearer Token
	client2 := httpclient.New(
		httpclient.WithBaseURL("https://api.example.com"),
		httpclient.WithBearerToken("your-token-here"),
	)

	resp2, _ := client2.Get("/protected").Do()
	defer resp2.Close()
}

// Example_formData 提交表单数据
func Example_formData() {
	client := httpclient.New(
		httpclient.WithBaseURL("https://api.example.com"),
	)

	resp, err := client.Post("/login").
		FormBody(map[string]string{
			"username": "user",
			"password": "pass",
		}).
		Do()

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()

	fmt.Println(resp.StatusCode)
}
