package interfaces

import "golang.org/x/oauth2"

type IAuth interface {
	IVendor
	IOAuth2

	GetUser(token *oauth2.Token) (*User, error)
}
