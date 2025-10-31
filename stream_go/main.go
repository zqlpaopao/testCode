// https://www.jianshu.com/p/3209edd28187

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

type message struct {
	Msg string
}

func (msg *message) String() string {
	return msg.Msg
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

func main() {
	// 运行简单的传感器数据生成示例
	//simpleSensorDataExample()

	// 运行官方滑动窗口示例
	officialSlidingWindowExample()

	// 等待足够长的时间让程序输出数据
	time.Sleep(250 * time.Second)
	println("程序结束")
}

func base() {
	// 构建输入源：该输入源每隔一段时间，输出当前时间的字符串
	source := ext.NewChanSource(tickerChan(time.Second * 1))

	// 对数据进行处理：时间字符后面加上"-UTC"
	flow := flow.NewMap(mapp, 1)

	// 结果输出：将结果输出到标准输出
	sink := ext.NewStdoutSink()

	// 将流式处理各个步骤串联起来
	source.Via(flow).To(sink)
}

var mapp = func(in interface{}) interface{} {
	msg := in.(*message)
	msg.Msg += "-UTC"
	return msg
}

// simpleSensorDataExample 简单的传感器数据生成示例
func simpleSensorDataExample() {
	println("=== 简单传感器数据生成示例开始 ===")

	// 预定义的数据中心、房间和机架选项，用于随机选择但保持重复性
	datacenters := []string{"d1", "d2", "d3", "d4"}
	rooms := []string{"r1", "r2", "r3", "r4", "r5"}
	racks := []string{"r1", "r2", "r3", "r4", "r5", "r6"}
	deviceTypes := []string{"temperature", "humidity", "pressure"}

	// 创建一个ticker，每秒生成一次数据
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	counter := 0

	go func() {
		for range ticker.C {
			counter++

			// 生成随机但重复性的传感器数据
			sensorData := &SensorData{
				DeviceID:       fmt.Sprintf("sensor_%02d", rand.Intn(10)+1), // sensor_01 到 sensor_10
				DeviceType:     deviceTypes[rand.Intn(len(deviceTypes))],
				Temperature:    20.0 + rand.Float64()*25.0, // 20.0 到 45.0 度之间的随机温度
				DatacenterName: datacenters[rand.Intn(len(datacenters))],
				Room:           rooms[rand.Intn(len(rooms))],
				Rack:           racks[rand.Intn(len(racks))],
			}

			// 输出JSON格式的数据
			fmt.Println(sensorData.String())

			if counter >= 20 { // 生成20个数据后停止
				ticker.Stop()
				println("=== 简单传感器数据生成示例结束 ===")
				return
			}
		}
	}()
}

// sensorDataExample 传感器数据生成示例
func sensorDataExample() {
	println("=== 传感器数据生成示例开始 ===")

	// 创建传感器数据源：每隔1秒生成一个传感器数据
	source := ext.NewChanSource(createSensorDataChan(time.Second * 1))

	// 创建一个简单的传递流（不做任何处理，直接传递数据）
	passThrough := flow.NewMap(func(in interface{}) interface{} {
		return in // 直接返回输入数据
	}, 1)

	// 直接输出到标准输出
	sink := ext.NewStdoutSink()

	// 构建处理管道
	source.Via(passThrough).To(sink)

	// 让程序运行一段时间以观察效果
	go func() {
		time.Sleep(10 * time.Second)
		println("=== 传感器数据生成示例结束 ===")
	}()
}

// createSensorDataChan 创建传感器数据源（专门用于直接输出JSON格式）
func createSensorDataChan(interval time.Duration) chan interface{} {
	ticker := time.NewTicker(interval)
	oc := ticker.C
	nc := make(chan interface{})

	// 预定义的数据中心、房间和机架选项，用于随机选择但保持重复性
	datacenters := []string{"d1", "d2", "d3", "d4"}
	rooms := []string{"r1", "r2", "r3", "r4", "r5"}
	racks := []string{"r1", "r2", "r3", "r4", "r5", "r6"}
	deviceTypes := []string{"temperature", "humidity", "pressure"}

	counter := 0

	go func() {
		defer close(nc)
		for range oc {
			counter++

			// 生成随机但重复性的传感器数据
			sensorData := &SensorData{
				DeviceID:       fmt.Sprintf("sensor_%02d", rand.Intn(10)+1), // sensor_01 到 sensor_10
				DeviceType:     deviceTypes[rand.Intn(len(deviceTypes))],
				Temperature:    20.0 + rand.Float64()*25.0, // 20.0 到 45.0 度之间的随机温度
				DatacenterName: datacenters[rand.Intn(len(datacenters))],
				Room:           rooms[rand.Intn(len(rooms))],
				Rack:           racks[rand.Intn(len(racks))],
			}

			nc <- sensorData

			if counter >= 10 { // 生成10个数据后停止
				ticker.Stop()
				return
			}
		}
	}()
	return nc
}

func tickerChan(repeat time.Duration) chan interface{} {
	ticker := time.NewTicker(repeat)
	oc := ticker.C
	nc := make(chan interface{})
	go func() {
		for range oc {
			nc <- &message{strconv.FormatInt(time.Now().UTC().UnixNano(), 10)}
		}
	}()
	return nc
}

// officialSlidingWindowExample 使用官方 SlidingWindow 的示例
func officialSlidingWindowExample() {
	println("=== 官方滑动窗口示例开始 ===")

	// 创建事件数据源：每隔0.5秒生成一个事件
	source := ext.NewChanSource(createEventChan(time.Millisecond * 500))

	// 创建滑动窗口：窗口大小3秒，滑动间隔1秒
	// 这意味着每1秒会输出一个包含过去3秒内所有传感器数据的窗口
	slidingWindow := flow.NewSlidingWindow[*SensorData](3*time.Second, 1*time.Second)

	// 处理窗口结果：将窗口内的事件进行汇总处理
	windowProcessor := flow.NewMap(processWindow, 1)

	// 输出结果
	sink := ext.NewStdoutSink()

	// 构建处理管道
	source.Via(slidingWindow).Via(windowProcessor).To(sink)

	// 让程序运行一段时间以观察效果
	go func() {
		time.Sleep(15 * time.Second)
		println("=== 官方滑动窗口示例结束 ===")
	}()
}

// createEventChan 创建事件数据源
func createEventChan(interval time.Duration) chan interface{} {
	ticker := time.NewTicker(interval)
	oc := ticker.C
	nc := make(chan interface{})

	// 预定义的数据中心、房间和机架选项，用于随机选择但保持重复性
	datacenters := []string{"d1", "d2", "d3", "d4"}
	rooms := []string{"r1", "r2", "r3", "r4", "r5"}
	racks := []string{"r1", "r2", "r3", "r4", "r5", "r6"}
	deviceTypes := []string{"temperature", "humidity", "pressure"}

	counter := 0

	go func() {
		defer close(nc)
		for range oc {
			counter++

			// 生成随机但重复性的传感器数据
			sensorData := &SensorData{
				DeviceID:       fmt.Sprintf("sensor_%02d", rand.Intn(10)+1), // sensor_01 到 sensor_10
				DeviceType:     deviceTypes[rand.Intn(len(deviceTypes))],
				Temperature:    20.0 + rand.Float64()*25.0, // 20.0 到 45.0 度之间的随机温度
				DatacenterName: datacenters[rand.Intn(len(datacenters))],
				Room:           rooms[rand.Intn(len(rooms))],
				Rack:           racks[rand.Intn(len(racks))],
			}

			nc <- sensorData

			if counter >= 50 { // 生成50个数据后停止，可以根据需要调整
				ticker.Stop()
				return
			}
		}
	}()
	return nc
}

// processWindow 处理滑动窗口的结果
func processWindow(windowData interface{}) interface{} {
	// 滑动窗口的输出是一个切片，包含窗口时间范围内的所有传感器数据
	sensorDataList, ok := windowData.([]*SensorData)
	if !ok {
		return fmt.Sprintf("窗口数据类型错误: %T", windowData)
	}

	if len(sensorDataList) == 0 {
		return "空窗口"
	}

	// 计算窗口统计信息
	var totalTemp float64
	var deviceIDs []string
	datacenterCount := make(map[string]int)

	for _, data := range sensorDataList {
		totalTemp += data.Temperature
		deviceIDs = append(deviceIDs, data.DeviceID)
		datacenterCount[data.DatacenterName]++
	}

	avgTemp := totalTemp / float64(len(sensorDataList))

	result := SensorWindowResult{
		SensorCount:     len(sensorDataList),
		DeviceIDs:       deviceIDs,
		AverageTemp:     avgTemp,
		DatacenterStats: datacenterCount,
	}

	return result
}

// SensorWindowResult 传感器窗口处理结果
type SensorWindowResult struct {
	SensorCount     int            `json:"sensor_count"`
	DeviceIDs       []string       `json:"device_ids"`
	AverageTemp     float64        `json:"average_temp"`
	DatacenterStats map[string]int `json:"datacenter_stats"`
}

func (swr SensorWindowResult) String() string {
	return fmt.Sprintf(
		"传感器滑动窗口 | 传感器数: %d | 设备IDs: %v | 平均温度: %.2f°C | 数据中心统计: %v",
		swr.SensorCount,
		swr.DeviceIDs,
		swr.AverageTemp,
		swr.DatacenterStats,
	)
}
