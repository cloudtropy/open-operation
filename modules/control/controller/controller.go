package controller

import (
  "encoding/json"
  "net/http"
  "time"

  "github.com/cloudtropy/open-operation/utils/rpc"
  "github.com/cloudtropy/open-operation/modules/control/ctx"
)

var (
  rpcMonitor *rpc.SingleConnRpcClient
)

type HttpResponseBody struct {
  RequestId string      `json:"requestId,omitempty"`
  Code      string      `json:"code"`
  Msg       string      `json:"msg,omitempty"`
  Data      interface{} `json:"data,omitempty"`
}

func Init() {
  LogInit()
  rpcMonitor = &rpc.SingleConnRpcClient{
    RpcServer: ctx.Cfg().Monitor.RpcHost,
    Timeout:   time.Second * 5,
  }
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

  ati := GetHttpActionTrailInfo()
  if ati == nil {
    httpResJson(w, HttpResponseBody{
      Code: code,
      Msg:  msg,
    })
    return
  }

  ati.Result = codemsg[0]
  if codemsg[1] != "" {
    ati.Result += ": " + codemsg[1]
  }
  httpResJson(w, HttpResponseBody{
    RequestId: ati.RequestId,
    Code: code,
    Msg:  msg,
  })
}

func HttpResData(w http.ResponseWriter, data interface{}) {
  httpResJson(w, HttpResponseBody{
    RequestId: GetRequestId(),
    Code:      "Success",
    Data:      data,
  })
}
