package common


type MetricValue struct {
  Endpoint  string      `json:"endpoint"`
  Metric    string      `json:"metric"`
  Value     interface{} `json:"value"`
  DataType  string      `json:"dataType"`
  Step      int64       `json:"step"`
  Type      string      `json:"counterType"`
  Tags      string      `json:"tags"`
  Timestamp int64       `json:"timestamp"`
}

type HostItem struct {
  ItemName  string
  Interval  int64
  Dst       string
  DataType  string
  Creator   string
  Command   string
  History   int64
  Timestamp int64
}

type TopicMsg struct {
  Topic     string      `json:"topic"`
  Data      interface{} `json:"data"`
}
