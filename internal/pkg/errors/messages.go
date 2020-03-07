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

package errors

import "errors"

var (
	NotLoggedMessage   = errors.New("you are currently not logged in to the system, to get proper access create a new user or login with an existing user")
	LoginErrorMessage  = errors.New("incorrect login or password")
	LogoutErrorMessage = errors.New("some problems with logout")
	UnknownMessage     = errors.New("oops, error occurred: Please provide bug to github: https://github.com/lastbackend/registry/issues")
)
