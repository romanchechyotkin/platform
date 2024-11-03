package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"

	"github.com/TakeAway-Inc/platform/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func NewClient(log *logger.Logger, cfg *ClientConfig, clientTLSConfig *tls.Config) (*grpc.ClientConn, error) {
	clientCreds := credentials.NewTLS(clientTLSConfig)

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(clientCreds),
		grpc.WithUnaryInterceptor(contextInterceptor(log)),
		grpc.WithUnaryInterceptor(clientLogInterceptor(log)),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func contextInterceptor(log *logger.Logger) grpc.UnaryClientInterceptor {
	log = log.With(slog.String("component", "grpc client"))

	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ginCtx, ok := ctx.(*gin.Context)
		if !ok {
			err := errors.New("failed to get request context")
			log.Error("failed to get request context", err, slog.String("method", method))
			return err
		}

		reqCtx := ginCtx.Request.Context()

		err := invoker(reqCtx, method, req, reply, cc, opts...)
		if err != nil {
			log.Error("failed to invoke", err, slog.String("method", method), slog.String("interceptor", "context interceptor"))
			return err
		}

		return nil
	}
}

func clientLogInterceptor(log *logger.Logger) grpc.UnaryClientInterceptor {
	log = log.With(slog.String("component", "grpc client"))
	tracer := otel.GetTracerProvider().Tracer("grpc client")

	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Info("[Restaurant Service Interceptor]", slog.String("method", method))

		ctx, span := tracer.Start(ctx, method)
		defer span.End()

		traceId := fmt.Sprintf("%s", span.SpanContext().TraceID())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceId)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			log.Error("failed to invoke", err, slog.String("method", method), slog.String("x-trace-id", traceId), slog.String("interceptor", "log/trace interceptor"))
			return err
		}

		log.Info("got reply", slog.Any("reply", reply))

		return nil
	}
}
