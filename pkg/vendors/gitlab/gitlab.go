//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package gitlab

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
	API_URL = "https://gitlab.com"
)

type GitLab struct {
	types.Vendor
}

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

func GetClient(token string) *GitLab {
	c := new(GitLab)
	c.Token = &oauth2.Token{AccessToken: token}
	c.Name = "gitlab"
	c.Host = "gitlab.com"
	return c
}

func (g *GitLab) VendorInfo() *types.Vendor {
	return &g.Vendor
}

func (g *GitLab) GetUser() (*types.User, error) {

	var err error

	payload := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		ID       int64  `json:"id"`
	}{}

	var uri = fmt.Sprintf("%s/api/v3/user?private_token=%s", API_URL, g.Token.AccessToken)

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

func (g *GitLab) ListRepositories(username string, org bool) (*types.VCSRepositories, error) {

	payload := []struct {
		Name          string  `json:"name"`
		Description   *string `json:"description"`
		Public        bool    `json:"public"`
		DefaultBranch string  `json:"default_branch"`
	}{}

	var uri = fmt.Sprintf("%s/api/v3/projects?private_token=%s", API_URL, g.Token.AccessToken)

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
		repository.Private = !repo.Public
		repository.DefaultBranch = repo.DefaultBranch
		if repo.Description != nil {
			repository.Description = *repo.Description
		}

		*repositories = append(*repositories, *repository)
	}

	return repositories, nil
}

func (g *GitLab) ListBranches(owner, repo string) (*types.VCSBranches, error) {

	owner = strings.ToLower(owner)
	repo = strings.ToLower(repo)

	payload := []struct {
		Name string `json:"name"`
	}{}

	var uri = fmt.Sprintf("%s/api/v3/projects/%s%%2F%s/repository/branches?private_token=%s", API_URL, owner, repo, g.Token.AccessToken)

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

func (g *GitLab) CreateHook(hookID, owner, repo, host string) (*string, error) {

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

	var uri = fmt.Sprintf("%s/api/v3/projects/%s/hooks?private_token=%s", API_URL, name, g.Token.AccessToken)

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

	id := strconv.FormatInt(int64(payload.ID), 10)

	return &id, nil
}

func (g *GitLab) RemoveHook(id, owner, repo string) error {

	var err error

	owner = strings.ToLower(owner)
	repo = strings.ToLower(repo)

	var uri = fmt.Sprintf("%s/api/v3/projects/%s%%2F%s/hooks/%s?private_token=%s", API_URL, owner, repo, id, g.Token.AccessToken)

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

func (g *GitLab) PushPayload(data []byte) (*types.VCSBranch, error) {

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

	var branch = new(types.VCSBranch)

	branch.Name = strings.Split(payload.Ref, "/")[2]
	branch.LastCommit = types.Commit{
		Username: commit.Committer.Username,
		Email:    commit.Committer.Email,
		Hash:     commit.Hash,
		Message:  commit.Message,
		Date:     commit.Date,
	}

	return branch, nil
}
