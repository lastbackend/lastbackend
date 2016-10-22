package utils

import (
	"encoding/base64"
	"strconv"
	"strings"
)

func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func StringToBool(s string) bool {
	s = strings.ToLower(s)
	if s == "true" || s == "1" || s == "t" {
		return true
	}
	return false
}

func Int64ToInt(i int64) int {
	return StringToInt(strconv.FormatInt(i, 10))
}

func DecodeBase64(s string) string {
	buf, _ := base64.StdEncoding.DecodeString(s)
	return string(buf)
}
