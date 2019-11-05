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

package utils

import (
	"bytes"
	"github.com/lastbackend/lastbackend/tools/log"
	"os/exec"
	"strings"
)

type FDBRule struct {
	Mac       string
	Device    string
	DST       string
	Vlan      string
	Master    string
	Self      bool
	Permanent bool
}

func BridgeFDBList() ([]FDBRule, error) {

	var rules []FDBRule

	fdblcmd := exec.Command("bridge", "fdb")

	var stdout, stderr bytes.Buffer
	fdblcmd.Stdout = &stdout
	fdblcmd.Stderr = &stderr

	if err := fdblcmd.Run(); err != nil {
		log.Errorf("cmd.Run() failed with %s\n", err.Error())
		return rules, err
	}

	fdbs := strings.Split(string(stdout.Bytes()), "\n")
	rules = make([]FDBRule, len(fdbs))
	for _, fdb := range fdbs {
		rule := BridgeFDBParse(fdb)
		if rule.Mac != "" {
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func BridgeFDBParse(r string) FDBRule {
	var rule FDBRule

	var rl = strings.Fields(r)
	if len(rl) == 0 {
		return rule
	}

	if len(strings.Split(rl[0], ":")) != 6 {
		return rule
	}

	rule.Mac = rl[0]

	for i := 1; i < len(rl); i++ {

		switch rl[i] {
		case "dev":
			rule.Device = rl[i+1]
			i++
			break
		case "dst":
			rule.DST = rl[i+1]
			i++
			break
		case "vlan":
			rule.Vlan = rl[i+1]
			i++
			break
		case "master":
			rule.Master = rl[i+1]
			i++
			break
		case "self":
			rule.Self = true
			break
		case "permanent":
			rule.Permanent = true
			break
		}
	}

	return rule
}
