package httpclient

import (
	"net/http"
	"time"
)

// Option 客户端配置选项
type Option func(*Client)

// WithBaseURL 设置基础URL
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHeader 设置默认请求头
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithHeaders 批量设置默认请求头
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithHTTPClient 设置自定义HTTP客户端
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTransport 设置自定义Transport
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.httpClient.Transport = transport
	}
}

// WithRetry 设置重试配置
func WithRetry(maxRetries int, retryDelay time.Duration, retryOn ...int) Option {
	return func(c *Client) {
		c.retryConfig = &RetryConfig{
			MaxRetries: maxRetries,
			RetryDelay: retryDelay,
			RetryOn:    retryOn,
		}
	}
}

// WithInterceptor 添加拦截器
func WithInterceptor(interceptor Interceptor) Option {
	return func(c *Client) {
		c.interceptors = append(c.interceptors, interceptor)
	}
}

// WithUserAgent 设置User-Agent
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.headers["User-Agent"] = userAgent
	}
}

// WithBasicAuth 设置Basic认证
func WithBasicAuth(username, password string) Option {
	return func(c *Client) {
		c.interceptors = append(c.interceptors, &basicAuthInterceptor{
			username: username,
			password: password,
		})
	}
}

// WithBearerToken 设置Bearer Token
func WithBearerToken(token string) Option {
	return func(c *Client) {
		c.headers["Authorization"] = "Bearer " + token
	}
}

// basicAuthInterceptor Basic认证拦截器
type basicAuthInterceptor struct {
	username string
	password string
}

func (i *basicAuthInterceptor) Before(req *http.Request) error {
	req.SetBasicAuth(i.username, i.password)
	return nil
}

func (i *basicAuthInterceptor) After(resp *http.Response) error {
	return nil
}
