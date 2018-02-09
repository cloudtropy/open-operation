package funcs

import (
  "strings"
  
  cc "github.com/cloudtropy/open-operation/utils/common"
)

func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *cc.MetricValue {
  mv := cc.MetricValue{
    Metric: metric,
    Value:  val,
    Type:   dataType,
  }

  size := len(tags)

  if size > 0 {
    mv.Tags = strings.Join(tags, ",")
  }

  return &mv
}

func GaugeValue(metric string, val interface{}, tags ...string) *cc.MetricValue {
  return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val interface{}, tags ...string) *cc.MetricValue {
  return NewMetricValue(metric, val, "COUNTER", tags...)
}
