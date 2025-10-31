package pkg

import (
	"fmt"
	"math/rand"
	"time"
)

// CreateEventChan 创建事件数据源
func CreateEventChan(interval time.Duration) chan interface{} {
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