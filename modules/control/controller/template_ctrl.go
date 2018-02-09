package controller

import (
  "encoding/json"
  "io/ioutil"
  "net/http"

  cc "github.com/cloudtropy/open-operation/utils/common"
)

func TemplateCtrl(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }
  defer r.Body.Close()

  var f struct {
    Topic string `json:"topic"`
  }
  err = json.Unmarshal(body, &f)
  if err != nil {
    log.Warn("json.Unmarshal", err)
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }
  switch f.Topic {
  case "addTemplate":
    addTemplate(w, body)
  case "delTemplate":
    delTemplate(w, body)
  case "updateTemplate":
    updateTemplate(w, body)
  case "queryTemplate":
    queryTemplate(w, body)
  case "queryTemplates":
    queryTemplates(w)
  default:
    HttpResMsg(w, "InvalidRequestParams", "topic")
    return
  }
}

func addTemplate(w http.ResponseWriter, body []byte) {
  args := cc.AddTemplateArgs{}
  var templateId int64
  err := json.Unmarshal(body, &args)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  err = rpcMonitor.Call("Template.AddTemplate", args, &templateId)
  if err != nil {
    log.Info("rpcMonitor.Call", "Template.AddTemplate", err)
    HttpResMsg(w, err.Error())
    return
  }
  HttpResData(w, map[string]int64{"templateId": templateId})
}

func delTemplate(w http.ResponseWriter, body []byte) {
  var f struct {
    TemplateId int64 `json:"templateId"`
  }
  err := json.Unmarshal(body, &f)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }
  var res string
  err = rpcMonitor.Call("Template.DelTemplate", f.TemplateId, &res)
  if err != nil {
    log.Info("rpcMonitor.Call", "Template.DelTemplate", err)
    HttpResMsg(w, err.Error())
    return
  }
  HttpResMsg(w, res)
}

func updateTemplate(w http.ResponseWriter, body []byte) {
  args := cc.UpdateTemplate{}
  err := json.Unmarshal(body, &args)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }
  var res string
  err = rpcMonitor.Call("Template.UpdateTemplate", args, &res)
  if err != nil {
    log.Info("rpcMonitor.Call", "Template.UpdateTemplate", err)
    HttpResMsg(w, err.Error())
    return
  }
  HttpResMsg(w, res)
}

func queryTemplate(w http.ResponseWriter, body []byte) {
  var f struct {
    TemplateId int64 `json:"templateId"`
  }
  err := json.Unmarshal(body, &f)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }
  res := cc.QueryTemplateRes{}
  err = rpcMonitor.Call("Template.QueryTemplate", f.TemplateId, &res)
  if err != nil {
    log.Info("rpcMonitor.Call", "Template.QueryTemplate", err)
    HttpResMsg(w, err.Error())
    return
  }
  HttpResData(w, res)
}

func queryTemplates(w http.ResponseWriter) {
  var args string
  res := make([]cc.QueryTemplateInfo, 0)
  err := rpcMonitor.Call("Template.QueryTemplates", args, &res)
  if err != nil {
    log.Info("rpcMonitor.Call", "Template.QueryTemplatesQueryTemplates", err)
    HttpResMsg(w, err.Error())
    return
  }
  HttpResData(w, res)
}
