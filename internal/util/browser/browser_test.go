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

package browser

import "testing"

func TestOpen(t *testing.T) {
	tmp := CommandWrapper
	CommandWrapper = func(name string, parameters ...string) error {
		return nil
	}

	err := Open("http://dummmy")

	if err != nil {
		t.Error("Unexpected error")
	}

	CommandWrapper = tmp
}

func TestOpenFail(t *testing.T) {
	tmp := Os
	Os = "Dummy"
	err := Open("http://dummmy")

	if err == nil {
		t.Error("Unexpected successfully url call")
	}

	Os = tmp
}
