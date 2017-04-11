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

package converter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type source struct {
	Resource string
	Hub      string
	Repo     string
	Owner    string
	Vendor   string
	Branch   string
}

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

// Parse incoming string git url in source type
// Ex:
// 	* https://github.com/lastbackend/lastbackend.git
// 	* git@github.com:lastbackend/lastbackend.git
func GitUrlParse(url string) (*source, error) {

	var (
		parse  = strings.Split(url, "#")
		branch = "master"
	)

	if len(parse) == 2 {
		branch = parse[1]
	}

	var match []string = regexp.MustCompile(`^(?:ssh|git|http(?:s)?)(?:@|:\/\/(?:.+@)?)((\w+)\.\w+)(?:\/|:)(.+)(?:\/)(.+)(?:\..+)$`).FindStringSubmatch(parse[0])

	if len(match) < 5 {
		return nil, errors.New("can't parse url")
	}

	return &source{
		Resource: match[0],
		Hub:      match[1],
		Vendor:   match[2],
		Owner:    match[3],
		Repo:     match[4],
		Branch:   branch,
	}, nil

}

func DockerNamespaceParse(namespace string) (*source, error) {

	var parsingNamespace *source = new(source)
	parsingNamespace.Vendor = "dockerhub"

	splitStr := strings.Split(namespace, "/")
	switch len(splitStr) {
	case 1:
		parsingNamespace.Repo = splitStr[0]
		return parsingNamespace, nil
	case 2:
		parsingNamespace.Owner = splitStr[0]
	case 3:
		parsingNamespace.Hub = splitStr[0]
		parsingNamespace.Owner = splitStr[1]
	default:
		return nil, errors.New("can't parse url")
	}
	repoAndTag := strings.Split(splitStr[len(splitStr)-1], ":")
	parsingNamespace.Repo = repoAndTag[0]
	if len(repoAndTag) == 2 {
		parsingNamespace.Branch = repoAndTag[1]
	}

	return parsingNamespace, nil

}

func EnforcePtr(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		if v.Kind() == reflect.Invalid {
			return reflect.Value{}, fmt.Errorf("expected pointer, but got invalid kind")
		}
		return reflect.Value{}, fmt.Errorf("expected pointer, but got %v type", v.Type())
	}
	if v.IsNil() {
		return reflect.Value{}, fmt.Errorf("expected pointer, but got nil")
	}
	return v.Elem(), nil
}
