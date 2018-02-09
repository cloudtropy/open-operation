package g

import (
	"encoding/json"
	"log"
)

type ServerInfo struct {
	Host       string `json:"host"`
	ExtHost    string `json:"ext_host"`
	RpcHost    string `json:"rpc_host"`
	ListenPort int    `json:"listen_port"`
}

type GlobalConfig struct {
	Self         *ServerInfo `json:"self"`
	Monitor      *ServerInfo `json:"monitor"`
	Shelver      *ServerInfo `json:"shelver"`
	Deploy       *ServerInfo `json:"deploy"`
	CookieDomain string      `json:"cookie_domain"`
	ElasticHost  string      `json:"elastic_host"`
}

var (
	config *GlobalConfig
)

func Config() *GlobalConfig {
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !PathIsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent.")
	}

	configContent, err := FileToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
