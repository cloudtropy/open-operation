package mysql

import (
  "github.com/pkg/errors"
)

type ActionTrail struct {
  Id         string `json:"id"`
  User       string `json:"user"`
  Module     string `json:"module"`
  Action     string `json:"action"`
  Result     string `json:"result"`
  Detail     string `json:"detail"`
  CreateTime string `json:"createTime"`
}

func AddActionTrail(at *ActionTrail) error {
  query := `INSERT INTO ops_action_trail (id,user,module,action,result,detail) 
    VALUES(?,?,?,?,?,?);`
  _, err := db.Exec(query, at.Id, at.User, at.Module, at.Action, at.Result, at.Detail)
  return errors.WithStack(err)
}

func GetActionTrails(start, end int64, pIndex, pCount int, searchBy, searchInfo string) ([]*ActionTrail, error) {
  query := "SELECT id,user,module,action,result,detail,create_time FROM ops_action_trail WHERE " +
    "UNIX_TIMESTAMP(create_time)>? AND UNIX_TIMESTAMP(create_time)<? "
  params := []interface{}{start, end}
  if searchBy != "" && searchInfo != "" {
    query += "AND " + searchBy + " LIKE ? "
    params = append(params, "%" + searchInfo + "%")
  }
  query += "ORDER BY create_time DESC LIMIT ?,?;"
  params = append(params, pCount * (pIndex - 1))
  params = append(params, pCount)

  rows, err := db.Query(query, params...)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]*ActionTrail, 0)
  for rows.Next() {
    at := &ActionTrail{}
    err = rows.Scan(&at.Id, &at.User, &at.Module, &at.Action, &at.Result, &at.Detail, &at.CreateTime)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, at)
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func GetActionTrailsCount(start, end int64, searchBy, searchInfo string) (int, error) {
  query := "SELECT COUNT(id) FROM ops_action_trail WHERE " +
    "UNIX_TIMESTAMP(create_time)>? AND UNIX_TIMESTAMP(create_time)<?"
  params := []interface{}{start, end}
  if searchBy != "" && searchInfo != "" {
    query += " AND " + searchBy + " LIKE ?"
    params = append(params, "%" + searchInfo + "%")
  }

  count := 0
  err := db.QueryRow(query, params...).Scan(&count)
  return count, errors.WithStack(err)
}
