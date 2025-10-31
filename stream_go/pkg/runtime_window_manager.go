package pkg

import (
	"fmt"
	"sync"

	"time"

	"github.com/reugn/go-streams"
	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

// RuntimeWindowManager 运行时窗口管理器
// 用于在程序运行时动态添加新的滑动窗口到现有数据源
type RuntimeWindowManager struct {
	originalSource streams.Outlet         // 原始数据源
	fanOutSource   []streams.Flow         // FanOut后的数据源分支
	windows        map[string]*WindowInfo // 窗口信息映射
	windowCount    int                    // 窗口计数器
	mutex          sync.RWMutex           // 读写锁
	maxFanOut      int                    // 最大分支数
	currentFanOut  int                    // 当前分支数
}

// WindowInfo 窗口信息
type WindowInfo struct {
	ID            string
	WindowSize    time.Duration
	SlideInterval time.Duration
	ProcessorFunc string
	Active        bool
	CreatedAt     time.Time
}

// NewRuntimeWindowManager 创建运行时窗口管理器
func NewRuntimeWindowManager(source streams.Outlet, maxFanOut int) *RuntimeWindowManager {
	return &RuntimeWindowManager{
		originalSource: source,
		windows:        make(map[string]*WindowInfo),
		maxFanOut:      maxFanOut,
		currentFanOut:  0,
	}
}

// AddNewWindow 在运行时添加新的滑动窗口
func (rwm *RuntimeWindowManager) AddNewWindow(windowSize, slideInterval time.Duration, processorType string) (string, error) {
	rwm.mutex.Lock()
	defer rwm.mutex.Unlock()

	// 检查是否还能添加新窗口
	if rwm.currentFanOut >= rwm.maxFanOut {
		return "", fmt.Errorf("已达到最大窗口数量限制: %d", rwm.maxFanOut)
	}

	// 生成窗口ID
	rwm.windowCount++
	windowID := fmt.Sprintf("runtime_window_%d_%s", rwm.windowCount, processorType)

	// 如果是第一次添加窗口，需要创建FanOut
	if rwm.currentFanOut == 0 {
		rwm.fanOutSource = flow.FanOut(rwm.originalSource, rwm.maxFanOut)
		fmt.Printf("🔄 创建FanOut，最大分支数: %d\n", rwm.maxFanOut)
	}

	// 使用下一个可用的分支
	branchIndex := rwm.currentFanOut
	rwm.currentFanOut++

	// 创建窗口信息
	windowInfo := &WindowInfo{
		ID:            windowID,
		WindowSize:    windowSize,
		SlideInterval: slideInterval,
		ProcessorFunc: processorType,
		Active:        true,
		CreatedAt:     time.Now(),
	}

	rwm.windows[windowID] = windowInfo

	// 启动新的窗口处理管道
	go rwm.startNewWindowPipeline(windowID, branchIndex, windowSize, slideInterval, processorType)

	fmt.Printf("✅ 成功添加运行时窗口: %s (窗口大小: %v, 滑动间隔: %v)\n",
		windowID, windowSize, slideInterval)

	return windowID, nil
}

// startNewWindowPipeline 启动新的窗口处理管道
func (rwm *RuntimeWindowManager) startNewWindowPipeline(windowID string, branchIndex int,
	windowSize, slideInterval time.Duration, processorType string) {

	// 创建滑动窗口
	slidingWindow := flow.NewSlidingWindow[*SensorData](windowSize, slideInterval)

	// 创建处理器
	processor := flow.NewMap(func(windowData interface{}) interface{} {
		return rwm.processWindowData(windowID, processorType, windowData)
	}, 1)

	// 创建输出
	sink := ext.NewStdoutSink()

	// 构建处理管道：数据源分支 -> 滑动窗口 -> 处理器 -> 输出
	rwm.fanOutSource[branchIndex].Via(slidingWindow).Via(processor).To(sink)

	fmt.Printf("🚀 窗口 %s 的处理管道已启动 (使用分支 %d)\n", windowID, branchIndex)
}

// processWindowData 处理窗口数据
func (rwm *RuntimeWindowManager) processWindowData(windowID, processorType string, windowData interface{}) interface{} {
	sensorDataList, ok := windowData.([]*SensorData)
	if !ok {
		return fmt.Sprintf("%s: 数据类型错误: %T", windowID, windowData)
	}

	if len(sensorDataList) == 0 {
		return fmt.Sprintf("%s: 空窗口", windowID)
	}

	currentTime := time.Now().Format("15:04:05")

	switch processorType {
	case "anomaly_detection":
		return rwm.processAnomalyDetection(windowID, currentTime, sensorDataList)
	case "performance_monitor":
		return rwm.processPerformanceMonitor(windowID, currentTime, sensorDataList)
	case "alert_system":
		return rwm.processAlertSystem(windowID, currentTime, sensorDataList)
	case "data_quality":
		return rwm.processDataQuality(windowID, currentTime, sensorDataList)
	default:
		return rwm.processCustom(windowID, currentTime, sensorDataList, processorType)
	}
}

// processAnomalyDetection 异常检测处理
func (rwm *RuntimeWindowManager) processAnomalyDetection(windowID, currentTime string, data []*SensorData) interface{} {
	var totalTemp float64
	var anomalyCount int

	for _, sensor := range data {
		totalTemp += sensor.Temperature
		// 简单的异常检测：温度超过40度或低于15度
		if sensor.Temperature > 40.0 || sensor.Temperature < 15.0 {
			anomalyCount++
		}
	}
	avgTemp := totalTemp / float64(len(data))

	result := fmt.Sprintf("🚨 [%s] %s: 异常检测 - %d条数据, 平均温度%.1f°C, 异常数据%d条",
		currentTime, windowID, len(data), avgTemp, anomalyCount)

	if anomalyCount > 0 {
		fmt.Printf("⚠️ %s\n", result)
	} else {
		fmt.Printf("✅ %s\n", result)
	}
	return result
}

// processPerformanceMonitor 性能监控处理
func (rwm *RuntimeWindowManager) processPerformanceMonitor(windowID, currentTime string, data []*SensorData) interface{} {
	deviceStats := make(map[string]int)
	datacenterStats := make(map[string]int)

	for _, sensor := range data {
		deviceStats[sensor.DeviceID]++
		datacenterStats[sensor.DatacenterName]++
	}

	result := fmt.Sprintf("📊 [%s] %s: 性能监控 - %d条数据, 活跃设备%d个, 数据中心%d个",
		currentTime, windowID, len(data), len(deviceStats), len(datacenterStats))
	fmt.Println(result)
	return result
}

// processAlertSystem 告警系统处理
func (rwm *RuntimeWindowManager) processAlertSystem(windowID, currentTime string, data []*SensorData) interface{} {
	var highTempCount, lowTempCount int
	var totalTemp float64

	for _, sensor := range data {
		totalTemp += sensor.Temperature
		if sensor.Temperature > 35.0 {
			highTempCount++
		}
		if sensor.Temperature < 20.0 {
			lowTempCount++
		}
	}
	avgTemp := totalTemp / float64(len(data))

	alertLevel := "正常"
	if highTempCount > len(data)/2 {
		alertLevel = "高温告警"
	} else if lowTempCount > len(data)/2 {
		alertLevel = "低温告警"
	}

	result := fmt.Sprintf("🔔 [%s] %s: 告警系统 - %d条数据, 平均温度%.1f°C, 状态:%s",
		currentTime, windowID, len(data), avgTemp, alertLevel)
	fmt.Println(result)
	return result
}

// processDataQuality 数据质量处理
func (rwm *RuntimeWindowManager) processDataQuality(windowID, currentTime string, data []*SensorData) interface{} {
	var validCount, invalidCount int
	deviceTypeStats := make(map[string]int)

	for _, sensor := range data {
		// 简单的数据质量检查
		if sensor.Temperature >= 0 && sensor.Temperature <= 100 &&
			sensor.DeviceID != "" && sensor.DatacenterName != "" {
			validCount++
		} else {
			invalidCount++
		}
		deviceTypeStats[sensor.DeviceType]++
	}

	qualityPercent := float64(validCount) / float64(len(data)) * 100

	result := fmt.Sprintf("🔍 [%s] %s: 数据质量 - %d条数据, 有效数据%.1f%%, 设备类型%d种",
		currentTime, windowID, len(data), qualityPercent, len(deviceTypeStats))
	fmt.Println(result)
	return result
}

// processCustom 自定义处理
func (rwm *RuntimeWindowManager) processCustom(windowID, currentTime string, data []*SensorData, processorType string) interface{} {
	result := fmt.Sprintf("⚙️ [%s] %s: 自定义处理(%s) - %d条数据",
		currentTime, windowID, processorType, len(data))
	fmt.Println(result)
	return result
}

// GetWindowInfo 获取窗口信息
func (rwm *RuntimeWindowManager) GetWindowInfo(windowID string) (*WindowInfo, bool) {
	rwm.mutex.RLock()
	defer rwm.mutex.RUnlock()

	info, exists := rwm.windows[windowID]
	return info, exists
}

// ListActiveWindows 列出所有活跃窗口
func (rwm *RuntimeWindowManager) ListActiveWindows() []*WindowInfo {
	rwm.mutex.RLock()
	defer rwm.mutex.RUnlock()

	var activeWindows []*WindowInfo
	for _, window := range rwm.windows {
		if window.Active {
			activeWindows = append(activeWindows, window)
		}
	}
	return activeWindows
}

// GetWindowCount 获取当前窗口数量
func (rwm *RuntimeWindowManager) GetWindowCount() int {
	rwm.mutex.RLock()
	defer rwm.mutex.RUnlock()
	return rwm.currentFanOut
}

// GetMaxCapacity 获取最大容量
func (rwm *RuntimeWindowManager) GetMaxCapacity() int {
	return rwm.maxFanOut
}

// DemonstrateRuntimeWindowAddition 演示如何在运行时添加新窗口
func DemonstrateRuntimeWindowAddition() {
	fmt.Println("=== 运行时动态添加滑动窗口演示 ===")

	// 创建数据源（模拟已经运行的程序）
	source := ext.NewChanSource(CreateEventChan(time.Millisecond * 500))

	// 创建运行时窗口管理器，最多支持8个窗口
	manager := NewRuntimeWindowManager(source, 8)

	fmt.Println("📡 程序已启动，数据源正在运行...")
	time.Sleep(2 * time.Second)

	// 模拟运行时添加新窗口的场景
	go func() {
		// 场景1：2秒后，需要添加异常检测窗口
		time.Sleep(2 * time.Second)
		fmt.Println("\n🔥 业务需求：需要添加异常检测功能")
		windowID1, err := manager.AddNewWindow(5*time.Second, 1*time.Second, "anomaly_detection")
		if err != nil {
			fmt.Printf("❌ 添加窗口失败: %v\n", err)
		} else {
			fmt.Printf("✅ 异常检测窗口已添加: %s\n", windowID1)
		}

		// 场景2：5秒后，需要添加性能监控窗口
		time.Sleep(3 * time.Second)
		fmt.Println("\n📊 业务需求：需要添加性能监控功能")
		windowID2, err := manager.AddNewWindow(8*time.Second, 2*time.Second, "performance_monitor")
		if err != nil {
			fmt.Printf("❌ 添加窗口失败: %v\n", err)
		} else {
			fmt.Printf("✅ 性能监控窗口已添加: %s\n", windowID2)
		}

		// 场景3：8秒后，需要添加告警系统窗口
		time.Sleep(3 * time.Second)
		fmt.Println("\n🔔 业务需求：需要添加告警系统功能")
		windowID3, err := manager.AddNewWindow(6*time.Second, 2*time.Second, "alert_system")
		if err != nil {
			fmt.Printf("❌ 添加窗口失败: %v\n", err)
		} else {
			fmt.Printf("✅ 告警系统窗口已添加: %s\n", windowID3)
		}

		// 场景4：12秒后，需要添加数据质量检查窗口
		time.Sleep(4 * time.Second)
		fmt.Println("\n🔍 业务需求：需要添加数据质量检查功能")
		windowID4, err := manager.AddNewWindow(10*time.Second, 3*time.Second, "data_quality")
		if err != nil {
			fmt.Printf("❌ 添加窗口失败: %v\n", err)
		} else {
			fmt.Printf("✅ 数据质量检查窗口已添加: %s\n", windowID4)
		}

		// 显示当前所有活跃窗口
		time.Sleep(2 * time.Second)
		fmt.Println("\n📋 当前所有活跃窗口:")
		activeWindows := manager.ListActiveWindows()
		for i, window := range activeWindows {
			fmt.Printf("  %d. %s (窗口大小: %v, 滑动间隔: %v, 创建时间: %s)\n",
				i+1, window.ID, window.WindowSize, window.SlideInterval,
				window.CreatedAt.Format("15:04:05"))
		}

		fmt.Printf("\n📊 当前窗口数量: %d/%d\n", manager.GetWindowCount(), manager.GetMaxCapacity())
	}()

	// 让程序运行足够长的时间来观察效果
	time.Sleep(25 * time.Second)
	fmt.Println("=== 运行时动态添加滑动窗口演示结束 ===")
}
