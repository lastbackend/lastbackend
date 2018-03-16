//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

import (
	"github.com/lastbackend/lastbackend/pkg/vendors/types"
)

type IVCS interface {
	VendorInfo() *types.Vendor
	GetUser() (*types.User, error)
	ListRepositories(username string, org bool) (*types.VCSRepositories, error)
	ListBranches(owner, repo string) (*types.VCSBranches, error)
	CreateHook(id, owner, repo, host string) (*string, error)
	RemoveHook(id, owner, repo string) error
	PushPayload(data []byte) (*types.VCSBranch, error)
}
