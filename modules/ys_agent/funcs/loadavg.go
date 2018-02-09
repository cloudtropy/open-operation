package funcs

import (
  "strconv"
  "strings"

  "github.com/pkg/errors"
  log "github.com/cloudtropy/open-operation/utils/logger"
  cc "github.com/cloudtropy/open-operation/utils/common"
  "github.com/cloudtropy/open-operation/utils/file"
)

type Loadavg struct {
  Avg1min  float64
  Avg5min  float64
  Avg15min float64
}

func LoadAvgMetrics() []*cc.MetricValue {
  load, err := LoadAvg()
  if err != nil {
    log.Println(err)
    return nil
  }

  return []*cc.MetricValue{
    GaugeValue("cpu.load.1min", load.Avg1min),
    GaugeValue("cpu.load.5min", load.Avg5min),
    GaugeValue("cpu.load.15min", load.Avg15min),
  }

}

func LoadAvgPercentMetrics() []*cc.MetricValue {
  load, err := LoadAvg()
  if err != nil {
    log.Println(err)
    return nil
  }

  cpuMetrics := CpuMetrics()
  if cpuMetrics == nil || len(cpuMetrics) == 0 {
    return nil
  }
  cpuCount := cpuMetrics[0].Value
  loadPercent := load.Avg1min * 100.0 / float64(cpuCount.(int))

  return []*cc.MetricValue{
    GaugeValue("cpu.load.1min.percent", loadPercent),
  }

}

func LoadAvg() (*Loadavg, error) {

  loadAvg := Loadavg{}

  data, err := file.ReadFileToTrimedString("/proc/loadavg")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  L := strings.Fields(data)
  if loadAvg.Avg1min, err = strconv.ParseFloat(L[0], 64); err != nil {
    return nil, errors.WithStack(err)
  }
  if loadAvg.Avg5min, err = strconv.ParseFloat(L[1], 64); err != nil {
    return nil, errors.WithStack(err)
  }
  if loadAvg.Avg15min, err = strconv.ParseFloat(L[2], 64); err != nil {
    return nil, errors.WithStack(err)
  }

  return &loadAvg, nil
}