package funcs

import (
  "bufio"
  "bytes"
  "io"
  "io/ioutil"
  "strconv"
  "strings"

  "github.com/pkg/errors"
  log "github.com/cloudtropy/open-operation/utils/logger"
  cc "github.com/cloudtropy/open-operation/utils/common"
  "github.com/cloudtropy/open-operation/utils/file"
)

type Mem struct {
  Buffers      uint64
  Cached       uint64
  MemTotal     uint64
  MemFree      uint64
  MemUsed      uint64
  SwapTotal    uint64
  SwapUsed     uint64
  SwapFree     uint64
  Shmem        uint64
  SReclaimable uint64
}

func MemMetrics() []*cc.MetricValue {
  m, err := MemInfo()
  if err != nil {
    log.Println(err)
    return nil
  }

  memUsed := m.MemUsed
  memUsed -= m.Buffers + m.Cached
  memFree := m.MemTotal - memUsed

  pmemUsed := 0.0
  if m.MemTotal != 0 {
    pmemUsed = float64(memUsed) * 100.0 / float64(m.MemTotal)
  }

  return []*cc.MetricValue{
    //GaugeValue("mem.memtotal", m.MemTotal),
    GaugeValue("mem.used", memUsed),
    GaugeValue("mem.free", memFree),
    // GaugeValue("mem.swaptotal", m.SwapTotal),
    // GaugeValue("mem.swapused", m.SwapUsed),
    // GaugeValue("mem.swapfree", m.SwapFree),
    // GaugeValue("mem.memfree.percent", pmemFree),
    GaugeValue("mem.used.percent", pmemUsed),
    // GaugeValue("mem.swapfree.percent", pswapFree),
    // GaugeValue("mem.swapused.percent", pswapUsed),
  }

}

func MemTotalMetrics() []*cc.MetricValue {
  memTotal, err := MemTotalInfo()
  if err != nil {
    log.Println(err)
    return nil
  }

  return []*cc.MetricValue{
    GaugeValue("mem.memtotal", memTotal),
  }
}

var Multi uint64 = 1

var WANT = map[string]struct{}{
  "Buffers:":      {},
  "Cached:":       {},
  "MemTotal:":     {},
  "MemFree:":      {},
  "SwapTotal:":    {},
  "SwapFree:":     {},
  "Shmem:":        {},
  "SReclaimable:": {},
}

func MemInfo() (*Mem, error) {
  contents, err := ioutil.ReadFile("/proc/meminfo")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  memInfo := &Mem{}

  reader := bufio.NewReader(bytes.NewBuffer(contents))

  for {
    line, err := file.ReadLine(reader)
    if err == io.EOF {
      err = nil
      break
    } else if err != nil {
      return nil, errors.WithStack(err)
    }

    fields := strings.Fields(string(line))
    fieldName := fields[0]

    _, ok := WANT[fieldName]
    if ok && len(fields) == 3 {
      val, numerr := strconv.ParseUint(fields[1], 10, 64)
      if numerr != nil {
        continue
      }
      val *= Multi
      switch fieldName {
      case "Buffers:":
        memInfo.Buffers = val
      case "Cached:":
        memInfo.Cached = val
      case "MemTotal:":
        memInfo.MemTotal = val
      case "MemFree:":
        memInfo.MemFree = val
      case "SwapTotal:":
        memInfo.SwapTotal = val
      case "SwapFree:":
        memInfo.SwapFree = val
      case "Shmem:":
        memInfo.Shmem = val
      case "SReclaimable:":
        memInfo.SReclaimable = val
      }
    }
  }

  memInfo.SwapUsed = memInfo.SwapTotal - memInfo.SwapFree
  memInfo.MemUsed = memInfo.MemTotal - memInfo.MemFree
  memInfo.Cached = memInfo.Cached + memInfo.SReclaimable - memInfo.Shmem

  return memInfo, nil
}

func MemTotalInfo() (uint64, error) {
  contents, err := ioutil.ReadFile("/proc/meminfo")
  if err != nil {
    return 0, errors.WithStack(err)
  }

  reader := bufio.NewReader(bytes.NewBuffer(contents))

  for {
    line, err := file.ReadLine(reader)
    if err == io.EOF {
      err = nil
      break
    } else if err != nil {
      return 0, errors.WithStack(err)
    }

    fields := strings.Fields(string(line))
    fieldName := fields[0]

    if fieldName == "MemTotal:" {
      val, numerr := strconv.ParseUint(fields[1], 10, 64)
      if numerr != nil {
        return 0, errors.WithStack(numerr)
      }
      return val * Multi, nil
    }
  }

  return 0, errors.New("not found MemTotal info from file /proc/meminfo")
}
