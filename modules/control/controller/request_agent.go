package controller

import (
  "net/http"
  "strings"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "fmt"

  "github.com/cloudtropy/open-operation/modules/control/ctx"
  "github.com/tidwall/gjson"
)


func HandleMonitor(w http.ResponseWriter, r *http.Request) {
  requestAgent(w, r, ctx.Cfg().Monitor.Host, nil)
}

func HandleMonitorWithUser(w http.ResponseWriter, r *http.Request) {
  requestAgent(w, r, ctx.Cfg().Monitor.Host, map[string]string{
    "username": GetUsername(r),
  })
}

func HandleMonitorWithRequestId(w http.ResponseWriter, r *http.Request) {
  requestAgent(w, r, ctx.Cfg().Monitor.Host, map[string]string{
    "requestId": GetRequestId(),
  })
}

func HandleMonitorWithUR(w http.ResponseWriter, r *http.Request) {
  requestAgent(w, r, ctx.Cfg().Monitor.Host, map[string]string{
    "username": GetUsername(r),
    "requestId": GetRequestId(),
  })
}

func requestAgent(w http.ResponseWriter, r *http.Request, host string, adds map[string]string) {

  addr := host + strings.Replace(r.RequestURI, "/api", "", 1)
  hasAdds := false
  for _, v := range adds {
    if v != "" {
      hasAdds = true
    }
  }

  client := &http.Client{}
  var request *http.Request
  var err error

  if r.Method == http.MethodGet {
    if hasAdds {
      arrAdds := make([]string, 0)
      for k, v := range adds {
        arrAdds = append(arrAdds, k + "=" + v)
      }
      strAdds := strings.Join(arrAdds, "&")
      if strings.Index(addr, "?") == -1 {
        addr += "?" + strAdds
      } else {
        addr += "&" + strAdds
      }
    }
    request, err = http.NewRequest(r.Method, addr, nil)

  } else if r.Method == http.MethodPost || r.Method == http.MethodPut {
    if hasAdds {
      bs, err := ioutil.ReadAll(r.Body)
      if err != nil {
        log.Error("ioutil.ReadAll", err)
        HttpResMsg(w, "InternalError", err.Error())
        return
      }

      var mss map[string]interface{}
      err = json.Unmarshal(bs, &mss)
      if err != nil {
        log.Error("json.Unmarshal", err)
        HttpResMsg(w, "InvalidRequestParams", err.Error())
        return
      }

      for k, v := range adds {
        mss[k] = v
      }

      bsBody, err := json.Marshal(mss)
      if err != nil {
        log.Error("json.Marshal", err)
        HttpResMsg(w, "InvalidRequestParams", err.Error())
        return
      }

      request, err = http.NewRequest(r.Method, addr, bytes.NewBuffer(bsBody))
    } else {
      request, err = http.NewRequest(r.Method, addr, r.Body)
    }

    request.Header.Set("Content-Type", "application/json; charset=utf-8")
  } else {
    HttpResMsg(w, "InvalidHttpMethod")
    return
  }

  resp, err := client.Do(request)
  if err != nil {
    log.Error("Http Request:", err)
    HttpResMsg(w, "InternalError", err.Error())
    return
  } else if resp.StatusCode != http.StatusOK {
    log.Error("Error status code:", resp.StatusCode, addr)
    HttpResMsg(w, "InternalError", fmt.Sprintf("%d", resp.StatusCode))
    return
  }

  defer resp.Body.Close()
  result, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Error("ioutil.ReadAll", err)
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  code := gjson.GetBytes(result, "code").String()
  if code != "Success" {
    msg := gjson.GetBytes(result, "msg").String()
    HttpResMsg(w, code, msg)
    return
  }

  data := gjson.GetBytes(result, "data").Value()
  HttpResData(w, data)
}
