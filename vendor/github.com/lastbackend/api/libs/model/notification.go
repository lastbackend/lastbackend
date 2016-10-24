package model

type Notification struct {
	UUID          string
	UserID        string
	ComponentType string
	ComponentID   string
	Service       string
	Channel       string
	Level1        int64
	Level2        int64
	Level3        int64
	Active        bool
}

type Notifications []Notification
