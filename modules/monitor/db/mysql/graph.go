package mysql

import (
  "github.com/pkg/errors"
)

func InsertGraphData(graphName, graphType, hostId string, screenId int64) (int64, error) {

  prepareClause := `insert into ops_graph(name, graph_type, host_id, screen_id) values(?,?,?,?)`

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return -1, errors.WithStack(err)
  }
  res, err := stmt.Exec(graphName, graphType, hostId, screenId)
  if err != nil {
    return -1, errors.WithStack(err)
  }

  lastId, err := res.LastInsertId()
  if err != nil {
    return -1, errors.WithStack(err)
  }

  return lastId, nil
}


func InsertGraphItemData(graph_id int64, item_id int64) (int64, error) {

  prepareClause := `insert into ops_graph_item(graph_id, item_id) values(?,?)`
  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return -1, errors.WithStack(err)
  }
  res, err := stmt.Exec(graph_id, item_id)
  if err != nil {
    return -1, errors.WithStack(err)
  }

  lastId, err := res.LastInsertId()
  if err != nil {
    return -1, errors.WithStack(err)
  }

  return lastId, nil
}

func DeleteGraphById(id int64) error {
  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_graph_item WHERE graph_id = ?;", id)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_graph WHERE id = ?;", id)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  err = tx.Commit()
  if err != nil {
    tx.Rollback()
  }
  return errors.WithStack(err)
}

func DeleteGraphData(id int64) error {

  prepareClause := `delete from ops_graph where id = ?`

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return errors.WithStack(err)
  }
  _, err = stmt.Exec(id)
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}


func DeleteGraphItemData(id int64) error {

  if db == nil {
    return errors.New("mysql nil pointer")
  }

  prepareClause := `delete from ops_graph_item where graph_id = ?`

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return errors.WithStack(err)
  }
  _, err = stmt.Exec(id)
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}


func QueryGraphIdByName(name ,hostId string, screenId int64) (int64, error) {

  var id int64
  query_string := `select id from ops_graph 
    where name = ? 
    and host_id = ? 
    and screenId = ?"`

  if err := db.QueryRow(query_string, name, hostId, screenId).Scan(&id); err != nil {
    return 0, errors.WithStack(err)
  } else {
    return id, nil
  }
}


func QueryGraphItem(id int64) ([]string, error) {

  query_string := `select a.name from ops_item a, ops_graph_item b
    where b.graph_id = ? and a.id = b.item_id;`
  
  rows, err := db.Query(query_string, id)
  if err != nil {
      return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]string, 0)

  for rows.Next() {
    var name string

    err = rows.Scan(&name)

    if err != nil {
      return nil, errors.WithStack(err)
    }
    
    res = append(res, name) 
  }

  return res, nil
}

func QueryGraph(graphId int64) (map[string]interface{}, error) {
  queryStr := `SELECT a.name, a.graph_type, a.host_id, b.ip, c.name 
    FROM ops_graph a, ops_host b, ops_screen c WHERE a.id=? AND 
    a.host_id=b.host_id AND a.screen_id=c.id;`
  var graphName, graphType, hostId, hostIp, screenName string
  err := db.QueryRow(queryStr, graphId).Scan(&graphName, 
    &graphType, &hostId, &hostIp, &screenName)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return map[string]interface{}{
    "graphId":      graphId,
    "graphName":    graphName,
    "graphType":    graphType,
    "hostId":       hostId,
    "hostIp":       hostIp,
    "screenName":   screenName,
  }, nil
}

func UpdateGraph(graphId int64, graphName, graphType string) error {
  queryStr := `UPDATE ops_graph SET name=?,graph_type=? 
    WHERE id=?;`
  stmt, err := db.Prepare(queryStr)
  if err != nil {
    return errors.WithStack(err)
  }
  defer stmt.Close()

  _, err = stmt.Exec(graphName, graphType, graphId)
  return errors.WithStack(err)
}

func UpdateGraphItems(graphId int64, hostId string, addItems, delItems []string) error {
  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }

  for _, item := range addItems {
    _, err = tx.Exec(`INSERT INTO ops_graph_item (graph_id, item_id) 
      SELECT ?,id FROM ops_item WHERE name=? AND id IN 
      (SELECT a.item_id FROM ops_template_item a, ops_job_template b 
        WHERE a.template_id=b.template_id AND b.job_id IN 
        (SELECT job_id FROM ops_host_server WHERE host_id=?));`, 
      graphId, item, hostId)
    if err != nil {
      tx.Rollback()
      return errors.WithStack(err)
    }
  }

  for _, item := range delItems {
    _, err = tx.Exec(`DELETE FROM ops_graph_item WHERE graph_id=? 
      AND item_id IN (SELECT id FROM ops_item WHERE name=?);`, 
      graphId, item)
    if err != nil {
      tx.Rollback()
      return errors.WithStack(err)
    }
  }

  err = tx.Commit()
  if err != nil {
    tx.Rollback()
  }
  return errors.WithStack(err)
}

