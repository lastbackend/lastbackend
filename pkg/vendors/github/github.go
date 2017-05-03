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

package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/vendors/types"
	"github.com/lastbackend/lastbackend/pkg/vendors/utils"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	API_URL = "https://api.github.com"
)

type GitHub struct {
	types.Vendor
}

func GetClient(token string) *GitHub {
	c := new(GitHub)
	c.Token = &oauth2.Token{AccessToken: token, TokenType: "Bearer"}
	c.Name = "github"
	c.Host = "github.com"
	return c
}

func (g *GitHub) VendorInfo() *types.Vendor {
	return &g.Vendor
}

func (g *GitHub) GetUser() (*types.User, error) {

	var err error

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(g.Token))

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

	var user = new(types.User)
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

func (g *GitHub) ListRepositories(username string, org bool) (*types.VCSRepositories, error) {

	var res *http.Response
	var err error

	username = strings.ToLower(username)

	payload := []struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		Private       bool   `json:"private"`
		DefaultBranch string `json:"default_branch"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(g.Token))

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

	var repositories = new(types.VCSRepositories)

	for _, repo := range payload {
		repository := new(types.VCSRepository)

		repository.Name = repo.Name
		repository.Description = repo.Description
		repository.Private = repo.Private
		repository.DefaultBranch = repo.DefaultBranch

		*repositories = append(*repositories, *repository)
	}

	return repositories, nil
}

func (g *GitHub) ListBranches(owner, repo string) (*types.VCSBranches, error) {

	var err error

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := []struct {
		Name string `json:"name"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(g.Token))

	res, err := client.Get(fmt.Sprintf("%s/repos/%s/%s/branches", API_URL, owner, repo))

	if err != nil {
		return nil, nil
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, nil
	}

	var branches = new(types.VCSBranches)

	for _, br := range payload {
		branch := new(types.VCSBranch)

		branch.Name = br.Name
		*branches = append(*branches, *branch)
	}

	return branches, nil
}

func (g *GitHub) GetLastCommitOfBranch(owner, repo, branch string) (*types.Commit, error) {

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

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(g.Token))

	var uri = fmt.Sprintf("%s/repos/%s/%s/branches/%s", API_URL, owner, repo, branch)

	res, err := client.Get(uri)

	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&branchResponse); err != nil {
		return nil, err
	}

	var commit = new(types.Commit)

	commit.Hash = branchResponse.Commit.Hash
	commit.Date = branchResponse.Commit.Commit.Committer.Date
	commit.Message = branchResponse.Commit.Commit.Message
	commit.Username = branchResponse.Commit.Committer.Login
	commit.Email = branchResponse.Commit.Commit.Committer.Email

	return commit, nil
}

func (g *GitHub) GetReadme(owner string, repo string) (string, error) {

	var err error

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := struct {
		Content string `json:"content"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(g.Token))

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

func (g *GitHub) PushPayload(data []byte) (*types.VCSBranch, error) {

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

	var branch = new(types.VCSBranch)

	branch.Name = strings.Split(payload.Ref, "/")[2]
	branch.LastCommit = types.Commit{
		Username: payload.Commit.Committer.Username,
		Email:    payload.Commit.Committer.Email,
		Hash:     payload.Commit.ID,
		Message:  payload.Commit.Message,
		Date:     payload.Commit.Date,
	}

	return branch, nil
}

func (g *GitHub) CreateHook(id, owner, repo, host string) (*string, error) {

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

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(g.Token))

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

func (g *GitHub) RemoveHook(id, owner, repo string) error {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(g.Token))

	var uri = fmt.Sprintf(`%s/repos/%s/%s/hooks/%s`, API_URL, owner, repo, id)

	req, err := http.NewRequest("DELETE", uri, nil)
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
