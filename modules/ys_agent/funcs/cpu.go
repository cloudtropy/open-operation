package funcs

import (
  "bufio"
  "bytes"
  "io"
  "io/ioutil"
  "strings"

  log "github.com/cloudtropy/open-operation/utils/logger"
  cc "github.com/cloudtropy/open-operation/utils/common"
  "github.com/cloudtropy/open-operation/utils/file"
)


func CpuMetrics() []*cc.MetricValue {
  contents, err := ioutil.ReadFile("/proc/cpuinfo")
  if err != nil {
    log.Error(err)
    return nil
  }

  reader := bufio.NewReader(bytes.NewBuffer(contents))
  cpuCount := 0
  for {
    line, err := file.ReadLine(reader)
    if err == io.EOF {
      err = nil
      break
    } else if err != nil {
      return nil
    }
    if i := strings.Index(string(line), "processor"); i != -1 {
      cpuCount += 1
    }
  }
  if cpuCount == 0 {
    cpuCount = 1
  }

  return []*cc.MetricValue{
    GaugeValue("cpu.count", cpuCount),
  }
}
