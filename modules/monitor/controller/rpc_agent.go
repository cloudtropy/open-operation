package controller

import (
  "strings"

  log "github.com/cloudtropy/open-operation/utils/logger"
  cc "github.com/cloudtropy/open-operation/utils/common"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
  "github.com/cloudtropy/open-operation/modules/monitor/db/rrdtool"
)

type Agent struct {}

func (a *Agent) HostItems(hostId string, res *[]cc.HostItem) error {
  var err error
  *res, err = mysql.GetItemsForAgent(hostId)
  if err != nil {
    log.Println("mysql.GetItemsForAgent", hostId, err)
    return err
  }
  return nil
}

func (a *Agent) HandleReportData(args cc.MetricValue, res *string) error {
  HandleMetrics([]*cc.MetricValue{&args})
  *res = "Success"
  return nil
}

func HandleMetrics(metricValues []*cc.MetricValue) {
  mapMetrics := make(map[string]*cc.MetricValue)
  for _, mv := range metricValues {
    if mv.DataType == "string" || mv.DataType == "none" {
      continue
    } else if strings.Index(mv.Tags, "NoHistory") != -1 {
      mapMetrics[mv.Metric] = mv
      continue
    }

    var err error
    if mv.Type == "GAUGE" {
      err = rrdtool.UpdateData(mv.Endpoint, mv.Metric, mv.Step, mv.Value)
    } else if mv.Type == "COUNTER" {
      tmpVI, ok := mv.Value.(int64)
      if !ok {
        log.Debug("metric value type not int64:", mv)
        tmpVF, ok := mv.Value.(float64)
        if !ok {
          log.Warn("metric value type not int64 or float64:", mv)
          continue
        }
        err = rrdtool.UpdateData(mv.Endpoint, mv.Metric, mv.Step, int64(tmpVF))
      } else {
        err = rrdtool.UpdateData(mv.Endpoint, mv.Metric, mv.Step, tmpVI)
      }
    }
    if err != nil {
      log.Printf("rrdtool.UpdateData(%s, %s, %d, %v) error:%s.\n", 
        mv.Endpoint, mv.Metric, mv.Step, mv.Value, err.Error())
      continue
    }

    mapMetrics[mv.Metric] = mv
  }

  // TriggerCheck(metricValues[0].Endpoint, mapMetrics)
}
