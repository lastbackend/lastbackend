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

package resource

import (
	"fmt"
	"github.com/docker/go-units"
	"math/big"
)

const (
	// m
	mili = 100
	// MB - MegaByte size
	MB = 1000 * 1000
	// MIB - MegaByte size
	MIB = 1024 * 1024
	// GB - GigaByte size
	GB = 1000 * 1000 * 1000
	// GIB - GibiByte size
	GIB = 1024 * 1024 * 1024
)

// parseResource - parse resource size string
// m,mb,mib,gb,gib,kb,kib,
func DecodeMemoryResource(value string) (int64, error) {
	return units.RAMInBytes(value)
}

func EncodeMemoryResource(res int64) string {
	return units.BytesSize(float64(res))
}

// returns microseconds
func DecodeCpuResource(value string) (int64, error) {
	cpu, ok := new(big.Rat).SetString(value)
	if !ok {
		return 0, fmt.Errorf("failed to parse %v as a rational number", value)
	}
	nano := cpu.Mul(cpu, big.NewRat(1e9, 1))
	if !nano.IsInt() {
		return 0, fmt.Errorf("value is too precise")
	}
	return nano.Num().Int64(), nil
}

func EncodeCpuResource(res int64) string {
	return big.NewRat(res, 1e9).FloatString(3)
}