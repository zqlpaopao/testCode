// https://www.jianshu.com/p/3209edd28187

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"test/stream_go/pkg"
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

	// 运行动态滑动窗口示例
	pkg.DemonstrateRuntimeWindowAddition()

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

// officialSlidingWindowExample 使用官方 SlidingWindow 的示例 - 多个不同长度的滑动窗口
func officialSlidingWindowExample() {
	println("=== 多个滑动窗口示例开始 ===")

	// 创建单一数据源
	source := ext.NewChanSource(createEventChan(time.Millisecond * 500))

	// 创建多个不同长度的滑动窗口
	// 窗口1：5秒窗口，每1秒滑动，用于短期分析
	slidingWindow5s := flow.NewSlidingWindow[*SensorData](5*time.Second, 1*time.Second)

	// 窗口2：10秒窗口，每2秒滑动，用于中期分析
	slidingWindow10s := flow.NewSlidingWindow[*SensorData](10*time.Second, 2*time.Second)

	// 窗口3：15秒窗口，每3秒滑动，用于长期分析
	slidingWindow15s := flow.NewSlidingWindow[*SensorData](15*time.Second, 3*time.Second)

	// 为每个窗口创建不同的处理器
	processor5s := flow.NewMap(processWindow5s, 1) //该 Map 操作使用 1个 goroutine 来处理数据
	processor10s := flow.NewMap(processWindow10s, 1)
	processor15s := flow.NewMap(processWindow15s, 1)

	// 输出结果
	sink := ext.NewStdoutSink()

	// 使用 FanOut 将同一数据源分发到多个滑动窗口
	go func() {
		fanOut := flow.FanOut(source, 3) //3-分支数量

		// 创建三个处理管道
		pipeline1 := fanOut[0].Via(slidingWindow5s).Via(processor5s)
		pipeline2 := fanOut[1].Via(slidingWindow10s).Via(processor10s)
		pipeline3 := fanOut[2].Via(slidingWindow15s).Via(processor15s)

		// 合并所有管道的输出到同一个 sink
		flow.Merge(pipeline1, pipeline2, pipeline3).To(sink)
	}()

	// 让程序运行一段时间以观察效果
	go func() {
		time.Sleep(30 * time.Second)
		println("=== 多个滑动窗口示例结束 ===")
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

			fmt.Printf("🔄 生成数据 #%d: %s (DeviceID: %s, Temp: %.1f°C)\n",
				counter, sensorData.DeviceID, sensorData.DeviceID, sensorData.Temperature)
			nc <- sensorData

			if counter >= 50 { // 生成50个数据后停止，可以根据需要调整
				fmt.Println("✅ 数据生成完成，共生成50条数据")
				ticker.Stop()
				return
			}
		}
	}()
	return nc
}

func processWindow1(windowData interface{}) interface{} {

	fmt.Println("第二个窗口", windowData)
	return "第二个窗口"
}

// 全局窗口计数器，用于跟踪窗口序号
var windowCounter int = 0

// processWindow 处理滑动窗口的结果 - 显示窗口内数据并实现分层分类统计
func processWindow(windowData interface{}) interface{} {
	windowCounter++
	// 滑动窗口的输出是一个切片，包含窗口时间范围内的所有传感器数据
	sensorDataList, ok := windowData.([]*SensorData)
	if !ok {
		return fmt.Sprintf("窗口数据类型错误: %T", windowData)
	}

	if len(sensorDataList) == 0 {
		return "空窗口"
	}

	// 显示当前滑动窗口内的所有原始数据
	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("\n🔍 ===== 第%d个滑动窗口数据详情 [%s] (共%d条数据) =====\n", windowCounter, currentTime, len(sensorDataList))

	// 显示窗口时间范围信息
	windowEndTime := time.Now()
	windowStartTime := windowEndTime.Add(-5 * time.Second) // 窗口大小是5秒
	fmt.Printf("📅 窗口时间范围: %s ~ %s (窗口大小: 5秒)\n",
		windowStartTime.Format("15:04:05"), windowEndTime.Format("15:04:05"))

	fmt.Println("📋 窗口内数据列表:")
	for i, data := range sensorDataList {
		fmt.Printf("  [%02d] DeviceID: %-10s | Type: %-11s | Temp: %5.1f°C | DC: %s | Room: %s | Rack: %s\n",
			i+1, data.DeviceID, data.DeviceType, data.Temperature,
			data.DatacenterName, data.Room, data.Rack)
	}
	fmt.Println("🔍 ===== 窗口数据详情结束 =====\n")

	// 第一层：按 rack 分类统计
	rackStats := make(map[string][]*SensorData)
	for _, data := range sensorDataList {
		rackStats[data.Rack] = append(rackStats[data.Rack], data)
	}

	// 分层分类结果
	result := &HierarchicalWindowResult{
		TotalSensors:    len(sensorDataList),
		RackStats:       make(map[string]int),
		RoomStats:       make(map[string]int),
		DatacenterStats: make(map[string]int),
		Details:         make(map[string]interface{}),
	}

	// 处理每个 rack
	for rack, rackData := range rackStats {
		rackCount := len(rackData)
		result.RackStats[rack] = rackCount

		// 如果 rack 数据大于2个，按 room 进一步分类
		if rackCount > 2 {
			roomStats := make(map[string][]*SensorData)
			for _, data := range rackData {
				roomKey := fmt.Sprintf("%s-%s", rack, data.Room)
				roomStats[roomKey] = append(roomStats[roomKey], data)
			}

			// 处理每个 room
			for roomKey, roomData := range roomStats {
				roomCount := len(roomData)
				result.RoomStats[roomKey] = roomCount

				// 如果 room 数据大于2个，按 datacenter 进一步分类
				if roomCount > 2 {
					datacenterStats := make(map[string][]*SensorData)
					for _, data := range roomData {
						dcKey := fmt.Sprintf("%s-%s", roomKey, data.DatacenterName)
						datacenterStats[dcKey] = append(datacenterStats[dcKey], data)
					}

					// 处理每个 datacenter
					for dcKey, dcData := range datacenterStats {
						result.DatacenterStats[dcKey] = len(dcData)
					}
				}
			}
		}
	}

	// 简化输出，专注于窗口数据验证
	fmt.Printf("📊 第%d个窗口统计: 总数据%d条，涉及%d个不同的Rack\n",
		windowCounter, result.TotalSensors, len(result.RackStats))
	fmt.Println("==================")

	// 暂时注释掉详细的分层统计结果，专注于验证窗口数据
	// return result
	return fmt.Sprintf("第%d个滑动窗口处理完成 - 数据条数: %d", windowCounter, len(sensorDataList))
}

// HierarchicalWindowResult 分层窗口处理结果
type HierarchicalWindowResult struct {
	TotalSensors    int                    `json:"total_sensors"`
	RackStats       map[string]int         `json:"rack_stats"`       // rack维度统计
	RoomStats       map[string]int         `json:"room_stats"`       // room维度统计 (只有rack>2时才有)
	DatacenterStats map[string]int         `json:"datacenter_stats"` // datacenter维度统计 (只有room>2时才有)
	Details         map[string]interface{} `json:"details"`          // 详细信息
}

func (hwr *HierarchicalWindowResult) String() string {
	result := fmt.Sprintf("=== 分层统计结果 ===\n")
	result += fmt.Sprintf("总传感器数: %d\n", hwr.TotalSensors)

	result += fmt.Sprintf("\n【Rack 维度统计】:\n")
	for rack, count := range hwr.RackStats {
		result += fmt.Sprintf("  Rack %s: %d个传感器", rack, count)
		if count > 2 {
			result += " (>2个，进入Room维度分析)"
		}
		result += "\n"
	}

	if len(hwr.RoomStats) > 0 {
		result += fmt.Sprintf("\n【Room 维度统计】(Rack数据>2的情况):\n")
		for room, count := range hwr.RoomStats {
			result += fmt.Sprintf("  %s: %d个传感器", room, count)
			if count > 2 {
				result += " (>2个，进入Datacenter维度分析)"
			}
			result += "\n"
		}
	}

	if len(hwr.DatacenterStats) > 0 {
		result += fmt.Sprintf("\n【Datacenter 维度统计】(Room数据>2的情况):\n")
		for datacenter, count := range hwr.DatacenterStats {
			result += fmt.Sprintf("  %s: %d个传感器\n", datacenter, count)
		}
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

// processWindow5s 处理5秒滑动窗口的结果 - 短期分析
func processWindow5s(windowData interface{}) interface{} {
	sensorDataList, ok := windowData.([]*SensorData)
	if !ok {
		return fmt.Sprintf("5秒窗口数据类型错误: %T", windowData)
	}

	if len(sensorDataList) == 0 {
		return "5秒窗口为空"
	}

	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("\n🔥 ===== 5秒滑动窗口 [%s] (共%d条数据) - 短期实时分析 =====\n", currentTime, len(sensorDataList))

	// 计算平均温度
	var totalTemp float64
	for _, data := range sensorDataList {
		totalTemp += data.Temperature
	}
	avgTemp := totalTemp / float64(len(sensorDataList))

	fmt.Printf("📊 短期分析: 平均温度 %.1f°C, 数据点 %d个\n", avgTemp, len(sensorDataList))
	fmt.Println("🔥 ===== 5秒窗口结束 =====\n")

	return fmt.Sprintf("5秒窗口: %d条数据, 平均温度%.1f°C", len(sensorDataList), avgTemp)
}

// processWindow10s 处理10秒滑动窗口的结果 - 中期分析
func processWindow10s(windowData interface{}) interface{} {
	sensorDataList, ok := windowData.([]*SensorData)
	if !ok {
		return fmt.Sprintf("10秒窗口数据类型错误: %T", windowData)
	}

	if len(sensorDataList) == 0 {
		return "10秒窗口为空"
	}

	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("\n🌊 ===== 10秒滑动窗口 [%s] (共%d条数据) - 中期趋势分析 =====\n", currentTime, len(sensorDataList))

	// 统计设备类型分布
	deviceTypeStats := make(map[string]int)
	var totalTemp float64
	for _, data := range sensorDataList {
		deviceTypeStats[data.DeviceType]++
		totalTemp += data.Temperature
	}
	avgTemp := totalTemp / float64(len(sensorDataList))

	fmt.Printf("📈 中期分析: 平均温度 %.1f°C, 设备类型分布: %v\n", avgTemp, deviceTypeStats)
	fmt.Println("🌊 ===== 10秒窗口结束 =====\n")

	return fmt.Sprintf("10秒窗口: %d条数据, 平均温度%.1f°C, 设备类型%d种", len(sensorDataList), avgTemp, len(deviceTypeStats))
}

// processWindow15s 处理15秒滑动窗口的结果 - 长期分析
func processWindow15s(windowData interface{}) interface{} {
	sensorDataList, ok := windowData.([]*SensorData)
	if !ok {
		return fmt.Sprintf("15秒窗口数据类型错误: %T", windowData)
	}

	if len(sensorDataList) == 0 {
		return "15秒窗口为空"
	}

	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("\n🏔️ ===== 15秒滑动窗口 [%s] (共%d条数据) - 长期统计分析 =====\n", currentTime, len(sensorDataList))

	// 详细统计分析
	datacenterStats := make(map[string]int)
	roomStats := make(map[string]int)
	var totalTemp, minTemp, maxTemp float64
	minTemp = 999.0
	maxTemp = -999.0

	for _, data := range sensorDataList {
		datacenterStats[data.DatacenterName]++
		roomStats[data.Room]++
		totalTemp += data.Temperature
		if data.Temperature < minTemp {
			minTemp = data.Temperature
		}
		if data.Temperature > maxTemp {
			maxTemp = data.Temperature
		}
	}
	avgTemp := totalTemp / float64(len(sensorDataList))

	fmt.Printf("📊 长期分析: 平均温度 %.1f°C (范围: %.1f°C - %.1f°C)\n", avgTemp, minTemp, maxTemp)
	fmt.Printf("🏢 数据中心分布: %v\n", datacenterStats)
	fmt.Printf("🏠 房间分布: %v\n", roomStats)
	fmt.Println("🏔️ ===== 15秒窗口结束 =====\n")

	return fmt.Sprintf("15秒窗口: %d条数据, 温度范围%.1f-%.1f°C, 数据中心%d个",
		len(sensorDataList), minTemp, maxTemp, len(datacenterStats))
}
