package table

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	_MINWIDTH int  = 0
	_TABWIDTH int  = 0
	_PADDING  int  = 2
	_PADCHAR  byte = ' '
	_FLAGS    uint = tabwriter.TabIndent
	// Indentation columns
	_IC string = "	"
)

func PrintTable(header []string, data [][]string, footer []string) {
	tw := tabwriter.NewWriter(os.Stdout, _MINWIDTH, _TABWIDTH, _PADDING, _PADCHAR, _FLAGS)

	fmt.Fprintln(tw, strings.ToUpper(parse(header)))
	for _, de := range data {
		fmt.Fprintln(tw, parse(de))
	}
	fmt.Fprintln(tw, parse(footer))

	tw.Flush()
}

func parse(values []string) string {
	var retStr string

	for _, v := range values {
		retStr += checkNullStr(v) + _IC
	}
	return retStr
}

func checkNullStr(str string) string {
	if str == "" {
		return "-"
	}
	return str
}
