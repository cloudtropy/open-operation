package controller

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/http"

  cc "github.com/cloudtropy/open-operation/utils/common"
  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
)

func HandleHostInfo(w http.ResponseWriter, r *http.Request) {

  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Error("ioutil.ReadAll", err)
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  var metricValues []*cc.MetricValue
  err = json.Unmarshal(body, &metricValues)
  if err != nil {
    log.Error("json.Unmarshal", err)
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  var hostInfo mysql.HostInfo
  GetMetricValues(&hostInfo, metricValues)
  err = mysql.UpsertHostInfo(&hostInfo)
  if err != nil {
    log.Error("mysql.UpsertHostInfo", err)
    HttpResMsg(w, "MysqlError", err.Error())
  } else {
    HttpResMsg(w, "Success")
  }
}

func GetHostList(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  hs, err := mysql.GetHostInfos(0)
  if err != nil {
    log.Error("mysql.GetHostInfos", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  res := make([]map[string]interface{}, 0)
  for _, hi := range hs {
    isOnline := "online"
    if !GetHostsStatus(hi.HostId) {
      isOnline = "offline"
    }
    res = append(res, map[string]interface{}{
      "host_id": hi.HostId,
      "host_ip": hi.HostIp,
      "hostname": hi.Hostname,
      "host_os": hi.HostOs,
      "cpu_count": hi.CpuCount,
      "mem_capacity": hi.MemCapacity,
      "disk_capacity": hi.DiskCapacity,
      "comment": hi.Comment,
      "location": hi.Location,
      "is_online": isOnline,
    })
  }
  HttpResData(w, res)
}

func GetRemovedHostList(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  hs, err := mysql.GetHostInfos(1)
  if err != nil {
    log.Error("mysql.GetHostInfos", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  res := make([]map[string]interface{}, 0)
  for _, hi := range hs {
    res = append(res, map[string]interface{}{
      "host_id": hi.HostId,
      "host_ip": hi.HostIp,
      "hostname": hi.Hostname,
      "host_os": hi.HostOs,
      "cpu_count": hi.CpuCount,
      "mem_capacity": hi.MemCapacity,
      "disk_capacity": hi.DiskCapacity,
      "comment": hi.Comment,
      "location": hi.Location,
      "is_online": "offline",
    })
  }
  HttpResData(w, res)
}

func UpdateHostInfo(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Error("ioutil.ReadAll", err)
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  mss := make(map[string]string)
  err = json.Unmarshal(body, &mss)
  if err != nil {
    log.Error("json.Unmarshal", err)
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  if mss["update_key"] == "" || mss["host_id"] == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  err = mysql.UpdateOneHostInfo(mss)
  if err != nil {
    log.Error("mysql.UpdateOneHostInfo", mss, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }
  HttpResMsg(w, "Success")
}

func GetMetricValues(hostInfo *mysql.HostInfo, metricValues []*cc.MetricValue) {

  for i := range metricValues {
    hostInfo.HostId = metricValues[i].Endpoint
    switch metricType := metricValues[i].Metric; metricType {
    case "host.os":
      hostInfo.HostOs = InterfaceToStr(metricValues[i].Value)
    case "host.ip":
      hostInfo.HostIp = InterfaceToStr(metricValues[i].Value)
    case "host.hostname":
      hostInfo.Hostname = InterfaceToStr(metricValues[i].Value)
    case "mem.memtotal":
      hostInfo.MemCapacity = uint64(metricValues[i].Value.(float64))
    case "df.statistics.total":
      hostInfo.DiskCapacity = uint64(metricValues[i].Value.(float64))
    case "cpu.count":
      hostInfo.CpuCount = uint8(metricValues[i].Value.(float64))
    default:
      continue
    }
  }
}

func InterfaceToStr(i interface{}) string {
  return fmt.Sprintf("%v", i)
}
