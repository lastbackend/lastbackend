package utils

import (
	"fmt"
	"github.com/satori/go.uuid"
	"math"
)

func GetUUIDV4() string {
	return uuid.NewV4().String()
}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func FilterQuery(filters map[string]string) string {

	var filter string = fmt.Sprint(` WHERE `)

	for k, v := range filters {
		if v != "" {
			filter += fmt.Sprintf("%s = %s, ", k, v)
		}
	}

	if string(filter[len(filter)-2:]) == ", " {
		filter = filter[:len(filter)-2]
	}

	if filter == ` WHERE ` {
		filter = ""
	}

	return filter
}
