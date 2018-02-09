package redis

import (
  "time"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/go-redis/redis"
  "github.com/cloudtropy/open-operation/modules/control/ctx"
)

var (
  client *redis.Client
  UserCookieKey = "CONSOLE_USER_SESSION_"
  UserCookieExpire  time.Duration
  UCRememberExpire  time.Duration
)

func Init() {
  if client != nil {
    return
  }

  UserCookieExpire = time.Minute * time.Duration(ctx.Cfg().User.LoginTimeout)
  UCRememberExpire = time.Hour * time.Duration(ctx.Cfg().User.RememberTimeout * 24)

  client = redis.NewClient(&redis.Options{
    Addr: ctx.Cfg().Redis.Host,
    Password: "",
    DB: 0,
    PoolSize: 50,
  })

  if pong, err := client.Ping().Result(); err != nil || pong != "PONG" {
    log.Fatal("redis init error:", err.Error())
  } else {
    log.Println("redis init successfully.")
  }
}
