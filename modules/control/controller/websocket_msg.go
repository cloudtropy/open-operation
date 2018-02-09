package controller

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "sync"

  "github.com/gorilla/websocket"
)

type Msg struct {
  Topic string      `json:"topic"`
  Data  interface{} `json:"data"`
}

var (
  WsMsgChans   = make(map[*websocket.Conn]chan Msg)
  WsMsgChansMu = new(sync.RWMutex)
  upgrader     = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
  }
)

func SetWsMsgChan(c *websocket.Conn, mc chan Msg) {
  WsMsgChansMu.Lock()
  defer WsMsgChansMu.Unlock()
  WsMsgChans[c] = mc
}

func DelWsMsgChan(c *websocket.Conn) {
  WsMsgChansMu.Lock()
  defer WsMsgChansMu.Unlock()
  delete(WsMsgChans, c)
}

func AddWsMsg(msg Msg) {
  WsMsgChansMu.RLock()
  defer WsMsgChansMu.RUnlock()
  chanCount := 0
  for _, v := range WsMsgChans {
    chanCount++
    go func(v chan Msg) {
      v <- msg
    }(v)
  }
  log.Println("Added WsMsg", chanCount, msg.Topic, msg.Data)
}

func HandleWsMsg(w http.ResponseWriter, r *http.Request) {
  c, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println("upgrade:", err)
    return
  }
  msgChan := make(chan Msg, 5)
  SetWsMsgChan(c, msgChan)
  defer c.Close()
  for msg := range msgChan {
    //msg := <-msgChan
    bmsg, err := json.Marshal(msg)
    if err != nil {
      log.Println("json.Marshal:", err)
      continue
    }

    err = c.WriteMessage(websocket.TextMessage, bmsg)
    if err != nil {
      log.Println("websocket write:", err)
      DelWsMsgChan(c)
      break
    }
  }
}

func HandleReport(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }
  if len(body) < 2 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  var msg Msg
  err = json.Unmarshal(body, &msg)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  if msg.Topic != "" {
    AddWsMsg(msg)
    HttpResMsg(w, "Success")
  } else {
    HttpResMsg(w, "InvalidRequestParams", "topic")
  }
}
