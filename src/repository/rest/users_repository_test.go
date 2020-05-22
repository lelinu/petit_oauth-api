package rest

import (
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.petit.com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"hello.world@gmail.com", "password": "123456"}`,
		RespHTTPCode: http.StatusInternalServerError,
	})

	repo := NewRepository()
	user, err := repo.Login("hello.world@gmail.com", "123456")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.EqualValues(t, http.StatusInternalServerError, err.HttpStatusCode)
}

func TestLoginInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.petit.com/users/login",
		HTTPMethod:   http.MethodPost,
		RespHTTPCode: http.StatusInternalServerError,
		RespBody:     `{"message": "invalid login credentials", "http_status_code": "404", "error": "not_found"}`,
	})

	repo := NewRepository()
	user, err := repo.Login("hello.world@gmail.com", "123456")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.EqualValues(t, http.StatusInternalServerError, err.HttpStatusCode)
	assert.EqualValues(t, "invalid rest error interface", err.Message)
}

func TestLoginInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.petit.com/users/login",
		HTTPMethod:   http.MethodPost,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "http_status_code": 404, "error": "not_found"}`,
	})

	repo := NewRepository()
	user, err := repo.Login("hello.world@gmail.com", "123456")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.EqualValues(t, http.StatusNotFound, err.HttpStatusCode)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.petit.com/users/login",
		HTTPMethod:   http.MethodPost,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{
							"id": "1",
							"first_name": "lelinu",
							"last_name": "mercieca",
							"email": "man.mercieca@gmail.com",
							"status": "Active",
							"date_created": "2020-05-11 16:42:05"
						}`,
	})

	repo := NewRepository()
	user, err := repo.Login("hello.world@gmail.com", "123456")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.EqualValues(t, http.StatusInternalServerError, err.HttpStatusCode)
	assert.EqualValues(t, "invalid rest user interface", err.Message)
}

func TestLoginNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.petit.com/users/login",
		HTTPMethod:   http.MethodPost,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{
							"id": 1,
							"first_name": "lelinu",
							"last_name": "mercieca",
							"email": "man.mercieca@gmail.com"
						}`,
	})

	repo := NewRepository()
	user, err := repo.Login("hello.world@gmail.com", "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "lelinu", user.FirstName)
	assert.EqualValues(t, "mercieca", user.LastName)
	assert.EqualValues(t, "man.mercieca@gmail.com", user.Email)
}
