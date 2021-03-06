package db

import (
	"github.com/gocql/gocql"
	"github.com/lelinu/api_utils/utils/error_utils"
	"github.com/lelinu/petit_oauth-api/src/clients/cassandra"
	"github.com/lelinu/petit_oauth-api/src/domain/access_token"
)

const(
	queryGetAccessToken = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryInsertAccessToken = "INSERT INTO access_tokens (access_token, user_id, client_id, expires) VALUES (?,?,?,?);"
	queryUpdateAccessTokenExpires = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

func NewDbRepository() IDbRepository {
	return &dbRepository{}
}

type IDbRepository interface{
	GetById(string) (*access_token.AccessToken, *error_utils.ApiError)
	Create(access_token.AccessToken) *error_utils.ApiError
	UpdateExpirationTime(access_token.AccessToken) *error_utils.ApiError
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *error_utils.ApiError){
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires); err != nil{

			if err == gocql.ErrNotFound {
				return nil, error_utils.NewNotFoundError("no access token found with given id")
			}

		return nil, error_utils.NewInternalServerError(err.Error())
	}

	return nil, nil
}

func (r *dbRepository) Create(token access_token.AccessToken) *error_utils.ApiError {

	if err := cassandra.GetSession().Query(queryInsertAccessToken,
		token.AccessToken,
		token.UserId,
		token.ClientId,
		token.Expires).Exec(); err != nil{
		return error_utils.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(token access_token.AccessToken) *error_utils.ApiError {

	if err := cassandra.GetSession().Query(queryUpdateAccessTokenExpires,
		token.Expires,
		token.AccessToken).Exec(); err != nil{
		return error_utils.NewInternalServerError(err.Error())
	}

	return nil
}