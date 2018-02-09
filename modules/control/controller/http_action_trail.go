package controller

import (
  "net/http"
  "time"
  "sync"
  "strconv"

  "github.com/cloudtropy/open-operation/utils/fun"
  "github.com/cloudtropy/open-operation/modules/control/db/mysql"
  "github.com/huandu/goroutine"
)

type ActionTrailInfo struct {
  RequestId  string
  ActionName string
  Result     string
  Timer      *time.Timer
}

var (
  httpActionTrailInfoCache = make(map[int64]*ActionTrailInfo)
  httpActionTrailInfoMu    sync.RWMutex
)

func NewHttpActionTrailInfo() *ActionTrailInfo {
  atInfo := &ActionTrailInfo{
    RequestId: fun.NewRequestId(),
    Result:    "Success",
  }

  gId := goroutine.GoroutineId()

  atInfo.Timer = time.AfterFunc(3*time.Minute, func() {
    httpActionTrailInfoMu.Lock()
    defer httpActionTrailInfoMu.Unlock()
    if _, isExist := httpActionTrailInfoCache[gId]; isExist {
      delete(httpActionTrailInfoCache, gId)
    }
  })

  httpActionTrailInfoMu.Lock()
  httpActionTrailInfoCache[gId] = atInfo
  httpActionTrailInfoMu.Unlock()
  return atInfo
}

func GetHttpActionTrailInfo() *ActionTrailInfo {
  httpActionTrailInfoMu.RLock()
  defer httpActionTrailInfoMu.RUnlock()
  return httpActionTrailInfoCache[goroutine.GoroutineId()]
}

func GetRequestId() string {
  httpActionTrailInfoMu.RLock()
  defer httpActionTrailInfoMu.RUnlock()
  if c, isExist := httpActionTrailInfoCache[goroutine.GoroutineId()]; isExist {
    return c.RequestId
  } else {
    return ""
  }
}

func DelHttpActionTrailInfo() {
  gId := goroutine.GoroutineId()
  httpActionTrailInfoMu.RLock()
  hati, isExist := httpActionTrailInfoCache[gId]
  if !isExist {
    httpActionTrailInfoMu.RUnlock()
    return
  }

  httpActionTrailInfoMu.RUnlock()
  hati.Timer.Stop()
  httpActionTrailInfoMu.Lock()
  delete(httpActionTrailInfoCache, gId)
  httpActionTrailInfoMu.Unlock()
}


func GetActionTrailList(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  start := r.FormValue("start")
  end := r.FormValue("end")
  pageIndex := r.FormValue("pageIndex")
  pageCount := r.FormValue("pageCount")
  searchBy := r.FormValue("searchBy")
  searchInfo := r.FormValue("searchInfo")

  iStart, err := strconv.ParseInt(start, 10, 64)
  if err != nil || iStart < 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  iEnd, err := strconv.ParseInt(end, 10, 64)
  if err != nil || iEnd < 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  pIndex, err := strconv.Atoi(pageIndex)
  if err != nil || pIndex < 1 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  pCount, err := strconv.Atoi(pageCount)
  if err != nil || pCount <= 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  if !((searchBy == "user" || searchBy == "module" || searchBy == "action") || 
    (searchBy == "" && searchInfo == "")) {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  totalCount, err := mysql.GetActionTrailsCount(iStart, iEnd, searchBy, searchInfo)
  if err != nil {
    log.Println("mysql.GetActionTrailsCount", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  actionTrails, err := mysql.GetActionTrails(iStart, iEnd, pIndex, pCount, searchBy, searchInfo)
  if err != nil {
    log.Println("mysql.GetActionTrails", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }
  HttpResData(w, map[string]interface{}{
    "totalCount": totalCount,
    "values": actionTrails,
  })
}
