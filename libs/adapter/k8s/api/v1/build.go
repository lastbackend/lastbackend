package v1

import (
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"time"
)

type Build struct {
	unversioned.TypeMeta `json:",inline"`
	// metadata for Build
	api.ObjectMeta `json:"metadata,omitempty"`
	// build Spec
	Spec BuildSpec `json:"spec,omitempty"`
}

type BuildSpec struct {
	// CommonSpec is the information that represents a build
	CommonSpec `json:",inline"`
	// Request called and created this build
	Request BuildRequest `json:"request,omitempty"`
}

type CommonSpec struct {
	// Build number, incremented automatically
	Number int `json:"number, omitempty"`
	// Build executing status
	Status BuildStatus `json:"status"`
	// Build sources used for build
	Source BuildSource `json:"source,omitempty"`
	// Build image output information
	Image BuildImage `json:"image, omitempty"`
	// Build created time
	Created time.Time `json:"created, omitempty"`
	// Build updated time
	Updated time.Time `json:"updated, omitempty"`
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
	Auth *api.LocalObjectReference `json:"sourceAuth,omitempty"`
}

type BuildImage struct {
	// Build image repo name
	Repo string `json:"repo,omitempty"`
	// Build image tag name
	Tag string `json:"tag,omitempty"`
	// Build image registry reference
	Registry *api.LocalObjectReference `json:"registry,omitempty"`
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

// BuildList is a collection of Builds.
type BuildList struct {
	unversioned.TypeMeta `json:",inline"`
	// metadata for BuildList.
	unversioned.ListMeta `json:"metadata,omitempty"`
	// items is a list of builds
	Items []Build `json:"items"`
}

func (obj *Build) GetObjectKind() unversioned.ObjectKind { return &obj.TypeMeta }
