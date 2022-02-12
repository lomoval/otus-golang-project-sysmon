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
	maxLines          = 15
	valueColumnMinLen = 7
)

type table struct {
	header           string
	lines            *list.List
	horizontalBorder string
	columnsIndexes   map[string]int
	lineTemplate     string
}

func newTable(group *api.MetricGroup) (*table, error) {
	t := &table{lines: list.New()}

	firstColumnLen := len(time.Now().Format(time.RFC3339))
	if len(group.GetName()) > firstColumnLen {
		firstColumnLen = len(group.GetName())
	}

	t.lines = list.New()
	t.columnsIndexes = make(map[string]int)

	header := strings.Builder{}
	if _, err := header.WriteString("|" + strings.Repeat(" ", firstColumnLen-len(group.GetName())) + group.GetName()); err != nil {
		return nil, err
	}
	lineTemplate := strings.Builder{}
	if _, err := lineTemplate.WriteString("|%" + strconv.Itoa(firstColumnLen) + "s"); err != nil {
		return nil, err
	}

	for i, m := range group.GetMetrics() {
		name := m.GetName()
		columnMinLen := valueColumnMinLen
		if len(name) > columnMinLen {
			columnMinLen = len(name)
		}
		if _, err := header.WriteString("|" + strings.Repeat(" ", columnMinLen-len(name)) + name); err != nil {
			return nil, err
		}

		if _, err := lineTemplate.WriteString("|%" + strconv.Itoa(columnMinLen) + ".2f"); err != nil {
			return nil, err
		}

		t.columnsIndexes[m.GetName()] = i
	}
	if _, err := header.WriteString("|\n"); err != nil {
		return nil, err
	}
	if _, err := lineTemplate.WriteString("|\n"); err != nil {
		return nil, err
	}

	t.header = header.String()
	t.lineTemplate = lineTemplate.String()
	t.horizontalBorder = strings.Repeat("-", len(t.header)-1) + "\n"

	return t, nil
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
