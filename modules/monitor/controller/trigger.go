package controller

import (
  "encoding/json"
  "io/ioutil"
  "strings"
  "net/http"
  "database/sql"
  "regexp"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
)

func HandleTrigger(w http.ResponseWriter, r *http.Request) {

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

  var ctrlBody map[string]string
  err = json.Unmarshal(body, &ctrlBody)
  if err != nil {
    log.Error("json.Unmarshal", err)
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  switch ctrlBody["topic"] {
  case "addTrigger":
    addTrigger(w, ctrlBody)
  case "delTrigger":
    delTrigger(w, ctrlBody)
  case "updateTrigger":
    updateTrigger(w, ctrlBody)
  case "queryTrigger":
    queryTrigger(w, ctrlBody)
  case "queryTriggers":
    queryTriggers(w, ctrlBody)
  default:
    HttpResMsg(w, "InvalidRequestParams", "topic")
  }
}

func addTrigger(w http.ResponseWriter, body map[string]string) {

  if body["templateName"] == "" || body["triggerName"] == "" ||
    body["severity"] == "" || body["noticePerson"] == "" ||
    body["noticeBy"] == "" || body["ruleItem"] == "" ||
    body["ruleTime"] == "" || body["ruleType"] == "" ||
    body["ruleOperator"] == "" || body["ruleValue"] == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  match, _ := regexp.MatchString(`^[1-9]\d{0,2}(分钟)?$`, body["ruleTime"])
  if !match {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  if body["ruleItem"] == "host.heartbeat" {
    HttpResMsg(w, "InvalidRequestParams", "Item host.heartbeat will not be used by users.")
    return
  }

  templateId, err := mysql.QueryIdByName(body["templateName"], "ops_template")
  if err != nil {
    if err == sql.ErrNoRows {
      HttpResMsg(w, "NotFound", "template")
      return
    }
    log.Error("mysql.QueryIdByName", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  itemId, err := mysql.QueryIdByName(body["ruleItem"], "ops_item")
  if err != nil {
    if err == sql.ErrNoRows {
      HttpResMsg(w, "NotFound", "item")
      return
    }
    log.Error("mysql.QueryIdByName", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  triggerInfo := &mysql.TriggerInfo{
    Name:          body["triggerName"],
    Severity:      body["severity"],
    NoticeMessage: body["noticeMessage"],
    NoticePerson:  body["noticePerson"],
    NoticeBy:      body["noticeBy"],
    TemplateId:    templateId,
    ItemId:        itemId,
    RuleTime:      body["ruleTime"],
    RuleType:      body["ruleType"],
    RuleOperator:  body["ruleOperator"],
    RuleValue:     body["ruleValue"],
  }

  // tId, err := mysql.InsertTriggerData(triggerInfo)
  _, err = mysql.InsertTriggerData(triggerInfo)
  if err != nil {
    log.Printf("mysql.InsertTriggerData(%v) error: %s.\n", *triggerInfo, err.Error())
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success")
  // todo
  // UpdateTriggers(tId, body["topic"])
}

func delTrigger(w http.ResponseWriter, body map[string]string) {

  if body["triggerName"] == "" || body["templateName"] == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  // The trigger named 'HostHeartbeat' of template 'Basic' will not be deleted.
  if body["triggerName"] == "HostHeartbeat" && body["templateName"] == "Basic" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  templateId, err := mysql.QueryIdByName(body["templateName"], "ops_template")
  if err != nil {
    if err == sql.ErrNoRows {
      HttpResMsg(w, "NotFound", "template:" + body["templateName"])
      return
    }
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  // tId, err := mysql.QueryIdByName(body["triggerName"], "ops_trigger")
  _, err = mysql.QueryIdByName(body["triggerName"], "ops_trigger")
  if err != nil {
    if err == sql.ErrNoRows {
      HttpResMsg(w, "NotFound", "trigger:" + body["triggerName"])
      return
    }
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  err = mysql.DeleteTriggerData(body["triggerName"], templateId)
  if err != nil {
    log.Error("mysql.DeleteTriggerData", body["triggerName"], templateId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success")

  // todo
  // UpdateTriggers(tId, body["topic"])
}

func updateTrigger(w http.ResponseWriter, body map[string]string) {
  if body["templateName"] == "" || body["triggerName"] == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  } else if tmpStr := body["newTriggerName"] + body["severity"] + body["noticeMessage"] +
    body["noticePerson"] + body["noticeBy"] + body["ruleItem"] + body["ruleTime"] +
    body["ruleType"] + body["ruleOperator"] + body["ruleValue"] + body["enabled"]; tmpStr == "" {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }
  if body["ruleTime"] != "" {
    match, _ := regexp.MatchString(`^[1-9]\d{0,2}(分钟)?$`, body["ruleTime"])
    if !match {
      HttpResMsg(w, "InvalidRequestParams")
      return
    }
  }

  templateId, err := mysql.QueryIdByName(body["templateName"], "ops_template")
  if err != nil {
    if err == sql.ErrNoRows {
      HttpResMsg(w, "NotFound", "template:" + body["templateName"])
      return
    }
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  t, err := mysql.GetTriggerData(templateId, body["triggerName"])
  if err != nil {
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if t == nil {
    HttpResMsg(w, "NotFound", "trigger:" + body["triggerName"])
    return
  }

  var needUpdate = false

  if body["ruleItem"] != "" {
    itemId, err := mysql.QueryIdByName(body["ruleItem"], "ops_item")
    if err != nil {
      if err == sql.ErrNoRows {
        HttpResMsg(w, "NotFound", "item:" + body["ruleItem"])
        return
      }
      log.Error("mysql.QueryIdByName", body["ruleItem"], "ops_item", err)
      HttpResMsg(w, "MysqlError", err.Error())
      return
    }
    if t.ItemId != itemId {
      needUpdate = true
      t.ItemId = itemId
    }
  }

  if body["newTriggerName"] != "" && body["newTriggerName"] != t.Name {
    needUpdate = true
    t.Name = body["newTriggerName"]
  }
  if body["severity"] != "" && body["severity"] != t.Severity {
    needUpdate = true
    t.Severity = body["severity"]
  }
  if body["noticeMessage"] != "" && body["noticeMessage"] != t.NoticeMessage {
    needUpdate = true
    t.NoticeMessage = body["noticeMessage"]
  }
  if body["noticePerson"] != "" && body["noticePerson"] != t.NoticePerson {
    needUpdate = true
    t.NoticePerson = body["noticePerson"]
  }
  if body["noticeBy"] != "" && body["noticeBy"] != t.NoticeBy {
    needUpdate = true
    t.NoticeBy = body["noticeBy"]
  }
  if body["ruleTime"] != "" && body["ruleTime"] != t.RuleTime {
    needUpdate = true
    t.RuleTime = body["ruleTime"]
  }
  if body["ruleType"] != "" && body["ruleType"] != t.RuleType {
    needUpdate = true
    t.RuleType = body["ruleType"]
  }
  if body["ruleOperator"] != "" && body["ruleOperator"] != t.RuleOperator {
    needUpdate = true
    t.RuleOperator = body["ruleOperator"]
  }
  if body["ruleValue"] != "" && body["ruleValue"] != t.RuleValue {
    needUpdate = true
    t.RuleValue = body["ruleValue"]
  }
  if body["enabled"] == "0" && t.Enabled != 0 {
    needUpdate = true
    t.Enabled = 0
  } else if body["enabled"] == "1" && t.Enabled != 1 {
    needUpdate = true
    t.Enabled = 1
  }

  if !needUpdate {
    HttpResMsg(w, "Success")
    return
  }

  err = mysql.UpdateTriggerById(t)
  if err != nil {
    log.Error("mysql.UpdateTriggerById", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success")

  // todo
  // UpdateTriggers(t.Id, body["topic"])
}

func queryTriggers(w http.ResponseWriter, body map[string]string) {
  if (body["templateName"] == "" && body["hostId"] == "") ||
    (body["templateName"] != "" && body["hostId"] != "") {
    HttpResMsg(w, "InvalidRequestParams")
    return
  }

  triggersInfo, err := mysql.QueryTriggersList(body["templateName"], body["hostId"])
  if err != nil {
    log.Error("mysql.QueryTriggersList", body["templateName"], body["hostId"], err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  var templateHostCount int = -1
  if body["templateName"] != "" {
    templateHostCount, err = mysql.QueryHostCountOfTemplate(body["templateName"])
    if err != nil {
      log.Error("mysql.QueryHostCountOfTemplate", body["templateName"], err)
      HttpResMsg(w, "MysqlError", err.Error())
      return
    }
  }

  resBody := make([]map[string]interface{}, 0)
  for i, t := range triggersInfo {
    hostCount := templateHostCount
    if hostCount == -1 {
      hostCount, err = mysql.QueryHostCountOfTemplate(t["templateName"])
      if err != nil {
        log.Error("mysql.QueryHostCountOfTemplate", t["templateName"], err)
        HttpResMsg(w, "MysqlError", err.Error())
        return
      }
    }

    rule := map[string]string{
      "ruleTime":     t["ruleTime"],
      "ruleType":     t["ruleType"],
      "ruleOperator": t["ruleOperator"],
      "ruleValue":    t["ruleValue"],
    }

    resBody = append(resBody, map[string]interface{}{
      "key":          i,
      "triggerName":  t["triggerName"],
      "severity":     t["severity"],
      "itemName":     t["itemName"],
      "hostCount":    hostCount,
      "rule":         rule,
      "enabled":      t["enabled"],
      "noticePerson": t["noticePerson"],
      "noticeBy":     t["noticeBy"],
    })
  }

  HttpResData(w, resBody)
}

func queryTrigger(w http.ResponseWriter, body map[string]string) {

  templateId, err := mysql.QueryIdByName(body["templateName"], "ops_template")
  if err != nil {
    if err == sql.ErrNoRows {
      HttpResMsg(w, "NotFound", "template:" + body["templateName"])
      return
    }
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  t, err := mysql.GetTriggerData(templateId, body["triggerName"])
  if err != nil {
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if t == nil {
    HttpResMsg(w, "NotFound", "trigger:" + body["triggerName"])
    return
  }

  item, err := mysql.GetItemById(t.ItemId)
  if err != nil {
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  ruleWay := "time"
  if index := strings.Index(t.RuleTime, "分钟"); index == -1 {
    ruleWay = "count"
  }

  res := map[string]interface{}{
    "triggerId":      t.Id,
    "triggerName":    t.Name,
    "severity":       t.Severity,
    "noticeMessage":  t.NoticeMessage,
    "noticePerson":   t.NoticePerson,
    "noticeBy":       t.NoticeBy,
    "ruleItem":       item.Name,
    "ruleTime":       t.RuleTime,
    "ruleType":       t.RuleType,
    "ruleOperator":   t.RuleOperator,
    "ruleValue":      t.RuleValue,
    "ruleWay":        ruleWay,
    "enabled":        t.Enabled,
  }

  HttpResData(w, res)
}
