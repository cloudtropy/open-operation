package controller

import (
  "net/http"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
)



func GetDetailItems(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  items, err := mysql.GetItems("born")
  if err != nil {
    log.Error("mysql.GetItems", "born", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  res := make([]map[string]interface{}, 0)
  for _, item := range items {
    templates, err := mysql.GetTemplatesOfItem(item.Name)
    if err != nil {
      log.Error("mysql.GetTemplatesOfItem", item.Name, err)
      HttpResMsg(w, "MysqlError", err.Error())
      return
    }
    res = append(res, map[string]interface{}{
      "itemName":    item.Name,
      "dataType":    item.DataType,
      "unit":        item.Unit,
      "interval":    item.Interval,
      "history":     item.History,
      "dst":         item.Dst,
      "description": item.Description,
      "creator":     item.Creator,
      "createTime":  item.CreateTime,
      "templates":   templates,
    })
  }

  HttpResData(w, res)
}

func GetItemList(w http.ResponseWriter, r *http.Request) {

  if r.Method != "GET" {
    http.NotFound(w, r)
    return
  }

  r.ParseForm()
  templateName := r.FormValue("template")
  hostId := r.FormValue("hostId")
  if hostId != "" && templateName != "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  itemList, err := mysql.QueryItemData(templateName, hostId)
  if err != nil {
    log.Error("mysql.QueryItemData", templateName, hostId, err)
    HttpResMsg(w, "MysqlError", err.Error())
  } else {
    HttpResData(w, itemList)
  }
}
