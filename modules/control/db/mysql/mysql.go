package mysql

import (
  "database/sql"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/control/ctx"
  _ "github.com/go-sql-driver/mysql"
)

var (
  db       *sql.DB
)

func Init() {
  if db != nil {
    return
  }
  sqlInfo := ctx.Cfg().Mysql
  if sqlInfo == nil {
    log.Fatal("mysql init error: no mysql config")
  }

  dataSourceName := sqlInfo.Username + ":" + sqlInfo.Passwd + "@tcp(" + sqlInfo.Host + ")/" + sqlInfo.Database
  var err error
  db, err = sql.Open("mysql", dataSourceName)
  if err != nil {
    log.Fatal("mysql init sql.Open:", err)
  }

  err = db.Ping()
  if err != nil {
    log.Fatal("mysql init db.Ping:", err)
  } else {
    log.Println("mysql init successfully.")
  }
}
