package utils

import (
	"encoding/json"
	"strconv"
)

func ToString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		j, _ := json.Marshal(t)
		return string(j)
	}
}

// any to map[string]any
func ToMap(v any) (map[string]any, error) {
	m := make(map[string]any)
	j, err := json.Marshal(v)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}
