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

package types

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {

	namespace := "483705e8-0ae0-453f-a973-8b79457ca23a"
	id := "483705e8-0ae0-453f-a973-8b79457ca22a"


	var asset = Secret{}
	asset.Meta.ID = id
	asset.Meta.NamespaceID = namespace

	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s:%s", namespace, id)))
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	assert.Equal(t, hash, asset.GetHash(), "hash generation validation")
}