package metric

import (
	"errors"
	"time"
)

var (
	ErrCollectorNotAvailable = errors.New("metrics not available for collection")
	ErrParseFailed           = errors.New("failed to parse collector output")
)

type Config struct {
	Cpu *bool
}

type Metric struct {
	Time  time.Time
	Name  string
	Value float64
}

type Group struct {
	Time    time.Time
	Name    string
	Metrics []Metric
}

type Collector interface {
	Name() string
	GroupName() string
	Available() bool
	GetMetrics() (Group, error)
}

type UnavailableCollector struct {
}

func (c UnavailableCollector) Name() string {
	return ""
}

func (c UnavailableCollector) GroupName() string {
	return ""
}

func (c UnavailableCollector) Available() bool {
	return false
}

func (c UnavailableCollector) GetMetrics() (Group, error) {
	return Group{}, ErrCollectorNotAvailable
}
