package access_token

import (
	"github.com/lelinu/api_utils/errors"
	at "github.com/lelinu/petit_oauth-api/src/domain/access_token"
	"github.com/lelinu/petit_oauth-api/src/repository/rest"
	"strings"
)

type IRepository interface {
	GetById(string) (*at.AccessToken, *errors.ApiError)
	Create(at.AccessToken) *errors.ApiError
	UpdateExpirationTime(at.AccessToken) *errors.ApiError
}

type IService interface {
	GetById(string) (*at.AccessToken, *errors.ApiError)
	Create(request at.AccessTokenRequest) (*at.AccessToken, *errors.ApiError)
	UpdateExpirationTime(at.AccessToken) *errors.ApiError
}

type service struct {
	restUsersRepo rest.IRestUsersRepository
	dbRepo        IRepository
}

func NewService(restUsersRepo rest.IRestUsersRepository, dbRepo IRepository) IService {
	return &service{
		restUsersRepo: restUsersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*at.AccessToken, *errors.ApiError) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request at.AccessTokenRequest) (*at.AccessToken, *errors.ApiError) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var username = request.Username
	var password = request.Password

	if request.GrantType == at.GrandTypeClientCredentials {
		username = request.ClientId
		password = request.ClientSecret
	}

	user, err := s.restUsersRepo.Login(username, password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token:
	at := at.GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at at.AccessToken) *errors.ApiError {
	if err := at.Validate(); err != nil {
		return err
	}

	return s.dbRepo.UpdateExpirationTime(at)
}
