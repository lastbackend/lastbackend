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
	"encoding/json"
	"sync"
	"time"
)

type RepoList []*Repo

type Repo struct {
	lock sync.RWMutex
	// Repo meta
	Meta RepoMeta `json:"meta"`
	// Repo state
	State RepoState `json:"state"`
	// Repo stats
	Stats RepoStats `json:"stats"`
	// Repo source
	Source RepoSource `json:"source"`
	// Repo tags
	Rules map[string]*RepoBuildRule `json:"rules"`
	// Repo tags
	Tags map[string]*RepoTag `json:"tags"`
	// Repo created time
	Registry string `json:"registry"`
	// Repo meta readme
	Readme string `json:"readme"`
	// Remote repo option
	Remote bool `json:"remote"`
	// Private repo option
	Private bool `json:"private"`
	// Repo created time
	Created time.Time `json:"created"`
	// Repo updated time
	Updated time.Time `json:"updated"`
}

type RepoMeta struct {
	Meta
	// Repo meta user
	Owner string `json:"owner"`
	// Repo readme
	Technology string `json:"technology"`
}

type RepoSource struct {
	// Repo source hub
	Hub string `json:"hub"`
	// Repo source owner
	Owner string `json:"owner"`
	// Repo source name
	Name string `json:"name"`
}

type RepoState struct {
	// Repo state
	State string `json:"state"`
	// Repo state status
	Status string `json:"status"`
	// Meta deleted
	Deleted bool `json:"deleted"`
	// Meta liked
	Liked bool `json:"liked"`
}

type RepoStats struct {
	// Repo stats pulls quantity
	PullsQuantity int64 `json:"pulls_quantity"`
	// Repo stats builds quantity
	BuildsQuantity int64 `json:"builds_quantity"`
	// Repo stats services quantity
	ServicesQuantity int64 `json:"services_quantity"`
	// Repo stats stars quantity
	StarsQuantity int64 `json:"stars_quantity"`
	// Repo stats views quantity
	ViewsQuantity int64 `json:"views_quantity"`
}

type RepoBuildRule struct {
	SystemMeta struct {
		Insert bool
		Update bool
		Delete bool
	} `json:"-"`
	// Repo rule id
	ID string `json:"id"`
	// Repo rule branch
	Branch string `json:"branch"`
	// Repo rule filepath
	FilePath string `json:"filepath"`
	// Repo rule tag
	Tag string `json:"tag"`
	// Repo rule registry
	Registry string `json:"registry"`
	// Repo rule autobuild
	AutoBuild bool `json:"autobuild"`
	// Repo rule disabled
	Disabled bool `json:"disable"`
	// Repo rule config
	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

type RepoTag struct {
	ID        string `json:"uuid"`
	RepoID    string `json:"repo_id"`
	Name      string `json:"name"`
	BuildSize int64  `json:"size"`
	Layers    struct {
		Count int64 `json:"count"`
		Size  struct {
			Average int64 `json:"average"`
			Max     int64 `json:"max"`
		} `json:"size"`
	} `json:"layers"`
	Build0   BuildInfo `json:"build_0"`
	Build1   BuildInfo `json:"build_1"`
	Build2   BuildInfo `json:"build_2"`
	Build3   BuildInfo `json:"build_3"`
	Build4   BuildInfo `json:"build_4"`
	Disabled bool      `json:"disabled"`
	Updated  time.Time `json:"updated"`
	Created  time.Time `json:"created"`
}

type BuildInfo struct {
	ID     string `json:"id"`
	Number int64  `json:"number"`
	Status string `json:"status"`
}

type RepoDeployTemplate struct {
	// Repo deploy template unique identificator
	ID string `json:"id"`
	// Repo deploy template unique repo identificator
	RepoID string `json:"repo"`
	// Repo deploy template name
	Name string `json:"name"`
	// Repo deploy template description
	Description string `json:"description"`
	// Repo deploy template entrypoint
	Entrypoint StringSlice `json:"entrypoint"`
	// Repo deploy template command
	Ports RepoDeployTemplatePortList `json:"ports"`
	// Repo deploy template memory
	Memory int64 `json:"memory"`
	// Repo deploy template option for mark as default
	Main bool `json:"main"`
	// Repo deploy template shared option
	Shared bool `json:"shared"`
	// Repo deploy template created time
	Created time.Time `json:"created"`
	// Repo deploy template updated time
	Updated time.Time `json:"updated"`
}

type RepoDeployTemplatePort struct {
	Protocol  string `json:"protocol"`
	Container int    `json:"internal"`
	Host      int    `json:"external"`
	Published bool   `json:"published"`
}

type RepoDeployTemplatePortList []RepoDeployTemplatePort

func (s *RepoDeployTemplatePortList) ToJson() string {
	if s == nil {
		return EmptyStringSlice
	}
	res, err := json.Marshal(s)
	if err != nil {
		return EmptyStringSlice
	}
	return string(res)
}
