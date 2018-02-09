package controller

import (
  "time"
  "sync"

  cc "github.com/cloudtropy/open-operation/utils/common"
  "github.com/cloudtropy/open-operation/utils/ticker"
  "github.com/cloudtropy/open-operation/modules/ys_agent/ctx"
  "github.com/cloudtropy/open-operation/modules/ys_agent/funcs"
)

type ReportFuncs struct {
  Title           string
  ReportSrv       *ctx.ServerInfo
  Fs              []func() []*cc.MetricValue
  IgnoreMetrics   map[string]bool
  BeginReport     bool
}

var (
  Mappers []ReportFuncs
  collectTickers = make(map[string]*ticker.Ticker)
  collectTickersMu sync.RWMutex
)

func BuildMappers() {

  Mappers = []ReportFuncs{
    ReportFuncs{
      Title: "HOST_INFO",
      ReportSrv: ctx.Cfg().MonitorSrv,
      Fs: []func() []*cc.MetricValue{
        funcs.OsMetrics,
        funcs.IpMetrics,
        funcs.HostnameMetrics,
        funcs.MemTotalMetrics,
        funcs.DeviceMetrics,
        funcs.CpuMetrics,
      },
      IgnoreMetrics: map[string]bool{
        "df.used": true,
        "df.used.percent": true,
      },
      BeginReport: true,
    },
  }
}

func StartMultiReport() {

  BuildMappers()

  for _, v := range Mappers {
    report := ctx.Cfg().Reports[v.Title]
    if report == nil || !report.Enabled {
      continue
    }

    var addrs []string
    for _, addr := range v.ReportSrv.Addrs {
      addrs = append(addrs, addr + report.Path)
    }

    if len(addrs) == 0 {
      continue
    }

    go collect(int64(report.Interval), v, addrs)
  }
}


func collect(sec int64, rfns ReportFuncs, addrs []string) {
  fns := rfns.Fs
  t := ticker.NewTicker(time.Second * time.Duration(sec))
  collectTickersMu.Lock()
  collectTickers[rfns.Title] = t
  collectTickersMu.Unlock()

  for {
    if !t.Running() {
      break
    }
    if !rfns.BeginReport {
      <-t.C
    }

    hostId := ctx.HostId()

    mvs := []*cc.MetricValue{}
    for _, fn := range fns {
      items := fn()
      if items == nil {
        continue
      }

      if len(items) == 0 {
        continue
      }

      for _, mv := range items {
        if b, ok := rfns.IgnoreMetrics[mv.Metric]; ok && b {
          continue
        } else {
          mvs = append(mvs, mv)
        }
      }
    }

    now := time.Now().Unix()
    for j := 0; j < len(mvs); j++ {
      mvs[j].Step = sec
      mvs[j].Endpoint = hostId
      mvs[j].Timestamp = now
    }

    //report mvs to control srv
    SendToConsole(mvs, addrs)

    if rfns.BeginReport {
      <-t.C
    }
  }
}


func CollectorDoNow(title string) {
  collectTickersMu.RLock()
  defer collectTickersMu.RUnlock()
  t, isExist := collectTickers[title]
  if !isExist {
    return
  }

  t.DoNow()
}
