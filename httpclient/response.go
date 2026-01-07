package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response HTTP响应
type Response struct {
	*http.Response
	body []byte
}

// newResponse 创建响应对象
func newResponse(resp *http.Response) *Response {
	return &Response{
		Response: resp,
	}
}

// Bytes 获取响应体字节
func (r *Response) Bytes() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.body = body
	return body, nil
}

// String 获取响应体字符串
func (r *Response) String() (string, error) {
	body, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// JSON 解析JSON响应
func (r *Response) JSON(v interface{}) error {
	body, err := r.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// IsSuccess 判断请求是否成功
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// Error 返回错误信息
func (r *Response) Error() error {
	if r.IsSuccess() {
		return nil
	}

	body, _ := r.String()
	return fmt.Errorf("HTTP %d: %s", r.StatusCode, body)
}

// Close 关闭响应体
func (r *Response) Close() error {
	if r.Response != nil && r.Response.Body != nil {
		return r.Response.Body.Close()
	}
	return nil
}
