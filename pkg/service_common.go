package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
)

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

	// do something else?

	return nil
}
