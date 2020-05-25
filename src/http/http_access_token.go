package http

import (
	"github.com/gin-gonic/gin"
	"github.com/lelinu/api_utils/utils/error_utils"
	at "github.com/lelinu/petit_oauth-api/src/domain/access_token"
	"github.com/lelinu/petit_oauth-api/src/services/access_token"
	"net/http"
	"strings"
)

type IAccessTokenHandler interface{
	GetById(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct{
	service access_token.IService
}

func NewHandler(service access_token.IService) IAccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (h *accessTokenHandler) GetById(g *gin.Context){
	accessToken, err := h.service.GetById(strings.TrimSpace(g.Param("access_token_id")))
	if err != nil{
		g.JSON(err.HttpStatusCode, err)
		return
	}
	g.JSON(http.StatusOK, accessToken)
}

func (h *accessTokenHandler) Create(g *gin.Context){
	var request at.AccessTokenRequest
	if err := g.ShouldBindJSON(&request); err != nil{
		restErr := error_utils.NewBadRequestError("invalid JSON body")
		g.JSON(restErr.HttpStatusCode, restErr)
		return
	}

	result, err := h.service.Create(request)
	if err != nil{
		g.JSON(err.HttpStatusCode, err)
		return
	}

	g.JSON(http.StatusCreated, result)
}

