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

package deployment

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

// Provision deployment
// Remove deployment or cancel if deployment is market for destroy
// Remove deployment if no active pods present and deployment is marked for destroy
func Provision(d *types.Deployment) error {

	var (

	)

	log.Debugf("Deployment Controller: provision deployment: %s/%s", d.Meta.Namespace, d.Meta.Name)

	// Check deployment is marked for destroy

	// Get all pods per deployment

	return nil
}

func Create () error {
	return nil
}

func Cancel () error {
	return nil
}

func Remove () error {
	return nil
}

func Destroy () error {
	return nil
}