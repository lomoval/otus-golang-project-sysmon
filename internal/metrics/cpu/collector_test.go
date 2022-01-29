package cpu_test

import (
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics/cpu"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCpuMetric(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		var a []int
		var i int
		ch <- struct{}{}
		for {
			i++
			a = append(a, i)
			_ = a
			select {
			case <-ch:
				return
			default:
			}
		}
	}()
	<-ch
	metric, err := cpu.Collector{}.GetMetrics()
	ch <- struct{}{}
	require.NoError(t, err)
	require.Greater(t, metric.Metrics[0].Value+metric.Metrics[1].Value+metric.Metrics[2].Value, float64(0))
	<-ch
}
