package router

import (
  "net/http"
  "strconv"
  "time"

  ctl "github.com/cloudtropy/open-operation/modules/monitor/controller"
  "github.com/cloudtropy/open-operation/modules/monitor/ctx"
  log "github.com/cloudtropy/open-operation/utils/logger"
)

func ListenHttp() {

  http.HandleFunc("/host/heartbeat", ctl.HandleHostHeartbeat)
  http.HandleFunc("/host/info/update", ctl.UpdateHostInfo)
  http.HandleFunc("/host/info", ctl.HandleHostInfo)
  http.HandleFunc("/host/list", ctl.GetHostList)
  http.HandleFunc("/removed/host/list", ctl.GetRemovedHostList)
  http.HandleFunc("/host/monitor/data", ctl.GetHostMonitorData)

  http.HandleFunc("/graph/ctrl", ctl.HandleGraph)
  http.HandleFunc("/screen/ctrl", ctl.HandleScreen)
  http.HandleFunc("/screen/graphs", ctl.GetScreenGraphData)
  http.HandleFunc("/preview/graph", ctl.HandPreviewGraph)
  http.HandleFunc("/trigger/ctrl", ctl.HandleTrigger)

  http.HandleFunc("/detail/items", ctl.GetDetailItems)
  http.HandleFunc("/items", ctl.GetItemList)
  http.HandleFunc("/host/items", ctl.GetHostItemsForAgent)
  http.HandleFunc("/graph/valid/items", ctl.GetHostValidItemsForGraph)

  var listenPort = ctx.Cfg().HttpListenPort
  s := &http.Server{
    Addr:           ":" + strconv.Itoa(listenPort),
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20, //if not set, use (DefaultMaxHeaderBytes = 1 << 20) // 1 MB
    //ErrorLog *log.Logger
  }

  log.Info("Server listen http on:", listenPort)
  log.Fatal(s.ListenAndServe())
}
