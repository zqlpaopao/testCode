package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	ext "github.com/reugn/go-streams/extension"
	"test/stream_go/pkg"
)

// 演示如何在接收到系统信号时优雅关闭滑动窗口
func main() {
	fmt.Println("=== 优雅关闭滑动窗口演示 ===")

	// 创建数据源
	source := ext.NewChanSource(pkg.CreateEventChan(time.Millisecond * 500))

	// 创建运行时窗口管理器
	manager := pkg.NewRuntimeWindowManager(source, 5)

	// 添加一些窗口
	fmt.Println("📝 添加初始窗口...")
	manager.AddNewWindow(3*time.Second, 1*time.Second, "anomaly_detection")
	manager.AddNewWindow(5*time.Second, 2*time.Second, "performance_monitor")
	manager.AddNewWindow(4*time.Second, 1500*time.Millisecond, "alert_system")

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个goroutine来监听信号
	go func() {
		sig := <-sigChan
		fmt.Printf("\n🔔 收到信号: %v\n", sig)
		fmt.Println("开始优雅关闭...")

		// 优雅关闭，给10秒超时
		err := manager.ShutdownWithTimeout(10 * time.Second)
		if err != nil {
			fmt.Printf("❌ 关闭失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("👋 程序已安全退出")
		os.Exit(0)
	}()

	fmt.Println("🚀 程序运行中... 按 Ctrl+C 来测试优雅关闭")
	fmt.Println("💡 程序会等待所有滑动窗口处理完剩余数据后再退出")

	// 模拟长时间运行
	for {
		time.Sleep(1 * time.Second)
		if manager.IsShuttingDown() {
			break
		}
	}

	// 等待关闭完成
	time.Sleep(2 * time.Second)
}