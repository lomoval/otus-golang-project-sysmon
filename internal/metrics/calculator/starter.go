package metriccalc

import (
	"context"
	"github.com/lomoval/otus-golang-project-sysmon/internal/metrics"
	log "github.com/sirupsen/logrus"
	"time"
)

func Start(
	ctx context.Context,
	collectors []metric.Collector,
	notifyTimeout time.Duration,
	avgInterval time.Duration,
) <-chan []metric.Group {
	calcGroupChan := make(chan []metric.Group)

	calculators := make(map[string]*groupCalc)
	for _, collector := range collectors {
		calculators[collector.GroupName()] = NewGroupCalc(collector.GroupName(), avgInterval)
	}

	tmpNotifyInterval := notifyTimeout
	if notifyTimeout < avgInterval {
		notifyTimeout = avgInterval
	}

	go func() {
		defer close(calcGroupChan)
		for {
			select {
			case <-ctx.Done():
				log.Debug("notifier goroutine done")
				return
			case <-time.After(notifyTimeout):
				notifyTimeout = tmpNotifyInterval
				metrics := make([]metric.Group, 0, len(calculators))
				for _, calc := range calculators {
					metrics = append(metrics, calc.Average())
				}
				calcGroupChan <- metrics
			}
		}
	}()

	for _, collector := range collectors {
		go func(collector metric.Collector) {
			for {
				select {
				case <-ctx.Done():
					log.Debug("calc goroutine done")
					return
				default:
					start := time.Now()
					metricGroup, err := collector.GetMetrics()

					switch {
					case err != nil:
						log.Errorf("failed to get metrics: %s", err.Error())
					default:
						for _, m := range metricGroup.Metrics {
							calculators[metricGroup.Name].Add(m)
						}
					}

					select {
					case <-ctx.Done():
						log.Debug("calc goroutine done")
						return
					case <-time.After(time.Second - time.Since(start)):
					}
				}
			}
		}(collector)
	}

	return calcGroupChan
}
