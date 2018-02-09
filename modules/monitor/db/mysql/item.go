package mysql

import (
  "strconv"
  "database/sql"

  "github.com/pkg/errors"
)

type ItemInfo struct {
  Name      string `json:"name"`
  AliasName string `json:"alias"`
  Unit      string `json:"unit"`
  Creator   string `json:"creator"`
  Interval  uint64 `json:"interval"`
  History   uint64 `json:"history"`
  Key       int    `json:"key"`
}

type Item struct {
  Id          int64  `json:"id"`
  Name        string `json:"name"`
  Interval    int64  `json:"interval"`
  History     int64  `json:"history"`
  Description string `json:"description"`
  DataType    string `json:"dataType"`
  Unit        string `json:"unit"`
  Dst         string `json:"dst"`
  Creator     string `json:"creator"`
  Command     string `json:"command"`
  CreateTime  string `json:"createTime"`
}

/*
 creator: system / born / [user]
 */
func GetItems(creator string) ([]*Item, error) {

  queryStr := "select id,name,`interval`,history,description," +
    "data_type,unit,dst,creator,create_time from ops_item where creator=?;"

  rows, err := db.Query(queryStr, creator)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]*Item, 0)
  for rows.Next() {
    var item = Item{}
    err = rows.Scan(&item.Id, &item.Name, &item.Interval,
      &item.History, &item.Description, &item.DataType,
      &item.Unit, &item.Dst, &item.Creator, &item.CreateTime)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, &item)
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func GetItemByName(name string) (*Item, error) {
  var item = Item{}
  err := db.QueryRow("select id,name,`interval`,history,description,command,"+
    "data_type,unit,dst,creator,create_time from ops_item where name=?;", name).Scan(&item.Id,
    &item.Name, &item.Interval, &item.History, &item.Description, &item.Command,
    &item.DataType, &item.Unit, &item.Dst, &item.Creator, &item.CreateTime)
  return &item, errors.WithStack(err)
}

func GetItemById(id int64) (*Item, error) {
  var item = Item{}
  err := db.QueryRow("select id,name,`interval`,history,description,command,"+
    "data_type,unit,dst,creator,create_time from ops_item where id=?;", id).Scan(&item.Id,
    &item.Name, &item.Interval, &item.History, &item.Description, &item.Command,
    &item.DataType, &item.Unit, &item.Dst, &item.Creator, &item.CreateTime)
  return &item, errors.WithStack(err)
}

func GetTemplatesOfItem(name string) ([]map[string]interface{}, error) {
  queryStr := `SELECT id,name FROM ops_template WHERE id IN 
    (SELECT a.template_id FROM ops_template_item a, ops_item b
      WHERE a.item_id=b.id AND b.name=?);`
  rows, err := db.Query(queryStr, name)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]interface{}, 0)
  for rows.Next() {
    var templateName string
    var templateId int64
    err = rows.Scan(&templateId, &templateName)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, map[string]interface{}{
      "templateId":   templateId,
      "templateName": templateName,
    })
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}


func QueryItemData(templateName, hostId string) (itemList []*ItemInfo, err error) {

  var rows *sql.Rows
  if templateName != "" {
    query_string := `select id, name, a.interval, history, unit, creator from ops_item a where id in 
      (select a.item_id from ops_template_item a, ops_template b where 
        a.template_id = b.id and b.name = ?)`
    rows, err = db.Query(query_string, templateName)

  } else if hostId != "" {
    query_string := `select id, name, a.interval, history, unit, creator from ops_item a where id in 
      (select a.item_id from ops_template_item a, ops_job_template b where a.template_id=b.template_id 
      and b.job_id in (select job_id from ops_host_server where host_id=?));`
    rows, err = db.Query(query_string, hostId)
  } else {
    query_string := `select id, name, a.interval, history, unit, creator from ops_item a`
    rows, err = db.Query(query_string)
  }

  if err != nil {
    err = errors.WithStack(err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    var name string
    var interval uint64
    var history uint64
    var unit string
    var creator string
    var id int

    err = rows.Scan(&id, &name, &interval, &history, &unit, &creator)
    if err != nil {
      err = errors.WithStack(err)
      return
    }

    itemMetric := &ItemInfo{
      Name:      name,
      AliasName: creator + "_" + name + "_" + strconv.FormatUint(interval, 10) + "_" + unit,
      Interval:  interval,
      History:   history,
      Unit:      unit,
      Creator:   creator,
      Key:       id,
    }

    itemList = append(itemList, itemMetric)
  }
  return
}

func GetGraphValidItems(hostId string) ([]map[string]interface{}, error) {
  query := "select name, `interval`, unit, creator from ops_item where " +
    "history>0 and data_type not in ('string','none') and id in " +
    "(select a.item_id from ops_template_item a, ops_job_template b " +
    "where a.template_id=b.template_id and b.job_id in " +
    "(select job_id from ops_host_server where host_id=?));"
  rows, err := db.Query(query, hostId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]interface{}, 0)
  for rows.Next() {
    var name, unit, creator string
    var interval int64
    err = rows.Scan(&name, &interval, &unit, &creator)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, map[string]interface{}{
      "itemName":   name,
      "aliasName": creator + "_" + name + "_" + strconv.FormatInt(interval, 10) + "_" + unit,
    })
  }

  err = rows.Err()
  return res, errors.WithStack(err)
}
