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

package vendors

import (
	"github.com/lastbackend/lastbackend/pkg/vendors/bitbucket"
	"github.com/lastbackend/lastbackend/pkg/vendors/github"
	"github.com/lastbackend/lastbackend/pkg/vendors/gitlab"
	"github.com/lastbackend/lastbackend/pkg/vendors/slack"
	"github.com/lastbackend/lastbackend/pkg/vendors/wechat"
)

func GetGitHub(clientID, clientSecretID, redirectURI string) *github.GitHub {
	return github.GetClient(clientID, clientSecretID, redirectURI)
}

func GetBitBucket(clientID, clientSecretID, redirectURI string) *bitbucket.BitBucket {
	return bitbucket.GetClient(clientID, clientSecretID, redirectURI)
}

func GetGitLab(clientID, clientSecretID, redirectURI string) *gitlab.GitLab {
	return gitlab.GetClient(clientID, clientSecretID, redirectURI)
}

func GetSlack(clientID, clientSecretID, redirectURI string) *slack.Slack {
	return slack.GetClient(clientID, clientSecretID, redirectURI)
}

func GetWeChat(clientID, clientSecretID, redirectURI string) *wechat.WeChat {
	return wechat.GetClient(clientID, clientSecretID, redirectURI)
}
