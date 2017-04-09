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

package types

import (
	"time"
)

type BuildList []Build

type Build struct {
	// Build Meta
	Meta Meta `json:"meta"`
	// Build executing status
	Status BuildStatus `json:"status"`
	// Build sources used for build
	Source BuildSource `json:"source"`
}

type BuildStatus struct {
	// Build current step
	Step BuildStep `json:"step"`
	// Is build cancelled
	Cancelled bool `json:"cancelled"`
	// Build executing message
	Message string `json:"message"`
	// Build error information
	Error string `json:"error"`
	// Build status updated time
	Updated time.Time `json:"updated"`
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
	Hub string `json:"hub"`
	// Build sources owner
	Owner string `json:"owner"`
	// Build sources repo
	Repo string `json:"repo"`
	// Build source tag (branch, tag)
	Tag string `json:"tag"`
	// Build commit information
	Commit GitSourceCommit `json:"commit"`
	// Build sources auth reference
}

type GitSourceCommit struct {
	// Git commit information hash
	Commit string `json:"commit"`
	// Git committer gravatar
	Committer string `json:"committer"`
	// Git committer username
	Author string `json:"author"`
	// Git committer username
	Email string `json:"email"`
	// Git commit message
	Message string `json:"message"`
	// Git commit date
	Date time.Time `json:"date"`
}
