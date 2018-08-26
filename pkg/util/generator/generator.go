//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package generator

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"strings"
	"time"
)

const RANDOM_PASS_LEN = 10

func GetUUIDV4() string {
	return uuid.NewV4().String()
}

func UnixTimestamp() int {
	return int(time.Now().Unix())
}

func GenerateSalt(password string) (string, error) {
	buf := make([]byte, 10, 10+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		fmt.Printf("random read failed: %v", err)
		return "", err
	}

	hash := sha1.New()
	_, err = hash.Write(buf)
	if err != nil {
		return "", err
	}

	_, err = hash.Write([]byte(password))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(buf)), nil
}

func GeneratePassword(password string, salt string) (string, error) {
	pass := []byte(password + salt)
	// Hashing the password with the default cost of 10
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func GenerateGravatar(email string) string {
	m := md5.New()
	if _, err := io.WriteString(m, strings.ToLower(email)); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", m.Sum(nil))
}

func GenerateRandomString(n int) string {

	var str string
	var index int

	for len(str)-(index*4) < n {
		str += GetUUIDV4()
		index++
	}

	str = strings.Replace(str, "-", "", -1)

	return str[:n]
}
