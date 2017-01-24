package model

import (
	"encoding/json"
	"time"
)

type DockerRepository struct {
	StarCount int64  `json:"name"`
	PullCount int64  `json:"pull_count"`
	Hub       string `json:"hub"`
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Automated bool   `json:"automated"`
	Official  bool   `json:"official"`
}

type DockerRepositoryList []DockerRepository

type DockerTag struct {
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

type DockerTagList []DockerTag

func (s *DockerRepository) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *DockerRepositoryList) ToJson() ([]byte, error) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *DockerTag) ToJson() ([]byte, error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *DockerTagList) ToJson() ([]byte, error) {

	if s == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
