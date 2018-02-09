package controller

import (
  "fmt"
  "os"

  "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/control/ctx"
)

var (
  lg *logger.Logger
  log Log
)

type Log uint64

func LogInit() {
  logConfig := ctx.Cfg().LogConfig
  if logConfig != nil {
    logConfig.Withprefix = true
    logConfig.LogFlag = logger.LstdFlags | logger.Lshortfile
    logConfig.AddCallDepth = 2
  } else {
    logConfig = &logger.LoggerConf{
      ConsoleAppender: true,
      Withprefix: true,
      LogFlag: logger.LstdFlags | logger.Lshortfile,
      AddCallDepth: 2,
    }
  }

  var err error
  lg, err = logger.NewLogger(logConfig)
  if err != nil {
    fmt.Printf("%+v\n", err)
    os.Exit(1)
  }
}

func withRequestId(v ...interface{}) ([]interface{}) {
  rId := GetRequestId()
  if rId == "" {
    return v
  } else {
    v = append(v, fmt.Sprintf("(%s)", rId))
    return v
  }
}

func output(lvl string, v ...interface{}) {
  if lg != nil {
    v = withRequestId(v...)
    switch lvl {
    case "Debug":
      lg.Debug(v...)
    case "Info":
      lg.Info(v...)
    case "Warn":
      lg.Warn(v...)
    case "Error":
      lg.Error(v...)
    }
  } else {
    fmt.Println(v...)
  }
}

func (l *Log) Debug(v ...interface{}) {
  output("Debug", v...)
}

func (l *Log) Info(v ...interface{}) {
  output("Info", v...)
}

func (l *Log) Warn(v ...interface{}) {
  output("Warn", v...)
}

func (l *Log) Error(v ...interface{}) {
  output("Error", v...)
}

func (l *Log) Printf(format string, v ...interface{}) {
  if lg != nil && len(format) > 0 {
    lenS := len(v)
    v = withRequestId(v...)
    if lenS + 1 == len(v) {
      if format[len(format)-1] == '\n' {
        format = format[:len(format)-1] + "%s\n"
      } else {
        format += "%s"
      }
    }
    lg.Printf(format, v...)
  } else {
    fmt.Printf(format, v...)
  }
}

func (l *Log) Println(v ...interface{}) {
  output("Debug", v...)
}
