package mysql

import (
  "database/sql"

  "github.com/pkg/errors"
)

type TriggerInfo struct {
  Id            int64       `json:"id"`
  Name          string      `json:"name"`
  Severity      string      `json:"severity"`
  NoticeMessage string      `json:"noticeMessage"`
  NoticePerson  string      `json:"noticePerson"`
  NoticeBy      string      `json:"noticeBy"`
  TemplateId    int64       `json:"templateId"`
  ItemId        int64       `json:"itemId"`
  RuleTime      string      `json:"ruleTime"`
  RuleType      string      `json:"ruleType"`
  RuleOperator  string      `json:"ruleOperator"`
  RuleValue     string      `json:"ruleValue"`
  Enabled       int         `json:"enabled"`
}

func InsertTriggerData(t *TriggerInfo) (int64, error) {

  prepareClause := `insert into ops_trigger(name, severity, notice_message, 
    notice_person, notice_by, template_id, item_id, rule_time, rule_type,
    rule_operator, rule_value) values(?,?,?,?,?,?,?,?,?,?,?)`

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return -1, errors.WithStack(err)
  }
  res, err := stmt.Exec(t.Name, t.Severity, t.NoticeMessage,
    t.NoticePerson, t.NoticeBy, t.TemplateId, t.ItemId,
    t.RuleTime, t.RuleType, t.RuleOperator, t.RuleValue)
  if err != nil {
    return -1, errors.WithStack(err)
  }

  lastId, err := res.LastInsertId()
  if err != nil {
    return -1, errors.WithStack(err)
  }

  return lastId, nil
}

func DeleteTriggerData(name string, templateId int64) error {

  prepareClause := `delete from ops_trigger where name = ? and 
    template_id = ?`

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return errors.WithStack(err)
  }
  _, err = stmt.Exec(name, templateId)
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}

func GetTriggerData(template_id int64, name string) (*TriggerInfo, error) {
  queryStr := `select id, name, severity, notice_message, 
    notice_person, notice_by, template_id, item_id, rule_time, 
    rule_type, rule_operator, rule_value, enabled from ops_trigger 
    where template_id=? and name=?;`

  var t TriggerInfo
  err := db.QueryRow(queryStr, template_id, name).Scan(&t.Id,
    &t.Name, &t.Severity, &t.NoticeMessage, &t.NoticePerson,
    &t.NoticeBy, &t.TemplateId, &t.ItemId, &t.RuleTime,
    &t.RuleType, &t.RuleOperator, &t.RuleValue, &t.Enabled)

  if err != nil {
    if err == sql.ErrNoRows {
      err = nil
    }
    return nil, errors.WithStack(err)
  }
  return &t, errors.WithStack(err)
}

func GetAllTriggers() (map[int64]*TriggerInfo, error) {
  queryStr := `select id, name, severity, notice_message, 
    notice_person, notice_by, template_id, item_id, rule_time, 
    rule_type, rule_operator, rule_value, enabled from ops_trigger;`

  rows, err := db.Query(queryStr)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make(map[int64]*TriggerInfo)
  for rows.Next() {
    var t TriggerInfo
    err = rows.Scan(&t.Id,
      &t.Name, &t.Severity, &t.NoticeMessage, &t.NoticePerson,
      &t.NoticeBy, &t.TemplateId, &t.ItemId, &t.RuleTime,
      &t.RuleType, &t.RuleOperator, &t.RuleValue, &t.Enabled)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res[t.Id] = &t
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func UpdateTriggerById(t *TriggerInfo) error {
  queryStr := `update ops_trigger set name=?, severity=?, notice_message=?, 
    notice_person=?, notice_by=?, item_id=?, rule_time=?, rule_type=?, 
    rule_operator=?, rule_value=?, enabled=? where id=?;`
  stmt, err := db.Prepare(queryStr)
  if err != nil {
    return errors.WithStack(err)
  }

  _, err = stmt.Exec(t.Name, t.Severity, t.NoticeMessage,
    t.NoticePerson, t.NoticeBy, t.ItemId, t.RuleTime, t.RuleType, 
    t.RuleOperator, t.RuleValue, t.Enabled, t.Id)
  return errors.WithStack(err)
}

func QueryTriggersList(template, hostId string) ([]map[string]string, error) {
  var queryStr, arg string
  if hostId != "" {
    queryStr = `select a.name, a.severity, a.notice_person, a.notice_by, a.enabled, 
      a.rule_time, a.rule_type, a.rule_operator, a.rule_value, b.name, c.name from 
      ops_trigger a, ops_item b, ops_template c where a.item_id=b.id 
      and a.template_id=c.id and c.id in (
      select d.template_id from ops_job_template d, ops_host_server e
      where d.job_id=e.job_id and e.host_id=?);`
    arg = hostId
  } else {
    queryStr = `select a.name, a.severity, a.notice_person, a.notice_by, a.enabled, 
      a.rule_time, a.rule_type, a.rule_operator, a.rule_value, b.name, c.name from 
      ops_trigger a, ops_item b, ops_template c where a.item_id=b.id 
      and a.template_id=c.id and c.name=?;`
    arg = template
  }

  rows, err := db.Query(queryStr, arg)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]string, 0)
  for rows.Next() {
    var triggerName, severity, noticePerson, noticeBy, ruleTime,
      ruleType, ruleOperator, ruleValue, itemName, enabled, templateName string
    err = rows.Scan(&triggerName, &severity, &noticePerson, &noticeBy, &enabled,
      &ruleTime, &ruleType, &ruleOperator, &ruleValue, &itemName, &templateName)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    res = append(res, map[string]string{
      "triggerName":  triggerName,
      "severity":     severity,
      "noticePerson": noticePerson,
      "noticeBy":     noticeBy,
      "itemName":     itemName,
      "ruleTime":     ruleTime,
      "ruleType":     ruleType,
      "ruleOperator": ruleOperator,
      "ruleValue":    ruleValue,
      "enabled":      enabled,
      "templateName": templateName,
    })
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}


/*

 ops_host_server job_id host_id
 ops_job_template job_id template_id
 ops_trigger template_id item_id
 ops_item name interval id
 ops_host host_id ip
 */

func GetHostTriggers(triggerId int64) ([]map[string]string, error) {
  var rows *sql.Rows
  var err error
  if triggerId == -1 {
    queryStr := `SELECT a.host_id, e.ip, c.id, d.name, d.dst, d.interval 
      FROM ops_host_server a, ops_job_template b,
      ops_trigger c, ops_item d, ops_host e WHERE a.job_id=b.job_id 
      AND a.host_id=e.host_id AND b.template_id=c.template_id AND 
      c.item_id=d.id;`
    rows, err = db.Query(queryStr)
  } else {
    queryStr := `SELECT a.host_id, e.ip, c.id, d.name, d.dst, d.interval 
      FROM ops_host_server a, ops_job_template b,
      ops_trigger c, ops_item d, ops_host e WHERE a.job_id=b.job_id 
      AND a.host_id=e.host_id AND b.template_id=c.template_id AND 
      c.item_id=d.id AND c.id=?;`
    rows, err = db.Query(queryStr, triggerId)
  }

  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]string, 0)
  for rows.Next() {
    var hostId, hostIp, triggerId, itemName, dst, interval string
    err = rows.Scan(&hostId, &hostIp, &triggerId, &itemName, &dst, &interval)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, map[string]string{
      "hostId":     hostId,
      "hostIp":     hostIp,
      "triggerId":  triggerId,
      "itemName":   itemName,
      "interval":   interval,
      "dst":        dst,
    })
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func GetHostTriggersOfHostId(hostId string) ([]map[string]string, error) {
  queryStr := `SELECT distinct(c.id), a.host_id, e.ip, d.name, d.dst, d.interval 
    FROM ops_host_server a, ops_job_template b,
    ops_trigger c, ops_item d, ops_host e WHERE a.host_id=? AND 
    a.job_id=b.job_id AND a.host_id=e.host_id AND b.template_id=c.template_id 
    AND c.item_id=d.id;`

  rows, err := db.Query(queryStr, hostId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]string, 0)
  for rows.Next() {
    var hostId, hostIp, triggerId, itemName, dst, interval string
    err = rows.Scan(&triggerId, &hostId, &hostIp, &itemName, &dst, &interval)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, map[string]string{
      "hostId":     hostId,
      "hostIp":     hostIp,
      "triggerId":  triggerId,
      "itemName":   itemName,
      "interval":   interval,
      "dst":        dst,
    })
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func GetTriggerIdsOfTemplateId(tId int64) ([]int64, error) {
  rows, err := db.Query("SELECT id from ops_trigger where template_id=?;", tId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]int64, 0)
  for rows.Next() {
    var id int64
    err = rows.Scan(&id)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, id)
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func GetTriggersOfHostId(hostId string) (map[int64]string, error) {
  queryStr := `SELECT c.id, c.name 
    FROM ops_host_server a, ops_job_template b, ops_trigger c 
    WHERE a.host_id=? AND a.job_id=b.job_id AND b.template_id=c.template_id;`
  rows, err := db.Query(queryStr, hostId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make(map[int64]string)
  for rows.Next() {
    var id int64
    var name string
    err = rows.Scan(&id, &name)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res[id] = name
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}