package controller

import (
  "bytes"
  "io/ioutil"
  "net/http"
  "regexp"
  "strings"

  "github.com/cloudtropy/open-operation/modules/control/db/mysql"
  "github.com/cloudtropy/open-operation/utils/fun"
)

var pathnameActions = map[string][]string{
  "/api/authority": {"CreateUser", "DeleteUser", "ResetPassword", "GetUserInfo", "UpdateUserInfo",
    "GetUserList", "GetActionTrailList"},
  "/api/basic/monitor": {"GetHostMonitorData", "GetHostList", "GetRemovedHostList", "UpdateHostInfo",
    "AddMonitorTemplate", "DeleteMonitorTemplate", "UpdateMonitorTemplate",
    "GetMonitorTemplateInfo", "GetMonitorTemplateList", "AddMonitorGraph",
    "DeleteMonitorGraph", "UpdateMonitorGraph", "GetMonitorGraphInfo",
    "AddMonitorScreen", "DeleteMonitorScreen", "UpdateMonitorScreen",
    "GetMonitorScreenInfo", "AddMonitorTrigger", "DeleteMonitorTrigger",
    "UpdateMonitorTrigger", "GetMonitorTriggerInfo", "GetMonitorTriggerList",
    "GetMonitorItemDetailList", "GetMonitorGraphPreviewData", "GetMonitorScreenData",
    "GetMonitorHostConfigList", "GetMonitorItemList", "GetMonitorItemsOfHost",
    "GetMonitorItemsForGraph"},
}

var pathnameModule = map[string]string{
  "/api/authority":     "访问控制",
  "/api/basic/monitor": "基础监控",
}

type ActionInfo struct {
  Alias   string
  Method  string
  ActFunc http.HandlerFunc
  Path    string
}

var Actions = map[string]*ActionInfo{
  "CreateUser":         {"创建用户", "POST", CreateUser, ""},
  "DeleteUser":         {"删除用户", "POST", DeleteUser, ""},
  "ResetPassword":      {"重置密码", "POST", ResetUserPassword, ""},
  "GetUserInfo":        {"获取用户信息", "GET", GetUserInfo, ""},
  "GetUserList":        {"获取用户信息列表", "GET", GetUserList, ""},
  "UpdateUserInfo":     {"更新用户信息", "POST", UpdateUser, ""},
  "GetActionTrailList": {"", "GET", GetActionTrailList, ""},


  "GetHostList":                {"获取所有未下架主机信息列表", "GET", HandleMonitor, "/api/host/list"},
  "GetRemovedHostList":         {"获取下架主机信息列表", "GET", HandleMonitor, "/api/removed/host/list"},
  "GetHostMonitorData":         {"获取主机基础监控数据", "GET", HandleMonitor, "/api/host/monitor/data"},
  "UpdateHostInfo":             {"更新主机信息", "POST", HandleMonitor, "/api/host/info/update"},

  "AddMonitorTemplate":         {"添加监控模板", "POST", TemplateCtrl, ""},
  "DeleteMonitorTemplate":      {"删除监控模板", "POST", TemplateCtrl, ""},
  "UpdateMonitorTemplate":      {"更新监控模板", "POST", TemplateCtrl, ""},
  "GetMonitorTemplateInfo":     {"获取监控模板信息", "GET", TemplateCtrl, ""},
  "GetMonitorTemplateList":     {"获取监控模板信息列表", "GET", TemplateCtrl, ""},

  "AddMonitorGraph":            {"添加监控图表", "POST", HandleMonitor, "/api/graph/ctrl"},
  "DeleteMonitorGraph":         {"删除监控图表", "POST", HandleMonitor, "/api/graph/ctrl"},
  "UpdateMonitorGraph":         {"更新监控图表", "POST", HandleMonitor, "/api/graph/ctrl"},
  "GetMonitorGraphInfo":        {"获取监控图表信息", "GET", HandleMonitor, "/api/graph/ctrl"},
  "AddMonitorScreen":           {"添加监控图表面板", "POST", HandleMonitor, "/api/screen/ctrl"},
  "DeleteMonitorScreen":        {"删除监控图表面板", "POST", HandleMonitor, "/api/screen/ctrl"},
  "UpdateMonitorScreen":        {"更新监控图表面板信息", "POST", HandleMonitor, "/api/screen/ctrl"},
  "GetMonitorScreenInfo":       {"获取监控图表面板信息", "GET", HandleMonitor, "/api/screen/ctrl"},
  "GetMonitorGraphPreviewData": {"获取监控图表的预览数据", "GET", HandleMonitor, "/api/preview/graph"},
  "GetMonitorScreenData":       {"获取监控图表面板的数据", "GET", HandleMonitor, "/api/screen/graphs"},

  "AddMonitorTrigger":          {"添加监控告警规则", "POST", HandleMonitor, "/api/trigger/ctrl"},
  "DeleteMonitorTrigger":       {"删除监控告警规则", "POST", HandleMonitor, "/api/trigger/ctrl"},
  "UpdateMonitorTrigger":       {"更新监控告警规则", "POST", HandleMonitor, "/api/trigger/ctrl"},
  "GetMonitorTriggerInfo":      {"获取监控告警规则信息", "GET", HandleMonitor, "/api/trigger/ctrl"},
  "GetMonitorTriggerList":      {"获取监控告警规则信息列表", "GET", HandleMonitor, "/api/trigger/ctrl"},

  "GetMonitorItemDetailList":   {"获取监控项详情列表", "GET", HandleMonitor, "/api/detail/items"},

  // "GetMonitorHostConfigList":   {"获取主机监控配置列表", "GET", HandleMonitor, "/api/config/hosts"},

  "GetMonitorItemList":         {"获取监控项列表", "GET", HandleMonitor, "/api/items"},
  "GetMonitorItemsOfHost":      {"获取主机的监控项列表", "GET", HandleMonitor, "/api/host/items"},
  "GetMonitorItemsForGraph":    {"获取图表可用监控项列表", "GET", HandleMonitor, "/api/graph/valid/items"},
}

func HandleAction(w http.ResponseWriter, r *http.Request) {
  user := GetUsername(r)
  if user == "" {
    HttpResMsg(w, "NotLogin")
    return
  }

  pathname := getUrlPathname(r.RequestURI)
  pathActions, isExist := pathnameActions[pathname]
  if !isExist {
    http.NotFound(w, r)
    return
  }

  act := getUrlQueryParam(r.RequestURI, "action")
  if act == "" || Actions[act] == nil {
    HttpResMsg(w, "InvalidRequestParams", "action")
    return
  } else if fun.IndexOfS(pathActions, act) == -1 {
    HttpResMsg(w, "InvalidRequestParams", "action")
    return
  }

  // if Actions[act].Method != r.Method {
  //   HttpResMsg(w, "InvalidHttpMethod", r.Method)
  //   return
  // }

  if Actions[act].Path != "" {
    r.RequestURI = Actions[act].Path + getUrlQueryParams(r.RequestURI)
  }

  var body []byte
  if Actions[act].Method == http.MethodPost {
    body, _ = ioutil.ReadAll(r.Body)
    if len(body) > 0 {
      resetRequestBody(r, body)
    }
  }

  Actions[act].ActFunc(w, r)

  if Actions[act].Method != http.MethodPost {
    return
  }

  ati := GetHttpActionTrailInfo()
  if ati == nil {
    log.Warn("Not found action trail info.")
    return
  }
  at := &mysql.ActionTrail{
    Id:     ati.RequestId,
    User:   user,
    Module: pathnameModule[pathname],
    Action: act,
    Result: ati.Result,
  }
  detailArr := []string{
    "RequestId: " + GetRequestId(),
    "Action: " + act,
    "Body: " + string(body),
    "Result: " + ati.Result,
  }
  at.Detail = "-----------------\n" + strings.Join(detailArr, "\n\n-----------------\n")
  err := mysql.AddActionTrail(at)
  if err != nil {
    log.Error("mysql.AddActionTrail", at, err)
  }
}

func getUrlPathname(path string) string {
  return path[:strings.Index(path, "?")]
}

func getUrlQueryParams(path string) string {
  querys := path[strings.Index(path, "?"):]
  re := regexp.MustCompile("(\\?|&)action=([^&]+)")
  trimedQuerys := re.ReplaceAllString(querys, "")
  if len(trimedQuerys) == 0 {
    return ""
  }
  if trimedQuerys[0] == '&' {
    trimedQuerys = "?" + trimedQuerys[1:]
  }
  return trimedQuerys
}

func getUrlQueryParam(path, key string) string {
  re := regexp.MustCompile("(\\?|&)" + key + "=([^&]+)")
  resArr := re.FindStringSubmatch(path)
  if resArr == nil {
    return ""
  } else {
    return resArr[2]
  }
}

func resetRequestBody(r *http.Request, body []byte) {
  r.Body = ioutil.NopCloser(bytes.NewReader(body))
}
