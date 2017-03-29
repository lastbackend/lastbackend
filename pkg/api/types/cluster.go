package types

import "time"

type ClusterList []Cluster

type Cluster struct {
	// Cluster uuid, generated automatically
	ID string `json:"id"`
	// Cluster owner username
	Owner string `json:"owner"`
	// Cluster name
	Name string `json:"name"`
	// Cluster region
	Region string `json:"name"`
	// Cluster labels lists
	Labels map[string]string `json:"labels"`
	// Cluster created time
	Created time.Time `json:"created"`
	// Cluster updated time
	Updated time.Time `json:"updated"`
}
