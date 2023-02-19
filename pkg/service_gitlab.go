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
	pushEvent, ok := webhook.(*gitlab.PushEvent)
	if !ok {
		c.JSON(http.StatusBadRequest, "not a push event")
		return
	}

	repoId := pushEvent.Project.PathWithNamespace
	revHash := pushEvent.CheckoutSHA
	revRef := pushEvent.Ref
	// ssh or http?
	pullUrl := pushEvent.Repository.GitHTTPURL

	// ok, upload
	logger.Infof("valid upload request: %s/%s from %s", repoId, revHash, pullUrl)
	go func() {
		err := doUpload(repoId, revHash, revRef, pullUrl)
		if err != nil {
			logger.Errorf("error when upload: %v", err)
		}
	}()
}
