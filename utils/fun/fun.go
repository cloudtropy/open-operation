package fun

import (
  "time"
  "strings"
  "crypto/md5"
  "io"
  "fmt"
  "regexp"

  "github.com/satori/go.uuid"
)


func NowTimestampS() int64 {
  return time.Now().Unix()
}

func NewRequestId() string {
  u4 := uuid.NewV4()
  return strings.Replace(u4.String(), "-", "", -1)
}

func GetUUIDV4() string {
  u4 := uuid.NewV4()
  return strings.Replace(u4.String(), "-", "", -1)
}

func GetMd5(s string) string {
  h := md5.New()
  io.WriteString(h, s)
  return fmt.Sprintf("%X", h.Sum(nil))
}

func IndexOfS(s []string, v string) int {
  for i, ts := range s {
    if ts == v {
      return i
    }
  }
  return -1
}

func IsSqlId(id string) bool {
  match, _ := regexp.MatchString(`^[1-9]\d{0,8}$`, id)
  return match
}

