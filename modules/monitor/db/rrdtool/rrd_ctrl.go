package rrdtool

import (
  "math"
  "os"
  "path"
  "strings"
  "sync"
  "time"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/ctx"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
  "github.com/cloudtropy/open-operation/utils/file"
  "github.com/ziutek/rrd"
  "github.com/pkg/errors"
)

type RrdUpdater struct {
  updater *rrd.Updater
  lock    sync.RWMutex
}

type RrdData struct {
  Timestamp int64       `json:"timestamp"`
  Value     interface{} `json:"value"`
}

var (
  updaters   = make(map[string]*RrdUpdater)
  updatersMu = new(sync.RWMutex)
)

func setUpdatersMap(f string, u *RrdUpdater) {
  updatersMu.Lock()
  defer updatersMu.Unlock()
  updaters[f] = u
}

func getUpdater(f string) (*RrdUpdater, bool) {
  updatersMu.RLock()
  defer updatersMu.RUnlock()
  u, isExist := updaters[f]
  return u, isExist
}

func closeUpdater(f string) {
  updatersMu.Lock()
  defer updatersMu.Unlock()
  delete(updaters, f)
}

func GetRrdFilePath(hostId, metric string) string {
  return path.Join(ctx.Cfg().RrdPath, hostId, metric+".rrd")
}

func GetRrdFileDir(hostId string) string {
  return path.Join(ctx.Cfg().RrdPath, hostId)
}

func CreateRrdFile(hostId, metric string) error {
  filePath := GetRrdFilePath(hostId, metric)
  if file.PathIsExist(filePath) {
    // To solve the problem when rrd file size is 0.
    f, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
    if err != nil {
      log.Println("os.OpenFile", filePath, err)
      return err
    }
    defer f.Close()
    fileInfo, err := f.Stat()
    if err != nil {
      log.Println("f.Stat", err)
      return err
    }
    if fileInfo.Size() > 0 {
      return nil
    } else {
      if err = os.Remove(filePath); err != nil {
        log.Println("os.Remove", filePath, err)
        return err
      }
    }

  } else if !file.PathIsExist(GetRrdFileDir(hostId)) {
    if err := os.MkdirAll(GetRrdFileDir(hostId), 0777); err != nil {
      return err
    }
  }

  itemInfo, err := mysql.GetItemByName(metric)
  if err != nil {
    return err
  }
  if itemInfo.History == 0 {
    return errors.New("InvalidItem:item neednt to record in rrd.")
  }
  step := itemInfo.Interval
  dst := itemInfo.Dst
  cdp := int64(itemInfo.History * 24 * 60 * 60 / step)
  // dsName := strings.Replace(metric, ".", "_", -1)

  c := rrd.NewCreator(filePath, time.Unix(time.Now().Unix()-step, 0), uint(step))
  c.DS("metric", dst, step*2, 0, "U")
  c.RRA("MAX", 0.5, 1, cdp)

  return c.Create(false)
}

func UpdateData(hostId, metric string, step int64, value interface{}) error {
  filePath := GetRrdFilePath(hostId, metric)
  rrdUpdater, isExist := getUpdater(filePath)
  if !isExist {
    u := RrdUpdater{}
    // lock to avoid two creating at the same time
    u.lock.Lock()
    defer u.lock.Unlock()
    err := CreateRrdFile(hostId, metric)
    if err != nil {
      return err
    }

    u.updater = rrd.NewUpdater(filePath)
    rrdUpdater = &u
    setUpdatersMap(filePath, &u)
  }

  var updateTime = adjustTimestamp(time.Now().Unix(), step)
  err := rrdUpdater.updater.Update(updateTime, value)
  if err != nil && strings.Index(err.Error(), "last update time") == -1 {
    closeUpdater(filePath)
  }
  return err
}

func FetchData(hostId, metric string, start, end int64) ([]*RrdData, error) {
  startTime := time.Unix(start, 0)
  endTime := time.Unix(end, 0)
  filePath := GetRrdFilePath(hostId, metric)

  fetchRes, err := rrd.Fetch(filePath, "MAX", startTime, endTime, 1*time.Second)
  if err != nil {
    return nil, err
  }
  defer fetchRes.FreeValues()

  row := 0
  dataSlice := make([]*RrdData, 0)
  for ti := fetchRes.Start.Add(fetchRes.Step); ti.Before(endTime) || ti.Equal(endTime); ti = ti.Add(fetchRes.Step) {
    for i := 0; i < len(fetchRes.DsNames); i++ {
      if v := fetchRes.ValueAt(i, row); !math.IsNaN(v) {
        rrdData := RrdData{
          Timestamp: ti.Unix(),
          Value:     v,
        }
        dataSlice = append(dataSlice, &rrdData)
      } else {
        rrdData := RrdData{
          Timestamp: ti.Unix(),
          Value:     -1,
        }
        dataSlice = append(dataSlice, &rrdData)
      }
    }
    row++
  }
  return dataSlice, nil
}

func adjustTimestamp(timestamp, interval int64) time.Time {
  var resT int64
  if interval == 0 {
    return time.Unix(0, 0)
  }
  if remainder := timestamp % interval; remainder >= interval/2 {
    resT = timestamp - remainder + interval
  } else {
    resT = timestamp - remainder
  }

  return time.Unix(resT, 0)
}
