package http

import (
  "net/http"
  "time"
  "strconv"
  "encoding/json"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/ys_agent/ctx"
)


func ListenHttp() {

  http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
    bs, _ := json.Marshal(map[string]string{
      "version": ctx.VERSION,
    })
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Write(bs)
  })


  var listenPort = ctx.Cfg().ListenPort
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
