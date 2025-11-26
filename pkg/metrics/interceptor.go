package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Add these to your pkg/metrics/metrics.go

var (
	// gRPC server metrics
	grpcServerRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_total",
			Help: "Total number of gRPC server requests",
		},
		[]string{"service", "method", "code"},
	)

	grpcServerRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_request_duration_seconds",
			Help:    "gRPC server request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "code"},
	)

	grpcServerRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_server_requests_in_flight",
			Help: "Current number of gRPC server requests in flight",
		},
		[]string{"service", "method"},
	)

	// gRPC client metrics (for HTTP service calling auth/search)
	grpcClientRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_client_requests_total",
			Help: "Total number of gRPC client requests",
		},
		[]string{"caller_service", "target_service", "method", "code"},
	)

	grpcClientRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_client_request_duration_seconds",
			Help:    "gRPC client request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"caller_service", "target_service", "method", "code"},
	)

	grpcClientRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_client_requests_in_flight",
			Help: "Current number of gRPC client requests in flight",
		},
		[]string{"caller_service", "target_service", "method"},
	)
)

// gRPC server interceptor for auth and search services
func GRPCServerMetricsInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Extract method name
		method := info.FullMethod

		// Record in-flight request
		grpcServerRequestsInFlight.WithLabelValues(serviceName, method).Inc()
		defer grpcServerRequestsInFlight.WithLabelValues(serviceName, method).Dec()

		// Call the handler
		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()

		// Convert error to gRPC code
		code := status.Code(err)
		codeStr := code.String()

		// Record metrics
		grpcServerRequestsTotal.WithLabelValues(serviceName, method, codeStr).Inc()
		grpcServerRequestDuration.WithLabelValues(serviceName, method, codeStr).Observe(duration)

		return resp, err
	}
}

// gRPC client interceptor for HTTP service calling auth/search
func GRPCClientMetricsInterceptor(callerService, targetService string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		// Record in-flight request
		grpcClientRequestsInFlight.WithLabelValues(callerService, targetService, method).Inc()
		defer grpcClientRequestsInFlight.WithLabelValues(callerService, targetService, method).Dec()

		// Make the call
		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start).Seconds()

		// Convert error to gRPC code
		code := status.Code(err)
		codeStr := code.String()

		// Record metrics
		grpcClientRequestsTotal.WithLabelValues(callerService, targetService, method, codeStr).Inc()
		grpcClientRequestDuration.WithLabelValues(callerService, targetService, method, codeStr).Observe(duration)

		return err
	}
}

// Record gRPC specific errors
func RecordGRPCError(serviceType, serviceName, method, code, errorType string) {
	if serviceType == "server" {
		// You can add specific gRPC error tracking here
		businessErrorsTotal.WithLabelValues(serviceName, errorType, method).Inc()
	}
}
