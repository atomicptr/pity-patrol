package util

import (
	"encoding/json"
)

func ToPrettyString(data any) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}
