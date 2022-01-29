package metriccalc

import (
	"container/list"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type metricCalc struct {
	SumValue metric.Metric
	Values   *list.List
	Period   time.Duration
}

func newCalc(metricName string, period time.Duration) *metricCalc {
	return &metricCalc{
		SumValue: metric.Metric{Name: metricName},
		Values:   list.New(),
		Period:   period,
	}
}

func (c *metricCalc) removeObsolete() {
	for e := c.Values.Front(); e != nil; e = c.Values.Front() {
		if time.Since(e.Value.(metric.Metric).Time) <= c.Period {
			return
		}
		c.SumValue.Value -= e.Value.(metric.Metric).Value
		c.Values.Remove(e)
	}
}

func (c *metricCalc) Add(metric metric.Metric) {
	c.Values.PushBack(metric)
	c.SumValue.Time = metric.Time
	c.SumValue.Value += metric.Value
	log.Printf("ADD %v %f", c.SumValue, metric.Value)
	c.removeObsolete()

}

func (c *metricCalc) Avg() metric.Metric {
	log.Printf("AVG %v", c.SumValue)
	return metric.Metric{Name: c.SumValue.Name, Time: c.SumValue.Time, Value: c.SumValue.Value / float64(c.Values.Len())}
}

type groupCalc struct {
	groupName    string
	calculators  map[string]*metricCalc
	calcInterval time.Duration
	mutex        sync.Mutex
}

func NewGroupCalc(group string, calcInterval time.Duration) *groupCalc {
	return &groupCalc{
		groupName:    group,
		calcInterval: calcInterval,
		calculators:  map[string]*metricCalc{},
	}
}

func (mc *groupCalc) Add(metric metric.Metric) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	calc, ok := mc.calculators[metric.Name]
	if !ok {
		calc = newCalc(metric.Name, mc.calcInterval)
		mc.calculators[metric.Name] = calc
	}
	calc.Add(metric)
}

func (mc *groupCalc) Average() metric.Group {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	g := metric.Group{Name: mc.groupName, Metrics: make([]metric.Metric, 0, len(mc.calculators))}
	for _, calc := range mc.calculators {
		g.Metrics = append(g.Metrics, calc.Avg())
	}
	return g
}
