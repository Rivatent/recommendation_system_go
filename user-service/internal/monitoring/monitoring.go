package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RequestsTotal)
}

func CollectMetrics(start time.Time, c *gin.Context) {
	method := c.Request.Method
	elapsed := time.Since(start).Seconds()
	RequestsTotal.WithLabelValues(method).Inc()
	RequestDuration.WithLabelValues(method).Observe(elapsed)
}
