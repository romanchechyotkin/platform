package minio

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/TakeAway-Inc/platform"
)

type Tracer struct {
	tracer trace.Tracer
}

type tracerConfig struct {
	tp    trace.TracerProvider
	attrs []attribute.KeyValue
}

func setTracing() (*Tracer, error) {
	cfg := &tracerConfig{
		tp: otel.GetTracerProvider(),
		attrs: []attribute.KeyValue{
			attribute.Key("service.name").String("minio"),
		},
	}

	return &Tracer{tracer: cfg.tp.Tracer(tracerName)}, nil
}
