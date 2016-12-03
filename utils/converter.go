package utils

import (
	"encoding/base64"
	"errors"
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
	Brunch   string
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
// 	* https://github.com/lastbackend/vendors.git
// 	* git@github.com:lastbackend/vendors.git
func GitUrlParse(url string) (*source, error) {

	var parsingUrl []string = regexp.MustCompile(`^(?:ssh|git|http(?:s)?)(?:@|:\/\/(?:.+@)?)((\w+)\.\w+)(?:\/|:)(.+)(?:\/)(.+)(?:\..+)$`).FindStringSubmatch(url)

	if len(parsingUrl) < 5 {
		return nil, errors.New("can't parse url")
	}

	return &source{
		Resource: parsingUrl[0],
		Hub:      parsingUrl[1],
		Vendor:   parsingUrl[2],
		Owner:    parsingUrl[3],
		Repo:     parsingUrl[4],
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
		parsingNamespace.Brunch = repoAndTag[1]
	}

	return parsingNamespace, nil

}
