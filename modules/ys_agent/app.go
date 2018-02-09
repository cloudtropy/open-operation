package main

import (
  "fmt"
  "flag"
  "os"

  log "github.com/cloudtropy/open-operation/utils/logger"
  ctl "github.com/cloudtropy/open-operation/modules/ys_agent/controller"
  "github.com/cloudtropy/open-operation/modules/ys_agent/ctx"
  "github.com/cloudtropy/open-operation/modules/ys_agent/http"
)

func main() {
  cfgPath := flag.String("c", "configure.json", "set configuration file")
  ver := flag.Bool("v", false, "show version and exit")
  flag.Parse()

  if *ver {
    fmt.Println("version:", ctx.VERSION)
    os.Exit(0)
  }

  ctx.ParseJsonConfigFile(*cfgPath)
  ctx.InitLog()
  ctx.InitHostId()

  ctl.HeartbeatReport()
  ctl.StartMultiReport()
  ctl.StartMonitorReport()

  go http.ListenHttp()

  log.Println("Server version:", ctx.VERSION)
  log.Println("Server pid:", os.Getpid())

  ctx.HandleSignals()
}
