package mon

import (
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"code", "method", "path"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency distribution",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 3, 5, 10},
		},
		[]string{"method", "path"},
	)

	processResourceAlloc = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_alloc_mem_MB",
			Help: "Process allocated memory",
		})

	processResourceTotalAlloc = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_total_alloc_mem_MB",
			Help: "Process allocated total memory",
		})

	processResourceGoroutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "process_total_goroutines_NUM",
			Help: "Process allocated total goroutines",
		})
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		processResourceAlloc,
		processResourceTotalAlloc,
		processResourceGoroutines,
	)
}

func PromMonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.FullPath()

		c.Next()

		duration := time.Since(start).Seconds()
		statusStr := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			statusStr,
			c.Request.Method,
			path,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}

func processMon() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	alloc := m.Alloc / (1024 * 1024)
	totalAlloc := m.TotalAlloc / (1024 * 1024)
	goroutines := runtime.NumGoroutine()

	processResourceAlloc.Set(float64(alloc))
	processResourceTotalAlloc.Set(float64(totalAlloc))
	processResourceGoroutines.Set(float64(goroutines)) 
}

func StartPromMetricsUpdate(interval time.Duration) {
	go func() {
		for {
			processMon()
			time.Sleep(interval)
		}
	}()
}
