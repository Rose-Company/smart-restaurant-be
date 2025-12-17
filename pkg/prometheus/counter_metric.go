package prometheus

import (
	"context"
	prm "github.com/prometheus/client_golang/prometheus"
)

type RequestCounterMetricInterface interface {
	Get(ctx context.Context) *prm.CounterVec
	Add(ctx context.Context, methodName string, method string, status string, num float64)
}

type requestCounterMetric struct {
	counterVec *prm.CounterVec
}

func NewRequestCounterMetric(name string) *requestCounterMetric {
	opts := prm.CounterOpts{
		Name: name,
	}
	ct := prm.NewCounterVec(opts, []string{"method_name", "method", "status"})
	return &requestCounterMetric{counterVec: ct}
}

func (c *requestCounterMetric) Get(ctx context.Context) *prm.CounterVec {
	return c.counterVec
}

func (c *requestCounterMetric) Add(ctx context.Context, methodName string, method string, status string, num float64) {
	c.counterVec.WithLabelValues(methodName, method, status).Add(num)
}
