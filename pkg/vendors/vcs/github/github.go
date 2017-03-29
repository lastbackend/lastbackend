package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/vendors/interfaces"
	"github.com/lastbackend/vendors/utils"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// const

const (
	API_URL   = "https://api.github.com"
	TOKEN_URL = "https://github.com/login/oauth/access_token"
)

// Model

type GitHub struct {
	proto interfaces.OAuth2
}

// IVendor

func GetClient(clientID, clientSecretID, redirectURI string) *GitHub {
	return &GitHub{proto: interfaces.OAuth2{ClientID: clientID, ClientSecret: clientSecretID, RedirectUri: redirectURI}}
}

func (GitHub) GetVendorInfo() *interfaces.Vendor {
	return &interfaces.Vendor{Vendor: "github", Host: "github.com"}
}

// IOAuth2 func

func (g GitHub) GetToken(code string) (*oauth2.Token, error) {

	token, err := g.getOAuth2Config().Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (g GitHub) RefreshToken(token *oauth2.Token) (*oauth2.Token, bool, error) {

	var err error

	if !token.Expiry.Before(time.Now()) || token.RefreshToken == "" {
		return token, false, nil
	}

	restoredToken := &oauth2.Token{
		RefreshToken: token.RefreshToken,
	}

	c := g.getOAuth2Config().Client(oauth2.NoContext, restoredToken)

	newToken, err := c.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		return nil, false, err
	}

	return newToken, true, nil
}

// IOAuth2 - Private functions

func (g GitHub) getOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     g.proto.ClientID,
		ClientSecret: g.proto.ClientSecret,
		RedirectURL:  g.proto.RedirectUri,
		Endpoint: oauth2.Endpoint{
			TokenURL: TOKEN_URL,
		},
	}
}

// IVCS func

func (GitHub) GetUser(token *oauth2.Token) (*interfaces.User, error) {

	var err error

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	payload := struct {
		Username string `json:"login"`
		ID       int64  `json:"id"`
	}{}

	var uri = fmt.Sprintf("%s/user", API_URL)

	resUser, err := client.Get(uri)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resUser.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var user = new(interfaces.User)
	user.Username = payload.Username
	user.ServiceID = strconv.FormatInt(payload.ID, 10)

	emailsResponse := []struct {
		Email     string `json:"email"`
		Confirmed bool   `json:"verified"`
		Primary   bool   `json:"primary"`
	}{}

	uri = fmt.Sprintf("%s/user/emails", API_URL)

	resEmails, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resEmails.Body).Decode(&emailsResponse)
	if err != nil {
		return nil, err
	}

	for _, email := range emailsResponse {
		if email.Confirmed == true && email.Primary == true {
			user.Email = email.Email
			break
		}
	}

	return user, nil
}

func (GitHub) ListRepositories(token *oauth2.Token, username string, org bool) (*interfaces.VCSRepositories, error) {

	var res *http.Response
	var err error

	username = strings.ToLower(username)

	payload := []struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		Private       bool   `json:"private"`
		DefaultBranch string `json:"default_branch"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/user/repos?per_page=99&type=owner", API_URL)
	if org {
		uri = fmt.Sprintf("%s/orgs/%s/repos?per_page=99", API_URL, username)
	}

	res, err = client.Get(uri)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	var repositories = new(interfaces.VCSRepositories)

	for _, repo := range payload {
		repository := new(interfaces.VCSRepository)

		repository.Name = repo.Name
		repository.Description = repo.Description
		repository.Private = repo.Private
		repository.DefaultBranch = repo.DefaultBranch

		*repositories = append(*repositories, *repository)
	}

	return repositories, nil
}

func (GitHub) ListBranches(token *oauth2.Token, owner, repo string) (*interfaces.VCSBranches, error) {

	var err error

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := []struct {
		Name string `json:"name"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	res, err := client.Get(fmt.Sprintf("%s/repos/%s/%s/branches", API_URL, owner, repo))

	if err != nil {
		return nil, nil
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, nil
	}

	var branches = new(interfaces.VCSBranches)

	for _, br := range payload {
		branch := new(interfaces.VCSBranch)

		branch.Name = br.Name
		*branches = append(*branches, *branch)
	}

	return branches, nil
}

func (GitHub) GetLastCommitOfBranch(token *oauth2.Token, owner, repo, branch string) (*interfaces.Commit, error) {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)
	branch = strings.ToLower(branch)

	branchResponse := struct {
		Commit struct {
			Hash   string `json:"sha"`
			Commit struct {
				Message   string `json:"message"`
				Committer struct {
					Date  time.Time `json:"date"`
					Email string    `json:"email"`
				} `json:"committer"`
			} `json:"commit"`
			Committer struct {
				Login string `json:"login"`
			} `json:"committer"`
		} `json:"commit"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/repos/%s/%s/branches/%s", API_URL, owner, repo, branch)

	res, err := client.Get(uri)

	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&branchResponse); err != nil {
		return nil, err
	}

	var commit = new(interfaces.Commit)

	commit.Hash = branchResponse.Commit.Hash
	commit.Date = branchResponse.Commit.Commit.Committer.Date
	commit.Message = branchResponse.Commit.Commit.Message
	commit.Username = branchResponse.Commit.Committer.Login
	commit.Email = branchResponse.Commit.Commit.Committer.Email

	return commit, nil
}

func (GitHub) GetReadme(token *oauth2.Token, owner string, repo string) (string, error) {

	var err error

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := struct {
		Content string `json:"content"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf(`%s/repos/%s/%s/readme`, API_URL, owner, repo)

	res, err := client.Get(uri)
	if err != nil {
		return "", nil
	}

	res.Header.Add("Accept", "application/vnd.github.VERSION.raw")

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", nil
	}

	if payload.Content != "" {
		payload.Content = utils.DecodeBase64(payload.Content)
	}

	return payload.Content, nil
}

func (GitHub) PushPayload(data []byte) (*interfaces.VCSBranch, error) {

	var err error

	payload := struct {
		Ref    string `json:"ref"`
		Commit struct {
			ID        string    `json:"id"`
			Message   string    `json:"message"`
			Date      time.Time `json:"timestamp"`
			Committer struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"committer"`
		} `json:"head_commit"`
	}{}

	if err = json.Unmarshal(data, &payload); err != nil {
		return nil, nil
	}

	var branch = new(interfaces.VCSBranch)

	branch.Name = strings.Split(payload.Ref, "/")[2]
	branch.LastCommit = interfaces.Commit{
		Username: payload.Commit.Committer.Username,
		Email:    payload.Commit.Committer.Email,
		Hash:     payload.Commit.ID,
		Message:  payload.Commit.Message,
		Date:     payload.Commit.Date,
	}

	return branch, nil
}

func (GitHub) CreateHook(token *oauth2.Token, id, owner, repo, host string) (*string, error) {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := struct {
		ID    int64  `json:"id"`
		Error string `json:"message,omitempty"`
	}{}

	body := struct {
		Name   string                 `json:"name"`
		Active bool                   `json:"active"`
		Events []string               `json:"events"`
		Config map[string]interface{} `json:"config"`
	}{"web", true, []string{"push", "pull_request"}, map[string]interface{}{"url": host + "/hook/github/process/" + id, "content_type": "json"}}

	var buf io.ReadWriter
	buf = new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, nil
	}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/repos/%s/%s/hooks", API_URL, owner, repo)

	res, err := client.Post(uri, "application/json", buf)
	if err != nil {
		return nil, nil
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if payload.Error != "" {
		return nil, errors.New(payload.Error)
	}

	id = strconv.FormatInt(int64(payload.ID), 10)

	return &id, nil
}

func (GitHub) RemoveHook(token *oauth2.Token, id, owner, repo string) error {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf(`%s/repos/%s/%s/hooks/%s`, API_URL, owner, repo, id)

	req, err := http.NewRequest("DELETE", uri, nil)
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
