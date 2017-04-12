package v1

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"time"
)

type Event struct {
	Meta      types.NodeMeta       `json:"meta"`
	State     types.NodeState      `json:"state"`
	Pods      []types.PodNodeState `json:"pods"`
	Timestamp time.Time            `json:"timestamp"`
}
