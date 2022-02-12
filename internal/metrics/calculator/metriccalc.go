package metriccalc

import (
	metric "github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	log "github.com/sirupsen/logrus"
	"time"
)

type metricCalc struct {
	Period   time.Duration
	SumValue metric.Metric
	count    int
}

func newMetricCalc(metricName string, period time.Duration) *metricCalc {
	return &metricCalc{
		SumValue: metric.Metric{Name: metricName},
		Period:   period,
	}
}

func (c *metricCalc) Add(metric metric.Metric) {
	c.SumValue.Time = metric.Time
	c.SumValue.Value += metric.Value
	c.count++
	log.Debugf("add %v %f", c.SumValue, metric.Value)
}

func (c *metricCalc) Subtract(metric metric.Metric) {
	c.count--
	c.SumValue.Value -= metric.Value
}

func (c *metricCalc) Average() metric.Metric {
	log.Debugf("avg count %d %v", c.count, c.SumValue)
	return metric.Metric{Name: c.SumValue.Name, Time: c.SumValue.Time, Value: c.SumValue.Value / float64(c.count)}
}
