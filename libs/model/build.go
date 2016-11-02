package model

import "time"

type BuildList []Build

type Build struct {
	// Build number, incremented automatically
	ID int `json:"id, omitempty" gorethink:"id,omitempty"`
	// Build executing status
	Status BuildStatus `json:"status,omitempty" gorethink:"status,omitempty"`
	// Build sources used for build
	Source BuildSource `json:"source,omitempty" gorethink:"source,omitempty"`
	// Build image output information
	Image BuildImage `json:"image, omitempty" gorethink:"image,omitempty"`
	// Build created time
	Created time.Time `json:"created, omitempty" gorethink:"created,omitempty"`
	// Build updated time
	Updated time.Time `json:"updated, omitempty" gorethink:"updated,omitempty"`
}

type BuildStatus struct {
	// Build current step
	Step BuildStep `json:"step,omitempty" gorethink:"step,omitempty"`
	// Is build cancelled
	Cancelled bool `json:"cancelled,omitempty" gorethink:"cancelled,omitempty"`
	// Build executing message
	Message string `json:"message,omitempty" gorethink:"message,omitempty"`
	// Build error information
	Error string `json:"error,omitempty" gorethink:"error,omitempty"`
	// Build status updated time
	Updated time.Time `json:"updated,omitempty" gorethink:"updated,omitempty"`
}

type BuildStep string

const (
	//BuildStepCreate - The first step after build creating
	BuildStepCreate = "create"
	//BuildStepFetch - Fetch sources step
	BuildStepFetch = "fetch"
	//BuildStepBuild - Build executing step
	BuildStepBuild = "build"
	//BuildStepUpload - Upload builded docker image step
	BuildStepUpload = "upload"
)

type BuildSource struct {
	// Build sources hub
	Hub string `json:"hub,omitempty" gorethink:"hub,omitempty"`
	// Build sources repo
	Repo string `json:"repo,omitempty" gorethink:"repo,omitempty"`
	// Build source tag (branch, tag)
	Tag string `json:"tag,omitempty" gorethink:"tag,omitempty"`
	// Build commit information
	Commit GitSourceCommit `json:"commit,omitempty" gorethink:"commit,omitempty"`
	// Build sources auth reference
}

type BuildImage struct {
	// Build image repo name
	Repo string `json:"repo,omitempty" gorethink:"repo,omitempty"`
	// Build image tag name
	Tag string `json:"tag,omitempty" gorethink:"tag,omitempty"`
	// Build image registry reference
	Registry string `json:"registry,omitempty" gorethink:"registry,omitempty"`
}

type GitSourceCommit struct {
	// Git commit information hash
	Commit string `json:"commit,omitempty" gorethink:"id,omitempty"`
	// Git committer gravatar
	Committer string `json:"committer,omitempty" gorethink:"id,omitempty"`
	// Git committer email
	Author string `json:"author,omitempty" gorethink:"id,omitempty"`
	// Git commit message
	Message string `json:"message,omitempty" gorethink:"id,omitempty"`
}
