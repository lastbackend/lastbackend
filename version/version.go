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

package version

var (
	// Package is filled at linking time
	Package = "github.com/lastbackend/lastbackend"

	// Version holds the complete version number. Filled in at linking time.
	Version = "1.0.0"

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""
)

