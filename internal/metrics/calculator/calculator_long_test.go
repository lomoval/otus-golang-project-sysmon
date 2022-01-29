// +build longtest

package metriccalc

import (
	"context"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testMetricsCollector struct {
	counter int
	sum     float64
	name    string
}

func (c testMetricsCollector) Name() string {
	return "test"
}

func (c testMetricsCollector) GroupName() string {
	return "testGroup"
}

func (tm *testMetricsCollector) Available() bool {
	return true
}

func (tm *testMetricsCollector) GetMetrics() (metric.Group, error) {
	tm.counter++
	tm.sum += float64(tm.counter)
	return metric.Group{
			Name:    tm.GroupName(),
			Metrics: []metric.Metric{{Time: time.Now(), Name: tm.name, Value: float64(tm.counter)}},
		},
		nil
}

func TestMetricExecutor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	testMetric := &testMetricsCollector{name: "test"}
	ch := Start(
		ctx,
		[]metric.Collector{testMetric},
		time.Millisecond*2100,
		time.Millisecond*2100,
	)

	m := <-ch
	require.NotEmpty(t, m)
	require.True(t, (2.0+3.0)/2.0 == m[0].Metrics[0].Value || (1.0+2.0+3.0)/3.0 == m[0].Metrics[0].Value)
}

func TestMetricExecutorFirstBiggerNotifyInterval(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	testMetric := &testMetricsCollector{name: "test"}
	ch := Start(
		ctx,
		[]metric.Collector{testMetric},
		time.Millisecond*500,
		time.Millisecond*1100,
	)

	m := <-ch
	require.NotEmpty(t, m)
	require.Equal(t, testMetric.counter, 2)
	require.GreaterOrEqual(t, m[0].Metrics[0].Value, 1.5)
}

func TestMetricExecutorSeveralMetrics(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	testMetric1 := &testMetricsCollector{name: "test1"}
	testMetric2 := &testMetricsCollector{name: "test2"}
	testMetric3 := &testMetricsCollector{name: "test3"}
	testMetrics := map[string]*testMetricsCollector{
		testMetric1.name: testMetric1,
		testMetric2.name: testMetric2,
		testMetric3.name: testMetric3,
	}
	ch := Start(
		ctx,
		[]metric.Collector{testMetric1, testMetric2, testMetric3},
		time.Millisecond*1500,
		time.Millisecond*2500,
	)

	groups := <-ch
	require.NotEmpty(t, groups)
	for _, m := range groups[0].Metrics {
		testMetric := testMetrics[m.Name]
		require.Equal(t, testMetric.counter, 3)
		require.True(t, (2.0+3.0)/2.0 == m.Value || (1.0+2.0+3.0)/3.0 == m.Value)
	}

}
