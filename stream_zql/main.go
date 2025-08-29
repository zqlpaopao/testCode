package main

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/rulego/streamsql"
)

// https://mp.weixin.qq.com/s/H_Z7IDOwL2gfmDJoObpRqw
func main() {
	if err := recover(); err != nil {
		fmt.Println(err)
		fmt.Println()
		fmt.Println(string(debug.Stack()))
	}

	test()
	time.Sleep(300 * time.Second)

}

func test() {
	sql := streamsql.New()
	//defer sql.Stop()

	var sqls string
	//定义处理逻辑
	// 普通逻辑
	sqls = `
	SELECT 
      deviceId,
      UPPER(deviceType) as device_type,
      temperature,
      NOW() as process_time
    FROM stream 
    WHERE temperature > 30
	
`

	//-- 	窗口计算
	//-- 统计每个设备最近10秒的平均温度（滚动窗口）

	sqls = `
	
SELECT 
  deviceId,
  AVG(temperature) as avg_temp,
  window_start() as window_start,
  window_end() as window_end
FROM stream 
GROUP BY deviceId, TumblingWindow('10s')
`

	// 3. 执行 SQL，加载处理逻辑
	if err := sql.Execute(sqls); err != nil {
		panic(err) // 处理 SQL 解析或执行错误
	}
	fmt.Println(2)

	// 5. 定义结果输出逻辑（Sink）
	sql.AddSink(func(results []map[string]interface{}) {
		fmt.Println(3)
		// results 是处理后的结果集（多条记录）
		for _, result := range results {
			fmt.Printf("处理后结果: %+v\n", result)
		}
	})

	go func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println()
			fmt.Println(string(debug.Stack()))
		}
		for {
			// 4、模拟3条传感器数据（其中2条温度>30℃）
			sensorData := []map[string]interface{}{
				{"deviceId": "sensor_01", "deviceType": "temperature", "temperature": 25.6},
				{"deviceId": "sensor_02", "deviceType": "temperature", "temperature": 32.1},
				{"deviceId": "sensor_03", "deviceType": "temperature", "temperature": 35.8},
			}
			// 逐个发送数据到流引擎
			for _, data := range sensorData {
				fmt.Println(1)
				sql.Emit(data) // 将数据注入流处理逻辑
				//time.Sleep(500 * time.Millisecond) // 模拟实时数据间隔
			}
			time.Sleep(3 * time.Second)
		}

	}()

	// 等待所有数据处理完成
	//time.Sleep(7 * time.Second)

	fmt.Println("\n\n所有数据处理完毕")

}
