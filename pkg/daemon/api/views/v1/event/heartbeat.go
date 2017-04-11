package event

import "time"

type Heartbeat struct {
	Timestamp time.Time `json:"timestamp"`
}

func NewHeartBeatEvent() *Heartbeat {
	return &Heartbeat{}
}
