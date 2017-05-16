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
	c.Token = &oauth2.Token{AccessToken: token}
	c.Name = "github"
	c.Host = "github.com"
	return c
}

func (g *GitHub) VendorInfo() *types.Vendor {
	return &g.Vendor
}

func (g *GitHub) GetUser() (*types.User, error) {

	var err error

	payload := struct {
		Username string `json:"login"`
		ID       int64  `json:"id"`
	}{}

	var uri = fmt.Sprintf("%s/user?access_token=%s", API_URL, g.Token.AccessToken)

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var user = new(types.User)
	user.Username = payload.Username
	user.ServiceID = strconv.FormatInt(payload.ID, 10)

	return user, nil
}

func (g *GitHub) ListRepositories(username string, org bool) (*types.VCSRepositories, error) {

	var err error

	username = strings.ToLower(username)

	payload := []struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		Private       bool   `json:"private"`
		DefaultBranch string `json:"default_branch"`
	}{}

	var uri = fmt.Sprintf("%s/user/repos?access_token=%s&per_page=99&type=owner", API_URL, g.Token.AccessToken)

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
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

	var uri = fmt.Sprintf("%s/repos/%s/%s/branches?access_token=%s", API_URL, owner, repo, g.Token.AccessToken)

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var branches = new(types.VCSBranches)

	for _, br := range payload {
		branch := new(types.VCSBranch)

		branch.Name = br.Name
		*branches = append(*branches, *branch)
	}

	return branches, nil
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

	var uri = fmt.Sprintf("%s/repos/%s/%s/hooks?access_token=%s", API_URL, owner, repo, g.Token.AccessToken)

	res, err := http.Post(uri, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
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

	var uri = fmt.Sprintf(`%s/repos/%s/%s/hooks/%s?access_token=%s`, API_URL, owner, repo, id, g.Token.AccessToken)

	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
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
