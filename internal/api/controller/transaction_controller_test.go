package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	httpService "github.com/iamviniciuss/casino-transactions/internal/api/http"
	"github.com/iamviniciuss/casino-transactions/internal/repository"
	"github.com/iamviniciuss/casino-transactions/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

// func getFreeListener() (net.Listener, int, error) {
// 	listener, err := net.Listen("tcp", ":0")
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	port := listener.Addr().(*net.TCPAddr).Port
// 	return listener, port, nil
// }

func TestTransactionController(t *testing.T) {
	dbConn, teardown := test_utils.SetupPostgres(t)
	defer teardown()

	repo := repository.NewTransactionRepository(dbConn)

	app := httpService.NewFiberHttp()
	app.Get("/transactions", NewTransactionController(repo).GetTransactions)

	port := 5274

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

	t.Run("GetTransactions - No Transactions", func(t *testing.T) {
		resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/transactions?user_id=62d54d96-88f4-4111-8564-c043d710bdcd&limit=10&offset=0")
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response PaginatedResponse
		err = json.Unmarshal(bodyBytes, &response)
		assert.NoError(t, err)
		assert.NotNil(t, response)

		assert.Equal(t, 0, response.Total)
		assert.Equal(t, 0, response.Offset)
		assert.Equal(t, 10, response.Limit)
		assert.Len(t, response.Items, 0)
	})

	t.Run("GetTransactions - Require user_id", func(t *testing.T) {
		resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/transactions")
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		var integrationError httpService.IntegrationError
		err = json.Unmarshal(bodyBytes, &integrationError)
		assert.NoError(t, err)
		assert.Equal(t, ErrUserIDRequired, integrationError)
		assert.Equal(t, 400, integrationError.StatusCode)
	})

	t.Run("GetTransactions - Error on repository when user_id is invalid", func(t *testing.T) {
		resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/transactions?user_id=invalid")
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		bodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		var integrationError httpService.IntegrationError
		err = json.Unmarshal(bodyBytes, &integrationError)
		assert.NoError(t, err)
		assert.Equal(t, 400, integrationError.StatusCode)
	})
}
