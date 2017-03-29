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

package types

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	"time"
)

var (
	ErrUnexpectedSigninMethod = errors.New("UNEXPECTED_SIGNIN_METHD")
	ErrSessionTokenHasNoEXP   = errors.New("NO_EXP_IN_TOKEN")
	ErrSessionTokenHasNoJTI   = errors.New("NO_JTI_IN_TOKEN")
)

// Generate new Session pointer structure
func NewSession(username, email string) *Session {
	return &Session{
		Username: username,
		Email:    email,
	}
}

type Session struct {
	Username string // session username
	Email    string // session email
}

// Decode - decode token to session object
func (s *Session) Decode(token string) error {

	var cfg = config.Get()

	payload, err := jwt.Parse(token, func(payload *jwt.Token) (interface{}, error) {
		result := []byte(cfg.TokenSecret)

		err := func(token *jwt.Token) error {

			claims := token.Claims.(jwt.MapClaims)

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return ErrUnexpectedSigninMethod
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

	if _, ok := claims["em"]; ok {
		s.Email = claims["em"].(string)
	}

	if _, ok := claims["un"]; ok {
		s.Username = claims["un"].(string)
	}

	return nil
}

// Encode - encode session structure to jwt token
func (s *Session) Encode() (string, error) {

	var cfg = config.Get()

	context := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"em":  s.Email,
		"un":  s.Username,
		"jti": time.Now().Add(time.Hour * 2232).Unix(),
		"exp": time.Now().Add(time.Hour * 2232).Unix(),
	})

	return context.SignedString([]byte(cfg.TokenSecret))
}
