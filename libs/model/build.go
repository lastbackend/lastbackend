package model

import "time"

type BuildList []Build

type Build struct {
	// Build number, incremented automatically
	ID string `json:"id" gorethink:"id"`
	// Build number, incremented automatically
	User string `json:"user" gorethink:"id"`
	// Build executing status
	Status BuildStatus `json:"status" gorethink:"status"`
	// Build sources used for build
	Source BuildSource `json:"source" gorethink:"source"`
	// Build image output information
	Image BuildImage `json:"image" gorethink:"image"`
	// Build created time
	Created time.Time `json:"created" gorethink:"created"`
	// Build updated time
	Updated time.Time `json:"updated" gorethink:"updated"`
}

type BuildStatus struct {
	// Build current step
	Step BuildStep `json:"step" gorethink:"step"`
	// Is build cancelled
	Cancelled bool `json:"cancelled" gorethink:"cancelled"`
	// Build executing message
	Message string `json:"message" gorethink:"message"`
	// Build error information
	Error string `json:"error" gorethink:"error"`
	// Build status updated time
	Updated time.Time `json:"updated" gorethink:"updated"`
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
	Hub string `json:"hub" gorethink:"hub"`
	// Build sources owner
	Owner string `json:"owner" gorethink:"owner"`
	// Build sources repo
	Repo string `json:"repo" gorethink:"repo"`
	// Build source tag (branch, tag)
	Tag string `json:"tag" gorethink:"tag"`
	// Build commit information
	Commit GitSourceCommit `json:"commit" gorethink:"commit"`
	// Build sources auth reference
}

type BuildImage struct {
	// Build image repo name
	Repo string `json:"repo" gorethink:"repo"`
	// Build image tag name
	Tag string `json:"tag" gorethink:"tag"`
	// Build image registry reference
	Registry string `json:"registry" gorethink:"registry"`
}

type GitSourceCommit struct {
	// Git commit information hash
	Commit string `json:"commit" gorethink:"id"`
	// Git committer gravatar
	Committer string `json:"committer" gorethink:"id"`
	// Git committer email
	Author string `json:"author" gorethink:"id"`
	// Git commit message
	Message string `json:"message" gorethink:"id"`
}
