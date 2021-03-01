package wrsvc

import (
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type wrSVCMetrics struct {
	requestCount   metrics.Counter
	requestSize    metrics.Counter
	requestLatency metrics.Histogram
}

func makeWRSVCMetrics() wrSVCMetrics {
	fieldKeys := []string{"method", "error"}

	return wrSVCMetrics{
		requestCount: kitprometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "dyslav",
			Subsystem: "wr_svc",
			Name:      "request_count",
			Help:      "Total number of requests received.",
		}, fieldKeys),
		requestSize: kitprometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "dyslav",
			Subsystem: "wr_svc",
			Name:      "request_size",
			Help:      "Size of requests recieved",
		}, fieldKeys),
		requestLatency: kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "dyslav",
			Subsystem: "wr_svc",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
	}
}
