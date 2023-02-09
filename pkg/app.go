package pkg

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanzy/go-gitlab"
)

func handleGitlabWebhook(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to read request body")
		return
	}
	webhook, err := gitlab.ParseWebhook(gitlab.EventTypePush, raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid gitlab request body")
		return
	}
	_ = webhook.(*gitlab.PushEvent)
}

func CreateServer() *gin.Engine {
	engine := gin.Default()
	apiV1Group := engine.Group("/api/v1")
	apiV1Group.Handle(http.MethodPost, "/gitlab", handleGitlabWebhook)
	return engine
}
