package go_trace_lib

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/trace"

	logging "github.com/lordtor/go-logging"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	Log = logging.Log
)

func AgentEndpoint(config ProviderConfig) (*jaeger.Exporter, error) {
	Log.Info("Jaeger provaider type set as AgentEndpoint")
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(config.JaegerHost),
		jaeger.WithAgentPort(config.JaegerPort),
	))
	// jaeger.WithAgentEndpoint(
	// jaeger.AgentEndpointOption(jaeger.WithAgentHost(config.JaegerHost)), jaeger.AgentEndpointOption(jaeger.WithAgentPort(config.JaegerPort))))
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return exp, nil
}
func CollectorEndpoint(config ProviderConfig) (*jaeger.Exporter, error) {
	Log.Info("Jaeger provaider type set as CollectorEndpoint")
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerEndpoint)))
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return exp, nil
}

// New возвращает новый тип "Поставщик`. Он использует экспортер Jaeger и глобально устанавливает поставщика трассировщика,
// а также глобальный трассировщик для пролетов.
func NewProvider(config ProviderConfig) (Provider, error) {
	Log.Info("Create new jaeger provider")
	if config.Disabled {
		return Provider{provider: trace.NewNoopTracerProvider()}, nil
	}
	exp := &jaeger.Exporter{}
	err := errors.New("Error sets jaeger exporter")
	if config.JaegerEndpoint != "" {
		exp, err = CollectorEndpoint(config)
	} else if config.JaegerHost != "" && config.JaegerPort != "" {
		exp, err = AgentEndpoint(config)
	}
	if err != nil {
		Log.Error("Error set jaeger provider")
		return Provider{}, err
	}
	Log.Info("Setting jaeger provider")
	prv := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(sdkresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(config.Environment),
		)),
	)
	Log.Debugf("Jaeger provider sets ServiceName: %s, ServiceVersion: %s, Environment: %s", config.ServiceName, config.ServiceVersion, config.Environment)
	otel.SetTracerProvider(prv)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return Provider{provider: prv}, nil
}

// Закрыть поставщика трассировщика только в том случае, если это не была версия "без операций".
func (p Provider) Close(ctx context.Context) error {
	if prv, ok := p.provider.(*sdktrace.TracerProvider); ok {
		return prv.Shutdown(ctx)
	}
	return nil
}
