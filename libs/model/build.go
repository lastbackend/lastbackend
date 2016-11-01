package model

import "time"

type Build struct {
	// Build number, incremented automatically
	Number int `json:"number, omitempty"`
	// Build executing status
	Status BuildStatus `json:"status,omitempty"`
	// Build sources used for build
	Source BuildSource `json:"source,omitempty"`
	// Build image output information
	Image BuildImage `json:"image, omitempty"`
	// Build created time
	Created time.Time `json:"created, omitempty"`
	// Build updated time
	Updated time.Time `json:"updated, omitempty"`
	// Request called and created this build
	Request BuildRequest `json:"request,omitempty"`
}

type BuildRequest struct {
	// Build request type
	Type string `json:"type, omitempty"`
	// Build owner
	Owner string `json:"owner, omitempty"`
}

type BuildStatus struct {
	// Build current step
	Step BuildStep `json:"step,omitempty"`
	// Is build cancelled
	Cancelled bool `json:"cancelled,omitempty"`
	// Build executing message
	Message string `json:"message,omitempty"`
	// Build error information
	Error string `json:"error,omitempty"`
	// Build status updated time
	Updated time.Time `json:"updated,omitempty"`
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
	Hub string `json:"hub,omitempty"`
	// Build sources repo
	Repo string `json:"repo,omitempty"`
	// Build source tag (branch, tag)
	Tag string `json:"tag,omitempty"`
	// Build commit information
	Commit GitSourceCommit `json:"commit,omitempty"`
	// Build sources auth reference
	// Auth *api.LocalObjectReference `json:"sourceAuth,omitempty"`
}

type BuildImage struct {
	// Build image repo name
	Repo string `json:"repo,omitempty"`
	// Build image tag name
	Tag string `json:"tag,omitempty"`
	// Build image registry reference
	Registry BuildStorage `json:"registry"`
}

type BuildStorage struct {
	// Storage host
	Host string `json:"host"`
	// Storage auth token
	Token string `json:"token"`
}

type GitSourceCommit struct {
	// Git commit information hash
	Commit string `json:"commit,omitempty"`
	// Git committer gravatar
	Committer string `json:"committer,omitempty"`
	// Git committer email
	Author string `json:"author,omitempty"`
	// Git commit message
	Message string `json:"message,omitempty"`
}
