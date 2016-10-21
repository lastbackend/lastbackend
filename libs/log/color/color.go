package color

import (
	"github.com/deployithq/deployit/libs/log/filesystem"
)

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
