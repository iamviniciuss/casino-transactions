package telemetry

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// FiberTraceMiddleware creates OpenTelemetry middleware for Fiber
func FiberTraceMiddleware(serviceName string) fiber.Handler {
	tracer := otel.Tracer(serviceName)
	
	return func(c *fiber.Ctx) error {
		// Extract context from headers
		ctx := c.UserContext()
		ctx = otel.GetTextMapPropagator().Extract(ctx, &fiberCarrier{c: c})
		
		// Create span
		spanName := fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
		if spanName == " " {
			spanName = fmt.Sprintf("%s %s", c.Method(), c.Path())
		}
		
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(
				semconv.HTTPMethod(c.Method()),
				semconv.HTTPURL(c.OriginalURL()),
				semconv.HTTPRoute(c.Route().Path),
				semconv.HTTPScheme(c.Protocol()),
				semconv.UserAgentOriginal(c.Get("User-Agent")),
				semconv.HTTPRequestContentLength(len(c.Body())),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()
		
		// Set context in Fiber
		c.SetUserContext(ctx)
		
		// Track request start time
		start := time.Now()
		
		// Continue to next handler
		err := c.Next()
		
		// Track request duration and set final attributes
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()
		
		span.SetAttributes(
			semconv.HTTPStatusCode(statusCode),
			semconv.HTTPResponseContentLength(len(c.Response().Body())),
			attribute.Int64("http.request.duration_ms", duration.Milliseconds()),
		)
		
		// Set span status based on HTTP status code
		if statusCode >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
		}
		
		// Record error if any
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		
		return err
	}
}

// fiberCarrier implements TextMapCarrier for Fiber context
type fiberCarrier struct {
	c *fiber.Ctx
}

func (fc *fiberCarrier) Get(key string) string {
	return fc.c.Get(key)
}

func (fc *fiberCarrier) Set(key, value string) {
	fc.c.Set(key, value)
}

func (fc *fiberCarrier) Keys() []string {
	keys := make([]string, 0)
	fc.c.Request().Header.VisitAll(func(key, _ []byte) {
		keys = append(keys, string(key))
	})
	return keys
}