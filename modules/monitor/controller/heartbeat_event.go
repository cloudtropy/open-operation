package controller

import (
  "net/http"
  "sync"
  "time"
  "io/ioutil"
  "encoding/json"

  log "github.com/cloudtropy/open-operation/utils/logger"
  cc "github.com/cloudtropy/open-operation/utils/common"
)


var (
  eventCacheMu = new(sync.Mutex)
  eventCache = make(map[string][]map[string]string)
  eventCacheTimeout = 60
  eventCacheTimers = make(map[string]*time.Timer)
)


func PostEvent(w http.ResponseWriter, r *http.Request) {
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

  var m cc.TopicMsg
  err = json.Unmarshal(body, &m)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }
  // var bData []byte
  // if m.Data != nil {
  //   bData, err = json.Marshal(m.Data)
  //   if err != nil {
  //     HttpResMsg(w, "InvalidRequestParams", err.Error())
  //     return
  //   }
  // }
  log.Info("Receive Event:", m.Topic, m.Data)
  switch m.Topic {
  case "job_add_host":
    // EventJobAddHost(w, bData)
  case "job_del_host":
    // EventJobDelHost(w, bData)
  default:
    HttpResMsg(w, "InvalidRequestParams")
  }
}

// func EventJobAddHost(w http.ResponseWriter, b []byte) {
//   var data map[string]string    // key: hostId
//   err := json.Unmarshal(b, &data)
//   if err != nil {
//     HttpResMsg(w, "InvalidRequestParams:" + err.Error())
//     return
//   } else if data["hostId"] == "" {
//     HttpResMsg(w, "InvalidRequestParams")
//     return
//   }

//   SetEventCache(data["hostId"], map[string]string{
//     "topic": "update_servers",
//   })
//   SetEventCache(data["hostId"], map[string]string{
//     "topic": "update_items",
//   })

//   err = AddHostTriggersOfHostId(data["hostId"])
//   if err != nil {
//     HttpResMsg(w, err.Error())
//     return
//   }
//   HttpResMsg(w, "Success")
// }

// func EventJobDelHost(w http.ResponseWriter, b []byte) {
//   var data map[string]string
//   err := json.Unmarshal(b, &data)
//   if err != nil {
//     HttpResMsg(w, "InvalidRequestParams:" + err.Error())
//     return
//   } else if data["hostId"] == "" {
//     HttpResMsg(w, "InvalidRequestParams")
//     return
//   }

//   SetEventCache(data["hostId"], map[string]string{
//     "topic": "update_servers",
//   })
//   SetEventCache(data["hostId"], map[string]string{
//     "topic": "update_items",
//   })

//   err = RemoveHostTriggersOfHostId(data["hostId"])
//   if err != nil {
//     HttpResMsg(w, err.Error())
//     return
//   }
//   HttpResMsg(w, "Success")
// }

func SetEventCache(hostId string, data map[string]string) {
  eventCacheMu.Lock()
  defer eventCacheMu.Unlock()
log.Debug("SetEventCache", hostId, data)
  if events, isExist := eventCache[hostId]; isExist {
    for _, d := range events {
      if d["topic"] == data["topic"] {
        if data["data"] != "" {
          d["data"] = data["data"]
        }
        return
      }
    }
    eventCache[hostId] = append(events, data)
  } else {
    eventCache[hostId] = []map[string]string{
      data,
    }
  }

  if timer, isExist := eventCacheTimers[hostId]; isExist {
    timer.Reset(time.Duration(eventCacheTimeout) * time.Second)
    return
  }

  eventCacheTimers[hostId] = time.AfterFunc(time.Duration(eventCacheTimeout) * time.Second, func () {
    eventCacheMu.Lock()
    defer eventCacheMu.Unlock()

    if _, isExist := eventCache[hostId]; isExist {
      delete(eventCache, hostId)
    }
  })
}


func GetEventCache(hostId string) []map[string]string {
  eventCacheMu.Lock()
  defer eventCacheMu.Unlock()

  events, isExist := eventCache[hostId]
  if !isExist {
    return nil
  }

  delete(eventCache, hostId)

  if timer, isExist := eventCacheTimers[hostId]; isExist {
    timer.Stop()
    delete(eventCacheTimers, hostId)
  }

  return events
}
