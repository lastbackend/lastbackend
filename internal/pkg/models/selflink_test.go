//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package models_test

import (
	"encoding/json"
	"testing"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestNamespaceSelfLink_MarshalJSON(t *testing.T) {

	ns := models.Namespace{}

	ns.Meta.Name = "test"
	ns.Meta.SelfLink = *models.NewNamespaceSelfLink("test")

	nsj, _ := json.Marshal(ns)

	ns2 := models.Namespace{}

	_ = json.Unmarshal(nsj, &ns2)

	assert.Equal(t, ns.SelfLink().String(), ns2.SelfLink().String(), "equal")
}

func TestServiceSelfLink_MarshalJSON(t *testing.T) {

	ns := models.Service{}

	ns.Meta.Name = "test"
	ns.Meta.SelfLink = *models.NewServiceSelfLink("demo", "test")

	nsj, _ := json.Marshal(ns)
	ns2 := models.Service{}

	_ = json.Unmarshal(nsj, &ns2)

	assert.Equal(t, ns.SelfLink().String(), ns2.SelfLink().String(), "equal")
}
