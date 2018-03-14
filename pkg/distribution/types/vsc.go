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

package types

import "time"

const GithubHost = "github.com"
const BitbucketHost = "bitbucket.org"
const GitlabHost = "gitlab.com"
const DockerHost = "index.docker.io"

const GithubType = "github"
const BitbucketType = "bitbucket"
const GitlabType = "gitlab"

type VCSRepo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
}

type VCSRepoList []VCSRepo

type VCSBranch struct {
	Name       string    `json:"name"`
	LastCommit VCSCommit `json:"commit,omitempty"`
}

type VCSBranchList []VCSBranch

type VCSCommit struct {
	Username string    `json:"username"`
	Hash     string    `json:"hash"`
	Message  string    `json:"message"`
	Email    string    `json:"email"`
	Date     time.Time `json:"date"`
}
