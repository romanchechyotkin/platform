package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func New(serviceName string) (trace.Tracer, error) {
	jaegerExporter, err := NewJaegerExporter()
	if err != nil {
		return nil, err
	}

	traceProvider, err := NewTraceProvider(jaegerExporter, serviceName)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(traceProvider)

	return traceProvider.Tracer("main service"), nil
}
