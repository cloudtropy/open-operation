package controller

import (
  "bytes"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "time"

  "github.com/cloudtropy/open-operation/modules/ys_agent/ctx"
  log "github.com/cloudtropy/open-operation/utils/logger"
)

type HeartbeatBody struct {
  Code string              `json:"code"`
  Data []map[string]string `json:"data"`
}

func HeartbeatReport() {
  report := ctx.Cfg().Reports["HEARTBEAT"]
  if report == nil || !report.Enabled {
    return
  }

  if hostId := ctx.HostId(); hostId == "" {
    return
  }

  var addrs []string
  for _, addr := range ctx.Cfg().MonitorSrv.Addrs {
    addrs = append(addrs, addr+ctx.Cfg().Reports["HEARTBEAT"].Path)
  }
  if addrs == nil || len(addrs) == 0 {
    return
  }

  go heartbeat(addrs)
}

func heartbeat(addrs []string) {
  heartbeatBody := map[string]string{
    "host_id": ctx.HostId(),
  }
  bs, err := json.Marshal(heartbeatBody)
  if err != nil {
    log.Println("json err:", err)
    return
  }

  t := time.NewTicker(time.Second * time.Duration(int64(ctx.Cfg().Reports["HEARTBEAT"].Interval))).C
  for {
    <-t

    for _, addr := range addrs {
      func() {
        body := bytes.NewBuffer(bs)
        resp, err := http.Post(addr, "application/json;charset=utf-8", body)
        if err != nil {
          log.Println(addr, err)
          return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
          log.Println(addr, resp.StatusCode)
          return
        }

        resBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
          log.Println(addr, err)
          return
        }
        var f HeartbeatBody
        err = json.Unmarshal(resBody, &f)
        if err != nil {
          log.Println(addr, err)
          return
        }
        if f.Code != "Success" {
          log.Println(addr, f.Code)
          return
        }

        if f.Data == nil || len(f.Data) == 0 {
          return
        }

        go HandleHeartbeatEvent(f.Data)
      }()
    }
  }
}

func HandleHeartbeatEvent(events []map[string]string) {
  for _, event := range events {
    log.Info("---receive event---", event["topic"])
    switch event["topic"] {
    case "update_items":
      // StopCollectors()
      // StartMonitorReport()
      go HandleItems()
    default:
    }
  }
}
