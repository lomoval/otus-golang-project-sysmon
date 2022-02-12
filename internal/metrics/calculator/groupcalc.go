package metriccalc

import (
	"container/list"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	"sync"
	"time"
)

type groupCalc struct {
	groupName    string
	calculators  map[string]*metricCalc
	calcInterval time.Duration
	mutex        sync.Mutex
	lastElement  *list.Element
}

func NewGroupCalc(group string, calcInterval time.Duration) *groupCalc {
	return &groupCalc{
		groupName:    group,
		calcInterval: calcInterval,
		calculators:  map[string]*metricCalc{},
	}
}

func (mc *groupCalc) Add(element *list.Element) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	group := element.Value.(metric.Group)
	for _, m := range group.Metrics {
		calc, ok := mc.calculators[m.Name]
		if !ok {
			calc = newMetricCalc(m.Name, mc.calcInterval)
			mc.calculators[m.Name] = calc
		}
		calc.Add(m)
		if mc.lastElement == nil {
			mc.lastElement = element
		}
		mc.moveToActualElem()
	}
}

func (mc *groupCalc) Average() metric.Group {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.moveToActualElem()
	g := metric.Group{Name: mc.groupName, Time: time.Now(), Metrics: make([]metric.Metric, 0, len(mc.calculators))}
	for _, calc := range mc.calculators {
		g.Metrics = append(g.Metrics, calc.Average())
	}
	return g
}

func (mc *groupCalc) moveToActualElem() {
	for e := mc.lastElement; e != nil; e = e.Next() {
		for _, m := range e.Value.(metric.Group).Metrics {
			if time.Since(m.Time) <= mc.calcInterval {
				return
			}
			mc.calculators[m.Name].Subtract(m)
			mc.lastElement = e.Next()
			if mc.lastElement == nil {
				return
			}
		}
	}
}
