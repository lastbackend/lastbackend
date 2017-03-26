package model

import (
	"encoding/json"
	"github.com/lastbackend/vendors/interfaces"
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
