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

// +build linux

package ipvs

const (
	SvcFlagPersist   = 0x1
	SvcFlagHashed    = 0x2
	SvcFlagOnePacket = 0x4

	DstFlagFwdMask   = 0x7
	DstFlagFwdMasq   = 0x0
	DstFlagFwdLocal  = 0x1
	DstFlagFwdTunnel = 0x2
	DstFlagFwdRoute  = 0x3
	DstFlagFwdBypass = 0x4
	DstFlagSync      = 0x20
	DstFlagHashed    = 0x40
	DstFlagNoOutput  = 0x80
	DstFlagInactive  = 0x100
	DstFlagOutSeq    = 0x200
	DstFlagInSeq     = 0x400
	DstFlagSeqMask   = 0x600
	DstFlagNoCPort   = 0x800
	DstFlagTemplate  = 0x1000
	DstFlagOnePacket = 0x2000
)

const (
	proxyTCPProto = "tcp"
	proxyUDPProto = "udp"
)

type Service struct {
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	Type        string    `json:"type"`
	Scheduler   string    `json:"scheduler"`
	Persistence int       `json:"persistence"`
	Netmask     string    `json:"netmask"`
	Backends    map[string]Backend `json:"backends"`
}

type Backend struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Forwarder      string `json:"forwarder"`
	Weight         int    `json:"weight"`
	UpperThreshold int    `json:"upper_threshold"`
	LowerThreshold int    `json:"lower_threshold"`
}

type DestinationFlags uint32

type DestinationStats struct {
	ActiveConns   uint32
	InactiveConns uint32
	PersistConns  uint32
}
