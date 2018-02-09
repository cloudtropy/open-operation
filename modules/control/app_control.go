package main

import (
  "flag"
  "os"
  "fmt"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/control/controller"
  "github.com/cloudtropy/open-operation/modules/control/ctx"
  "github.com/cloudtropy/open-operation/modules/control/router"
  "github.com/cloudtropy/open-operation/modules/control/db/redis"
  "github.com/cloudtropy/open-operation/modules/control/db/mysql"
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

  mysql.Init()
  redis.Init()
  controller.Init()

  log.Println("Server version:", ctx.VERSION)
  log.Println("Server pid:", os.Getpid())

  router.ListenHttp()
}
