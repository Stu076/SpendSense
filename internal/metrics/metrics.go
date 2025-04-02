package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// HTTP Request Counter
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests received",
		},
		[]string{"method", "path"},
	)

	// Request Duration Histogram
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(RequestDuration)
}
