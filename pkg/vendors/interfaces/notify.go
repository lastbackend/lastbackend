//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
