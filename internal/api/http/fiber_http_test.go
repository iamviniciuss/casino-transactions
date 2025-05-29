package http

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// getFreePortListener returns a listener bound to an available port
func getFreePortListener(t *testing.T) (net.Listener, int) {
	t.Helper()

	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)

	port := listener.Addr().(*net.TCPAddr).Port
	return listener, port
}

type fibberHTTPMethod string
const (
	FIBER_GET fibberHTTPMethod = "GET"
	FIBER_POST fibberHTTPMethod = "POST"
)

func setupTestServer(t *testing.T, method fibberHTTPMethod, route string, handler FibberHandlerFunc) (*FiberHttp, int) {
	t.Helper()

	app := NewFiberHttp()


	if method == FIBER_GET {
		app.Get(route, handler)
	} else if method == FIBER_POST {
		app.Post(route, handler)
	} else {
		t.Fatalf("Unsupported HTTP method: %s", method)
	}

	// app.Get("/fiber", func(ctx context.Context, m map[string]string, b []byte, qp QueryParams, lf LocalsFunc) (interface{}, *IntegrationError) {
	// 	return "GET - OK", nil
	// })

	// app.Post("/fiber", func(ctx context.Context, m map[string]string, b []byte, qp QueryParams, lf LocalsFunc) (interface{}, *IntegrationError) {
	// 	return "POST - OK", nil
	// })

	listener, port := getFreePortListener(t)

	go func() {
		if err := app.Listen(listener); err != nil {
			panic(err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	t.Cleanup(func() {
		_ = app.Shutdown()
	})

	return app, port
}

func TestFiberGET_HTTP(t *testing.T) {
	t.Run("GET /fiber - success", func(t *testing.T) {
		_, port := setupTestServer(t, FIBER_GET, "/fiber", func(ctx context.Context, m map[string]string, b []byte, qp QueryParams, lf LocalsFunc) (interface{}, *IntegrationError) {
			return "GET - OK", nil
		})

		resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/fiber")

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		assert.NoError(t, err)
		assert.Equal(t, "\"GET - OK\"", string(body))
	})
	

	t.Run("GET /fiber - fail", func(t *testing.T) {
		_, port := setupTestServer(t, FIBER_GET, "/fiber", func(ctx context.Context, m map[string]string, b []byte, qp QueryParams, lf LocalsFunc) (interface{}, *IntegrationError) {
			return "", &IntegrationError{
				StatusCode: 500,
				Message:    "Internal Server Error",
			}
		})

		resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/fiber")

		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		assert.NoError(t, err)

		var response IntegrationError
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, response.Message, "Internal Server Error")
	})
}


func TestFiberPOST_HTTP(t *testing.T) {
	t.Run("POST /fiber - success", func(t *testing.T) {
		_, port := setupTestServer(t, FIBER_POST, "/fiber", func(ctx context.Context, m map[string]string, b []byte, qp QueryParams, lf LocalsFunc) (interface{}, *IntegrationError) {
			return "POST - OK", nil
		})

		resp, err := http.Post("http://localhost:" + strconv.Itoa(port) + "/fiber", "application/json", nil)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		assert.NoError(t, err)
		assert.Equal(t, "\"POST - OK\"", string(body))
	})

	t.Run("POST /fiber - fail", func(t *testing.T) {
		_, port := setupTestServer(t, FIBER_POST, "/fiber", func(ctx context.Context, m map[string]string, b []byte, qp QueryParams, lf LocalsFunc) (interface{}, *IntegrationError) {
			return "", &IntegrationError{
				StatusCode: 500,
				Message:    "unexpected error",
			}
		})

	resp, err := http.Post("http://localhost:" + strconv.Itoa(port) + "/fiber", "application/json", nil)

		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		assert.NoError(t, err)

		var response IntegrationError
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		assert.Equal(t, response.Error(), "Status 500: Integration Error: unexpected error")
	})
}