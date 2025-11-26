package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "path", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path", "status_code"},
	)

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"service", "method", "path"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"service", "method", "path", "status_code"},
	)

	// Error metrics
	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"service", "method", "path", "status_code", "error_type"},
	)

	// Business logic metrics
	businessErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "business_errors_total",
			Help: "Total number of business logic errors",
		},
		[]string{"service", "error_type", "operation"},
	)

	// Database metrics
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "operation", "table"},
	)

	dbErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"service", "operation", "table", "error_type"},
	)

	// Cache metrics
	cacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"service", "cache_type"},
	)

	cacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"service", "cache_type"},
	)

	cacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"service", "operation", "cache_type"},
	)
)

// Service name constants
const (
	ServiceHTTP   = "http"
	ServiceAuth   = "auth"
	ServiceSearch = "search"
)

// HTTPMiddleware creates middleware for HTTP metrics
func HTTPMiddleware(serviceName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Record request size
		if r.ContentLength > 0 {
			httpRequestSize.WithLabelValues(
				serviceName,
				r.Method,
				r.URL.Path,
			).Observe(float64(r.ContentLength))
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(rw.statusCode)

		// Record metrics
		httpRequestsTotal.WithLabelValues(
			serviceName,
			r.Method,
			r.URL.Path,
			statusCode,
		).Inc()

		httpRequestDuration.WithLabelValues(
			serviceName,
			r.Method,
			r.URL.Path,
			statusCode,
		).Observe(duration)

		httpResponseSize.WithLabelValues(
			serviceName,
			r.Method,
			r.URL.Path,
			statusCode,
		).Observe(float64(rw.bytes))

		// Record error if status code indicates error
		if rw.statusCode >= 400 {
			errorType := getErrorType(rw.statusCode)
			httpErrorsTotal.WithLabelValues(
				serviceName,
				r.Method,
				r.URL.Path,
				statusCode,
				errorType,
			).Inc()
		}
	})
}

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func getErrorType(statusCode int) string {
	switch {
	case statusCode >= 400 && statusCode < 500:
		return "client_error"
	case statusCode >= 500:
		return "server_error"
	default:
		return "unknown"
	}
}

// RecordBusinessError records business logic errors
func RecordBusinessError(serviceName, errorType, operation string) {
	businessErrorsTotal.WithLabelValues(serviceName, errorType, operation).Inc()
}

// RecordDBQuery records database query metrics
func RecordDBQuery(serviceName, operation, table string, duration time.Duration) {
	dbQueryDuration.WithLabelValues(serviceName, operation, table).Observe(duration.Seconds())
}

// RecordDBError records database errors
func RecordDBError(serviceName, operation, table, errorType string) {
	dbErrorsTotal.WithLabelValues(serviceName, operation, table, errorType).Inc()
}

// RecordCacheHit records cache hits
func RecordCacheHit(serviceName, cacheType string) {
	cacheHitsTotal.WithLabelValues(serviceName, cacheType).Inc()
}

// RecordCacheMiss records cache misses
func RecordCacheMiss(serviceName, cacheType string) {
	cacheMissesTotal.WithLabelValues(serviceName, cacheType).Inc()
}

// RecordCacheOperation records cache operations
func RecordCacheOperation(serviceName, operation, cacheType string) {
	cacheOperationsTotal.WithLabelValues(serviceName, operation, cacheType).Inc()
}
