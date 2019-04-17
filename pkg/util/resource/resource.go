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

package resource

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

const (
	KiB = 1024 << (10 * iota)
	MiB
	GiB
	TiB
	PiB
	EiB
	ZiB
	YiB
)

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB
	EB = 1000 * PB
	ZB = 1000 * EB
	YB = 1000 * ZB
)

type unitMap map[string]float64

var (
	decimalMapSize = unitMap{"kb": KB, "mb": MB, "gb": GB, "tb": TB, "pb": PB, "eb": EB, "zb": ZB, "yb": YB}
	binaryMapSize  = unitMap{"kib": KiB, "mib": MiB, "gib": GiB, "tib": TiB, "pib": PiB, "eib": EiB, "zib": ZiB, "yib": YiB}
)

var (
	decimalSizeName = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	binarySizeName  = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
)

// parseResource - parse resource size string
func DecodeMemoryResource(value string) (int64, error) {
	return ToBytes(value)
}

func EncodeMemoryResource(res int64) string {
	return BytesSize(float64(res))
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

var invalidByteQuantityError = errors.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")

// HumanSizeWithPrecision allows the size to be in any precision.
func HumanSizeWithPrecision(size float64, precision int) string {
	size, unit := getSizeAndUnit(size, 1000.0, decimalSizeName)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

// HumanSize returns a human-readable approximation of a size
// capped at 4 valid numbers (eg. "2.752 MB", "128 KB").
func HumanSize(size float64) string {
	return HumanSizeWithPrecision(size, 4)
}

// BytesSize returns a human-readable size in bytes, kibibytes (eg. "12kiB", "32MiB").
func BytesSize(size float64) string {
	size, unit := getSizeAndUnit(size, 1024.0, binarySizeName)
	return fmt.Sprintf("%.4g%s", size, unit)
}

func ToBytes(s string) (int64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	i := strings.IndexFunc(s, unicode.IsLetter)

	if i == -1 {
		return 0, invalidByteQuantityError
	}

	bytesString, multiple := s[:i], s[i:]
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes <= 0 {
		return 0, invalidByteQuantityError
	}

	if multiple == "b" {
		return int64(bytes), nil
	}

	if len(multiple) == 3 {
		if v, ok := binaryMapSize[multiple]; ok {
			return int64(bytes * v), nil
		}

		return 0, invalidByteQuantityError
	}

	if len(multiple) == 2 {
		if v, ok := decimalMapSize[multiple]; ok {
			return int64(bytes * v), nil
		}

		return 0, invalidByteQuantityError
	}

	return 0, invalidByteQuantityError
}

func getSizeAndUnit(size float64, base float64, dictionary []string) (float64, string) {
	i := 0
	limits := len(dictionary) - 1
	for size >= base && i < limits {
		size = size / base
		i++
	}
	return size, dictionary[i]
}


