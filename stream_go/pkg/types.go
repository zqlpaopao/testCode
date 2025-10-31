package pkg

import (
	"encoding/json"
	"fmt"
	"time"
)

// SensorData 传感器数据结构
type SensorData struct {
	DeviceID       string  `json:"deviceId"`
	DeviceType     string  `json:"deviceType"`
	Temperature    float64 `json:"temperature"`
	DatacenterName string  `json:"datacenter_name"`
	Room           string  `json:"room"`
	Rack           string  `json:"rack"`
}

func (s *SensorData) String() string {
	jsonData, _ := json.Marshal(s)
	return string(jsonData)
}

// Event 事件结构，用于滑动窗口处理
type Event struct {
	ID        int       `json:"id"`
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *Event) String() string {
	return fmt.Sprintf("Event{ID:%d, Value:%d, Time:%s}", e.ID, e.Value, e.Timestamp.Format("15:04:05"))
}