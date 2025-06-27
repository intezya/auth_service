package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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

func Init(config Config, logger Logger) func() {
	ctx := context.Background()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(config.Endpoint),
		otlptracegrpc.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		logger.Warnf("error creating OTLP trace exporter: %v", err)
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
		logger.Warnf("error creating resource: %v", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	logger.Infoln("Tracer initialized successfully:", config.Endpoint)

	return func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logger.Warnf("Error shutting down tracer provider: %v", err)
		}
	}
}
