package http

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type FiberHttp struct {
	Port         int
	app          *fiber.App
	CustomParams QueryParams
}

func NewFiberHttp() *FiberHttp {
	f := new(FiberHttp)
	f.app = fiber.New()

	f.app.Use(cors.New())

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
	return f.app.Shutdown()
}
