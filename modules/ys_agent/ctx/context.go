package ctx

import (
  "encoding/json"
  "io/ioutil"

  "github.com/cloudtropy/open-operation/utils/file"
  "github.com/cloudtropy/open-operation/utils/fun"
  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/pkg/errors"
)

type ReportConfig struct {
  Enabled  bool   `json:"enabled"`
  Path     string `json:"path"`
  Interval int    `json:"interval"`
}

type ServerInfo struct {
  Addrs   []string `json:"addrs"`
  RpcHost string   `json:"rpc_host"`
  Timeout int      `json:"timeout"`
}

type ContextConfig struct {
  ListenPort int                      `json:"listen_port"`
  HostIdPath string                   `json:"host_id_path"`
  MonitorSrv *ServerInfo              `json:"monitor_srv"`
  Reports    map[string]*ReportConfig `json:"reports"`
  LogConfig  *log.LoggerConf          `json:"log_config"`
}

var (
  localIp string
  hostId  string
  config  *ContextConfig
  cfgPath string
)

func HostId() string {
  return hostId
}

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

func InitHostId() {
  hostIdPath := Cfg().HostIdPath

  if !file.PathIsExist(hostIdPath) {
    log.Printf("host_id file(%s) is not existent.\n", hostIdPath)

    initHostId := fun.GetUUIDV4()

    err := ioutil.WriteFile(hostIdPath, []byte(initHostId), 0444)
    if err != nil {
      log.Fatalln("create host_id failed: ", err)
    }
    log.Println("create host_id: ", initHostId)
    hostId = initHostId
    return
  }

  // todo: hostId format check
  var err error
  hostId, err = file.ReadFileToTrimedString(hostIdPath)
  if err != nil {
    log.Fatalf("read host_id file(%s) fail: %v\n", hostIdPath, err)
  }
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
