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

package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

func AddGlobalFlags(fs *pflag.FlagSet) {

	// lookup flags in global flag set and re-register the values with our flagset
	global := pflag.CommandLine
	local := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	pflagRegister(global, local, "verbose", "v")
	pflagRegister(global, local, "config", "c")

	fs.AddFlagSet(local)
}

func pflagRegister(global, local *pflag.FlagSet, globalName string, shortName string) {
	if f := global.Lookup(globalName); f != nil {
		f.Name = normalize(f.Name)
		f.Shorthand = shortName
		local.AddFlag(f)
	} else {
		panic(fmt.Sprintf("failed to find flag in global flagset (pflag): %s", globalName))
	}
}

func normalize(s string) string {
	return strings.Replace(s, "_", "-", -1)
}
