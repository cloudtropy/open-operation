package controller

import (
  "net/http"
  "encoding/json"

  "github.com/cloudtropy/open-operation/modules/monitor/db/redis"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
)

var (
  HostOnlineCacheKey = "HOSTS_ONLINE"
  HostOfflineCacheKey = "HOSTS_OFFLINE"
)

func Init() {
  mysql.Init()
  redis.Init()
  go DoOnlineHostMonitor()
}

type HttpResponseBody struct {
  Code      string      `json:"code"`
  Msg       string      `json:"msg,omitempty"`
  Data      interface{} `json:"data,omitempty"`
}

func httpResJson(w http.ResponseWriter, v interface{}) {
  bs, _ := json.Marshal(v)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.Write(bs)
}

func HttpResMsg(w http.ResponseWriter, codemsg ...string) {
  var code, msg string
  if len(codemsg) == 2 {
    code = codemsg[0]
    msg = codemsg[1]
  } else if len(codemsg) == 1 {
    code = codemsg[0]
  }
  httpResJson(w, HttpResponseBody{
    Code:      code,
    Msg:       msg,
  })
}

func HttpResData(w http.ResponseWriter, data interface{}) {
  httpResJson(w, HttpResponseBody{
    Code:      "Success",
    Data:      data,
  })
}
