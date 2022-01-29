package cpu

import (
	"fmt"
	"github.com/lomoval/otus-golang-project-sysmon/internal/executor"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"strconv"
	"strings"
	"time"
)

const command = "top"

var args = []string{"-b", "-n1"}

func (c Collector) Available() bool {
	return true
}

func (c Collector) GetMetrics() (metric.Group, error) {
	output, err := executor.Exec(command, args)
	if err != nil {
		return metric.Group{}, err
	}

	var m metricData
	var t time.Time = time.Now()
	if err := parse(output, &m); err != nil {
		return metric.Group{}, err
	}
	return toGroup(t, m), nil
}

func parse(output string, m *metricData) error {
	// Output examples:
	// top - 22:40:00 up 14:56,  1 user,  load average: 0.05, 0.35, 0.40
	// Tasks: 140 total,   1 running, 139 sleeping,   0 stopped,   0 zombie
	// %Cpu(s):  2.0 us,  4.0 sy,  0.0 ni, 94.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
	//
	// On some OS output has no ',' and has percentage symbols:
	// Mem: 1848560K used, 194504K free, 20064K shrd, 22332K buff, 479880K cached
	// CPU:  15% usr   0% sys   0% nic  85% idle   0% io   0% irq   0% sirq

	var err error
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("incorrect output from '%s': %w", command, metric.ErrParseFailed)
	}

	var elems []string
	const maxLines = 5
	for i := 0; i < maxLines && i < len(lines); i++ {
		if index := strings.Index(lines[i], ":"); index > 0 {
			if strings.Contains(strings.ToLower(lines[i][:index]), "cpu") {
				elems = strings.Fields(lines[i])
				break
			}
		}
	}

	if len(elems) <= 1 {
		return fmt.Errorf("incorrect output from '%s': %w", command, metric.ErrParseFailed)
	}

	for i := 1; i < len(elems); i++ {
		elems[i] = strings.Replace(elems[i], ",", "", 1)
		elems[i] = strings.Replace(elems[i], "%", "", 1)
		switch {
		case elems[i] == "us" || elems[i] == "usr":
			m.UserTime, err = strconv.ParseFloat(elems[i-1], 64)
			if err != nil {
				return fmt.Errorf("failed to parse metric '%s': %w", elems[i-1], metric.ErrParseFailed)
			}
		case elems[i] == "sy" || elems[i] == "sys":
			m.SystemTime, err = strconv.ParseFloat(elems[i-1], 64)
			if err != nil {
				return fmt.Errorf("failed to parse metric '%s': %w", elems[i-1], metric.ErrParseFailed)
			}
		case elems[i] == "id" || elems[i] == "idle":
			m.IdleTime, err = strconv.ParseFloat(elems[i-1], 64)
			if err != nil {
				return fmt.Errorf("failed to parse metric '%s': %w", elems[i-1], metric.ErrParseFailed)
			}
		}
	}

	return nil
}
