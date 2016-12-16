package table

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
)

type table struct {
	fields        []string
	rows          []map[string]string
	fieldSizes    map[string]int
	VisibleHeader bool
}

func New(fields []string) *table {
	return &table{
		fields:     fields,
		rows:       make([]map[string]string, 0),
		fieldSizes: make(map[string]int),
	}
}

func (t *table) AddRow(row map[string]interface{}) {

	newRow := make(map[string]string)

	for _, key := range t.fields {
		v := row[key]

		newRow[key] = ""
		if v != nil {
			newRow[key] = fmt.Sprintf("%v", v)
		}
	}

	t.calculate(newRow)

	if len(newRow) > 0 {
		t.rows = append(t.rows, newRow)
	}
}

func PrintHorizontal(m map[string]interface{}) {
	table := New([]string{"Key", "Value"})
	table.VisibleHeader = false
	for key, value := range m {
		row := map[string]interface{}{}
		row["Key"] = strings.Title(key)
		row["Value"] = value
		table.AddRow(row)
	}

	table.Print()
}

func (t *table) Print() {
	if len(t.rows) == 0 {
		return
	}

	t.printHeader()

	for _, r := range t.rows {
		t.printRow(r)
	}
}

func (t *table) printHeader() {
	var s string

	if !t.VisibleHeader {
		return
	}

	for _, name := range t.fields {
		s += t.fieldToString(name, strings.Title(name))
	}

	fmt.Println(s)
}

func (t *table) printRow(row map[string]string) {
	var s string

	for _, name := range t.fields {
		value := row[name]
		s += t.fieldToString(name, value)
	}

	fmt.Println(s)
}

func (t *table) fieldToString(name, value string) string {
	value = fmt.Sprintf(" %s ", value)
	spacesLeft := t.fieldSizes[name] - runewidth.StringWidth(value)

	if spacesLeft > 0 {
		for i := 0; i < spacesLeft; i++ {
			value += " "
		}
	}

	return value
}

func (t *table) calculate(row map[string]string) {
	for _, k := range t.fields {
		if v, ok := row[k]; !ok {
			continue
		} else {
			vlen := runewidth.StringWidth(v)

			// align to field name length
			if klen := runewidth.StringWidth(k); vlen < klen {
				vlen = klen
			}

			vlen += 2 // + 2 spaces

			if t.fieldSizes[k] < vlen {
				t.fieldSizes[k] = vlen
			}
		}
	}
}
