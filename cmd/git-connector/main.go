package main

import (
	"fmt"
	"log"

	"github.com/opensibyl/git-connector/pkg"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var port int
	var sibylUrl string
	var gitlabUser string
	var gitlabPwd string

	var rootCmd = &cobra.Command{
		Use:   "git-connector",
		Short: "git-connector cmd",
		Long:  `git-connector cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			config := pkg.DefaultServerConfig()
			if port != 0 {
				config.Port = port
			}
			if sibylUrl != "" {
				config.SibylUrl = sibylUrl
			}
			if gitlabUser != "" {
				config.GitlabConf.Username = gitlabUser
			}
			if gitlabPwd != "" {
				config.GitlabConf.Password = gitlabPwd
			}

			engine := pkg.CreateServer(*config)
			err := engine.Run(fmt.Sprintf(":%d", config.Port))
			if err != nil {
				return
			}
		},
	}
	rootCmd.PersistentFlags().IntVar(&port, "port", 0, "port")
	rootCmd.PersistentFlags().StringVar(&sibylUrl, "url", "", "sibyl server url")
	rootCmd.PersistentFlags().StringVar(&gitlabUser, "gitlab_user", "", "gitlab_user for http clone")
	rootCmd.PersistentFlags().StringVar(&gitlabPwd, "gitlab_pwd", "", "gitlab_pwd for http clone")
	return rootCmd
}

func main() {
	if err := createCmd().Execute(); err != nil {
		log.Fatalln(err)
	}
}
