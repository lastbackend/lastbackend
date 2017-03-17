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

package color

import "github.com/lastbackend/lastbackend/pkg/util/filesystem"

func Black(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[30m" + s + "\x1b[39m"
	}
	return s
}

func Red(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[31m" + s + "\x1b[39m"
	}
	return s
}

func Green(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[32m" + s + "\x1b[39m"
	}
	return s
}

func Yellow(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[33m" + s + "\x1b[39m"
	}
	return s
}

func Blue(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[34m" + s + "\x1b[39m"
	}
	return s
}

func Magenta(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[35m" + s + "\x1b[39m"
	}
	return s
}

func Cyan(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[36m" + s + "\x1b[39m"
	}
	return s
}

func White(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[37m" + s + "\x1b[39m"
	}
	return s
}

func Default(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[39m" + s + "\x1b[39m"
	}
	return s
}

func LightGray(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[90m" + s + "\x1b[39m"
	}
	return s
}

func LightRed(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[91m" + s + "\x1b[39m"
	}
	return s
}

func LightGreen(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[92m" + s + "\x1b[39m"
	}
	return s
}

func LightYellow(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[93m" + s + "\x1b[39m"
	}
	return s
}

func LightBlue(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[94m" + s + "\x1b[39m"
	}
	return s
}

func LightMagenta(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[95m" + s + "\x1b[39m"
	}
	return s
}

func LightCyan(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[96m" + s + "\x1b[39m"
	}
	return s
}

func LightWhite(s string) string {
	if !filesystem.IsWindows() {
		return "\x1b[97m" + s + "\x1b[39m"
	}
	return s
}
