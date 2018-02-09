package mysql

import (
  "github.com/pkg/errors"
  cc "github.com/cloudtropy/open-operation/utils/common"
)

type HostInfo struct {
  HostId       string `json:"host_id"`
  HostOs       string `json:"host_os"`
  HostIp       string `json:"host_ip"`
  HostStatus   int    `json:"host_status"`
  Hostname     string `json:"hostname"`
  Comment      string `json:"comment"`
  Location     string `json:"location"`
  CpuCount     uint8  `json:"cpu_count"`
  MemCapacity  uint64 `json:"mem_capacity"`
  DiskCapacity uint64 `json:"disk_capacity"`
  CreateTime   string `json:"create_time"`
  UpdateTime   string `json:"update_time,omitempty"`
}

func UpsertHostInfo(h *HostInfo) error {
  query := `INSERT INTO ops_host (host_id,ip,hostname,cpu_count,
    mem_capacity,disk_capacity,os) VALUES(?,?,?,?,?,?,?) 
    ON DUPLICATE KEY UPDATE ip=?,hostname=?,cpu_count=?,mem_capacity=?,
    disk_capacity=?,os=?,update_time=now();`
  _, err := db.Exec(query, h.HostId, h.HostIp, h.Hostname, h.CpuCount,
    h.MemCapacity, h.DiskCapacity, h.HostOs, h.HostIp, h.Hostname, h.CpuCount,
    h.MemCapacity, h.DiskCapacity, h.HostOs)
  return errors.WithStack(err)
}

/*
 sFlag: 0 online
        1 offline
 */
func GetHostInfos(sFlag int) ([]*HostInfo, error) {
  query := `SELECT host_id, ip, status, hostname, cpu_count, 
    mem_capacity, disk_capacity, os, comment, create_time, 
    location FROM ops_host WHERE status=? 
    ORDER BY create_time DESC;`
  rows, err := db.Query(query, sFlag)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer rows.Close()

  hs := make([]*HostInfo, 0)
  for rows.Next() {
    hi := HostInfo{}
    err = rows.Scan(&hi.HostId, &hi.HostIp, &hi.HostStatus, &hi.Hostname, &hi.CpuCount,
      &hi.MemCapacity, &hi.DiskCapacity, &hi.HostOs, &hi.Comment, &hi.CreateTime, &hi.Location)
    if err != nil {
      return nil, errors.WithStack(err)
    }
    hs = append(hs, &hi)
  }
  err = rows.Err()
  return hs, errors.WithStack(err)
}

func UpdateOneHostInfo(mss map[string]string) error {
  query := "UPDATE ops_host SET " + mss["update_key"] + "=? WHERE host_id=?"
  _, err := db.Exec(query, mss["new_value"], mss["host_id"])
  return errors.WithStack(err)
}

func GetItemsOfHostId(hostId string) ([]cc.HostItem, error) {
  rows, err := db.Query("select name,`interval`,dst,creator,history,data_type,"+
    "UNIX_TIMESTAMP(update_time) from ops_item where id in (select a.item_id " +
    "from ops_template_item a, ops_job_template b where a.template_id=b.template_id " +
    "and b.job_id in (select job_id from ops_host_server where host_id=?));", hostId)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  res := make([]cc.HostItem, 0)
  for rows.Next() {
    hi := cc.HostItem{}
    err = rows.Scan(&hi.ItemName, &hi.Interval, &hi.Dst, &hi.Creator, 
      &hi.History, &hi.DataType, &hi.Timestamp)
    if err != nil {
      return res, err
    }
    res = append(res, hi)
  }
  err = rows.Err()
  return res, err
}

func GetItemsForAgent(hostId string) ([]cc.HostItem, error) {
  rows, err := db.Query("select name,`interval`,dst,creator,history,data_type,"+
    "UNIX_TIMESTAMP(update_time) from ops_item where creator=\"born\" or id in (select a.item_id " +
    "from ops_template_item a, ops_job_template b where a.template_id=b.template_id " +
    "and b.job_id in (select job_id from ops_host_server where host_id=?));", hostId)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  res := make([]cc.HostItem, 0)
  for rows.Next() {
    hi := cc.HostItem{}
    err = rows.Scan(&hi.ItemName, &hi.Interval, &hi.Dst, &hi.Creator, 
      &hi.History, &hi.DataType, &hi.Timestamp)
    if err != nil {
      return res, err
    }
    res = append(res, hi)
  }
  err = rows.Err()
  return res, err
}


func QueryIdByName(name, tableName string) (int64, error) {
  var id int64
  query := "SELECT id FROM " + tableName + " WHERE name=?;"
  if err := db.QueryRow(query, name).Scan(&id); err != nil {
    return 0, errors.WithStack(err)
  } else {
    return id, nil
  }
}

