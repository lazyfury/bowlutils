package utils

import (
	"reflect"
	"strings"
)

func IsZero(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)

	// 如果是指针，需要判断指针是否为 nil 或指向的值是否为 zero
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
		if rv.IsNil() {
			return true
		}
	}

	return reflect.ValueOf(v).IsZero()
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t) == ""
	}
	return false
}
