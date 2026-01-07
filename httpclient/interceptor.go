package httpclient

import (
	"net/http"
)

// Interceptor 拦截器接口
type Interceptor interface {
	// Before 请求前拦截
	Before(req *http.Request) error
	// After 响应后拦截
	After(resp *http.Response) error
}

// LogInterceptor 日志拦截器
type LogInterceptor struct {
	Logger func(format string, args ...interface{})
}

func (l *LogInterceptor) Before(req *http.Request) error {
	if l.Logger != nil {
		l.Logger("[HTTP] Request: %s %s", req.Method, req.URL.String())
	}
	return nil
}

func (l *LogInterceptor) After(resp *http.Response) error {
	if l.Logger != nil {
		l.Logger("[HTTP] Response: %s %d", resp.Request.URL.String(), resp.StatusCode)
	}
	return nil
}
