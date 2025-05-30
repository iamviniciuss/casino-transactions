package controller

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	httpService "github.com/iamviniciuss/casino-transactions/internal/api/http"
	"github.com/stretchr/testify/assert"
)

func getFreeListener() (net.Listener, int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	return listener, port, nil
}

func TestHealthCheck(t *testing.T) {
	app := httpService.NewFiberHttp()
	app.Get("/health", NewHealthCheckController().Check)
	
	port := 5271

	go func() {
		listener, portList, _ := getFreeListener()
		port = portList
		err := app.Listen(listener)
		if err != nil {
			panic(err)
		}
	}()

	defer app.Shutdown()

	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/health")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyString := string(bodyBytes)
	assert.Equal(t, "\"OK\"", bodyString)
}
