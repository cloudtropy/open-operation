package rrdtool

import (
  "errors"
  "os"
  "path"

  "github.com/cloudtropy/open-operation/modules/monitor/ctx"
  "github.com/cloudtropy/open-operation/utils/file"
)

const (
  step      = 60
  heartbeat = step + 5
)

func GetDbFilePath(dbType string, fileName string) string {
  return path.Join(ctx.Cfg().RrdPath, dbType, fileName+".rrd")
}

func GetDbFileDir(dbType string) string {
  return path.Join(ctx.Cfg().RrdPath, dbType)
}

func CreateDbFilePath(dbFile string) error {
  if file.PathIsExist(dbFile) {
    return errors.New("exist")
  } else if !file.PathIsExist(path.Dir(dbFile)) {
    if err := os.MkdirAll(path.Dir(dbFile), 0777); err != nil {
      return err
    }
  }
  return nil
}
