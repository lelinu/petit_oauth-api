package rest

import (
	"encoding/json"
	"github.com/lelinu/api_utils/errors"
	"github.com/lelinu/petit_oauth-api/src/domain/users"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
)

var (
	restClient = rest.RequestBuilder{
		BaseURL: "https://api.petit.com",
		Timeout: 100 * time.Millisecond,
	}
)

type IRestUsersRepository interface {
	Login(string, string) (*users.User, *errors.ApiError)
}

type usersRepository struct{}

func NewRepository() IRestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) Login(email string, password string) (*users.User, *errors.ApiError) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := restClient.Post("/users/login", request)
	// timeout
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerError("invalid rest client response when trying to login user")
	}
	// error situation
	if response.StatusCode > 299 {
		var apiError errors.ApiError
		if err := json.Unmarshal(response.Bytes(), &apiError); err != nil {
			return nil, errors.NewInternalServerError("invalid rest error interface")
		}
		return nil, &apiError
	}

	// unmarshall user
	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError("invalid rest user interface")
	}

	return &user, nil
}
