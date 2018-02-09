package mysql

import (
  "database/sql"
  "fmt"
  "strconv"
  "strings"

  cc "github.com/cloudtropy/open-operation/utils/common"
  "github.com/pkg/errors"
)

func InsertTemplateData(name, description string) (int64, error) {

  prepareClause := "Insert into ops_template(name, description) values(?, ?)"

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  res, err := stmt.Exec(name, description)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    return 0, errors.WithStack(err)
  }

  return lastId, nil
}


func InsertItemTemplateData(templateId, itemId int64) (int64, error) {

  prepareClause := "Insert into ops_template_item(template_id, item_id) values(?, ?)"

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  res, err := stmt.Exec(templateId, itemId)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    return 0, errors.WithStack(err)
  }

  return lastId, nil
}

func InsertItemGraphData(graphId, itemId int64) (int64, error) {

  prepareClause := "Insert into ops_graph_item(graph_id, item_id) values(?, ?)"

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  res, err := stmt.Exec(graphId, itemId)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    return 0, errors.WithStack(err)
  }

  return lastId, nil
}

func InsertJobData(jobId, templateId int64) (int64, error) {

  prepareClause := "Insert into ops_job_template(job_id, template_id) values(?, ?)"

  stmt, err := db.Prepare(prepareClause)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  res, err := stmt.Exec(jobId, templateId)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    return 0, errors.WithStack(err)
  }

  return lastId, nil
}

func QueryNumbersByTemplateId(id int64, tableName string) (int, error) {

  var number int
  query_string := fmt.Sprintf("select count(*) as number from %s where template_id = ?", tableName)
  if err := db.QueryRow(query_string, id).Scan(&number); err != nil {
    return 0, errors.WithStack(err)
  } else {
    return number, nil
  }

}

func QueryHostCountOfTemplate(template string) (int, error) {
  queryStr := `select count(host_id) from ops_host_server where job_id in 
    (select a.job_id from ops_job_template a, ops_template b where 
    a.template_id=b.id and b.name=?);`
  var hostCount int
  err := db.QueryRow(queryStr, template).Scan(&hostCount)
  if err != nil {
    if err == sql.ErrNoRows {
      return 0, nil
    } else {
      return -1, errors.WithStack(err)
    }
  }
  return hostCount, nil
}

func DeleteTemplate(templateId int64) error {
  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_trigger WHERE template_id=?;", templateId)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_template_item WHERE template_id=?;", templateId)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_job_template WHERE template_id=?;", templateId)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_template WHERE id=?;", templateId)
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

func UpdateTemplate(id int64, name, description string) error {
  _, err := db.Exec(`UPDATE ops_template SET name=?,description=? 
    WHERE id=?;`, name, description, id)
  return errors.WithStack(err)
}

func DeleteItemsOfTemplate(itemId, templateId int64) error {
  tx, err := db.Begin()
  if err != nil {
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_trigger WHERE item_id=? AND template_id=?;", 
    itemId, templateId)
  if err != nil {
    tx.Rollback()
    return errors.WithStack(err)
  }
  _, err = tx.Exec("DELETE FROM ops_template_item WHERE item_id=? AND template_id=?;", 
    itemId, templateId)
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

func QueryItemsOfTemplate(templateId int64) (map[string]int64, error) {
  queryStr := `SELECT a.id, a.name FROM ops_item a, ops_template_item b 
    WHERE a.id=b.item_id AND b.template_id=?;`
  rows, err := db.Query(queryStr, templateId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make(map[string]int64)
  for rows.Next() {
    var itemName string
    var itemId int64
    err := rows.Scan(&itemId, &itemName)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res[itemName] = itemId
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func QueryTemplateInfo(templateId int64) (string, string, error) {
  var name, description string
  err := db.QueryRow("SELECT name,description FROM ops_template WHERE id=?;", 
    templateId).Scan(&name, &description)
  return name, description, errors.WithStack(err)
}

func QueryTemplateItems(templateId int64, isBlong bool) ([]ItemInfo, error) {
  
  queryStr := `SELECT id, name, a.interval, history, unit, creator FROM ops_item a WHERE id IN 
    (SELECT item_id FROM ops_template_item WHERE template_id=?);`
  if !isBlong {
    queryStr = strings.Replace(queryStr, " IN ", " NOT IN ", 1)
  }
  rows, err := db.Query(queryStr, templateId)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  res := make([]ItemInfo, 0)
  for rows.Next() {
    var name, unit, creator string
    var interval, history uint64
    var id int
    err = rows.Scan(&id, &name, &interval, &history, &unit, &creator)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    res = append(res, ItemInfo{
      Name:      name,
      AliasName: creator + "_" + name + "_" + strconv.FormatUint(interval, 10) + "_" + unit,
      Interval:  interval,
      History:   history,
      Unit:      unit,
      Creator:   creator,
      Key:       id,
    })
  }
  err = rows.Err()
  return res, errors.WithStack(err)
}

func QueryTemplates(res *[]cc.QueryTemplateInfo) error {
  rows, err := db.Query("SELECT id,name,description FROM ops_template;")
  if err != nil {
    return errors.WithStack(err)
  }
  defer rows.Close()
  for rows.Next() {
    t := cc.QueryTemplateInfo{}
    err = rows.Scan(&t.TemplateId, &t.TemplateName, &t.Description)
    if err != nil {
      return errors.WithStack(err)
    }
    *res = append(*res, t)
  }
  err = rows.Err()
  return errors.WithStack(err)
}


func AddTemplateGroup(templateId, groupId int64) (int64, error) {
  sqlStr := `INSERT INTO ops_job_template (template_id, job_id) 
    SELECT ?,? FROM dual WHERE EXISTS(SELECT * FROM ops_job WHERE id=?);`
  r, err := db.Exec(sqlStr, templateId, groupId, groupId)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  return r.RowsAffected()
}

func AddTemplateJob(templateId, envJobId int64) (int64, error) {
  sqlStr := `INSERT INTO ops_job_template (template_id, job_id) 
    SELECT ?,job_id FROM ops_env_job WHERE id=?;`
  r, err := db.Exec(sqlStr, templateId, envJobId)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  return r.RowsAffected()
}

func GetJobIdOfEnvJobId(envJobId int64) (int64, error) {
  var jobId int64
  err := db.QueryRow(`SELECT job_id FROM ops_env_job WHERE id=?;`, envJobId).Scan(&jobId)
  return jobId, errors.WithStack(err)
}

type JobInfo struct {
  Name string `json:"name"`
  Key  int64  `json:"key"`
}

func QueryTemplateJobs(templateId int64, belongOrNot bool) ([]JobInfo, error) {

  var query_string string
  jobList := make([]JobInfo, 0)

  if belongOrNot == true {
    query_string = `select id, name from ops_job where job_type!="GROUP" and id in (
      select job_id from ops_job_template where template_id = ?);`
  } else {
    query_string = `select id, name from ops_job where job_type!="GROUP" and id not in (
      select job_id from ops_job_template where template_id = ?);`
  }

  rows, err := db.Query(query_string, templateId)
  if err != nil {
    return jobList, errors.WithStack(err)
  }
  defer rows.Close()

  for rows.Next() {
    var id int64
    var name string

    err = rows.Scan(&id, &name)
    if err != nil {
      return jobList, errors.WithStack(err)
    }

    jobMetric := JobInfo{
      Name: name,
      Key:  id,
    }
    jobList = append(jobList, jobMetric)
  }
  return jobList, nil
}

func QueryTemplateGroups(templateId int64, belongOrNot bool) ([]JobInfo, error) {

  var query_string string
  jobList := make([]JobInfo, 0)

  if belongOrNot == true {
    query_string = `select id, name from ops_job where job_type="GROUP" and id in (
      select job_id from ops_job_template where template_id=?);`
  } else {
    query_string = `select id, name from ops_job where job_type="GROUP" and id not in (
      select job_id from ops_job_template where template_id=?);`
  }

  rows, err := db.Query(query_string, templateId)
  if err != nil {
    return jobList, errors.WithStack(err)
  }
  defer rows.Close()

  for rows.Next() {
    var id int64
    var name string
    err = rows.Scan(&id, &name)
    if err != nil {
      return jobList, errors.WithStack(err)
    }
    jobList = append(jobList, JobInfo{
      Name: name,
      Key:  id,
    })
  }
  return jobList, nil
}

func QueryTemplateEnvJobs(templateId int64, belongOrNot bool) ([]JobInfo, error) {

  var query_string string
  jobList := make([]JobInfo, 0)

  if belongOrNot == true {
    query_string = `select id, job_alias from ops_env_job where job_id in (
      select job_id from ops_job_template where template_id = ?);`
  } else {
    query_string = `select id, job_alias from ops_env_job where job_id not in (
      select job_id from ops_job_template where template_id = ?);`
  }

  rows, err := db.Query(query_string, templateId)
  if err != nil {
    return jobList, errors.WithStack(err)
  }
  defer rows.Close()

  for rows.Next() {
    var id int64
    var name string
    err = rows.Scan(&id, &name)
    if err != nil {
      return jobList, errors.WithStack(err)
    }
    jobList = append(jobList, JobInfo{
      Name: name,
      Key:  id,
    })
  }

  return jobList, nil
}

func DelTemplateJob(templateId, jobId int64) error {
  sqlStr := `delete from ops_job_template where template_id=? and job_id=?;`
  stmt, err := db.Prepare(sqlStr)
  if err != nil {
    return errors.WithStack(err)
  }
  defer stmt.Close()

  _, err = stmt.Exec(templateId, jobId)
  return errors.WithStack(err)
}
