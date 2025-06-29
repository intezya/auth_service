package tracer

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Warnf(template string, args ...interface{})
}

type Config struct {
	Endpoint           string `env:"TRACE_COLLECTOR_ENDPOINT" env-default:"localhost:4317"`
	ServiceName        string `env:"SERVICE_NAME" env-default:"auth"`
	ServiceVersion     string `env:"SERVICE_VERSION" env-required:"true"`
	ServiceEnvironment string `env:"ENV" env-default:"dev"`
}

var (
	globalTracer   trace.Tracer
	tracerProvider *sdktrace.TracerProvider
	once           sync.Once
	initErr        error
)

func Init(config Config, logger Logger) error {
	once.Do(
		func() {
			initErr = initTracer(config, logger)
		},
	)
	return initErr
}

func InitOrIgnore(config Config, logger Logger) {
	_ = Init(config, logger)
}

func initTracer(config Config, logger Logger) error {
	ctx := context.Background()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(config.Endpoint),
		otlptracegrpc.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(config.ServiceEnvironment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	globalTracer = otel.Tracer(config.ServiceName)

	logger.Infoln("Tracer initialized successfully:", config.Endpoint)
	return nil
}

func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	if globalTracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}
	return globalTracer.Start(ctx, spanName)
}

func Shutdown(ctx context.Context) error {
	if tracerProvider == nil {
		return nil
	}
	return tracerProvider.Shutdown(ctx)
}

func GetTracer() trace.Tracer {
	return globalTracer
}

func IsInitialized() bool {
	return globalTracer != nil
}
