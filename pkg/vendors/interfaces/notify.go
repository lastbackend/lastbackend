package interfaces

import "golang.org/x/oauth2"

// Model

type Notify struct{}

// Types

type NotifyChannel struct {
	ID   string
	Name string
	Type string
}

type NotifyChannels []NotifyChannel

type NotifyGroup struct {
	ID   string
	Name string
	Type string
}

type NotifyGroups []NotifyGroup

// Interfaces

type INotify interface {
	IOAuth2

	ListChannels(token *oauth2.Token) (*NotifyChannels, error)
	ListGroups(token *oauth2.Token) (*NotifyGroups, error)
}

// Functions
