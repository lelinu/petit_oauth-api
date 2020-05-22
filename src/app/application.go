package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lelinu/petit_oauth-api/src/http"
	"github.com/lelinu/petit_oauth-api/src/repository/db"
	"github.com/lelinu/petit_oauth-api/src/repository/rest"
	"github.com/lelinu/petit_oauth-api/src/services/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewHandler(access_token.NewService(rest.NewRepository(), db.NewDbRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	router.Run(":8080")
}
