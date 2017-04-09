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

package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"github.com/lastbackend/lastbackend/pkg/vendors"
	"github.com/lastbackend/lastbackend/pkg/vendors/interfaces"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	ErrVendorAuth      = "vendor auth data not set"
	ErrVendorSupported = "vendor is not supported yet"
)

func Create(ctx context.Context, imageID uuid.UUID, source *types.ServiceSource) (*types.Build, error) {
	var (
		err            error
		clientID       string
		clientSecretID string
		redirectURI    string
		lctx           = c.Get()
		vendor         = strings.Split(source.Hub, ".")[0]
		commit         *types.GitSourceCommit
	)

	lctx.Log.Debug("Create build")

	var client interfaces.IVCS
	clientID, clientSecretID, redirectURI = config.Get().GetVendorConfig(vendor)

	if clientID == "" || clientSecretID == "" {
		commit, err = getLastCommit(vendor, source.Owner, source.Repo, source.Branch)
		if err != nil {
			return nil, err
		}
	} else {

		// Get client for github/bitbucket/gitlab (or anything if implement adapter.OAuth interface) by vendor in user or organization mode
		switch vendor {
		case "github":
			client = vendors.GetGitHub(clientID, clientSecretID, redirectURI)
		case "bitbucket":
			client = vendors.GetBitBucket(clientID, clientSecretID, redirectURI)
		case "gitlab":
			client = vendors.GetGitLab(clientID, clientSecretID, redirectURI)
		default:
			lctx.Log.Error(ErrVendorSupported)
			return nil, errors.New(ErrVendorSupported)
		}

		vendorInfo := client.GetVendorInfo()

		oauth, err := lctx.Storage.Vendor().Get(ctx, vendorInfo.Vendor)
		if err != nil {
			lctx.Log.Error(err)
			return nil, err
		}

		token, modify, err := client.RefreshToken(oauth.Token)
		if err != nil {
			lctx.Log.Error(err)
			return nil, err
		}

		u, err := client.GetUser(token)
		if err != nil {
			lctx.Log.Error(err)
			return nil, err
		}

		if modify {

			oauth.Host = vendorInfo.Host
			oauth.Vendor = vendorInfo.Vendor
			oauth.ServiceID = u.ServiceID
			oauth.Token = token
			oauth.Username = u.Username

			if err = lctx.Storage.Vendor().Update(ctx, oauth); err != nil {
				lctx.Log.Error(err)
				return nil, err
			}
		}

		info, err := client.GetLastCommitOfBranch(token, source.Owner, source.Repo, source.Branch)
		if err != nil {
			lctx.Log.Error(err)
			return nil, err
		}

		commit.Date = info.Date
		commit.Email = info.Email
		commit.Commit = info.Hash
		commit.Message = info.Message
		commit.Author = info.Username
		commit.Committer = generator.GenerateGravatar(info.Email)

	}

	bsource := &types.BuildSource{
		Hub:    source.Hub,
		Owner:  source.Owner,
		Repo:   source.Repo,
		Tag:    source.Branch,
		Commit: *commit,
	}

	bld, err := lctx.Storage.Build().Insert(ctx, imageID, bsource)
	if err != nil {
		return nil, err
	}

	return bld, nil
}

func getLastCommit(vendor, owner, repo, branch string) (*types.GitSourceCommit, error) {

	var commit = new(types.GitSourceCommit)

	switch vendor {
	case "github":

		payload := struct {
			Error  string `json:"message"`
			Commit struct {
				Hash   string `json:"sha"`
				Commit struct {
					Message string `json:"message"`
					Author  struct {
						Date  time.Time `json:"date"`
						Email string    `json:"email"`
					} `json:"author"`
				} `json:"commit"`
				Author struct {
					Login string `json:"login"`
				} `json:"author"`
			} `json:"commit"`
		}{}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", owner, repo, branch)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}

		if payload.Error != "" {
			return nil, errors.New(payload.Error)
		}

		commit.Date = payload.Commit.Commit.Author.Date
		commit.Email = payload.Commit.Commit.Author.Email
		commit.Commit = payload.Commit.Hash
		commit.Message = payload.Commit.Commit.Message
		commit.Author = payload.Commit.Author.Login
		commit.Committer = generator.GenerateGravatar(commit.Email)

	case "bitbucket":

		payload := struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
			Target struct {
				Hash   string `json:"hash"`
				Author struct {
					Raw  string `json:"raw"`
					User struct {
						Username string `json:"username"`
					} `json:"user"`
				} `json:"author"`
				Date    time.Time `json:"date"`
				Message string    `json:"message"`
			} `json:"target"`
		}{}

		url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/refs/branches/%s", owner, repo, branch)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}

		if payload.Error.Message != "" {
			return nil, errors.New(payload.Error.Message)
		}

		// Another regular expression: <(\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}\b)>$
		r, _ := regexp.Compile("<(.+)>$")

		commit.Date = payload.Target.Date
		commit.Email = r.FindStringSubmatch(payload.Target.Author.Raw)[1]
		commit.Commit = payload.Target.Hash
		commit.Message = payload.Target.Message
		commit.Author = payload.Target.Author.User.Username
		commit.Committer = generator.GenerateGravatar(commit.Email)
	case "gitlab":

		payload := struct {
			Error  string `json:"error"`
			Commit struct {
				Hash           string    `json:"id"`
				Message        string    `json:"message"`
				CommitterEmail string    `json:"committer_email"`
				CommitterName  string    `json:"committer_name"`
				CommitterDate  time.Time `json:"committed_date"`
			} `json:"commit"`
		}{}

		url := fmt.Sprintf("https://gitlab.com/repos/%s/%s/branches/%s", owner, repo, branch)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, err
		}

		if payload.Error != "" {
			return nil, errors.New(payload.Error)
		}

		commit.Date = payload.Commit.CommitterDate
		commit.Email = payload.Commit.CommitterEmail
		commit.Message = payload.Commit.Message
		commit.Commit = payload.Commit.Hash
		commit.Author = payload.Commit.CommitterName
		commit.Committer = generator.GenerateGravatar(commit.Email)

	default:
		return nil, errors.New(ErrVendorSupported)
	}

	return commit, nil
}
