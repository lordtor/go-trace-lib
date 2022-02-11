package go_trace_lib

import (
	"go.opentelemetry.io/otel/trace"
)

// ProviderConfig представляет конфигурацию поставщика и используется для создания нового типа "Поставщик".
type ProviderConfig struct {
	JaegerEndpoint string `json:"jaeger_endpoint,omitempty" yaml:"jaeger_endpoint"`
	JaegerHost     string `json:"jaeger_host,omitempty" yaml:"jaeger_host"`
	JaegerPort     string `json:"jaeger_port,omitempty" yaml:"jaeger_port"`
	ServiceName    string `json:"service_name,omitempty" yaml:"service_name"`
	ServiceVersion string `json:"service_version,omitempty" yaml:"service_version"`
	Environment    string `json:"environment,omitempty" yaml:"environment"`
	Disabled       bool   `json:"disabled,omitempty" yaml:"disabled"`
}

// Поставщик представляет поставщика трассировщика. В зависимости от конфигурации.Параметр "Отключен",
// он будет использовать либо "живого" поставщика, либо версию "без операций".
// "Никаких операций" означает, что отслеживание будет отключено глобально.
type Provider struct {
	provider trace.TracerProvider
}
