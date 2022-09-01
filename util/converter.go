/**
 * @Author: lzw5399
 * @Date: 2020/8/1 16:21
 * @Desc: 用于map jsonStr互相转换
 */
package util

import (
	"encoding/json"
)

func JsonStrToMap(jsonStr string) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func MapToJsonStr(m map[string]string) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		return "", nil
	}

	return string(jsonByte), nil
}