package stats

import "time"

type Collector interface {
	Increment(metric string)
	Timing(metric string, duration time.Duration)
	Gauge(metric string, value float64)
}