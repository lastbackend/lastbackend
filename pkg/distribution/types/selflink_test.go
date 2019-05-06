//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package types_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNamespaceSelfLink_MarshalJSON(t *testing.T) {

	ns := types.Namespace{}

	ns.Meta.Name = "test"
	ns.Meta.SelfLink = *types.NewNamespaceSelfLink("test")

	nsj, _ := json.Marshal(ns)

	ns2 := types.Namespace{}

	_ = json.Unmarshal(nsj, &ns2)

	assert.Equal(t, ns.SelfLink().String(), ns2.SelfLink().String(), "equal")
}

func TestServiceSelfLink_MarshalJSON(t *testing.T) {

	ns := types.Service{}

	ns.Meta.Name = "test"
	ns.Meta.SelfLink = *types.NewServiceSelfLink("demo", "test")

	nsj, _ := json.Marshal(ns)
	ns2 := types.Service{}

	_ = json.Unmarshal(nsj, &ns2)

	assert.Equal(t, ns.SelfLink().String(), ns2.SelfLink().String(), "equal")
}
