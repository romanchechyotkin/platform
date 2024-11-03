package tracing

import (
	"os"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/trace"
)

const defaultEndpoint = "http://localhost:14268/api/traces"

// NewJaegerExporter creates new jaeger exporter
func NewJaegerExporter() (trace.SpanExporter, error) {
	var endpoint string

	if val := os.Getenv("JAEGER_ENDPOINT"); val == "" {
		endpoint = defaultEndpoint
	} else {
		endpoint = val
	}

	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
}
