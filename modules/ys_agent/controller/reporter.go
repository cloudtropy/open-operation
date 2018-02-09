package controller

import (
  "bytes"
  "errors"
  "os/exec"
  "strconv"
  "strings"
  "sync"
  "time"

  "github.com/cloudtropy/open-operation/modules/ys_agent/ctx"
  "github.com/cloudtropy/open-operation/modules/ys_agent/funcs"
  cc "github.com/cloudtropy/open-operation/utils/common"
  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/utils/rpc"
  "github.com/cloudtropy/open-operation/utils/ticker"
)

type MetricFunc func() []*cc.MetricValue

var MFunc = map[string]MetricFunc{
  "cpu.load.1min":         funcs.LoadAvgMetrics,
  "cpu.load.1min.percent": funcs.LoadAvgPercentMetrics,
  "cpu.load.5min":         funcs.LoadAvgMetrics,
  "cpu.load.15min":        funcs.LoadAvgMetrics,
  "mem.free":              funcs.MemMetrics,
  "mem.used":              funcs.MemMetrics,
  "mem.used.percent":      funcs.MemMetrics,
  "net.if.in":             funcs.NetMetrics,
  "net.if.out":            funcs.NetMetrics,
  "net.if.total":          funcs.NetMetrics,
  "df.free":               funcs.DeviceMetrics,
  "df.used":               funcs.DeviceMetrics,
  "df.used.percent":       funcs.DeviceMetrics,
}

type ItemScheduler struct {
  Ticker *ticker.Ticker
  Item   cc.HostItem
}

func NewItemScheduler(item cc.HostItem) *ItemScheduler {
  is := ItemScheduler{Item: item}
  is.Ticker = ticker.NewTicker(time.Duration(item.Interval) * time.Second)
  return &is
}

func (is *ItemScheduler) Stop() {
  is.Ticker.StopTimer()
}

func (is *ItemScheduler) Schedule() {
  go func() {
    // var errCount = 0
    for {
      if !is.Ticker.Running() {
        break
      }

      var metricValue *cc.MetricValue
      if is.Item.Creator == "system" || is.Item.Creator == "born" {
        metricValues := MFunc[is.Item.ItemName]()
        for _, mv := range metricValues {
          if mv.Metric == is.Item.ItemName {
            metricValue = mv
            break
          }
        }
      } else {
        // todo
        mv, err := GetCommandMetric(is.Item)
        if err != nil {
          log.Println("GetCommandMetric", err)
        }
        metricValue = mv
      }

      if metricValue != nil {
        metricValue.Type = is.Item.Dst
        metricValue.Step = is.Item.Interval
        metricValue.Endpoint = ctx.HostId()
        metricValue.Timestamp = time.Now().Unix()
        metricValue.DataType = is.Item.DataType
        if is.Item.History == 0 {
          metricValue.Tags = "NoHistory"
        }

        var res string
        err := rpcMonitor.Call("Agent.HandleReportData", *metricValue, &res)
        if err != nil {
          log.Println("rpcMonitor.Call Agent.HandleReportData", err)
        }
      } /* else if errCount++; errCount > 5 {
         log.Println("Too many error of item report:", is.Item)
         is.Ticker.StopTimer()
       }*/

      <-is.Ticker.C
    }
  }()
}

var (
  rpcMonitor       *rpc.SingleConnRpcClient
  Items            = make(map[string]cc.HostItem)
  ItemSchedulers   = make(map[string]*ItemScheduler)
  UpdateMu         sync.Mutex
  handleItemsTimer *time.Timer
  cmdTimeout       = 3
)

func StartMonitorReport() {
  rpcMonitor = &rpc.SingleConnRpcClient{
    RpcServer: ctx.Cfg().MonitorSrv.RpcHost,
    Timeout:   time.Second * 30,
  }

  go HandleItems()
}

func HandleItems() {
  if handleItemsTimer != nil {
    handleItemsTimer.Stop()
    handleItemsTimer = nil
  }

  items := make([]cc.HostItem, 0)
  err := rpcMonitor.Call("Agent.HostItems", ctx.HostId(), &items)
  if err != nil {
    log.Error("rpcMonitor.Call Agent.HostItems", err)
    handleItemsTimer = time.AfterFunc(60*time.Second, HandleItems)
    return
  }

  UpdateMu.Lock()
  defer UpdateMu.Unlock()

  desiredItems := make(map[string]cc.HostItem, 0)
  for _, item := range items {
    _, ok := MFunc[item.ItemName]
    if !ok /*&& item.Creator == "system"*/ {
      log.Info("Item", item.ItemName, "has not handle function.")
      continue
    } else if item.Interval <= 0 {
      continue
    } else if item.Dst != "GAUGE" && item.Dst != "COUNTER" {
      continue
    } else if (item.Creator != "system" && item.Creator != "born") && item.Command == "" {
      continue
    }
    desiredItems[item.ItemName] = item
  }

  DelNoUseItems(desiredItems)
  AddNewItems(desiredItems)
}

func DelNoUseItems(newItems map[string]cc.HostItem) {
  for itemName, item := range Items {
    newItem, ok := newItems[itemName]
    if !ok || item.Timestamp != newItem.Timestamp {
      delItem(itemName)
    }
  }
}

func AddNewItems(newItems map[string]cc.HostItem) {
  for itemName, newItem := range newItems {
    item, ok := Items[itemName]
    if ok && item.Timestamp == newItem.Timestamp {
      continue
    }

    Items[itemName] = newItem
    itemScheduler := NewItemScheduler(newItem)
    ItemSchedulers[itemName] = itemScheduler
    itemScheduler.Schedule()
  }
}

func delItem(itemName string) {
  itemScheduler, ok := ItemSchedulers[itemName]
  if ok {
    itemScheduler.Stop()
    delete(ItemSchedulers, itemName)
  }
  delete(Items, itemName)
}

func GetCommandMetric(item cc.HostItem) (*cc.MetricValue, error) {

  cmd := exec.Command("sh", "-c", item.Command)
  var stdout bytes.Buffer
  cmd.Stdout = &stdout
  var stderr bytes.Buffer
  cmd.Stderr = &stderr
  cmd.Start()

  err, isTimeout := CmdRunWithTimeout(cmd, time.Duration(cmdTimeout)*time.Second)
  if isTimeout {
    return nil, errors.New("CommandTimeout")
  } else if err != nil {
    return nil, errors.New("CommandError:" + err.Error())
  }

  errStr := stderr.String()
  if errStr != "" {
    log.Println("exec.Command", errStr)
    return nil, errors.New("CommandError:" + errStr)
  }

  var value interface{}
  valueStr := strings.TrimSpace(stdout.String())
  if item.DataType == "float" {
    valueF, err := strconv.ParseFloat(valueStr, 64)
    if err != nil {
      log.Println("strconv.ParseFloat", item, valueStr)
      return nil, errors.New("InvalidDataType")
    }
    value = valueF
  } else {
    value = valueStr
  }

  return &cc.MetricValue{
    Metric: item.ItemName,
    Value:  value,
  }, nil
}

func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
  done := make(chan error)
  go func() {
    done <- cmd.Wait()
  }()

  var err error
  select {
  case <-time.After(timeout):
    log.Printf("timeout, process:%s will be killed", cmd.Path)

    go func() {
      <-done // allow goroutine to exit
    }()

    // timeout
    if err = cmd.Process.Kill(); err != nil {
      log.Printf("failed to kill: %s, error: %s", cmd.Path, err)
    }

    return err, true
  case err = <-done:
    return err, false
  }
}
