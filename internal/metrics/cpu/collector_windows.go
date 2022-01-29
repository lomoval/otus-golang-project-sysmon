package cpu

import (
	"fmt"
	"github.com/lomoval/otus-golang-project-sysmon/internal/executor"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"strconv"
	"strings"
	"time"
)

const command = "typeperf"

var args = []string{
	`Processor Information(_Total)\% Privileged Time`,
	`Processor Information(_Total)\% User Time`,
	`Processor Information(_Total)\% Idle Time`,
	"-sc", "1",
}

const (
	timestampColumns = iota
	systemTimeColumn
	userTimeColumn
	idleTimeColumn
)

func (c Collector) Available() bool {
	return true
}

func (c Collector) GetMetrics() (metric.Group, error) {
	output, err := executor.Exec(command, args)
	if err != nil {
		return metric.Group{}, err
	}

	var m metricData
	var t time.Time
	if err := parse(output, &t, &m); err != nil {
		return metric.Group{}, err
	}
	return toGroup(time.Now(), m), nil
}

func parse(output string, t *time.Time, m *metricData) error {
	// Output example:
	//
	// "(PDH-CSV 4.0)","\\LAPTOP\Processor Information(_Total)\% User Time","\\LAPTOP\Processor Information(_Total)\% Processor Time","\\LAPTOP-TQVSMR2A\Processor Information(_Total)\% Idle Time"
	// "01/25/2022 21:08:57.559","10.763132","15.816928","84.183072"
	var err error
	parts := strings.Split(output, "\r\n")
	if len(parts) < 3 {
		return fmt.Errorf("incorrect output from '%s': %w", command, metric.ErrParseFailed)
	}

	parts = strings.Split(parts[2], ",")

	*t, err = time.Parse("01/02/2006 15:04:05.000", strings.Trim(parts[timestampColumns], "\""))
	if err != nil {
		return fmt.Errorf("failed to parse timestamp (%s): %w", parts[timestampColumns], metric.ErrParseFailed)
	}
	m.SystemTime, err = strconv.ParseFloat(strings.Trim(parts[systemTimeColumn], "\""), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", systemTimeMetricName, parts[systemTimeColumn], metric.ErrParseFailed)
	}
	m.UserTime, err = strconv.ParseFloat(strings.Trim(parts[userTimeColumn], "\""), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", userTimeMetricName, parts[userTimeColumn], metric.ErrParseFailed)
	}
	m.IdleTime, err = strconv.ParseFloat(strings.ReplaceAll(parts[idleTimeColumn], "\"", ""), 64)
	if err != nil {
		return fmt.Errorf("failed to parse metric '%s' (%s): %w", idleTimeMetricName, parts[idleTimeColumn], metric.ErrParseFailed)
	}
	return nil
}
