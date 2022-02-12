package metriccalc

import (
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMetricCalcOne(t *testing.T) {
	c := newMetricCalc("test", time.Second)
	c.Add(metric.Metric{Time: time.Now(), Name: "test", Value: 1})

	m := c.Average()

	require.Equal(t, metric.Metric{Time: m.Time, Name: "test", Value: 1}, m)
}

func TestMetricCalcThree(t *testing.T) {
	c := newMetricCalc("test", time.Second*5)
	c.Add(metric.Metric{Time: time.Now(), Name: "test", Value: 1})
	c.Add(metric.Metric{Time: time.Now(), Name: "test", Value: 2})
	c.Add(metric.Metric{Time: time.Now(), Name: "test", Value: 3})

	m := c.Average()

	require.Equal(t, metric.Metric{Time: m.Time, Name: "test", Value: (1 + 2 + 3) / 3}, m)
}
