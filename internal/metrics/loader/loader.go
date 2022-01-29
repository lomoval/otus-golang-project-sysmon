package metricloader

import (
	"errors"
	"fmt"
	metric "github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics/cpu"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics/loadaverage"
)

type Config struct {
	IgnoreUnavailable bool
	Collect           Metric
}

type Metric struct {
	Cpu         *bool
	LoadAverage *bool
}

var ErrCollectorNotAvailable = errors.New("collector is not available")

func Load(config Config) ([]metric.Collector, error) {
	var collectors []metric.Collector
	var err error
	if config.Collect.Cpu == nil || *config.Collect.Cpu {
		collectors, err = appendCollector(collectors, cpu.Collector{}, config.IgnoreUnavailable)
		if err != nil {
			return nil, err
		}
	}
	if config.Collect.LoadAverage == nil || *config.Collect.LoadAverage {
		collectors, err = appendCollector(collectors, loadaverage.Collector{}, config.IgnoreUnavailable)
		if err != nil {
			return nil, err
		}
	}
	return collectors, nil
}

func appendCollector(collectors []metric.Collector, collector metric.Collector, ignoreNotAvailable bool) ([]metric.Collector, error) {
	if !collector.Available() {
		if !ignoreNotAvailable {
			return collectors, fmt.Errorf("failed load collector '%s': %w", collector.Name(), ErrCollectorNotAvailable)
		}
		return collectors, nil
	}
	return append(collectors, collector), nil
}
