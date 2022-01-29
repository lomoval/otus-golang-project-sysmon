package loadaverage

import (
	metric "github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"time"
)

const (
	collectorName      = "Average Load"
	groupName          = "LoadAverage"
	minute1MetricName  = "1-min"
	minute5MetricName  = "5-min"
	minute15MetricName = "15-min"
)

//nolint:deadcode,unused // ignore when collector is not available
type metricData struct {
	Minute1  float64
	Minute5  float64
	Minute15 float64
}

/* eslint-enable no-unused-vars */

type Collector struct {
	metric.UnavailableCollector
	// _ metricData // Just for `unused` linter
}

func (c Collector) Name() string {
	return collectorName
}

func (c Collector) GroupName() string {
	return groupName
}

//nolint:deadcode,unused // ignore when collector is not available
func toGroup(t time.Time, m metricData) metric.Group {
	return metric.Group{
		Name: groupName,
		Time: t,
		Metrics: []metric.Metric{
			{
				Time:  t,
				Name:  minute1MetricName,
				Value: m.Minute1,
			},
			{
				Time:  t,
				Name:  minute5MetricName,
				Value: m.Minute5,
			},
			{
				Time:  t,
				Name:  minute15MetricName,
				Value: m.Minute15,
			},
		}}
}
