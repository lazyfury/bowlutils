package httpclient

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client HTTP客户端
type Client struct {
	baseURL      string
	httpClient   *http.Client
	headers      map[string]string
	interceptors []Interceptor
	retryConfig  *RetryConfig
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries int           // 最大重试次数
	RetryDelay time.Duration // 重试延迟
	RetryOn    []int         // 需要重试的状态码
}

// New 创建新的HTTP客户端
func New(options ...Option) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers:      make(map[string]string),
		interceptors: make([]Interceptor, 0),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// Get 发送GET请求
func (c *Client) Get(url string) *Request {
	return c.NewRequest(http.MethodGet, url)
}

// Post 发送POST请求
func (c *Client) Post(url string) *Request {
	return c.NewRequest(http.MethodPost, url)
}

// Put 发送PUT请求
func (c *Client) Put(url string) *Request {
	return c.NewRequest(http.MethodPut, url)
}

// Delete 发送DELETE请求
func (c *Client) Delete(url string) *Request {
	return c.NewRequest(http.MethodDelete, url)
}

// Patch 发送PATCH请求
func (c *Client) Patch(url string) *Request {
	return c.NewRequest(http.MethodPatch, url)
}

// NewRequest 创建新请求
func (c *Client) NewRequest(method, path string) *Request {
	fullURL := c.buildURL(path)
	return &Request{
		client:  c,
		method:  method,
		url:     fullURL,
		headers: c.copyHeaders(),
		query:   url.Values{},
		ctx:     context.Background(),
	}
}

// buildURL 构建完整URL
func (c *Client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

// copyHeaders 复制headers
func (c *Client) copyHeaders() map[string]string {
	headers := make(map[string]string)
	for k, v := range c.headers {
		headers[k] = v
	}
	return headers
}

// do 执行HTTP请求
func (c *Client) do(req *http.Request) (*Response, error) {
	var resp *http.Response
	var err error

	// 应用拦截器
	for _, interceptor := range c.interceptors {
		if err := interceptor.Before(req); err != nil {
			return nil, err
		}
	}

	// 执行请求（包含重试逻辑）
	if c.retryConfig != nil && c.retryConfig.MaxRetries > 0 {
		resp, err = c.doWithRetry(req)
	} else {
		resp, err = c.httpClient.Do(req)
	}

	if err != nil {
		return nil, err
	}

	// 应用拦截器
	for _, interceptor := range c.interceptors {
		if err := interceptor.After(resp); err != nil {
			return nil, err
		}
	}

	return newResponse(resp), nil
}

// doWithRetry 带重试的请求执行
func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.retryConfig.MaxRetries; i++ {
		// 克隆请求以支持重试
		reqClone := req.Clone(req.Context())

		resp, err = c.httpClient.Do(reqClone)
		if err != nil {
			if i < c.retryConfig.MaxRetries {
				time.Sleep(c.retryConfig.RetryDelay)
				continue
			}
			return nil, err
		}

		// 检查是否需要重试
		if c.shouldRetry(resp.StatusCode) && i < c.retryConfig.MaxRetries {
			resp.Body.Close()
			time.Sleep(c.retryConfig.RetryDelay)
			continue
		}

		break
	}

	return resp, err
}

// shouldRetry 判断是否应该重试
func (c *Client) shouldRetry(statusCode int) bool {
	if c.retryConfig == nil || len(c.retryConfig.RetryOn) == 0 {
		return false
	}

	for _, code := range c.retryConfig.RetryOn {
		if code == statusCode {
			return true
		}
	}

	return false
}

// AddInterceptor 添加拦截器
func (c *Client) AddInterceptor(interceptor Interceptor) {
	c.interceptors = append(c.interceptors, interceptor)
}
