package ctx

import (
  "encoding/json"

  "github.com/cloudtropy/open-operation/utils/file"
  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/pkg/errors"
)

type RedisConn struct {
  Host string `json:"host"`
}

type MysqlConn struct {
  Host     string `json:"host"`
  Username string `json:"username"`
  Passwd   string `json:"passwd"`
  Database string `json:"database"`
}

type ServerInfo struct {
  Host string `json:"host"`
}

type ContextConfig struct {
  HttpListenPort      int             `json:"http_listen_port"`
  RpcListenPort       int             `json:"rpc_listen_port"`
  Mysql               *MysqlConn      `json:"mysql"`
  Redis               *RedisConn      `json:"redis"`
  Control             *ServerInfo     `json:"control"`
  LogConfig           *log.LoggerConf `json:"log_config"`
  RrdPath             string          `json:"rrd_path"`
  HostOfflineTimeoutS int             `json:"host_offline_timeout_s"`
}

var (
  config  *ContextConfig
  cfgPath string
)

func Cfg() *ContextConfig {
  return config
}

func ParseJsonConfigFile(filePath string) {
  if !file.PathIsExist(filePath) {
    log.Panicf("configuration file(%s) is not exist.\n", filePath)
  }

  // cache filePath to reload configuration file
  cfgPath = filePath

  cfgContentBytes, err := file.ReadFileToBytes(filePath)
  if err != nil {
    log.Panicf("read configuration file(%s) fail: %s\n", filePath, err.Error())
  }

  c := ContextConfig{}
  err = json.Unmarshal(cfgContentBytes, &c)
  if err != nil {
    log.Panicf("unmarshal configuration file(%s)'s json content fail: %s\n", filePath, err.Error())
  }

  config = &c

  // log.Printf("read configuration file(%s) successfully.", filePath)
}

func ReloadJsonConfigFile() (err error) {
  defer func() {
    e := recover()
    if e == nil {
      return
    }
    if panicErr, ok := e.(string); ok {
      err = errors.New(panicErr)
    }
  }()

  ParseJsonConfigFile(cfgPath)
  return nil
}

func InitLog() {
  logCfg := Cfg().LogConfig
  if logCfg == nil {
    return
  }
  logCfg.Withprefix = true
  logCfg.LogFlag = log.LstdFlags | log.Lshortfile

  err := log.InitLogger(logCfg)
  if err != nil {
    log.Fatalln(err)
  }
}
