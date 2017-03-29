package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/vendors/interfaces"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	API_URL   = "https://api.weixin.qq.com"
	TOKEN_URL = "https://api.weixin.qq.com/sns/oauth2/access_token"
)

type WeChat struct {
	clientID       string
	clientSecretID string
	redirectURI    string
	vendor         string
	host           string
}

type Error struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

type AccessToken struct {
	AccessToken  string   `json:"access_token"`
	CreateAt     int64    `json:"create_at"`
	ExpiresIn    int64    `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	OpenId       string   `json:"openid"`
	UnionId      string   `json:"unionid,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

type RefreshToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scopes       string `json:"scope,omitempty"`
}

func GetClient(clientID, clientSecretID, redirectURI string) *WeChat {
	return &WeChat{
		clientID:       clientID,
		clientSecretID: clientSecretID,
		redirectURI:    redirectURI,
		vendor:         "wechat",
		host:           "open.weixin.qq.com",
	}
}

func (w WeChat) GetToken(code string) (token *oauth2.Token, err error) {

	now := time.Now().Unix()

	query := make(url.Values)
	query.Set("grant_type", "authorization_code")
	query.Set("appid", w.clientID)
	query.Set("secret", w.clientSecretID)
	query.Set("code", code)

	var uri = fmt.Sprintf("%s/sns/oauth2/access_token?%s", API_URL, query.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
	}()

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		var e Error
		if err := decoder.Decode(&e); err != nil {
			return nil, err
		}

		return nil, errors.New(e.Message)
	}

	var t AccessToken
	if err := decoder.Decode(&t); err != nil {
		return nil, err
	}

	token = &oauth2.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       time.Unix(now+int64(t.ExpiresIn), 0),
	}

	return token, err
}

func (w WeChat) RefreshToken(token *oauth2.Token) (rt *oauth2.Token, _ bool, err error) {

	now := time.Now().Unix()

	if token.Expiry.Before(time.Now()) == false || token.RefreshToken == "" {
		return token, false, nil
	}

	query := make(url.Values)
	query.Set("grant_type", "refresh_token")
	query.Set("appid", w.clientID)
	query.Set("refresh_token", token.RefreshToken)

	var uri = fmt.Sprintf("%s/sns/oauth2/refresh_token?%s", API_URL, query.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return nil, false, err
	}
	defer func() {
		err = resp.Body.Close()
	}()

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		var e Error
		if err := decoder.Decode(&e); err != nil {
			return nil, false, err
		}

		return nil, false, errors.New(e.Message)
	}

	var t RefreshToken
	if err := decoder.Decode(&t); err != nil {
		return nil, false, err
	}

	rt = &oauth2.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       time.Unix(now+int64(t.ExpiresIn), 0),
	}

	return rt, true, err
}

func (w WeChat) GetUser(token *oauth2.Token) (*interfaces.User, error) {

	var err error

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	payload := struct {
		Username string `json:"nickname"`
		Email    string `json:"email"`
		ID       string `json:"openid"`
		UnionID  string `json:"unionid"`
	}{}

	user := new(interfaces.User)

	query := make(url.Values)
	query.Set("access_token", token.AccessToken)
	query.Set("openid", w.clientID)

	//if lang != "" {
	//	query.Set("lang", WCAuthClientSecret)
	//}

	var uri = fmt.Sprintf("%s/sns/userinfo?%s", API_URL, query.Encode())

	resUser, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resUser.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	username := strings.Replace(payload.Username, " ", "_", -1)

	user.Username = strings.ToLower(username)
	user.Email = payload.Email
	user.ServiceID = payload.ID

	return user, nil
}

func (w WeChat) GetVendorInfo() *interfaces.Vendor {
	return &interfaces.Vendor{
		Vendor: w.vendor,
		Host:   w.host,
	}
}

func (w WeChat) getOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     w.clientID,
		ClientSecret: w.clientSecretID,
		RedirectURL:  w.redirectURI,
		Endpoint: oauth2.Endpoint{
			TokenURL: TOKEN_URL,
		},
	}
}
