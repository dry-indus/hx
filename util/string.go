package util

import "fmt"

func ValueString(m map[interface{}]interface{}, key string) string {
	val := m[key]
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%s", val)
}

func DefaultString(s, def string) string {
	if len(s) != 0 {
		return s
	}

	return def
}
