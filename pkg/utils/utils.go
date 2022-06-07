package utils

import (
	"dingdong/pkg/textual"
)

func FilterFields(data map[string]any, fields []string) map[string]any {
	res := make(map[string]any)
	for k, v := range data {
		if !textual.InArray(k, fields) {
			continue
		}
		res[k] = v
	}
	return res
}

func ConvertMapAny(data map[string]string) map[string]any {
	res := make(map[string]any)
	for k, v := range data {
		res[k] = v
	}
	return res
}
