package controller

import (
  "time"
  "math/rand"
  "encoding/json"
  "net/http"
  "bytes"
  "io/ioutil"

  cc "github.com/cloudtropy/open-operation/utils/common"
  log "github.com/cloudtropy/open-operation/utils/logger"
)



func HttpPostJson(addr string, v interface{}) bool {
  bs, err := json.Marshal(v)
  if err != nil {
    log.Println("json err:", err)
    return false
  }
  body := bytes.NewBuffer(bs)

  client := &http.Client{}
  req, err := http.NewRequest(http.MethodPost, addr, body)
  if err != nil {
    log.Println(addr, err)
    return false
  }
  req.Header.Set("Content-Type","application/json;charset=utf-8")
  resp, err := client.Do(req)
  if err != nil {
    log.Println(addr, err)
    return false
  }

  /*resp, err := http.Post(addr, "application/json;charset=utf-8", body)
  if err != nil {
    log.Println(addr, err)
    return false
  }*/

  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    log.Printf("<= StatusCode: %d  (path:%s)\n", resp.StatusCode, addr)
    return false
  }

  result, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Error(addr, err)
    return false
  }

  log.Debug("<= ", addr, string(result))
  return true
}


func SendToConsole(metrics []*cc.MetricValue, addrs []string) {
  if len(metrics) == 0 {
    return
  }

  log.Debug("=> ", addrs, len(metrics), metrics[0])

  rand.Seed(time.Now().UnixNano())
  for _, i := range rand.Perm(len(addrs)) {
    addr := addrs[i]
    if HttpPostJson(addr, metrics) {
      break
    }
  }
}
