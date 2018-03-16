//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package docker

import (
	"encoding/json"
	"time"
)

type Repository struct {
	StarCount int64  `json:"star_count"`
	PullCount int64  `json:"pull_count"`
	Hub       string `json:"hub"`
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Automated bool   `json:"automated"`
	Official  bool   `json:"official"`
}

type RepositoryList []Repository

func (obj *Repository) ToJson() ([]byte, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (list *RepositoryList) ToJson() ([]byte, error) {

	if list == nil {
		return []byte(`[]`), nil
	}

	buf, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

type Tag struct {
	Name        string    `json:"name"`
	ID          int64     `json:"id"`
	Size        int64     `json:"size"`
	Repo        int64     `json:"repo"`
	Creator     int64     `json:"creator"`
	LastUpdater int64     `json:"last_updated"`
	V2          bool      `json:"v2"`
	ImageID     int64     `json:"image_id"`
	Platforms   []int64   `json:"platforms"`
	LastUpdated time.Time `json:"last_updated"`
}

type TagList struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Tags  []Tag  `json:"tags"`
}

func (obj *Tag) ToJson() ([]byte, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (list *TagList) ToJson() ([]byte, error) {

	if list == nil {
		return []byte(`[]`), nil
	}

	buf, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
