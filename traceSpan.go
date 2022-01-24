package go_trace_lib

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Новый диапазон возвращает новый диапазон из глобального трассировщика.
// В зависимости от аргумента "cus" диапазон может быть либо простым, либо настраиваемым.
// Каждый результирующий интервал должен быть завершен с помощью функции `defer span.End()` сразу после вызова.
func NewSpan(ctx context.Context, name string, cus SpanCustomiser) (context.Context, trace.Span) {
	Log.Debugf("Create new span: %s", name)
	if cus == nil {
		return otel.Tracer("").Start(ctx, name)
	}

	return otel.Tracer("").Start(ctx, name)
}

// Интервал из контекста возвращает текущий интервал из контекста.
// Если вы хотите избежать создания дочерних интервалов для каждой операции и просто полагаться на родительский интервал,
// используйте эту функцию во всем приложении. При такой практике вы получите более плоское дерево пролетов,
// в отличие от более глубокой версии. Вы всегда можете смешивать и сочетать обе функции.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// Добавить теги диапазона добавляет новые теги в диапазон. Он появится в разделе "Теги"
// выбранного диапазона. Используйте это, если вы считаете, что тег и его значение могут быть полезны при отладке.
func AddSpanTags(span trace.Span, tags map[string]string) {
	Log.Debugf("Add span: %s tags", span.SpanContext().SpanID())
	list := make([]attribute.KeyValue, len(tags))
	var i int
	for k, v := range tags {
		list[i] = attribute.Key(k).String(v)
		i++
	}

	span.SetAttributes(list...)
}

// AddSpanEvents добавляет новые события в диапазон. Он появится в разделе "Журналы" выбранного диапазона.
// Используйте это, если событие может означать что-то ценное во время отладки.
func AddSpanEvents(span trace.Span, name string, events map[string]string) {
	Log.Debugf("Add span: %s events: %s", span.SpanContext().SpanID(), name)
	list := make([]trace.EventOption, len(events))
	var i int
	for k, v := range events {
		list[i] = trace.WithAttributes(attribute.Key(k).String(v))
		i++
	}
	span.AddEvent(name, list...)
}

// Ошибка добавления диапазона добавляет новое событие в диапазон. Он появится в разделе "Журналы" выбранного диапазона.
// Это не приведет к тому, что промежуток будет помечен как "неудачный".
// Используйте это, если вы считаете, что вам следует регистрировать любые исключения, такие как критические, ошибки,
// предупреждения, предостережения и т.д. Избегайте регистрации конфиденциальных данных!
func AddSpanError(span trace.Span, err error) {
	Log.Debugf("Add span: %s error: %s", span.SpanContext().SpanID(), err.Error())
	span.RecordError(err)
}

// Интервал сбоя помечает интервал как "сбой" и добавляет метку "ошибка" в указанную трассировку.
// Используйте это после вызова функции `AddSpanError`, чтобы для нее было зарегистрировано какое-то соответствующее исключение.
func FailSpan(span trace.Span, msg string) {
	Log.Errorf("Span: %s is faild: %s", span.SpanContext().SpanID(), msg)
	span.SetStatus(codes.Error, msg)
}

// SpanCustomiser используется для применения пользовательских параметров диапазона.
// Любой пользовательский тип настройщика бетонных пролетов должен реализовывать этот интерфейс.
type SpanCustomiser interface {
	customise() []trace.SpanOption
}
