package http

import (
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/iamviniciuss/casino-transactions/pkg/shared/telemetry"
	"github.com/valyala/fasthttp"
)

type FiberHttp struct {
	Port         int
	app          *fiber.App
	CustomParams QueryParams
	shutdown     func()
}

func NewFiberHttp() *FiberHttp {
	f := new(FiberHttp)
	
	// Initialize OpenTelemetry
	serviceName := getEnvWithDefault("OTEL_SERVICE_NAME", "casino-api")
	serviceVersion := getEnvWithDefault("OTEL_SERVICE_VERSION", "1.0.0")
	f.shutdown = telemetry.InitOTEL(serviceName, serviceVersion)
	
	f.app = fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0", "127.0.0.1"},
	})

	// OpenTelemetry middleware (should be first)
	f.app.Use(telemetry.FiberTraceMiddleware(serviceName))
	
	// Logging middleware
	f.app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	
	f.app.Use(cors.New())
	f.app.Group("/", func(c *fiber.Ctx) error {
		f.CustomParams = NewFiberQueryParams(c.Context().QueryArgs())
		return c.Next()
	})

	f.app.Use(func(c *fiber.Ctx) error {
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")
		c.Set("Access-Control-Allow-Origin", "*")

		return c.Next()
	})

	log.Printf("OpenTelemetry initialized for service: %s", serviceName)
	return f
}

func (f *FiberHttp) Get(path string, callback FiberHandlerFunc) {
	f.app.Get(path, func(c *fiber.Ctx) error {
		localsWrapper := func(key interface{}) interface{} {
			return c.Locals(key)
		}

		result, err := callback(c.Context(), c.AllParams(), c.Body(), f.CustomParams, localsWrapper)

		if err != nil {
			c.Status(err.StatusCode)
			return c.JSON(err)
		}

		return c.JSON(result)
	})
}

func (f *FiberHttp) Post(path string, callback FiberHandlerFunc) {
	f.app.Post(path, func(c *fiber.Ctx) error {
		localsWrapper := func(key interface{}) interface{} {
			return c.Locals(key)
		}

		result, err := callback(c.Context(), c.AllParams(), c.Body(), f.CustomParams, localsWrapper)

		if err != nil {
			c.Status(err.StatusCode)
			return c.JSON(err)
		}

		return c.JSON(result)
	})
}

func (f *FiberHttp) ListenAndServe(port string) error {
	return f.app.Listen(port)
}

func (f *FiberHttp) Listen(listener net.Listener) error {
	return f.app.Listener(listener)
}

func (f *FiberHttp) Shutdown() error {
	// Shutdown OpenTelemetry first
	if f.shutdown != nil {
		f.shutdown()
	}
	return f.app.Shutdown()
}

type FiberQueryParams struct {
	Args *fasthttp.Args
}

func NewFiberQueryParams(args *fasthttp.Args) *FiberQueryParams {
	return &FiberQueryParams{
		Args: args,
	}
}

func (fqp *FiberQueryParams) GetParam(key string) []byte {
	if !fqp.Args.Has(key) {
		return []byte{}
	}
	return fqp.Args.Peek(key)
}

func (fqp *FiberQueryParams) AddParam(key string, value string) {
	fqp.Args.Add(key, value)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
