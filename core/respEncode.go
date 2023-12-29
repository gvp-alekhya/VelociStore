package core

import (
	"fmt"
)

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", value))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value.(string)), value))
	case int, int8, int16, int32, int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	case []string:
		{
			buf := make([]byte, 0)

			for _, val := range value.([]string) {
				buf = append(buf, Encode(val, false)...)
			}
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value.([]string)), string(buf)))
		}
	}

	return []byte{}
}
