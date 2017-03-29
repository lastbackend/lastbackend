package gitlab

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/vendors/interfaces"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// const

const (
	API_URL   = "https://gitlab.com"
	TOKEN_URL = "https://gitlab.com/oauth/token"
)

// Model

type GitLab struct {
	proto interfaces.OAuth2
}

// Types

type CommitResponse struct {
	ID        int64     `json:"project_id"`
	Hash      string    `json:"checkout_sha"`
	Message   string    `json:"message"`
	Date      time.Time `json:"timestamp"`
	Committer struct {
		Username string `json:"name"`
		Email    string `json:"email"`
	} `json:"author"`
}

// IVendor

func GetClient(clientID, clientSecretID, redirectURI string) *GitLab {
	return &GitLab{proto: interfaces.OAuth2{ClientID: clientID, ClientSecret: clientSecretID, RedirectUri: redirectURI}}
}

func (GitLab) GetVendorInfo() *interfaces.Vendor {
	return &interfaces.Vendor{Vendor: "gitlab", Host: "gitlab.com"}
}

// IOAuth2 func

func (g GitLab) GetToken(code string) (*oauth2.Token, error) {

	token, err := g.getOAuth2Config().Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (g GitLab) RefreshToken(token *oauth2.Token) (*oauth2.Token, bool, error) {

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

func (g GitLab) getOAuth2Config() *oauth2.Config {
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

func (GitLab) GetUser(token *oauth2.Token) (*interfaces.User, error) {

	var err error

	payload := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		ID       int64  `json:"id"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/api/v3/user", API_URL)

	resUser, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	x := []byte{}
	_, err = resUser.Body.Read(x)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resUser.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var user = new(interfaces.User)
	user.Username = payload.Username
	user.Email = payload.Email
	user.ServiceID = strconv.FormatInt(payload.ID, 10)

	return user, nil
}

func (GitLab) ListRepositories(token *oauth2.Token, username string, org bool) (*interfaces.VCSRepositories, error) {

	payload := []struct {
		Name          string  `json:"name"`
		Description   *string `json:"description"`
		Public        bool    `json:"public"`
		DefaultBranch string  `json:"default_branch"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/api/v3/projects", API_URL)

	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	var repositories = new(interfaces.VCSRepositories)

	for _, repo := range payload {
		repository := new(interfaces.VCSRepository)

		repository.Name = repo.Name
		repository.Private = !repo.Public
		repository.DefaultBranch = repo.DefaultBranch
		if repo.Description != nil {
			repository.Description = *repo.Description
		}

		*repositories = append(*repositories, *repository)
	}

	return repositories, nil
}

func (GitLab) ListBranches(token *oauth2.Token, owner, repo string) (*interfaces.VCSBranches, error) {

	owner = strings.ToLower(owner)
	repo = strings.ToLower(repo)

	payload := []struct {
		Name string `json:"name"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/api/v3/projects/%s%%2F%s/repository/branches", API_URL, owner, repo)

	res, err := client.Get(uri)
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

func (GitLab) GetLastCommitOfBranch(token *oauth2.Token, owner, repo, branch string) (*interfaces.Commit, error) {

	owner = strings.ToLower(owner)
	repo = strings.ToLower(repo)
	branch = strings.ToLower(branch)

	branchResponse := struct {
		Name   string `json:"name"`
		Commit struct {
			Hash           string    `json:"id"`
			Message        string    `json:"message"`
			CommitterEmail string    `json:"committer_email"`
			CommitterName  string    `json:"committer_name"`
			CommitterDate  time.Time `json:"committed_date"`
		} `json:"commit"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/api/v3/projects/%s%%2F%s/repository/branches/%s", API_URL, owner, repo, branch)

	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&branchResponse); err != nil {
		return nil, err
	}

	var commit = new(interfaces.Commit)

	commit.Hash = branchResponse.Commit.Hash
	commit.Date = branchResponse.Commit.CommitterDate
	commit.Message = branchResponse.Commit.Message
	commit.Username = branchResponse.Commit.CommitterName // TODO: Get username
	commit.Email = branchResponse.Commit.CommitterEmail

	return commit, nil
}

func (GitLab) GetReadme(token *oauth2.Token, owner string, repo string) (string, error) {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf(`%s/%s/%s/raw/master/README.md`, API_URL, owner, repo)

	res, err := client.Get(uri)
	if err != nil {
		return "", nil
	}

	var content string

	if res.StatusCode == 200 {
		buf, _ := ioutil.ReadAll(res.Body)
		content = string(buf)
	}

	return string(content), nil
}

func (GitLab) PushPayload(data []byte) (*interfaces.VCSBranch, error) {

	var err error

	payload := struct {
		Ref     string           `json:"ref"`
		Hash    string           `json:"checkout_sha"`
		Commits []CommitResponse `json:"commits"`
	}{}

	if err = json.Unmarshal(data, &payload); err != nil {
		return nil, nil
	}

	commit := CommitResponse{}

	for index := range payload.Commits {
		commit = payload.Commits[index]

		if commit.Hash == payload.Hash {
			break
		}
	}

	var branch = new(interfaces.VCSBranch)

	branch.Name = strings.Split(payload.Ref, "/")[2]
	branch.LastCommit = interfaces.Commit{
		Username: commit.Committer.Username,
		Email:    commit.Committer.Email,
		Hash:     commit.Hash,
		Message:  commit.Message,
		Date:     commit.Date,
	}

	return branch, nil
}

func (GitLab) CreateHook(token *oauth2.Token, hookID, owner, repo, host string) (*string, error) {

	owner = strings.ToLower(owner)
	repo = strings.ToLower(repo)
	name := owner + "%2F" + repo

	payload := struct {
		ID    int64  `json:"id"`
		Error string `json:"error,omitempty"`
	}{}

	body := struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}{name, fmt.Sprintf("%s/hook/gitlab/process/%s", host, hookID)}

	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, nil
	}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/api/v3/projects/%s/hooks", API_URL, name)

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

	id := strconv.FormatInt(int64(payload.ID), 10)

	return &id, nil
}

func (GitLab) RemoveHook(token *oauth2.Token, id, owner, repo string) error {

	var err error

	owner = strings.ToLower(owner)
	repo = strings.ToLower(repo)

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/api/v3/projects/%s%%2F%s/hooks/%s", API_URL, owner, repo, id)

	req, err := http.NewRequest("DELETE", uri, nil)
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
