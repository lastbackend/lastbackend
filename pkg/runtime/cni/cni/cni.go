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
// +build !linux

package cni

import (
	"github.com/lastbackend/lastbackend/pkg/runtime/cni"
	"github.com/lastbackend/lastbackend/pkg/runtime/cni/local"
	"github.com/spf13/viper"
)

func New(v *viper.Viper) (cni.CNI, error) {
	switch v.GetString("network.cni.type")  {
	default:
		return local.New()
	}
}
