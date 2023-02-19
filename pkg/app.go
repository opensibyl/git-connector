package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var serverConfig *ServerConfig

func CreateServer(config ServerConfig) *gin.Engine {
	// shared config
	serverConfig = &config
	// logger
	prod, _ := zap.NewProduction()
	logger = prod.Sugar()

	engine := gin.Default()
	apiV1Group := engine.Group("/api/v1")
	apiV1Group.Handle(http.MethodPost, "/gitlab", handleGitlabWebhook)
	return engine
}
