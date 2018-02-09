package main

import (
  "flag"
  "os"
  "fmt"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/ctx"
  "github.com/cloudtropy/open-operation/modules/monitor/router"
  "github.com/cloudtropy/open-operation/modules/monitor/controller"
)

func main() {
  cfgPath := flag.String("c", "configure.json", "set configuration file")
  // sig := flag.String("s", "", "send signal to a master process: stop, restart, reload")
  ver := flag.Bool("v", false, "show version and exit")
  flag.Parse()

  if *ver {
    fmt.Println("version:", ctx.VERSION)
    os.Exit(0)
  }

  ctx.ParseJsonConfigFile(*cfgPath)
  ctx.InitLog()
  controller.Init()

  log.Info("Server version:", ctx.VERSION)
  log.Info("Server pid:", os.Getpid())

  go router.ListenJsonRpc()
  router.ListenHttp()
}
