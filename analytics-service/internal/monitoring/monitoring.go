package monitoring

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	// RequestsTotal хранит общее количество HTTP запросов, сгруппированных по методу.
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method"},
	)
	// RequestDuration хранит распределение временных задержек HTTP запросов, сгруппированных по методу.
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.LinearBuckets(0, 10, 20),
		},
		[]string{"method"},
	)
)

// InitMetrics инициализирует метрики для мониторинга запросов.
// В этой функции регистрируются метрики RequestDuration и RequestsTotal
// в системе мониторинга Prometheus.
//
// Примечание: должна быть вызвана один раз при инициализации приложения.
func InitMetrics() {
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RequestsTotal)
}

// CollectMetrics собирает метрики для каждого запроса.
// Он рассчитывает продолжительность запроса и инкрементирует
// счетчик RequestsTotal
//
// Параметры:
//   - start: время начала запроса, используется для вычисления продолжительности.
//   - c: указатель на контекст Gin, содержащий информацию о текущем запросе.
//
// Примечание: эта функция должна вызываться в обработчике для сбора метрик
// после завершения обработки запроса.
func CollectMetrics(start time.Time, c *gin.Context) {
	method := c.Request.Method
	elapsed := time.Since(start).Seconds()
	RequestsTotal.WithLabelValues(method).Inc()
	RequestDuration.WithLabelValues(method).Observe(elapsed)
}
