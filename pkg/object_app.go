package pkg

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
