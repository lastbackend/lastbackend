package docker

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/libs/model"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetRepository(name string) (*model.DockerRepositoryList, error) {

	var page, size int64 = 1, 10
	var results = struct {
		Count    int64  `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			StarCount int64  `json:"star_count"`
			PullCount int64  `json:"pull_count"`
			Repo      string `json:"repo_name"`
			Owner     string `json:"repo_owner"`
			Desc      string `json:"short_description"`
			Automated bool   `json:"is_automated"`
			Official  bool   `json:"is_official"`
		} `json:"results"`
	}{}

	var url = fmt.Sprintf("https://%s/v2/search/repositories/?query=%s&page=%d&page_size=%d", "hub.docker.com", name, page, size)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &results)
	if err != nil {
		return nil, err
	}

	var repos = new(model.DockerRepositoryList)
	for _, item := range results.Results {

		var owner, name = "", ""
		var items = strings.Split(item.Repo, "/")

		if len(items) == 2 {
			owner = items[0]
			name = items[1]
		} else if len(items) == 1 {
			owner = "library"
			name = items[0]
		}

		*repos = append(*repos, model.DockerRepository{
			StarCount: item.StarCount,
			PullCount: item.PullCount,
			Hub:       "index.docker.io",
			Owner:     owner,
			Name:      name,
			Desc:      item.Desc,
			Automated: item.Automated,
			Official:  item.Official,
		})
	}

	return repos, nil
}

func ListTag(owner, name string) (*model.DockerTagList, error) {

	var page, size int64 = 1, 10

	var results = struct {
		Count    int64  `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name        string    `json:"name"`
			Size        int64     `json:"full_size"`
			ID          int64     `json:"id"`
			Repo        int64     `json:"repository"`
			Creator     int64     `json:"creator"`
			LastUpdater int64     `json:"last_updater"`
			LastUpdated time.Time `json:"last_updated"`
			ImageID     int64     `json:"image_id"`
			V2          bool      `json:"v2"`
			Platforms   []int64   `json:"platforms"`
		} `json:"results"`
	}{}

	var url = fmt.Sprintf("https://%s/v2/repositories/%s/%s/tags/?page=%d&page_size=%d", "hub.docker.com", owner, name, page, size)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &results)
	if err != nil {
		return nil, err
	}

	var tags = new(model.DockerTagList)
	for _, item := range results.Results {
		*tags = append(*tags, model.DockerTag{
			ID:          item.ID,
			Name:        item.Name,
			Size:        item.Size,
			Repo:        item.Repo,
			Creator:     item.Creator,
			LastUpdater: item.LastUpdater,
			LastUpdated: item.LastUpdated,
			ImageID:     item.ImageID,
			V2:          item.V2,
			Platforms:   item.Platforms,
		})
	}

	return tags, nil
}
