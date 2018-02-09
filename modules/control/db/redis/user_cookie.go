package redis

import (
  "time"
  "encoding/json"

  "github.com/pkg/errors"
)


func AddUserSession(sid string, remember bool, v interface{}) error {
  bs, err := json.Marshal(v)
  if err != nil {
    return errors.WithStack(err)
  }
  expire := UserCookieExpire
  if remember {
    expire = UCRememberExpire
  }
  err = client.Set(UserCookieKey + sid, bs, expire).Err()
  return errors.WithStack(err)
}


func ExpireUserSession(sid string) (time.Duration, error) {
  ttl, err := client.TTL(UserCookieKey + sid).Result()
  if err != nil {
    return ttl, errors.WithStack(err)
  }
  if ttl >= UserCookieExpire || ttl < 0 {
    return ttl, nil
  }
  _, err = client.Expire(UserCookieKey + sid, UserCookieExpire).Result()
  if err != nil {
    return ttl, errors.WithStack(err)
  } else {
    return ttl, nil
  }
}


func DelUserSession(sid string) error {
  err := client.Del(UserCookieKey + sid).Err()
  return errors.WithStack(err)
}

func GetUserSession(sid string) (map[string]string, error) {
  bs, err := client.Get(UserCookieKey + sid).Bytes()
  if err != nil {
    return nil, errors.WithStack(err)
  }
  var mss map[string]string
  err = json.Unmarshal(bs, &mss)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return mss, nil
}
