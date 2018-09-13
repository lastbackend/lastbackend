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

package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpecTemplateContainerPort_Parse(t *testing.T) {

	var tests =[]struct{
		name string
		want SpecTemplateContainerPort
		arg  string
	}{
		{
			name: "check full port map",
			want: SpecTemplateContainerPort{
				HostIP: "0.0.0.0",
				HostPort: 2967,
				ContainerPort: 2967,
				Protocol: "tcp",
			},
			arg: "0.0.0.0:2967:2967/tcp",
		},
		{
			name: "check port map without host port",
			want: SpecTemplateContainerPort{
				HostIP: "0.0.0.0",
				HostPort: 2967,
				ContainerPort: 2967,
				Protocol: "tcp",
			},
			arg: "0.0.0.0:2967/tcp",
		},
		{
			name: "check port map without host ip",
			want: SpecTemplateContainerPort{
				HostIP: "127.0.0.1",
				HostPort: 2967,
				ContainerPort: 2967,
				Protocol: "tcp",
			},
			arg: "2967:2967/tcp",
		},
		{
			name: "check port map without host ip and protocol",
			want: SpecTemplateContainerPort{
				HostIP: "127.0.0.1",
				HostPort: 2967,
				ContainerPort: 2967,
				Protocol: "tcp",
			},
			arg: "2967:2967",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			got := SpecTemplateContainerPort{}
			got.Parse(tc.arg)

			assert.Equal(t, tc.want.HostIP, got.HostIP, "host ip mismatch")
			assert.Equal(t, tc.want.HostPort, got.HostPort, "host port mismatch")
			assert.Equal(t, tc.want.ContainerPort, got.ContainerPort, "container port mismatch")
			assert.Equal(t, tc.want.Protocol, got.Protocol, "protocol mismatch")
		})
	}

}