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
	"github.com/lastbackend/lastbackend/pkg/vendors/interfaces"
)

type VCSRepository struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
}

func (r *VCSRepository) Convert(repository *interfaces.VCSRepository) {

	if repository == nil {
		return
	}

	r.Name = repository.Name
	r.Description = repository.DefaultBranch
	r.Private = repository.Private
	r.DefaultBranch = repository.DefaultBranch
}

func (r *VCSRepository) ToJson() ([]byte, error) {
	buf, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

type VCSRepositoryList []VCSRepository

func (r *VCSRepositoryList) Convert(repositories *interfaces.VCSRepositories) {

	if repositories == nil {
		return
	}

	for _, repo := range *repositories {
		item := VCSRepository{}
		item.Convert(&repo)
		*r = append(*r, item)
	}
}

func (r *VCSRepositoryList) ToJson() ([]byte, error) {
	if r == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

type VCSBranch struct {
	Name string `json:"name"`
}

func (b *VCSBranch) Convert(branch *interfaces.VCSBranch) {

	if branch == nil {
		return
	}

	b.Name = branch.Name
}

func (b *VCSBranch) ToJson() ([]byte, error) {
	buf, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

type VCSBranchList []VCSBranch

func (b *VCSBranchList) Convert(branches *interfaces.VCSBranches) {

	if branches == nil {
		return
	}

	for _, branch := range *branches {
		item := VCSBranch{}
		item.Convert(&branch)
		*b = append(*b, item)
	}
}

func (b *VCSBranchList) ToJson() ([]byte, error) {
	if b == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
