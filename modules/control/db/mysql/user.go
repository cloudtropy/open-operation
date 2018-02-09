package mysql

import (
  "database/sql"
  "strings"

  "github.com/pkg/errors"
)

type User struct {
  Id         int64  `json:"id"`
  User       string `json:"user"`
  Alias      string `json:"alias"`
  Passwd     string `json:"passwd"`
  Salt       string `json:"salt"`
  Email      string `json:"email"`
  Phone      string `json:"phone"`
  Wechat     string `json:"wechat"`
  Sex        string `json:"sex"`
  CreateTime string `json:"createTime`
}

func CreateUser(u *User) error {
  tx, err := db.Begin()
  if err != nil {
    return err
  }

  query := `SELECT id FROM ops_user WHERE user=? AND role!=-1;`
  tmpId := 0
  err = db.QueryRow(query, u.User).Scan(&tmpId)
  if err == nil {
    return errors.New("AlreadyExist")
  } else if err != sql.ErrNoRows {
    return errors.WithStack(err)
  }

  query = `INSERT INTO ops_user (user, alias, passwd, salt, 
    email, phone, wechat, sex, role) VALUES(?,?,?,?,?,?,?,?,0) 
    ON DUPLICATE KEY UPDATE alias=?,passwd=?,salt=?,email=?,
    phone=?,wechat=?,sex=?,create_time=now(),role=0;`
  _, err = tx.Exec(query, u.User, u.Alias, u.Passwd, u.Salt, 
    u.Email, u.Phone, u.Wechat, u.Sex, u.Alias, u.Passwd, u.Salt, 
    u.Email, u.Phone, u.Wechat, u.Sex)
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

func DeleteUser(id int64) error {
  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }

  query := `DELETE FROM ops_team_user WHERE uid=?;`
  _, err = tx.Exec(query, id)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }

  query = `UPDATE ops_user SET role=-1 WHERE id=? AND role=0;`
  _, err = tx.Exec(query, id)
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

func UpdatePassword(user, passwd string) error {
  query := `UPDATE ops_user SET passwd=? WHERE user=?;`
  _, err := db.Exec(query, passwd, user)
  return errors.WithStack(err)
}

func UpdateUser(u *User) error {
  query := `UPDATE ops_user SET alias=?,email=?,phone=?,wechat=?,sex=? WHERE id=?;`
  _, err := db.Exec(query, u.Alias, u.Email, u.Phone, u.Wechat, u.Sex, u.Id)
  return errors.WithStack(err)
}

func GetUser(user string) (*User, error) {
  query := `SELECT id,user,alias,passwd,salt,email,phone,wechat,sex,create_time 
    FROM ops_user WHERE user=? AND role!=-1;`
  u := &User{}
  err := db.QueryRow(query, user).Scan(&u.Id, &u.User, &u.Alias, &u.Passwd, 
    &u.Salt, &u.Email, &u.Phone, &u.Wechat, &u.Sex, &u.CreateTime)
  if err != nil && err == sql.ErrNoRows {
    return nil, nil
  }
  return u, errors.WithStack(err)
}

func GetUsers() ([]*User, error) {
  query := `SELECT id,user,alias,passwd,salt,email,phone,wechat,sex,create_time 
    FROM ops_user WHERE role=0;`
  rows, err := db.Query(query)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]*User, 0)
  for rows.Next() {
    u := &User{}
    err = rows.Scan(&u.Id, &u.User, &u.Alias, &u.Passwd, 
      &u.Salt, &u.Email, &u.Phone, &u.Wechat, &u.Sex, &u.CreateTime)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, u)
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func GetUsersOfTeam(teamId int64, isBelong bool) ([]map[string]string, error) {
  query := `SELECT user,alias FROM ops_user WHERE role=0 AND id IN (
    SELECT uid FROM ops_team_user WHERE tid=?);`
  if !isBelong {
    query = strings.Replace(query, " IN ", " NOT IN ", 1)
  }
  rows, err := db.Query(query, teamId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]map[string]string, 0)
  for rows.Next() {
    var user, alias string
    err = rows.Scan(&user, &alias)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, map[string]string{
      "user": user,
      "alias": alias,
    })
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func AddUsersToTeam(tid int64, users []string) error {
  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }

  query := `INSERT INTO ops_team_user (tid, uid) 
    SELECT ?,id FROM ops_user WHERE role=0 AND user=?;`
  for _, user := range users {
    _, err = tx.Exec(query, tid, user)
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

func RemoveUsersFromTeam(tid int64, users []string) error {
  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }

  query := `DELETE FROM ops_team_user WHERE tid=? 
    AND uid IN (SELECT id FROM ops_user WHERE role=0 AND user=?);`
  for _, user := range users {
    _, err = tx.Exec(query, tid, user)
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
