package controller

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "sync"
  "time"
  "bytes"

  log "github.com/cloudtropy/open-operation/utils/logger"
  cc "github.com/cloudtropy/open-operation/utils/common"
  "github.com/cloudtropy/open-operation/utils/fun"
  "github.com/cloudtropy/open-operation/modules/monitor/db/redis"
  "github.com/cloudtropy/open-operation/modules/monitor/ctx"
)

var (
  HostsStatus   = make(map[string]bool) // true online, false offline
  HostsStatusMu sync.RWMutex
)

func HandleHostHeartbeat(w http.ResponseWriter, r *http.Request) {

  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  var heartbeatBody struct {
    HostId string `json:"host_id"`
  }
  err = json.Unmarshal(body, &heartbeatBody)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  if len(heartbeatBody.HostId) == 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  isInsert, err := redis.UpsertHostStatus(HostOnlineCacheKey, heartbeatBody.HostId)
  if err != nil {
    HttpResMsg(w, "RedisError", err.Error())
    return
  }

  SetHostsStatus(heartbeatBody.HostId, true)

  events := GetEventCache(heartbeatBody.HostId)
  if events == nil || len(events) == 0 {
    HttpResMsg(w, "Success")
  } else {
    HttpResData(w, events)
  }

  if isInsert {
    ControlNewOnlineHost(heartbeatBody.HostId)
  }
}

func HandleHostStatus(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  mapSI, err := redis.GetAllHostStatus(HostOnlineCacheKey)
  if err != nil {
    HttpResMsg(w, "RedisError", err.Error())
    return
  }
  HttpResData(w, mapSI)
}

func SetHostsStatus(hostId string, status bool) {
  HostsStatusMu.Lock()
  defer HostsStatusMu.Unlock()
  HostsStatus[hostId] = status
}

func GetHostsStatus(hostId string) bool {
  HostsStatusMu.RLock()
  defer HostsStatusMu.RUnlock()
  status, isExist := HostsStatus[hostId]
  if !isExist {
    return false
  }
  return status
}

func DoOnlineHostMonitor() {
  var hostOfflineTimeoutS = int64(ctx.Cfg().HostOfflineTimeoutS)
  time.Sleep(time.Second * time.Duration(hostOfflineTimeoutS))
  t := time.NewTicker(time.Second).C
  for {
    mapSI, err := redis.GetAllHostStatus(HostOnlineCacheKey)
    if err != nil {
      log.Println("GetAllHostStatus:", err)
      continue
    }

    nowTimestampS := fun.NowTimestampS()
    for f, v := range mapSI {
      if nowTimestampS-v > hostOfflineTimeoutS {
        _, err := redis.DelHostStatus(HostOnlineCacheKey, f)
        if err != nil {
          log.Println("DelHostStatus:", err)
          continue
        }
        log.Println("Host online to offline:", f)
        //todo
        redis.UpsertHostStatus(HostOfflineCacheKey, f)
        go NotifyHostStatusChange(f, "offline")
      }
    }

    <-t
  }
}

func ControlNewOnlineHost(hostId string) {

  go NotifyHostStatusChange(hostId, "online")

  delCount, err := redis.DelHostStatus(HostOfflineCacheKey, hostId)
  if err != nil {
    log.Println("ControlNewOnlineHost:", err)
    return
  }

  if delCount == 0 {
    return
  }
  log.Println("Host offline to online:", hostId)
}

func NotifyHostStatusChange(hostId, changeTo string) {

  var value float64 = 0
  if changeTo == "offline" {
    SetHostsStatus(hostId, false)
    value = 1
  }
  mvs := make(map[string]*cc.MetricValue)
  mvs["host.heartbeat"] = &cc.MetricValue{
    Endpoint:  hostId,
    Metric:    "host.heartbeat",
    Value:     value,
    Step:      0,
    Type:      "GAUGE",
    Tags:      "",
    Timestamp: fun.NowTimestampS(),
  }
  // TriggerCheck(hostId, mvs) // todo

  mapSS := make(map[string]string)
  mapSS["hostId"] = hostId
  mapSS["hostStatus"] = changeTo

  NotifyControl("host_status_change", mapSS)
}


func NotifyControl(topic string, data interface{}) {

  bs, err := json.Marshal(cc.TopicMsg{Topic: topic, Data: data})
  if err != nil {
    log.Error("json.Marshal:", err)
    return
  }
  body := bytes.NewBuffer(bs)

  addr := ctx.Cfg().Control.Host + "/report"
  resp, err := http.Post(addr, "application/json;charset=utf-8", body)
  if err != nil {
    log.Printf("Report %s, error: %s\n", topic, err.Error())
    return
  }

  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    log.Printf("Report %s, http res StatusCode: %d\n", topic, resp.StatusCode)
    return
  }
}
