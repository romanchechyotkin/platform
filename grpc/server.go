package grpc

import (
	"context"
	"crypto/tls"
	"log/slog"

	"github.com/TakeAway-Inc/platform/logger"
	"github.com/TakeAway-Inc/platform/metrics"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func NewServer(log *logger.Logger, clientTLSConfig *tls.Config) *grpc.Server {
	creds := credentials.NewTLS(clientTLSConfig)

	rpcSrv := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			serverMetricInterceptor,
			serverLogInterceptor(log),
		),
	)

	return rpcSrv
}

func serverLogInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	log = log.With(slog.String("component", "grpc server"))
	tracer := otel.GetTracerProvider().Tracer("grpc server")

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Info("[Storage Service Interceptor]", slog.Any("info", info.FullMethod))

		// Extract TraceID from header
		md, _ := metadata.FromIncomingContext(ctx)
		traceIdString := md["x-trace-id"][0]
		// Convert string to byte array
		traceId, err := trace.TraceIDFromHex(traceIdString)
		if err != nil {
			return nil, err
		}

		// Creating a span context with a predefined trace-id
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceId,
		})

		// Embedding span config into the context
		ctx = trace.ContextWithSpanContext(ctx, spanContext)

		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		m, err := handler(ctx, req)

		log.Info("post proc message", slog.Any("msg", m))

		return m, err
	}
}

func serverMetricInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metrics.GRPCServerRequestsCount.With(map[string]string{
		"method": info.FullMethod,
	}).Inc()

	m, err := handler(ctx, req)

	return m, err
}
