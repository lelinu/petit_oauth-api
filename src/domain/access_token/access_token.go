package access_token

import (
	"fmt"
	"github.com/lelinu/api_utils/utils/crypto_utils"
	"github.com/lelinu/api_utils/utils/error_utils"
	"strings"
	"time"
)

const (
	expirationTime             = 24
	GrantTypePassword          = "password"
	GrandTypeClientCredentials = "client_credentials"
)

type AccessToken struct {
	AccessToken string `json:access_token`
	UserId      int64  `json:user_id`
	ClientId    int64  `json:client_id`
	Expires     int64  `json:expires`
}

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client_credentials grant type
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *error_utils.ApiError {
	switch at.GrantType {
	case GrantTypePassword:
		at.Username = strings.TrimSpace(at.Username)
		if len(at.Username) == 0 {
			return error_utils.NewBadRequestError("Invalid username")
		}
		at.Password =  strings.TrimSpace(at.Password)
		if len(at.Password) == 0 {
			return error_utils.NewBadRequestError("Invalid password")
		}
		break

	case GrandTypeClientCredentials:

		at.ClientId = strings.TrimSpace(at.ClientId)
		if len(at.ClientId) == 0 {
			return error_utils.NewBadRequestError("Invalid client id")
		}
		at.ClientSecret =  strings.TrimSpace(at.ClientSecret)
		if len(at.ClientSecret) == 0 {
			return error_utils.NewBadRequestError("Invalid client secret")
		}

		break

	default:
		return error_utils.NewBadRequestError("invalid grant_type parameter")
	}

	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at *AccessToken) Validate() *error_utils.ApiError {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if len(at.AccessToken) == 0 {
		return error_utils.NewBadRequestError("invalid access token id")
	}
	if at.UserId <= 0 {
		return error_utils.NewBadRequestError("invalid user id")
	}
	if at.ClientId <= 0 {
		return error_utils.NewBadRequestError("invalid client id")
	}
	if at.Expires <= 0 {
		return error_utils.NewBadRequestError("invalid expiration time")
	}

	return nil
}

func (at *AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	hash, _ := crypto_utils.GenerateHashFromString(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
	at.AccessToken = hash
}
