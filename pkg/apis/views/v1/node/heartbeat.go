package node

import "time"

type Heartbeat struct {
	Memory     HostMemory `json:"memory"`
	Pods       int        `json:"pods"`
	Containers int        `json:"containers"`
	Images     int        ` json:"images"`
	Timestamp  time.Time  `json:"timestamp"`
}
