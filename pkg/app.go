package pkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var serverConfig *ServerConfig

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

func doUpload(repoId string, revHash string, revRef string, pullUrl string) error {
	// clone to tmp dir
	// permission issues
	safeRepoId := strings.ReplaceAll(repoId, "/", "_")
	tmpSrcDir, err := os.MkdirTemp("", fmt.Sprintf("sibyl_%s_%s", safeRepoId, revHash))
	// already existed?
	if err != nil {
		logger.Errorf("create tmp dir failed")
		return err
	}
	defer func(path string) {
		logger.Infof("clean up: %s", path)
		err := os.RemoveAll(path)
		if err != nil {
			logger.Errorf("clean up failed: %v", err)
		}
	}(tmpSrcDir)

	logger.Infof("start cloning: %s", pullUrl)
	// OOM?
	// https://github.com/go-git/go-git/issues/315
	repo, err := git.PlainClone(tmpSrcDir, false, &git.CloneOptions{
		URL: pullUrl,
		// if another push happened after webhook, depth=1 is not enough
		Depth:         3,
		ReferenceName: plumbing.ReferenceName(revRef),
		SingleBranch:  true,
		Auth: &gitHttp.BasicAuth{
			Username: serverConfig.GitlabConf.Username,
			Password: serverConfig.GitlabConf.Password,
		},
	})
	if err != nil {
		logger.Errorf("clone failed: %v", err)
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(revHash),
	})
	if err != nil {
		logger.Errorf("checkout failed: %v", err)
		return err
	}

	// apply analysis and upload (sibyl2 client
	logger.Infof("start sibyl2 upload: %s", tmpSrcDir)
	conf := upload.DefaultConfig()
	conf.Src = tmpSrcDir
	conf.Url = serverConfig.SibylUrl
	conf.RepoId = repoId
	conf.RevHash = revHash
	upload.ExecWithConfig(conf)
	if err != nil {
		return fmt.Errorf("error when calling sibyl uploader: %w", err)
	}
	logger.Infof("upload finished: %s/%s", repoId, revHash)
	return nil
}

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

type gitlabConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Port       int         `json:"port"`
	SibylUrl   string      `json:"sibylUrl"`
	GitlabConf *gitlabConf `json:"gitlabConf"`
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:     9448,
		SibylUrl: "http://127.0.0.1:9876",
		GitlabConf: &gitlabConf{
			Username: "",
			Password: "",
		},
	}
}
