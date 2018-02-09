package router

import (
  "net/http"
  "time"
  "strconv"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/control/ctx"
  ctl "github.com/cloudtropy/open-operation/modules/control/controller"
)


func ListenHttp() {

  ctl.HandleWithMid("/api/login", ctl.UserLogin)
  ctl.HandleWithMid("/api/logout", ctl.UserLogout)

  // websocket
  http.HandleFunc("/wsapi/msg", ctl.HandleWsMsg)
  http.HandleFunc("/report", ctl.HandleReport)

  // ctl.HandleWithMid("/report", ctl.HandleReport)
  ctl.HandleWithMid("/api/user/info", ctl.GetUserInfoSelf, ctl.AuthSessionCheck)
  ctl.HandleWithMid("/api/user/update", ctl.UpdateUserSelf, ctl.AuthSessionCheck)

  ctl.HandleWithMid("/api/authority", ctl.HandleAction, ctl.AuthSessionCheck)
  ctl.HandleWithMid("/api/basic/monitor", ctl.HandleAction, ctl.AuthSessionCheck)


  var listenPort = ctx.Cfg().HttpListenPort
  s := &http.Server{
    Addr:           ":" + strconv.Itoa(listenPort),
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20, //if not set, use (DefaultMaxHeaderBytes = 1 << 20) // 1 MB
    //ErrorLog *log.Logger
  }

  log.Println("Server listen on:", listenPort)
  log.Fatal(s.ListenAndServe())
}
