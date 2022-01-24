package go_trace_lib

import (
	"go.opentelemetry.io/otel/trace"
)

// ProviderConfig представляет конфигурацию поставщика и используется для создания нового типа "Поставщик".
type ProviderConfig struct {
	JaegerEndpoint string
	JaegerHost     string
	JaegerPort     string
	ServiceName    string
	ServiceVersion string
	Environment    string
	// Set this to `true` if you want to disable tracing completly.
	Disabled bool
}

// Поставщик представляет поставщика трассировщика. В зависимости от конфигурации.Параметр "Отключен",
// он будет использовать либо "живого" поставщика, либо версию "без операций".
// "Никаких операций" означает, что отслеживание будет отключено глобально.
type Provider struct {
	provider trace.TracerProvider
}
