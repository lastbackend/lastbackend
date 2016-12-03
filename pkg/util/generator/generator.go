package generator

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	mathrand "math/rand"
	"strings"
	"time"
)

const RANDOM_PASS_LEN = 10

func GetUUIDV4() string {
	return uuid.NewV4().String()
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

func GenerateRandomPassword() (string, error) {
	mathrand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	result := make([]byte, RANDOM_PASS_LEN)
	for i := 0; i < RANDOM_PASS_LEN; i++ {
		result[i] = chars[mathrand.Intn(len(chars))]
	}
	return string(result), nil
}

func GenerateToken(n int) string {

	var str string
	var index int

	for len(str)-(index*4) < n {
		str += GetUUIDV4()
		index++
	}

	str = strings.Replace(str, "-", "", -1)

	return str[:n]
}
