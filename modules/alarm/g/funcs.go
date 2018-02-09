package g

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

func NowTimestampS() int64 {
	return time.Now().Unix()
}

func GetMd5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%X", h.Sum(nil))
}
