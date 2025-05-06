// Log utils
package log

import "encoding/json"

func GetJsonString(v any) string {
	j, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(j)
}
