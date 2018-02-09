package controller

import (
  "time"
  // "strings"
  "net/http"

  "github.com/cloudtropy/open-operation/modules/control/db/redis"
  // "github.com/cloudtropy/open-operation/modules/control/db/mysql"
)

type MidHandleFunc func(http.HandlerFunc) http.HandlerFunc

/*
 midHandlers: [f1, f2, f3], 在handler之前执行
 执行顺序：f1 f2 f3 依次执行
 */
func HandleWithMid(path string, handler http.HandlerFunc, midHandlers ...MidHandleFunc) {
  tmpHandler := handler

  for i := len(midHandlers) - 1; i >= 0; i-- {
    tmpHandler = midHandlers[i](tmpHandler)
  }

  tmpHandler = HandleActionTrail(tmpHandler)

  http.HandleFunc(path, tmpHandler)
}


func AuthSessionCheck(handler http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

    c, err := r.Cookie("sid")
    if err != nil {
      HttpResMsg(w, "NotLogin", "")
      return
    }

    ttl, err := redis.ExpireUserSession(c.Value)
    if err != nil {
      HttpResMsg(w, "RedisError", err.Error())
      return
    }
    if ttl < 0 {
      SetResCookieKV(w, "sid", c.Value, -1)
      HttpResMsg(w, "NotLogin", "")
      return
    }

    if ttl <= redis.UserCookieExpire {
      SetResCookieKV(w, "sid", c.Value, int(redis.UserCookieExpire / time.Second) - 2)
    }

    handler(w, r)
  })
}

func HandleActionTrail(handler http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    NewHttpActionTrailInfo()

    handler(w, r)

    if r.Method == http.MethodPost {
      // todo: record in mysql
    }

    DelHttpActionTrailInfo()
  })
}
