package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request HTTP请求
type Request struct {
	client  *Client
	method  string
	url     string
	headers map[string]string
	query   url.Values
	body    io.Reader
	ctx     context.Context
}

// Header 设置请求头
func (r *Request) Header(key, value string) *Request {
	r.headers[key] = value
	return r
}

// Headers 批量设置请求头
func (r *Request) Headers(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

// Query 设置查询参数
func (r *Request) Query(key, value string) *Request {
	r.query.Add(key, value)
	return r
}

// QueryParams 批量设置查询参数
func (r *Request) QueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.query.Add(k, v)
	}
	return r
}

// Body 设置请求体
func (r *Request) Body(body io.Reader) *Request {
	r.body = body
	return r
}

// JSONBody 设置JSON请求体
func (r *Request) JSONBody(v interface{}) *Request {
	data, err := json.Marshal(v)
	if err != nil {
		// 这里可以考虑返回error，但为了链式调用的流畅性，暂时忽略
		return r
	}
	r.body = bytes.NewReader(data)
	r.Header("Content-Type", "application/json")
	return r
}

// FormBody 设置表单请求体
func (r *Request) FormBody(data map[string]string) *Request {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	r.body = strings.NewReader(form.Encode())
	r.Header("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// Context 设置上下文
func (r *Request) Context(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// Do 执行请求
func (r *Request) Do() (*Response, error) {
	// 构建完整URL
	fullURL := r.url
	if len(r.query) > 0 {
		if strings.Contains(fullURL, "?") {
			fullURL += "&" + r.query.Encode()
		} else {
			fullURL += "?" + r.query.Encode()
		}
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(r.ctx, r.method, fullURL, r.body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	// 执行请求
	return r.client.do(req)
}

// DoJSON 执行请求并解析JSON响应
func (r *Request) DoJSON(v interface{}) error {
	resp, err := r.Do()
	if err != nil {
		return err
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		return resp.Error()
	}

	return resp.JSON(v)
}

// DoString 执行请求并返回字符串响应
func (r *Request) DoString() (string, error) {
	resp, err := r.Do()
	if err != nil {
		return "", err
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		return "", resp.Error()
	}

	return resp.String()
}

// DoBytes 执行请求并返回字节响应
func (r *Request) DoBytes() ([]byte, error) {
	resp, err := r.Do()
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		return nil, resp.Error()
	}

	return resp.Bytes()
}
