//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package validator

import (
	"encoding/base64"
	"fmt"
	"github.com/asaskevich/govalidator"
	"regexp"
	"strconv"
	"strings"
	"reflect"
)

func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}

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

func IsDockerNamespace(s string) bool {
	// Todo: check valid docker namespace
	return true
}

func IsServiceName(s string) bool {
	reg, _ := regexp.Compile("[a-z]+(?:[_-][a-z0-9]+)*")
	str := reg.FindStringSubmatch(s)
	if len(str) == 1 && str[0] == s && len(s) >= 4 && len(s) <= 64 {
		return true
	}
	return false
}

func IsProjectName(s string) bool {
	reg, _ := regexp.Compile("[a-z]+(?:[_-][a-z0-9]+)*")
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

// Check incoming string on git valid utl
// Ex:
// 	* https://github.com/lastbackend/lastbackend.git
// 	* git@github.com:lastbackend/lastbackend.git
func IsGitUrl(url string) bool {
	res, err := regexp.MatchString(`^(?:ssh|git|http(?:s)?)(?:@|:\/\/(?:.+@)?)((\w+)\.\w+)(?:\/|:)(.+)(?:\/)(.+)(?:\..+)$`, url)
	if err != nil {
		return false
	}

	return res
}

func validateDockerRepositoryName(repositoryName string) error {
	var (
		namespace string
		name      string
	)
	nameParts := strings.SplitN(repositoryName, "/", 2)
	if len(nameParts) < 2 {
		namespace = "library"
		name = nameParts[0]
	} else {
		namespace = nameParts[0]
		name = nameParts[1]
	}
	validNamespace := regexp.MustCompile(`^([a-z0-9_]{4,30})$`)
	if !validNamespace.MatchString(namespace) {
		return fmt.Errorf("Invalid namespace name (%s), only [a-z0-9_] are allowed, size between 4 and 30", namespace)
	}
	validRepo := regexp.MustCompile(`^([a-zA-Z0-9-_.]+)$`)
	if !validRepo.MatchString(name) {
		return fmt.Errorf("Invalid repository name (%s), only [a-zA-Z0-9-_.] are allowed", name)
	}
	return nil
}