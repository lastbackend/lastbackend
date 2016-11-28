package utils

import (
	"encoding/base64"
	"github.com/asaskevich/govalidator"
	"regexp"
	"strconv"
	"strings"
)

func IsBool(s string) bool {
	s = strings.ToLower(s)
	if s == "true" || s == "1" || s == "t" || s == "false" || s == "0" || s == "f" {
		return true
	}
	return false
}

func IsEmail(s string) bool {
	return govalidator.IsEmail(s)
}

func IsUsername(s string) bool {
	reg, _ := regexp.Compile("[A-Za-z0-9]+(?:[_-][A-Za-z0-9]+)*")
	str := reg.FindStringSubmatch(s)
	if len(str) == 1 && str[0] == s && len(s) >= 4 && len(s) <= 64 {
		return true
	}
	return false
}

func IsPassword(s string) bool {
	return len(s) > 6
}

func IsServiceName(s string) bool {
	reg, _ := regexp.Compile("[a-z0-9]+(?:[_-][a-z0-9]+)*")
	str := reg.FindStringSubmatch(s)
	if len(str) == 1 && str[0] == s && len(s) >= 4 && len(s) <= 64 {
		return true
	}
	return false
}

func IsProjectName(s string) bool {
	reg, _ := regexp.Compile("[a-z0-9]+(?:[_-][a-z0-9]+)*")
	str := reg.FindStringSubmatch(s)
	if len(str) == 1 && str[0] == s && len(s) >= 4 && len(s) <= 64 {
		return true
	}
	return false
}

func IsIP(ip string) bool {
	return govalidator.IsIP(ip)
}

func IsMac(mac string) bool {
	return govalidator.IsMAC(mac)
}

func IsUUID(uuid string) bool {
	return govalidator.IsUUIDv4(uuid)
}

func IsRole(role string) bool {
	switch role {
	case "member":
		return true
	case "admin":
		return true
	}
	return false
}

func IsPort(port int) bool {
	return govalidator.IsPort(strconv.Itoa(port))
}

func IsDomain(domain string) bool { // TODO domait Validator
	return true
}

func IsProtocol(protocol string) bool {
	correctProtocols := []string{"tcp", "udp"}
	for _, correctProtocol := range correctProtocols {
		if strings.EqualFold(correctProtocol, protocol) {
			return true
		}
	}
	return false
}

func IsPublicKey(key string) bool {

	var splited = strings.SplitN(key, " ", 3)
	if len(splited) < 2 {
		return false
	}

	var alg = strings.TrimSpace(splited[0])
	var cb64 = strings.TrimSpace(splited[1])

	_, err := base64.StdEncoding.DecodeString(cb64)
	if err != nil {
		return false
	}

	switch alg {
	case "ssh-rsa":
		return true
	case "ssh-dss":
		return true
	case "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521":
		return true
	}

	return false
}
