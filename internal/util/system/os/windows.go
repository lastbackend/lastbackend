// +build windows
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

package os

import (
	"bytes"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/util/system/types"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetInfo() *types.OsInfo {
	info := strings.Replace(info(), "\n", "", -1)
	info = strings.Replace(info, "\r\n", "", -1)
	index1 := strings.Index(info, "[Version")
	index2 := strings.Index(info, "]")
	ver := "unknown"
	if index1 != -1 && index2 != -1 {
		ver = info[index1+9 : index2]
	}
	hostname, _ := os.Hostname()

	return &types.OsInfo{
		Kernel:   "Windows",
		Core:     ver,
		Platform: "unknown",
		OS:       "Windows",
		GoOS:     runtime.GOOS,
		CPUs:     runtime.NumCPU(),
		Hostname: hostname,
	}
}

func info() string {
	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)
	cmd := exec.Command("cmd", "ver")
	cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("error get info %s", err.Error())
		return ""
	}
	return out.String()
}
