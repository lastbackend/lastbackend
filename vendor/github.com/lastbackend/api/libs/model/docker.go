package model

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base32"
	"github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"io/ioutil"
	"strings"
	"time"
)

var (
	SERVICE string
	ISSUER  string
	KEYPATH string
)

type DockerRepository struct {
	StarCount int64
	PullCount int64
	Hub       string
	Owner     string
	Name      string
	Desc      string
	Automated bool
	Official  bool
}

type DockerRepositories []DockerRepository

type DockerTag struct {
	Name        string
	ID          int64
	Size        int64
	Repo        int64
	Creator     int64
	LastUpdater int64
	V2          bool
	ImageID     int64
	Platforms   []int64
	LastUpdated time.Time
}

type DockerTags []DockerTag

type AccountInfo struct {
	Username string
	Password string
	Token    string
}

type Scope struct {
	Type      string
	Name      string
	Namespace string
	Actions   []string
}

type Scopes []Scope

type AccessItem struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Actions []string `json:"actions"`
}

type AccessItems []AccessItem

type JwtToken struct {
	Account    string
	Service    string
	Scope      *Scopes
	PrivateKey *rsa.PrivateKey
}

func NewJwtToken(account string, scope *Scopes) (*JwtToken, error) {

	token := new(JwtToken)
	token.Account = account
	token.Service = SERVICE
	token.Scope = scope

	pem, err := ioutil.ReadFile(KEYPATH)
	if err != nil {
		return token, err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return token, err
	}

	token.PrivateKey = privKey

	return token, err
}

func (t *JwtToken) Claim(account string, scopes Scopes) (map[string]interface{}, error) {

	claims := make(map[string]interface{})

	claims["iss"] = ISSUER
	claims["sub"] = account
	claims["aud"] = SERVICE

	claims["exp"] = time.Now().Add(time.Minute * 240).Unix()
	claims["nbf"] = time.Now().Add(time.Minute * -240).Unix()
	claims["iat"] = time.Now().Unix()

	u := uuid.NewV4()
	claims["jti"] = u.String()

	accessItems := AccessItems{}

	for _, s := range scopes {
		action := AccessItem{
			Type:    s.Type,
			Name:    s.Name,
			Actions: s.Actions,
		}

		accessItems = append(accessItems, action)
	}

	claims["access"] = accessItems

	return claims, nil
}

func (t *JwtToken) SignedString(claim map[string]interface{}, privateKey *rsa.PrivateKey) (string, error) {

	var err error

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = claim

	derBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return "", err
	}

	hasher := crypto.SHA256.New()
	_, err = hasher.Write(derBytes)
	if err != nil {
		return "", err
	}

	s := strings.TrimRight(base32.StdEncoding.EncodeToString(hasher.Sum(nil)[:30]), "=")
	var buf bytes.Buffer
	var i int
	for i = 0; i < len(s)/4-1; i++ {
		start := i * 4
		end := start + 4
		_, err = buf.WriteString(s[start:end] + ":")
		if err != nil {
			return "", err
		}
	}

	_, err = buf.WriteString(s[i*4:])
	if err != nil {
		return "", err
	}

	token.Header["kid"] = buf.String()

	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}
