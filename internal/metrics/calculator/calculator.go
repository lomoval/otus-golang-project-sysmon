package metriccalc

import (
	"container/list"
	"context"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Calculator struct {
	collectors       []metric.Collector
	calcsByGroup     map[string]*calculators
	intervals        map[time.Duration]int
	maxInterval      time.Duration
	calcId           int
	mutex            sync.Mutex
	cancelCollectors context.CancelFunc
}

type calculators struct {
	calculators map[int]*groupCalc
	m           sync.Mutex
}

func (cd *calculators) addCalc(key int, calc *groupCalc) {
	cd.m.Lock()
	defer cd.m.Unlock()
	cd.calculators[key] = calc
}

func (cd *calculators) removeCalc(key int) {
	cd.m.Lock()
	defer cd.m.Unlock()
	delete(cd.calculators, key)
}

func NewCalculator(collectors []metric.Collector) *Calculator {
	return &Calculator{collectors: collectors, intervals: make(map[time.Duration]int)}
}

func (c *Calculator) createCalcs(avgInterval time.Duration) {
	if len(c.intervals) == 0 {
		c.calcId = 0
		c.calcsByGroup = make(map[string]*calculators)
		for _, collector := range c.collectors {
			c.calcsByGroup[collector.GroupName()] = &calculators{calculators: map[int]*groupCalc{}}
			c.calcsByGroup[collector.GroupName()].addCalc(c.calcId, NewGroupCalc(collector.GroupName(), avgInterval))
		}
		var collectorsCtx context.Context
		collectorsCtx, c.cancelCollectors = context.WithCancel(context.Background())
		c.startCollectors(collectorsCtx, c.collectors)
		return
	}

	for s, cd := range c.calcsByGroup {
		cd.addCalc(c.calcId, NewGroupCalc(s, avgInterval))
	}
}

func (c *Calculator) Start(
	ctx context.Context,
	notifyTimeout time.Duration,
	avgInterval time.Duration,
) <-chan []metric.Group {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.calcId++
	c.createCalcs(avgInterval)

	c.intervals[avgInterval]++
	if avgInterval > c.maxInterval {
		c.maxInterval = avgInterval
	}

	go func(calcId int) {
		<-ctx.Done()
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.removeCalcs(calcId)
		c.intervals[avgInterval]--
		if c.intervals[avgInterval] == 0 {
			delete(c.intervals, avgInterval)
		}
		if avgInterval == c.maxInterval {
			c.maxInterval = 0
			for duration := range c.intervals {
				if duration > c.maxInterval {
					c.maxInterval = duration
				}
			}
		}
		if len(c.intervals) == 0 {
			c.cancelCollectors()
		}
	}(c.calcId)

	calcGroupChan := make(chan []metric.Group)

	tmpNotifyInterval := notifyTimeout
	notifyTicker := time.NewTicker(notifyTimeout)
	if notifyTimeout < avgInterval {
		notifyTimeout = avgInterval
		notifyTicker.Reset(notifyTimeout)
	}
	log.Debugf("notify %d", tmpNotifyInterval)

	go func(calcId int) {
		defer close(calcGroupChan)
		defer notifyTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Debug("notifier goroutine done")
				return
			case <-notifyTicker.C:
				if notifyTimeout != tmpNotifyInterval {
					notifyTimeout = tmpNotifyInterval
					notifyTicker.Reset(notifyTimeout)
				}
				metrics := make([]metric.Group, 0, len(c.calcsByGroup))
				for _, calc := range c.calcsByGroup {
					metrics = append(metrics, calc.calculators[calcId].Average())
				}
				calcGroupChan <- metrics
			}
		}
	}(c.calcId)

	return calcGroupChan
}

func (c *Calculator) startCollectors(ctx context.Context, collectors []metric.Collector) {
	for _, collector := range collectors {
		go func(collector metric.Collector, cd *calculators) {
			metrics := list.New()
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					log.Debug("collectors goroutine done")
					return
				default:
					metricGroup, err := collector.GetMetrics()
					switch {
					case err != nil:
						log.Errorf("failed to get metrics: %s", err.Error())
					default:
						elem := metrics.PushBack(metricGroup)
						cd.m.Lock()
						for _, calc := range cd.calculators {
							calc.Add(elem)
						}
						cd.m.Unlock()

						maxInterval := c.getMaxInterval()
						for e := metrics.Front(); e != nil; {
							g := e.Value.(metric.Group)
							if time.Since(g.Time) <= maxInterval {
								break
							}
							tmp := e.Next()
							metrics.Remove(e)
							e = tmp
							//							log.Debugf("removed old element (count: %d, maxInterval: %v, time: %v)", metrics.Len(), maxInterval, g.Time)
						}
					}

					select {
					case <-ctx.Done():
						log.Debug("calc goroutine done")
						return
					case <-ticker.C:
					}
				}
			}
		}(collector, c.calcsByGroup[collector.GroupName()])
	}
}

func (c *Calculator) getMaxInterval() time.Duration {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.maxInterval
}

func (c *Calculator) removeCalcs(key int) {
	for _, cd := range c.calcsByGroup {
		cd.removeCalc(key)
	}
}
