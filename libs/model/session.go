package model

import "errors"

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	SESSION_TOKEN_SECRET = "LB:TokeN_#2015#_SecreT:KEY"
)

var (
	ErrUnexpectedSigninMethod = errors.New("UNEXPECTED_SIGNIN_METHD")
	ErrSessionTokenHasNoUID   = errors.New("NO_UID_IN_TOKEN")
	ErrSessionTokenHasNoEXP   = errors.New("NO_EXP_IN_TOKEN")
	ErrSessionTokenHasNoJTI   = errors.New("NO_JTI_IN_TOKEN")
)

// Generate new Session pointer structure
func NewSession(uid, oid, username, email string) *Session {
	return &Session{
		Uid:      uid,
		Oid:      oid,
		Username: username,
		Email:    email,
	}
}

type Session struct {
	Uid      string // session user id
	Oid      string // session organization id
	Username string // session username
	Email    string // session email
}

// Decode - decode token to session object
func (s *Session) Decode(token string) error {

	payload, err := jwt.Parse(token, func(payload *jwt.Token) (interface{}, error) {
		result := []byte(SESSION_TOKEN_SECRET)

		err := func(token *jwt.Token) error {

			claims := token.Claims.(jwt.MapClaims)

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return ErrUnexpectedSigninMethod
			}

			if claims["uid"] == nil {
				return ErrSessionTokenHasNoUID
			}

			if claims["exp"] == nil {
				return ErrSessionTokenHasNoEXP
			}

			if claims["jti"] == nil {
				return ErrSessionTokenHasNoJTI
			}

			return nil
		}(payload)

		return result, err
	})

	if err != nil || !payload.Valid {
		return err
	}

	claims := payload.Claims.(jwt.MapClaims)

	if _, ok := claims["oid"]; ok {
		s.Oid = claims["oid"].(string)
	}

	if _, ok := claims["uid"]; ok {
		s.Uid = claims["uid"].(string)
	}

	if _, ok := claims["em"]; ok {
		s.Email = claims["em"].(string)
	}

	if _, ok := claims["user"]; ok {
		s.Username = claims["user"].(string)
	}

	return nil
}

// Encode - encode session structure to jwt token
func (s *Session) Encode() (string, error) {

	context := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  s.Uid,
		"oid":  s.Oid,
		"em":   s.Email,
		"user": s.Username,
		"jti":  time.Now().Add(time.Hour * 2232).Unix(),
		"exp":  time.Now().Add(time.Hour * 2232).Unix(),
	})

	return context.SignedString([]byte(SESSION_TOKEN_SECRET))
}
