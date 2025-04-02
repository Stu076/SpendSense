package middlewares

import (
	"github.com/Stu076/SpendSense/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// TODO: fix at a later date
func TrackRequest(method, endpoint string, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues(method, endpoint))

		defer timer.ObserveDuration()

		metrics.RequestCounter.WithLabelValues(method, endpoint).Inc()
		handler(c)
	}
}
