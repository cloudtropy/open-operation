package controller

import (
  "database/sql"
  "strings"

  cc "github.com/cloudtropy/open-operation/utils/common"
  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/db/mysql"
  "github.com/pkg/errors"
)

type Template struct{}

func (t *Template) AddTemplate(args cc.AddTemplateArgs, templateId *int64) error {
  var err error
  *templateId, err = mysql.InsertTemplateData(args.TemplateName, args.Description)
  if err != nil {
    log.Println("mysql.InsertTemplateData", args.TemplateName, args.Description, err)
    return err
  }

  for _, item := range args.Items {
    // The item host.heartbeat can not be used in other template.
    if item == "host.heartbeat" {
      continue
    }
    id, err := mysql.QueryIdByName(item, "ops_item")
    if err != nil {
      if err == sql.ErrNoRows {
        log.Println("mysql.QueryIdByName", item, "ops_item", err)
        continue
      }
      return err
    }

    if _, err := mysql.InsertItemTemplateData(*templateId, id); err != nil {
      return err
    }
  }

  for _, groupId := range args.GroupIds {
    rowsAffected, err := mysql.AddTemplateGroup(*templateId, groupId)
    if err != nil {
      if strings.Index(err.Error(), "Duplicate") != -1 {
        continue
      }
      log.Println("mysql.AddTemplateGroup", *templateId, groupId, err)
      return err
    }
    if rowsAffected != 1 {
      continue
    }
    // afterSetTemplateJobs(*templateId, groupId, "add")
  }

  for _, envJobId := range args.EnvJobIds {
    rowsAffected, err := mysql.AddTemplateJob(*templateId, envJobId)
    if err != nil {
      if strings.Index(err.Error(), "Duplicate") != -1 {
        continue
      }
      log.Println("mysql.AddTemplateJob", *templateId, envJobId, err)
      return err
    }
    if rowsAffected != 1 {
      continue
    }
    // jobId, err := mysql.GetJobIdOfEnvJobId(envJobId)
    _, err = mysql.GetJobIdOfEnvJobId(envJobId)
    if err != nil {
      log.Println("mysql.GetJobIdOfEnvJobId", envJobId, err)
      return err
    }
    // afterSetTemplateJobs(*templateId, jobId, "add")
  }
  return nil
}

func (t *Template) DelTemplate(templateId int64, res *string) error {
  // Delete the template named Basic will not be allowed.
  if templateId == 1 {
    return errors.New("InvalidOperation")
  }

  // jobs, err := mysql.QueryTemplateJobs(templateId, true)
  _, err := mysql.QueryTemplateJobs(templateId, true)
  if err != nil {
    log.Println("mysql.QueryTemplateJobs", templateId, err)
    return err
  }

  // trigger template_item template_job
  err = mysql.DeleteTemplate(templateId)
  if err != nil {
    return err
  }

  // for _, job := range jobs {
  //   afterSetTemplateJobs(templateId, job.Key, "del")
  // }

  *res = "Success"
  return nil
}

func (t *Template) UpdateTemplate(args cc.UpdateTemplate, res *string) error {
  // update ops_template
  err := mysql.UpdateTemplate(args.TemplateId, args.TemplateName, args.Description)
  if err != nil {
    log.Println("mysql.UpdateTemplate", args.TemplateId, args.TemplateName, args.Description, err)
    return err
  }

  items, err := mysql.QueryItemsOfTemplate(args.TemplateId)
  if err != nil {
    log.Println("mysql.QueryItemsOfTemplate", args.TemplateId, err)
    return err
  }

  itemChanged := false
  for _, itemName := range args.Items {
    itemId, isExist := items[itemName]
    if isExist {
      delete(items, itemName)
      continue
    }

    // The item host.heartbeat can not be used in other template.
    if itemName == "host.heartbeat" {
      continue
    }

    itemId, err = mysql.QueryIdByName(itemName, "ops_item")
    if err != nil {
      if err == sql.ErrNoRows {
        log.Println("mysql.QueryIdByName", itemName, "ops_item", err)
        continue
      }
      return err
    }

    if _, err = mysql.InsertItemTemplateData(args.TemplateId, itemId); err != nil {
      if strings.Index(err.Error(), "Duplicate") != -1 {
        continue
      }
      return err
    }
    itemChanged = true
  }
  for _, itemId := range items {
    err = mysql.DeleteItemsOfTemplate(itemId, args.TemplateId)
    if err != nil {
      log.Println("mysql.DeleteItemsOfTemplate", itemId, args.TemplateId, err)
      return err
    }
    itemChanged = true
  }

  groups, err := mysql.QueryTemplateGroups(args.TemplateId, true)
  if err != nil {
    log.Println("mysql.QueryTemplateGroups", args.TemplateId, err)
    return err
  }
  mIdGroups := make(map[int64]mysql.JobInfo)
  for _, group := range groups {
    mIdGroups[group.Key] = group
  }

  for _, groupId := range args.GroupIds {
    if _, isExist := mIdGroups[groupId]; isExist {
      delete(mIdGroups, groupId)
      if itemChanged {
        // afterSetTemplateJobs(args.TemplateId, groupId, "add")
      }
      continue
    }

    _, err = mysql.AddTemplateGroup(args.TemplateId, groupId)
    if err != nil {
      if strings.Index(err.Error(), "Duplicate") != -1 {
        if itemChanged {
          // afterSetTemplateJobs(args.TemplateId, groupId, "add")
        }
        continue
      }
      log.Println("mysql.AddTemplateJob", args.TemplateId, groupId, err)
      return err
    }
    // afterSetTemplateJobs(args.TemplateId, groupId, "add")
  }

  for groupId, _ := range mIdGroups {
    err = mysql.DelTemplateJob(args.TemplateId, groupId)
    if err != nil {
      log.Println("mysql.DelTemplateJob", args.TemplateId, groupId, err)
      return err
    }
    // afterSetTemplateJobs(args.TemplateId, groupId, "del")
  }

  jobs, err := mysql.QueryTemplateJobs(args.TemplateId, true)
  if err != nil {
    log.Println("mysql.QueryTemplateJobs", args.TemplateId, err)
    return err
  }
  mIdJobs := make(map[int64]mysql.JobInfo)
  for _, job := range jobs {
    mIdJobs[job.Key] = job
  }

  for _, envJobId := range args.EnvJobIds {
    jobId, err := mysql.GetJobIdOfEnvJobId(envJobId)
    if err != nil {
      log.Println("mysql.GetJobIdOfEnvJobId", envJobId, err)
      return err
    }

    if _, isExist := mIdJobs[jobId]; isExist {
      delete(mIdJobs, jobId)
      if itemChanged {
        // afterSetTemplateJobs(args.TemplateId, jobId, "add")
      }
      continue
    }

    _, err = mysql.AddTemplateGroup(args.TemplateId, jobId)
    if err != nil {
      if strings.Index(err.Error(), "Duplicate") != -1 {
        if itemChanged {
          // afterSetTemplateJobs(args.TemplateId, jobId, "add")
        }
        continue
      }
      log.Println("mysql.AddTemplateJob", args.TemplateId, jobId, err)
      return err
    }
    // afterSetTemplateJobs(args.TemplateId, jobId, "add")
  }

  for jobId, _ := range mIdJobs {
    err = mysql.DelTemplateJob(args.TemplateId, jobId)
    if err != nil {
      log.Println("mysql.DelTemplateJob", args.TemplateId, jobId, err)
      return err
    }
    // afterSetTemplateJobs(args.TemplateId, jobId, "del")
  }

  *res = "Success"
  return nil
}

func (t *Template) QueryTemplate(templateId int64, res *cc.QueryTemplateRes) error {
  var err error
  if templateId != -1 {
    res.TemplateName, res.Description, err = mysql.QueryTemplateInfo(templateId)
    if err != nil {
      log.Println("mysql.QueryTemplateInfo", templateId, err)
      return err
    }
  }
  res.TemplateId = templateId
  
  belongItems, err := mysql.QueryTemplateItems(templateId, true)
  if err != nil {
    log.Println("mysql.QueryTemplateItems", templateId, true, err)
    return err
  }
  notBelongItems, err := mysql.QueryTemplateItems(templateId, false)
  if err != nil {
    log.Println("mysql.QueryTemplateItems", templateId, false, err)
    return err
  }

  belongGroups, err := mysql.QueryTemplateGroups(templateId, true)
  if err != nil {
    log.Println("mysql.QueryTemplateGroups", templateId, true, err)
    return err
  }
  notBelongGroups, err := mysql.QueryTemplateGroups(templateId, false)
  if err != nil {
    log.Println("mysql.QueryTemplateGroups", templateId, false, err)
    return err
  }

  belongEnvJobs, err := mysql.QueryTemplateEnvJobs(templateId, true)
  if err != nil {
    log.Println("mysql.QueryTemplateEnvJobs", templateId, true, err)
    return err
  }
  notBelongEnvJobs, err := mysql.QueryTemplateEnvJobs(templateId, false)
  if err != nil {
    log.Println("mysql.QueryTemplateEnvJobs", templateId, false, err)
    return err
  }

  res.Items = map[string]interface{}{
    "belong": belongItems,
    "notBelong": notBelongItems,
  }
  res.HostGroups = map[string]interface{}{
    "belong": belongGroups,
    "notBelong": notBelongGroups,
  }
  res.HostJobs = map[string]interface{}{
    "belong": belongEnvJobs,
    "notBelong": notBelongEnvJobs,
  }

  return nil
}

func (t *Template) QueryTemplates(args string, res *[]cc.QueryTemplateInfo) error {
  err := mysql.QueryTemplates(res)
  if err != nil {
    return err
  }
  for i, t := range *res {
    (*res)[i].ItemCount, err = mysql.QueryNumbersByTemplateId(t.TemplateId, "ops_template_item")
    if err != nil {
      log.Println("mysql.QueryNumbersByTemplateId", t.TemplateId, "ops_template_item", err)
      return err
    }

    (*res)[i].TriggerCount, err = mysql.QueryNumbersByTemplateId(t.TemplateId, "ops_trigger")
    if err != nil {
      log.Println("mysql.QueryNumbersByTemplateId", t.TemplateId, "ops_trigger", err)
      return err
    }

    (*res)[i].Groups, err = mysql.QueryTemplateGroups(t.TemplateId, true)
    if err != nil {
      log.Println("mysql.QueryTemplateGroups", t.TemplateId, true, err)
      return err
    }

    (*res)[i].Jobs, err = mysql.QueryTemplateEnvJobs(t.TemplateId, true)
    if err != nil {
      log.Println("mysql.QueryTemplateEnvJobs", t.TemplateId, true, err)
      return err
    }
  }
  return nil
}
