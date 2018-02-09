package controller

import (
  "net/http"
  "strconv"
  "strings"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
  "github.com/cloudtropy/open-operation/modules/monitor/db/rrdtool"
)

type EChartData struct {
  ChartType string `json:"chartType"`
  Legend    struct {
    Data []string `json:"data"`
  } `json:"legend"`
  // Series []struct {
  //   Data []interface{} `json:"data"`
  //   Name string        `json:"name"`
  // } `json:"series"`
  Series []map[string]interface{} `json:"series"`
  Title  string                   `json:"title"`
  XAxis  struct {
    Data []int64 `json:"data"`
    Name string  `json:"name"`
    Type string  `json:"type"`
  } `json:"xAxis"`
}

func GetHostItemsForAgent(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  r.ParseForm()
  hostId := r.FormValue("hostId")
  if hostId == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  items, err := mysql.GetItemsOfHostId(hostId)
  if err != nil {
    log.Error("mysql.GetItemsOfHostId", hostId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResData(w, items)
}

func GetHostValidItemsForGraph(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  r.ParseForm()
  hostId := r.FormValue("hostId")
  if hostId == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  items, err := mysql.GetGraphValidItems(hostId)
  if err != nil {
    log.Println("mysql.GetGraphValidItems", hostId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResData(w, items)
}

func GetHostMonitorData(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  r.ParseForm()
  start := r.FormValue("start")
  end := r.FormValue("end")
  hostId := r.FormValue("hostId")
  if start == "" || end == "" || hostId == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  intStart, err1 := strconv.ParseInt(start, 10, 64)
  intEnd, err2 := strconv.ParseInt(end, 10, 64)
  if err1 != nil || err2 != nil || intEnd-intStart <= 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  } else if (intEnd - intStart) > 604800 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  var resData = []EChartData{
    {
      Title:     "cpu load",
      ChartType: "line",
      Legend: struct {
        Data []string `json:"data"`
      }{[]string{"cpu.load.1min", "cpu.load.5min", "cpu.load.15min"}},
      XAxis: struct {
        Data []int64 `json:"data"`
        Name string  `json:"name"`
        Type string  `json:"type"`
      }{[]int64{}, "timestamp per 60s", "timestamp"}},
    {
      Title:     "mem used",
      ChartType: "line",
      Legend: struct {
        Data []string `json:"data"`
      }{[]string{"mem.used.percent"}},
      XAxis: struct {
        Data []int64 `json:"data"`
        Name string  `json:"name"`
        Type string  `json:"type"`
      }{[]int64{}, "percent per 60s", "timestamp"}},
    {
      Title:     "net flow",
      ChartType: "line",
      Legend: struct {
        Data []string `json:"data"`
      }{[]string{"net.if.in", "net.if.out"}},
      XAxis: struct {
        Data []int64 `json:"data"`
        Name string  `json:"name"`
        Type string  `json:"type"`
      }{[]int64{}, "timestamp per 60s", "timestamp"}},
  }

  for i, chartData := range resData {
    timestampArr := make([]int64, 0)

    for _, itemName := range chartData.Legend.Data {
      itemDatas, err := rrdtool.FetchData(hostId, itemName, intStart, intEnd)
      if err != nil {
        log.Println("rrdtool.FetchData", err)
        if strings.Index(err.Error(), "No such file") != -1 {
          continue
        }
        log.Error("rrdtool.FetchData", hostId, itemName, intStart, intEnd, err)
        HttpResMsg(w, "InternalError", err.Error())
        return
      }

      values := make([]interface{}, 0)
      for _, itemData := range itemDatas {
        values = append(values, itemData.Value)
        timestampArr = append(timestampArr, itemData.Timestamp * 1000)
      }
      resData[i].Series = append(resData[i].Series, map[string]interface{}{
        "name": itemName,
        "data": values,
      })

      if len(resData[i].XAxis.Data) == 0 {
        resData[i].XAxis.Data = timestampArr
        timestampArr = make([]int64, 0)
      }
    }
  }

  HttpResData(w, resData)
}
