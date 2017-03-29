package interfaces

import (
	"golang.org/x/oauth2"
	"time"
)

type User struct {
	Username  string
	ServiceID string
	Email     string
}

type VCSRepository struct {
	Name          string
	Description   string
	Private       bool
	DefaultBranch string
	Permissions   struct {
		Admin bool
	}
}

type VCSRepositories []VCSRepository

type VCSBranch struct {
	Name       string
	LastCommit Commit
}

type VCSBranches []VCSBranch

type Commit struct {
	Username string
	Hash     string
	Message  string
	Date     time.Time
	Email    string
}

type Commits []Commit

// interface

type IVCS interface {
	IVendor
	IOAuth2

	GetUser(token *oauth2.Token) (*User, error)
	ListRepositories(token *oauth2.Token, username string, org bool) (*VCSRepositories, error)
	ListBranches(token *oauth2.Token, owner, repo string) (*VCSBranches, error)
	GetLastCommitOfBranch(token *oauth2.Token, owner, repo, branch string) (*Commit, error)
	GetReadme(token *oauth2.Token, owner string, repo string) (string, error)
	PushPayload(data []byte) (*VCSBranch, error)
	CreateHook(token *oauth2.Token, id, owner, repo, host string) (*string, error)
	RemoveHook(token *oauth2.Token, id, owner, repo string) error
}
