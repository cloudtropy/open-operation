package redis

import (
  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/ctx"
  "github.com/go-redis/redis"
)

var (
  client *redis.Client
)

func Init() {
  if client != nil {
    return
  }

  client = redis.NewClient(&redis.Options{
    Addr: ctx.Cfg().Redis.Host,
    Password: "",
    DB: 0,
    PoolSize: 50,
  })

  if pong, err := client.Ping().Result(); err != nil || pong != "PONG" {
    log.Fatal("redis init error:", err.Error())
  } else {
    log.Info("redis init successfully.")
  }
}
