package interfaces

import "golang.org/x/oauth2"

// Model

type OAuth2 struct {
	ClientID     string
	ClientSecret string
	State        string
	RedirectUri  string
}

// Types

// Interface

type IOAuth2 interface {
	GetToken(code string) (*oauth2.Token, error)
	RefreshToken(token *oauth2.Token) (*oauth2.Token, bool, error)
}
