package rest

import (
	"encoding/json"
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/petit_oauth-api/src/domain/users"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
)

var (
	restClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8081",
		Timeout: 100 * time.Millisecond,
	}
)

type IRestUsersRepository interface {
	Login(string, string) (*users.User, *error_utils.ApiError)
}

type usersRepository struct{}

func NewRepository() IRestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) Login(email string, password string) (*users.User, *error_utils.ApiError) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := restClient.Post("/users/login", request)
	// timeout
	if response == nil || response.Response == nil {
		return nil, error_utils.NewInternalServerError("invalid rest client response when trying to login user")
	}
	// error situation
	if response.StatusCode > 299 {
		var apiError error_utils.ApiError
		if err := json.Unmarshal(response.Bytes(), &apiError); err != nil {
			return nil, error_utils.NewInternalServerError("invalid rest error interface")
		}
		return nil, &apiError
	}

	// unmarshall user
	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, error_utils.NewInternalServerError("invalid rest user interface")
	}

	return &user, nil
}
