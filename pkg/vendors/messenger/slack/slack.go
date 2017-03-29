package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/vendors/interfaces"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// const

const (
	API_URL   = "https://slack.com/api"
	TOKEN_URL = "https://slack.com/api/oauth.access"
)

// Model

type Slack struct {
	clientID       string
	clientSecretID string
	redirectURI    string
	vendor         string
	host           string
	access         string
	mode           string
	locker         sync.Mutex
}

func GetClient(clientID, clientSecretID, redirectURI string) *Slack {
	return &Slack{
		clientID:       clientID,
		clientSecretID: clientSecretID,
		redirectURI:    redirectURI,
	}
}

// IVendor

func (Slack) GetVendorInfo() *interfaces.Vendor {
	return &interfaces.Vendor{
		Vendor: "slack",
		Host:   "slack.com",
	}
}

func (Slack) httpGet(url string, i interface{}) error {

	var e error

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {

		if err := decoder.Decode(&e); err != nil {
			return err
		}

		return errors.New(e.Error())
	}

	if err := decoder.Decode(&i); err != nil {
		return err
	}

	return nil
}

// IOAuth2

func (s Slack) GetToken(code string) (*oauth2.Token, error) {

	token, err := s.getOAuth2Config().Exchange(oauth2.NoContext, code)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (s Slack) RefreshToken(token *oauth2.Token) (*oauth2.Token, bool, error) {

	if !token.Expiry.Before(time.Now()) || token.RefreshToken == "" {
		return token, false, nil
	}

	restoredToken := &oauth2.Token{
		RefreshToken: token.RefreshToken,
	}

	c := s.getOAuth2Config().Client(oauth2.NoContext, restoredToken)

	newToken, err := c.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		return newToken, false, err
	}

	return newToken, true, nil
}

func (s Slack) getOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecretID,
		Endpoint: oauth2.Endpoint{
			TokenURL: TOKEN_URL,
		},
	}
}

// INotify

func (Slack) ListChannels(token *oauth2.Token) (*interfaces.NotifyChannels, error) {

	var err error

	payload := struct {
		Ok       bool `json:"ok"`
		Channels []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			IsChannel bool   `json:"is_channel"`
		} `json:"channels"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/channels.list?token=%s", API_URL, token.AccessToken)

	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var channels = new(interfaces.NotifyChannels)

	for _, ch := range payload.Channels {
		channel := interfaces.NotifyChannel{}

		if ch.IsChannel {
			channel.ID = ch.ID
			channel.Name = ch.Name
			channel.Type = "channel"

			*channels = append(*channels, channel)
		}
	}

	return channels, nil
}

func (Slack) ListGroups(token *oauth2.Token) (*interfaces.NotifyGroups, error) {

	var err error

	payload := struct {
		Ok     bool `json:"ok"`
		Groups []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			IsGroup bool   `json:"is_group"`
		} `json:"groups"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/groups.list?token=%s", API_URL, token.AccessToken)

	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var groups = new(interfaces.NotifyGroups)

	for _, gr := range payload.Groups {
		group := interfaces.NotifyGroup{}

		if gr.IsGroup {
			group.ID = gr.ID
			group.Name = gr.Name
			group.Type = "group"

			*groups = append(*groups, group)
		}
	}

	return groups, nil
}

func (s Slack) GetUser(token *oauth2.Token) (*interfaces.User, error) {

	payload := struct {
		ID       string `json:"user_id"`
		Username string `json:"user"`
	}{}

	query := make(url.Values)
	query.Set("token", token.AccessToken)

	var uri = fmt.Sprintf("%s/auth.test?%s", API_URL, query.Encode())

	err := s.httpGet(uri, &payload)
	if err != nil {
		return nil, err
	}

	userResponse := struct {
		Profile struct {
			Email string `json:"email"`
		} `json:"profile"`
	}{}

	query = make(url.Values)
	query.Set("token", token.AccessToken)
	query.Set("user", payload.ID)

	uri = fmt.Sprintf("%s/users.profile.get?%s", API_URL, query.Encode())
	if err := s.httpGet(uri, &userResponse); err != nil {
		return nil, err
	}

	var user = new(interfaces.User)

	user.Username = payload.Username
	user.Email = userResponse.Profile.Email
	user.ServiceID = payload.ID

	return user, nil
}
