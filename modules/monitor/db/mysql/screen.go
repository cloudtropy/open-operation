package mysql

import (
  "github.com/pkg/errors"
)

type ScreenInfo struct {
  Id   int64
  Name string
}

func InsertScreenData(s *ScreenInfo) (int64, error) {

  prepareClause := `insert into ops_screen(name) values(?)`

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return -1, errors.WithStack(err)
  }
  res, err := stmt.Exec(s.Name)
  if err != nil {
    return -1, errors.WithStack(err)
  }

  lastId, err := res.LastInsertId()
  if err != nil {
    return -1, errors.WithStack(err)
  }

  return lastId, nil
}

func DeleteScreenData(id int64) error {

  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }
  _, err = tx.Exec(`DELETE FROM ops_graph_item WHERE graph_id in 
    (SELECT a.id FROM ops_graph a, ops_screen b WHERE a.screen_id=b.id 
    AND b.id=?);`, id)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  _, err = tx.Exec(`DELETE FROM ops_graph WHERE screen_id in 
    (SELECT id FROM ops_screen WHERE id=?);`, id)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_screen WHERE id=?;", id)
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

func UpdateScreenData(s *ScreenInfo) error {
  queryStr := `UPDATE ops_screen SET name=? WHERE id=?;`
  _, err := db.Exec(queryStr, s.Name, s.Id)
  return errors.WithStack(err)
}

func QueryScreenList() ([]map[string]interface{}, error) {

  query_string := "select id, name from ops_screen"
  rows, err := db.Query(query_string)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]interface{}, 0)

  for rows.Next() {
    var name string
    var id int64

    err = rows.Scan(&id, &name)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    res = append(res, map[string]interface{}{
      "name": name,
      "key":  id,
    })
  }
  return res, nil
}

func QueryScreenGraph(screenId int64) ([]map[string]interface{}, error) {

  query_string := `select id, name, graph_type, host_id 
    from ops_graph where screen_id=?;`
  rows, err := db.Query(query_string, screenId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]interface{}, 0)
  i := 0
  for rows.Next() {
    var id int64
    var name string
    var graphType string
    var hostId string

    err = rows.Scan(&id, &name, &graphType, &hostId)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    res = append(res, map[string]interface{}{
      "id":        id,
      "title":     name,
      "graphType": graphType,
      "hostId":    hostId,
      "key":       i,
    })

    i++
  }
  return res, nil
}
