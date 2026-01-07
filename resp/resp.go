package resp

import (
	"encoding/json"
	"net/http"
)

var (
	SuccessCode = 200
	SuccessMsg  = "ok"
)

var (
	BusinessErrCode = 1000
	BusinessErrMsg  = "business error"

	RequirePhoneErrCode = 1001
	RequirePhoneErrMsg  = "require phone"
)

type resp[T any] struct {
	w      http.ResponseWriter
	status int
	code   int
	msg    string
	data   T
}

type option[T any] func(*resp[T])

func WithStatus[T any](s int) option[T] {
	return func(r *resp[T]) {
		r.status = s
	}
}

func WithCode[T any](x int) option[T] {
	return func(r *resp[T]) {
		r.code = x
	}
}

func WithMsg[T any](m string) option[T] {
	return func(r *resp[T]) {
		r.msg = m
	}
}

func WithData[T any](d T) option[T] {
	return func(r *resp[T]) {
		r.data = d
	}
}

func New[T any](w http.ResponseWriter, opts ...option[T]) *resp[T] {
	r := &resp[T]{w: w, status: 200, code: SuccessCode, msg: SuccessMsg}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(r)
	}
	return r
}

func (b *resp[T]) Send() error {
	b.w.Header().Set("Content-Type", "application/json")
	b.w.WriteHeader(b.status)
	err := json.NewEncoder(b.w).Encode(map[string]any{
		"code": b.code,
		"msg":  b.msg,
		"data": b.data,
	})
	return err
}

// ===== shortcuts =============================

// Ok 成功响应
func Ok[T any](w http.ResponseWriter, data T) error {
	return New(w, WithData(data), WithMsg[T](SuccessMsg), WithCode[T](SuccessCode), WithStatus[T](200)).Send()
}

// Fail 失败响应
func Fail[T any](w http.ResponseWriter, msg string, opts ...option[T]) error {
	opts = append([]option[T]{WithMsg[T](msg), WithCode[T](int(BusinessErrCode)), WithStatus[T](400)}, opts...)
	return New(w, opts...).Send()
}

// Error 错误响应
func Error[T any](w http.ResponseWriter, code int, err string, data T, opts ...option[T]) error {
	opts = append([]option[T]{WithMsg[T](err), WithCode[T](int(code)), WithStatus[T](500), WithData(data)}, opts...)
	return New(w, opts...).Send()
}

// NotFound 404 响应
func NotFound[T any](w http.ResponseWriter, msg string, opts ...option[T]) error {
	opts = append([]option[T]{WithMsg[T](msg), WithCode[T](int(BusinessErrCode)), WithStatus[T](404)}, opts...)
	return New(w, opts...).Send()
}

// Unauthorized 401 响应
func Unauthorized[T any](w http.ResponseWriter, msg string, opts ...option[T]) error {
	opts = append([]option[T]{WithMsg[T](msg), WithCode[T](int(BusinessErrCode)), WithStatus[T](401)}, opts...)
	return New(w, opts...).Send()
}

// Forbidden 403 响应
func Forbidden[T any](w http.ResponseWriter, msg string, opts ...option[T]) error {
	opts = append([]option[T]{WithMsg[T](msg), WithCode[T](int(BusinessErrCode)), WithStatus[T](403)}, opts...)
	return New(w, opts...).Send()
}
