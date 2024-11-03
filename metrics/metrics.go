package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(
		HTTPRequestsCount,
		GRPCServerRequestsCount,
	)
}

var HTTPRequestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests",
}, []string{"uri", "method"})

var GRPCServerRequestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "grpc_requests_total",
	Help: "Total number of gRPC requests",
}, []string{"method"})
