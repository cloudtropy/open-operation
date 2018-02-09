package controller

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strconv"
  "strings"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
  "github.com/cloudtropy/open-operation/modules/monitor/db/rrdtool"
)

type ScreenCtrlBody struct {
  Topic         string    `json:"topic"`
  ScreenName    string    `json:"screenName"`
  ScreenId      int64     `json:"screenId"`
}

func HandleScreen(w http.ResponseWriter, r *http.Request) {

  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  var ctrlBody ScreenCtrlBody
  err = json.Unmarshal(body, &ctrlBody)
  if err != nil {
    log.Error("json.Unmarshal", err)
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  switch ctrlBody.Topic {
  case "addScreen":
    addScreen(w, ctrlBody)
  case "delScreen":
    delScreen(w, ctrlBody)
  case "updateScreen":
    updateScreen(w, ctrlBody)
  case "queryScreen":
    queryScreen(w)
  default:
    HttpResMsg(w, "InvalidRequestParams", "topic")
  }
}

func GetScreenGraphData(w http.ResponseWriter, r *http.Request) {

  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }
  r.ParseForm()

  screenId := r.FormValue("screenId")
  start := r.FormValue("start")
  end := r.FormValue("end")

  if screenId == "" || start == "" || end == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  intScreenId, err0 := strconv.ParseInt(screenId, 10, 64)
  intStart, err1 := strconv.ParseInt(start, 10, 64)
  intEnd, err2 := strconv.ParseInt(end, 10, 64)

  if err0 != nil || err1 != nil || err2 != nil {
    HttpResMsg(w, "InvalidRequestParams")
    return
  } else if (intEnd - intStart) < 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  } else if (intEnd - intStart) > 604800 {
    HttpResMsg(w, "InvalidRequestParams", "Time interval too long")
    return
  }

  graphDetails, err := mysql.QueryScreenGraph(intScreenId)
  if err != nil {
    log.Error("mysql.QueryScreenGraph", intScreenId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  for _, graphMetric := range graphDetails {
    // graphId, err := strconv.ParseInt()
    itemMetric, err := mysql.QueryGraphItem(graphMetric["id"].(int64))
    if err != nil {
      log.Error("mysql.QueryGraphItem", err)
      HttpResMsg(w, "MysqlError", err.Error())
      return
    }

    for _, item := range itemMetric {
      itemData, err := rrdtool.FetchData(graphMetric["hostId"].(string), item, intStart, intEnd)
      if err != nil {
        if strings.Index(err.Error(), "No such file") != -1 {
          continue
        }
        log.Error("rrdtool.FetchData", graphMetric["hostId"].(string), item, intStart, intEnd, err)
        HttpResMsg(w, "InternalError", "Fetch rrdtool data error")
        return
      }
      graphMetric = FormatGraphItemData(graphMetric, item, itemData)
      // log.Printf("%#v\n", graphMetric)
    }

  }
  // log.Printf("%#v\n", graphDetails)
  HttpResData(w, graphDetails)
}

func addScreen(w http.ResponseWriter, body ScreenCtrlBody) {

  if body.ScreenName == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  screenInfo := &mysql.ScreenInfo{
    Name: body.ScreenName,
  }

  _, err := mysql.InsertScreenData(screenInfo)
  if err != nil {
    log.Error("mysql.InsertScreenData", screenInfo, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success")
}

func delScreen(w http.ResponseWriter, body ScreenCtrlBody) {

  if body.ScreenId < 1 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  err := mysql.DeleteScreenData(body.ScreenId)
  if err != nil {
    log.Error("mysql.DeleteScreenData", body.ScreenId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success")
}

func updateScreen(w http.ResponseWriter, body ScreenCtrlBody) {
  if body.ScreenId < 1 || body.ScreenName == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  screenInfo := &mysql.ScreenInfo{
    Id :  body.ScreenId,
    Name: body.ScreenName,
  }

  err := mysql.UpdateScreenData(screenInfo)
  if err != nil {
    log.Error("mysql.UpdateScreenData", screenInfo, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }
  HttpResMsg(w, "Success")
}

func queryScreen(w http.ResponseWriter) {

  allScreens, err := mysql.QueryScreenList()
  if err != nil {
    log.Error("mysql.QueryScreenList", err)
    HttpResMsg(w, "MysqlError", err.Error())
  } else {
    HttpResData(w, allScreens)
  }
}

func HandPreviewGraph(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  var f struct {
    HostId string   `json:"hostId"`
    Start  string   `json:"start"`
    End    string   `json:"end"`
    Items  []string `json:"items"`
  }

  err = json.Unmarshal(body, &f)
  if err != nil {
    log.Error("json.Unmarshal", err)
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  intStart, err1 := strconv.ParseInt(f.Start, 10, 64)
  intEnd, err2 := strconv.ParseInt(f.End, 10, 64)

  if err1 != nil || err2 != nil {
    HttpResMsg(w, "InvalidRequestParams")
    return
  } else if (intEnd - intStart) < 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  } else if (intEnd - intStart) > 604800 {
    HttpResMsg(w, "InvalidRequestParams", "Time interval too long")
    return
  }

  graphMetric := make(map[string]interface{})

  for _, item := range f.Items {
    itemData, err := rrdtool.FetchData(f.HostId, item, intStart, intEnd)
    if err != nil {
      log.Printf("%#v\n", err)
      if strings.Index(err.Error(), "No such file") != -1 {
        continue
      }
      log.Error("rrdtool.FetchData", f.HostId, item, intStart, intEnd, err)
      HttpResMsg(w, "InternalError", "Fetch rrdtool data error")
      return
    }
    graphMetric = FormatSingleItemData(graphMetric, item, itemData)
  }

  // // log.Println("%#v\n", graphMetric)
  HttpResData(w, graphMetric)
}
