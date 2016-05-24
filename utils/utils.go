package utils

import (
	"crypto/sha1"
	"encoding/base64"
)

func Hash(data []byte) string {
	hash := sha1.Sum(data)

	var hashString string

	for i := 0; i < len(hash); i++ {
		hashString += base64.URLEncoding.EncodeToString(hash[:i])
	}

	return hashString
}
