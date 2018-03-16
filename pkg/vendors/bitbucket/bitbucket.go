//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/vendors/types"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	API_URL = "https://api.bitbucket.org"
)

type BitBucket struct {
	types.Vendor
}

func GetClient(token string) *BitBucket {
	c := new(BitBucket)
	c.Token = &oauth2.Token{AccessToken: token}
	c.Name = "bitbucket"
	c.Host = "bitbucket.org"
	return c
}

func (b *BitBucket) VendorInfo() *types.Vendor {
	return &b.Vendor
}

func (b *BitBucket) GetUser() (*types.User, error) {

	var err error

	payload := struct {
		Username string `json:"username"`
		ID       string `json:"uuid"`
	}{}

	var uri = fmt.Sprintf("https://%s@api.bitbucket.org/2.0/user", b.Token.AccessToken)

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
	user.ServiceID = payload.ID

	return user, nil
}

func (b *BitBucket) ListRepositories(username string, org bool) (*types.VCSRepositories, error) {

	var err error

	username = strings.ToLower(username)

	payload := struct {
		Repos []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Private     bool   `json:"is_private"`
		} `json:"values"`
	}{}

	var uri = fmt.Sprintf("https://%s@api.bitbucket.org/2.0/repositories?role=owner", b.Token.AccessToken)

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

	for _, repo := range payload.Repos {
		repository := new(types.VCSRepository)

		repository.Name = repo.Name
		repository.Description = repo.Description
		repository.Private = repo.Private
		*repositories = append(*repositories, *repository)
	}

	return repositories, nil
}

func (b *BitBucket) ListBranches(owner, repo string) (*types.VCSBranches, error) {

	var branches = new(types.VCSBranches)

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := struct {
		Branches []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"values"`
	}{}

	var uri = fmt.Sprintf("https://%s@api.bitbucket.org/2.0/repositories/%s/%s/refs/branches?pagelen=100", b.Token.AccessToken, owner, repo)

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	for _, br := range payload.Branches {
		if br.Type == "branch" {
			branch := new(types.VCSBranch)

			branch.Name = br.Name
			*branches = append(*branches, *branch)
		}
	}

	return branches, nil
}

func (b *BitBucket) CreateHook(hookID, owner, repo, host string) (*string, error) {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := struct {
		ID    string `json:"uuid"`
		Error struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}{}

	body := struct {
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Active      bool     `json:"active"`
		Events      []string `json:"events"`
	}{"web", fmt.Sprintf(`%s/hook/bitbucket/process/%s`, host, hookID), true, []string{"repo:push", "pullrequest:approved", "pullrequest:created"}}

	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	var uri = fmt.Sprintf("https://%s@api.bitbucket.org/2.0/repositories/%s/%s/hooks", b.Token.AccessToken, owner, repo)

	res, err := http.Post(uri, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	if payload.Error.Message != "" {
		return nil, errors.New(payload.Error.Message)
	}

	id := payload.ID

	return &id, nil
}

func (b *BitBucket) RemoveHook(hookID, owner, repo string) error {

	var err error

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	var uri = fmt.Sprintf("https://%s@api.bitbucket.org/2.0/repositories/%s/%s/hooks/%s", b.Token.AccessToken, owner, repo, hookID)

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

func (b *BitBucket) PushPayload(data []byte) (*types.VCSBranch, error) {

	var err error

	payload := struct {
		Push struct {
			Changes []struct {
				New struct {
					Name string `json:"name"`
				} `json:"new"`
				Commits []struct {
					Hash    string    `json:"hash"`
					Message string    `json:"message"`
					Date    time.Time `json:"date"`
					Author  struct {
						User struct {
							Username string `json:"username"`
						} `json:"user"`
						Raw string `json:"raw"`
					} `json:"author"`
				} `json:"commits"`
			} `json:"changes"`
		} `json:"push"`
	}{}

	if err = json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	r, _ := regexp.Compile("<(.+)>$")

	branch := new(types.VCSBranch)
	branch.Name = payload.Push.Changes[0].New.Name
	branch.LastCommit = types.Commit{
		Hash:     payload.Push.Changes[0].Commits[0].Hash,
		Date:     payload.Push.Changes[0].Commits[0].Date,
		Username: payload.Push.Changes[0].Commits[0].Author.User.Username,
		Message:  payload.Push.Changes[0].Commits[0].Message,
		Email:    r.FindStringSubmatch(payload.Push.Changes[0].Commits[0].Author.Raw)[1],
	}

	return branch, nil
}
