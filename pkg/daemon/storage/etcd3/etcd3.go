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

package etcd3

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	s "github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"path"
	"time"
)

var debug bool

func New(client *clientv3.Client, codec serializer.Codec, prefix string) s.IStore {
	var result = &store{
		client:     client,
		codec:      codec,
		pathPrefix: path.Join("/", prefix),
	}
	result.opts = append(result.opts, clientv3.WithSerializable())
	return result
}

func SetDebug() {
	debug = true
}

// printf formats according to a format specifier and writes to standard output.
func printf(format string, a ...interface{}) {
	if debug {
		t := time.Now()
		d := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
		fmt.Printf(fmt.Sprintf("DEBU[%s] %s\n", d, format), a...)
	}
}
