package common

type AddTemplateArgs struct {
  TemplateName    string        `json:"templateName"`
  Items           []string      `json:"items"`
  GroupIds        []int64       `json:"groupIds"`
  EnvJobIds       []int64       `json:"envJobIds"`
  Description     string        `json:"description"`
}

type UpdateTemplate struct {
  TemplateId      int64         `json:"templateId"`
  TemplateName    string        `json:"templateName"`
  Items           []string      `json:"items"`
  GroupIds        []int64       `json:"groupIds"`
  EnvJobIds       []int64       `json:"envJobIds"`
  Description     string        `json:"description"`
}

type QueryTemplateRes struct {
  TemplateId      int64         `json:"templateId"`
  TemplateName    string        `json:"templateName"`
  Description     string        `json:"description"`
  Items           interface{}   `json:"items"`
  HostGroups      interface{}   `json:"hostGroups"`
  HostJobs        interface{}   `json:"hostJobs"`
}

type QueryTemplateInfo struct {
  TemplateId      int64         `json:"key"`
  TemplateName    string        `json:"name"`
  Description     string        `json:"description"`
  ItemCount       int           `json:"itemCount"`
  TriggerCount    int           `json:"triggerCount"`
  Groups          interface{}   `json:"groups"`
  Jobs            interface{}   `json:"jobs"`
}

