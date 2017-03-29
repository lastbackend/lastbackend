package utils

import (
	"encoding/base64"
)

func DecodeBase64(s string) string {
	buf, _ := base64.StdEncoding.DecodeString(s)
	return string(buf)
}
