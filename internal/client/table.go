package client

import (
	"container/list"
	"fmt"
	"github.com/lomoval/otus-golang-project-sysmon/api"
	"strconv"
	"strings"
	"time"
)

const (
	maxLines = 15
)

type table struct {
	header           string
	lines            *list.List
	horizontalBorder string
	columnsIndexes   map[string]int
	lineTemplate     string
}

func newTable(group *api.MetricGroup) *table {
	t := &table{lines: list.New()}

	firstColumnLen := len(time.Now().Format(time.RFC3339))
	if len(group.GetName()) > firstColumnLen {
		firstColumnLen = len(group.GetName())
	}

	t.lines = list.New()
	t.columnsIndexes = make(map[string]int)

	t.header = "|" + strings.Repeat(" ", firstColumnLen-len(group.GetName())) + group.GetName()

	t.lineTemplate = "|%" + strconv.Itoa(firstColumnLen) + "s"
	for i, m := range group.GetMetrics() {
		name := m.GetName()
		columnMinLen := 7
		if len(name) > columnMinLen {
			columnMinLen = len(name)
		}
		t.header += "|" + strings.Repeat(" ", columnMinLen-len(name)) + name

		t.lineTemplate += "|%" + strconv.Itoa(columnMinLen) + ".2f"
		t.columnsIndexes[m.GetName()] = i
	}
	t.header += "|\n"
	t.lineTemplate += "|\n"
	t.horizontalBorder = strings.Repeat("-", len(t.header)-1) + "\n"

	return t
}

func (t *table) buildLine(values []interface{}) string {
	return fmt.Sprintf(t.lineTemplate, values...)
}

func (t *table) addLine(line string) {
	t.lines.PushBack(line)
	if t.lines.Len() > maxLines {
		t.lines.Remove(t.lines.Front())
	}
}

func (t *table) height() int {
	return t.lines.Len() + 3 // borders + header
}

func (t *table) print() {
	fmt.Print(t.horizontalBorder)
	fmt.Print(t.header)
	for i := t.lines.Front(); i != nil; i = i.Next() {
		fmt.Print(i.Value)
	}
	fmt.Print(t.horizontalBorder)
}
