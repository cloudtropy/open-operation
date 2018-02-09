package funcs

import (
  "strings"
  "os"
  "os/exec"
  "regexp"
  "net"
  "time"

  "github.com/cloudtropy/open-operation/modules/ys_agent/ctx"
  cc "github.com/cloudtropy/open-operation/utils/common"
  log "github.com/cloudtropy/open-operation/utils/logger"
)

var Ip string


func OsMetrics() []*cc.MetricValue {
  cmd := exec.Command("sh", "-c", "cat /etc/redhat-release")
  out, err := cmd.CombinedOutput()
  if err != nil {
    log.Error(err)
    return nil
  }

  operatingSystem := strings.TrimSpace(string(out))
  return []*cc.MetricValue{
    GaugeValue("host.os", operatingSystem),
  }
}


func IpMetrics() []*cc.MetricValue {
  if len(Ip) > 0 {
    return []*cc.MetricValue{
      GaugeValue("host.ip", Ip),
    }
  }

  shelverSrvAddrs := ctx.Cfg().MonitorSrv.Addrs
  re := regexp.MustCompile("^https?://")

  for _, addr := range shelverSrvAddrs {
    addr = re.ReplaceAllString(addr, "")
    conn, err := net.DialTimeout("tcp", addr, time.Second*3)
    if err != nil {
      log.Error("get local ip failed ! shelver_srv:", addr, err)
    } else {
      Ip = strings.Split(conn.LocalAddr().String(), ":")[0]
      conn.Close()
      break
    }
  }

  if len(Ip) > 0 {
    return []*cc.MetricValue{
      GaugeValue("host.ip", Ip),
    }
  } else {
    return nil
  }
}


func HostnameMetrics() []*cc.MetricValue {
  hostname, err := os.Hostname()
  if err != nil {
    log.Error(err)
    return nil
  } else {
    return []*cc.MetricValue{
      GaugeValue("host.hostname", hostname),
    }
  }
}
