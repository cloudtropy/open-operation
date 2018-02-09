package redis

import (
  "strconv"
  
  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/utils/fun"
  "github.com/pkg/errors"
)


/*
 return: bool: insert is true, update is false
 */
func UpsertHostStatus(key, field string) (bool, error) {
  br, err := client.HSet(key, field, fun.NowTimestampS()).Result()
  return br, errors.WithStack(err)
}


/*
 return: int64, the count of deleted fields
 */
func DelHostStatus(key string, fields ...string) (int64, error) {
  ir, err := client.HDel(key, fields...).Result()
  return ir, errors.WithStack(err)
}


/*
 return: int64, the timestamp 
 */
func GetOneHostStatus(key, field string) (int64, error) {
  stringCmd := client.HGet(key, field)
  if stringCmd.Err() != nil {
    return 0, errors.WithStack(stringCmd.Err())
  }
  return stringCmd.Int64()
}


func GetAllHostStatus(key string) (map[string]int64, error) {

  mapSS, err := client.HGetAll(key).Result()
  if err != nil {
    return nil, errors.WithStack(err)
  }

  mapSI := make(map[string]int64)
  for f, v := range mapSS {
    i64Time, err := strconv.ParseInt(v, 10, 64)
    if err != nil {
      log.Debug("Redis hdel key:", key, ",field:", f, ",value:", v, err)
      DelHostStatus(key, f)
      continue
    }
    mapSI[f] = i64Time
  }

  return mapSI, nil
}
