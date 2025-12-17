package prometheus

import (
	"context"
	prm "github.com/prometheus/client_golang/prometheus"
)

type RequestDurationMetric interface {
	Get(ctx context.Context) *prm.HistogramVec
	ObserveDuration(ctx context.Context, methodName string, method string) *prm.Timer
}

type requestDurationMetric struct {
	histogramVec *prm.HistogramVec
}

func NewRequestHistogramMetric(name string) *requestDurationMetric {
	opts := prm.HistogramOpts{
		Name:    name,
		Buckets: prm.LinearBuckets(0.01, 0.05, 10),
	}
	ct := prm.NewHistogramVec(opts, []string{"method_name", "method"})
	return &requestDurationMetric{histogramVec: ct}
}

func (c *requestDurationMetric) Get(ctx context.Context) *prm.HistogramVec {
	return c.histogramVec
}

func (c *requestDurationMetric) ObserveDuration(ctx context.Context, methodName string, method string) *prm.Timer {
	timer := prm.NewTimer(prm.ObserverFunc(func(v float64) {
		c.histogramVec.WithLabelValues(methodName, method).Observe(v)
	}))
	return timer
}
