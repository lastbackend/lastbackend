//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package service

import (
	"errors"
	"strings"

	"github.com/lastbackend/lastbackend/internal/cli/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
)

type ClusterService struct {
	storage storage.IStorage
}

func NewClusterService(stg storage.IStorage) *ClusterService {
	s := new(ClusterService)
	s.storage = stg
	return s
}

func (s *ClusterService) List() ([]*models.Cluster, error) {
	return s.listLocalCluster()
}

func (s *ClusterService) AddLocalCluster(name, endpoint, token string, local bool) error {
	items, err := s.listLocalCluster()
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.Name == name {
			return errors.New("already exists")
		}
	}

	cluster := new(models.Cluster)
	cluster.Name = name
	cluster.Endpoint = endpoint
	cluster.Token = token

	items = append(items, cluster)

	return s.storage.Set("cluster", "list", items)
}

func (s *ClusterService) GetLocalCluster(name string) (*models.Cluster, error) {
	items, err := s.listLocalCluster()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}

	return nil, nil
}

func (s *ClusterService) DelLocalCluster(name string) error {

	items, err := s.listLocalCluster()
	if err != nil {
		return err
	}

	for i, item := range items {
		if item.Name == name {

			cl, err := s.GetCluster()
			if err != nil {
				return err
			}

			match := strings.Split(cl, ".")
			if match[0] == "l" && match[1] == name {
				if err := s.SetCluster(""); err != nil {
					return err
				}
			}

			items = append(items[:i], items[i+1:]...)
			break
		}
	}

	return s.storage.Set("cluster", "list", items)
}

func (s *ClusterService) SetCluster(cluster string) error {
	return s.storage.Set("cluster", "current", cluster)
}

func (s *ClusterService) GetCluster() (string, error) {
	var cluster string
	err := s.storage.Get("cluster", "current", &cluster)
	return cluster, err
}

func (s *ClusterService) listLocalCluster() ([]*models.Cluster, error) {
	items := make([]*models.Cluster, 0)
	err := s.storage.Get("cluster", "list", &items)
	return items, err
}
