package loadaverage

import (
	"fmt"
	"github.com/lomoval/otus-golang-project-sysmon/internal/executor"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"strconv"
	"strings"
	"time"
)

const (
	command = "uptime"
)

const (
	minute1Column  = 0
	minute5Column  = 1
	minute15Column = 2
)

func (c Collector) Available() bool {
	return true
}

func (c Collector) GetMetrics() (metric.Group, error) {
	output, err := executor.Exec(command, nil)
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
	// Output example:
	// 00:03:51 up 1 min,  1 user,  load average: 0.62, 0.30, 0.11
	var err error
	parts := strings.Split(output, "load average:")
	if len(parts) < 2 {
		return fmt.Errorf("incorrect output from '%s': %w", command, metric.ErrParseFailed)
	}

	parts = strings.Fields(parts[1])
	if len(parts) < 3 {
		return fmt.Errorf("incorrect output from '%s': %w", command, metric.ErrParseFailed)
	}
	for i, part := range parts {
		parts[i] = strings.Trim(part, ",")
	}

	m.Minute1, err = strconv.ParseFloat(strings.Trim(parts[minute1Column], " "), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", minute1MetricName, parts[minute1Column], metric.ErrParseFailed)
	}
	m.Minute5, err = strconv.ParseFloat(strings.Trim(parts[minute5Column], " "), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", minute5MetricName, parts[minute5Column], metric.ErrParseFailed)
	}
	m.Minute15, err = strconv.ParseFloat(strings.Trim(parts[minute15Column], " "), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", minute15MetricName, parts[minute15Column], metric.ErrParseFailed)
	}
	return nil
}
