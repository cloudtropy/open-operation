package controller

import (
  "encoding/json"
  "io/ioutil"
  "net/http"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/utils/fun"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
  "github.com/cloudtropy/open-operation/modules/monitor/db/rrdtool"
)

type GraphCtrl struct {
  Topic     string   `json:"topic"`
  GraphId   int64    `json:"graphId"`
  GraphName string   `json:"graphName"`
  ScreenId  int64    `json:"screenId"`
  HostId    string   `json:"hostId"`
  GraphType string   `json:"graphType"`
  Items     []string `json:"items"`
}

func HandleGraph(w http.ResponseWriter, r *http.Request) {

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

  var ctrlBody GraphCtrl
  err = json.Unmarshal(body, &ctrlBody)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  switch ctrlBody.Topic {
  case "addGraph":
    addGraph(w, ctrlBody)
  case "delGraph":
    delGraph(w, ctrlBody)
  case "updateGraph":
    updateGraph(w, ctrlBody)
  case "queryGraph":
    queryGraph(w, ctrlBody)
  default:
    HttpResMsg(w, "InvalidRequestParams")
  }
}

func addGraph(w http.ResponseWriter, body GraphCtrl) {

  if body.GraphName == "" || body.HostId == "" ||
    body.GraphType == "" || body.ScreenId <= 0 ||
    len(body.Items) == 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  // todo: some item can not be used by graph.
  // The item host.heartbeat can not be used by graph.
  if index := fun.IndexOfS(body.Items, "host.heartbeat"); index != -1 {
    HttpResMsg(w, "InvalidRequestParams", "Item host.heartbeat can't be used by graph.")
    return
  }

  gId, err := mysql.InsertGraphData(body.GraphName, body.GraphType, body.HostId, body.ScreenId)
  if err != nil {
    log.Error("mysql.InsertGraphData", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  for _, item := range body.Items {

    itemId, err := mysql.QueryIdByName(item, "ops_item")
    if err != nil {
      // if err == sql.ErrNoRows {
      //   HttpResMsg(w, "MysqlNotFound")
      //   return
      // }
      log.Error("mysql.QueryIdByName", item, "ops_item", err)
      HttpResMsg(w, "MysqlError", err.Error())
      return
    }

    _, err = mysql.InsertGraphItemData(gId, itemId)
    if err != nil {
      log.Printf("mysql.InsertGraphItemData(%v) error: %s.\n", body, err.Error())
      HttpResMsg(w, "MysqlError")
      return
    }
  }

  HttpResMsg(w, "Success")
}


func delGraph(w http.ResponseWriter, body GraphCtrl) {

  if body.GraphId <= 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  err := mysql.DeleteGraphById(body.GraphId)
  if err != nil {
    log.Error("mysql.DeleteGraphById", body.GraphId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success")
}

func updateGraph(w http.ResponseWriter, body GraphCtrl) {
  // todo: some item can not be used by graph.
  // The item host.heartbeat can not be used by graph.
  if index := fun.IndexOfS(body.Items, "host.heartbeat"); index != -1 {
    HttpResMsg(w, "InvalidRequestParams", "Item host.heartbeat can't be used by graph.")
    return
  }

  err := mysql.UpdateGraph(body.GraphId, body.GraphName, body.GraphType)
  if err != nil {
    log.Error("mysql.UpdateGraph", body.GraphId, body.GraphName, body.GraphType, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  items, err := mysql.QueryGraphItem(body.GraphId)
  if err != nil {
    log.Error("mysql.QueryGraphItem", body.GraphId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  itemsMap := make(map[string]string)
  for _, item := range items {
    itemsMap[item] = "placeholder"
  }

  itemsToAdd := make([]string, 0)
  for _, item := range body.Items {
    if itemsMap[item] != "" {
      delete(itemsMap, item)
      continue
    }
    itemsToAdd = append(itemsToAdd, item)
  }

  itemsToDel := make([]string, 0)
  for item, _ := range itemsMap {
    itemsToDel = append(itemsToDel, item)
  }

  if len(itemsToAdd) == 0 && len(itemsToDel) == 0 {
    HttpResMsg(w, "Success")
    return
  }

  err = mysql.UpdateGraphItems(body.GraphId, body.HostId, itemsToAdd, itemsToDel)
  if err != nil {
    log.Println("mysql.UpdateGraphItems", body.GraphId, itemsToAdd, itemsToDel, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success")
}

func queryGraph(w http.ResponseWriter, body GraphCtrl) {
  if body.GraphId <= 0 {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  graph, err := mysql.QueryGraph(body.GraphId)
  if err != nil {
    log.Println("mysql.QueryGraph", body.GraphId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  items, err := mysql.QueryGraphItem(body.GraphId)
  if err != nil {
    log.Println("mysql.QueryGraphItem", body.GraphId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  graph["items"] = items
  HttpResData(w, graph)
}


func FormatGraphItemData(graphDetail map[string]interface{}, item string, itemData []*rrdtool.RrdData) (map[string]interface{}) {
  
  if graphDetail["data"] == nil {
    // log.Println("data nil")
    graphData := make(map[string]interface{}, 0)
    timeStampData := make([]int64, 0)
    for _, data := range itemData {
      timeStampData = append(timeStampData, data.Timestamp)
    }
    graphData["timestamp"] = timeStampData
    graphDetail["data"] = graphData
  } 

  dataSlice := make([]interface{}, 0)
  for _, data := range itemData{
    dataSlice = append(dataSlice, data.Value)
  }
  graphData := graphDetail["data"].(map[string]interface{})
  graphData[item] = dataSlice

  return graphDetail
}


func FormatSingleItemData(graphDetail map[string]interface{}, item string, itemData []*rrdtool.RrdData) (map[string]interface{}) {
  
  if graphDetail["timestamp"] == nil {
    
    timeStampData := make([]int64, 0)
    for _, data := range itemData {
      timeStampData = append(timeStampData, data.Timestamp)
    }

    graphDetail["timestamp"] = timeStampData
    
  }

  dataSlice := make([]interface{}, 0)
  for _, data := range itemData{
    dataSlice = append(dataSlice, data.Value)
  }
  // graphData := graphDetail["data"].(map[string]interface{})
  graphDetail[item] = dataSlice

  return graphDetail
}
