package model

import "time"

type BuildList []Build

type Build struct {
	// Build number, incremented automatically
	ID string `json:"id" gorethink:"id,omitempty"`
	// Build number, incremented automatically
	User string `json:"user" gorethink:"id,omitempty"`
	// Build executing status
	Status BuildStatus `json:"status" gorethink:"status,omitempty"`
	// Build sources used for build
	Source BuildSource `json:"source" gorethink:"source,omitempty"`
	// Build image output information
	Image BuildImage `json:"image" gorethink:"image,omitempty"`
	// Build created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Build updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

type BuildStatus struct {
	// Build current step
	Step BuildStep `json:"step" gorethink:"step,omitempty"`
	// Is build cancelled
	Cancelled bool `json:"cancelled" gorethink:"cancelled,omitempty"`
	// Build executing message
	Message string `json:"message" gorethink:"message,omitempty"`
	// Build error information
	Error string `json:"error" gorethink:"error,omitempty"`
	// Build status updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

type BuildStep string

const (
	//BuildStepCreate - The first step after build creating
	BuildStepCreate = "create"
	//BuildStepFetch - Fetch sources step
	BuildStepFetch = "fetch"
	//BuildStepBuild - Build executing step
	BuildStepBuild = "build"
	//BuildStepUpload - Upload docker image step
	BuildStepUpload = "upload"
)

type BuildSource struct {
	// Build sources hub
	Hub string `json:"hub" gorethink:"hub,omitempty"`
	// Build sources owner
	Owner string `json:"owner" gorethink:"owner,omitempty"`
	// Build sources repo
	Repo string `json:"repo" gorethink:"repo,omitempty"`
	// Build source tag (branch, tag)
	Tag string `json:"tag" gorethink:"tag,omitempty"`
	// Build commit information
	Commit GitSourceCommit `json:"commit" gorethink:"commit,omitempty"`
	// Build sources auth reference
}

type BuildImage struct {
	// Build image repo name
	Repo string `json:"repo" gorethink:"repo,omitempty"`
	// Build image tag name
	Tag string `json:"tag" gorethink:"tag,omitempty"`
	// Build image registry reference
	Registry string `json:"registry" gorethink:"registry,omitempty"`
}

type GitSourceCommit struct {
	// Git commit information hash
	Commit string `json:"commit" gorethink:"id,omitempty"`
	// Git committer gravatar
	Committer string `json:"committer" gorethink:"id,omitempty"`
	// Git committer email
	Author string `json:"author" gorethink:"id,omitempty"`
	// Git commit message
	Message string `json:"message" gorethink:"id,omitempty"`
}
