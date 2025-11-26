package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// PluginSystem 插件系统主结构
type PluginSystem struct {
	dbManager    *DatabaseManager
	memoryCache  *MemoryCache
	config       *Config
}

// Config 系统配置
type Config struct {
	DatabaseDSN     string `json:"database_dsn"`
	CheckInterval   int    `json:"check_interval_seconds"`
	MaxCacheSize    int    `json:"max_cache_size"`
	EnableAutoSync  bool   `json:"enable_auto_sync"`
}

// NewPluginSystem 创建新的插件系统
func NewPluginSystem(config *Config) (*PluginSystem, error) {
	// 创建数据库管理器
	dbManager, err := NewDatabaseManager(config.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("创建数据库管理器失败: %v", err)
	}

	// 创建内存缓存管理器
	memoryCache, err := NewMemoryCache(dbManager)
	if err != nil {
		return nil, fmt.Errorf("创建内存缓存管理器失败: %v", err)
	}

	return &PluginSystem{
		dbManager:   dbManager,
		memoryCache: memoryCache,
		config:      config,
	}, nil
}

// Start 启动插件系统
func (ps *PluginSystem) Start() error {
	fmt.Println("🚀 MySQL 插件系统启动中...")
	
	// 从数据库加载所有插件到内存
	if err := ps.memoryCache.LoadAllPlugins(); err != nil {
		return fmt.Errorf("加载插件到内存失败: %v", err)
	}

	fmt.Println("✅ 插件系统启动完成")
	return nil
}

// Stop 停止插件系统
func (ps *PluginSystem) Stop() {
	fmt.Println("🛑 正在关闭插件系统...")
	
	ps.memoryCache.Stop()
	ps.dbManager.Close()
	
	fmt.Println("✅ 插件系统已关闭")
}

func main() {
	// 默认配置
	config := &Config{
		DatabaseDSN:    "root:password@tcp(localhost:3306)/plugin_system?charset=utf8mb4&parseTime=True&loc=Local",
		CheckInterval:  30,
		MaxCacheSize:   100,
		EnableAutoSync: true,
	}

	// 尝试从配置文件加载
	if configData, err := os.ReadFile("config.json"); err == nil {
		json.Unmarshal(configData, config)
	}

	// 创建插件系统
	ps, err := NewPluginSystem(config)
	if err != nil {
		log.Fatalf("❌ 创建插件系统失败: %v", err)
	}

	// 启动系统
	if err := ps.Start(); err != nil {
		log.Fatalf("❌ 启动插件系统失败: %v", err)
	}

	// 设置优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		fmt.Println("\n🔄 收到关闭信号...")
		ps.Stop()
		os.Exit(0)
	}()

	// 显示系统信息
	showSystemInfo(ps)
	
	// 交互式命令循环
	runInteractiveMode(ps)
}

// showSystemInfo 显示系统信息
func showSystemInfo(ps *PluginSystem) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎯 MySQL 动态插件系统")
	fmt.Println("   - 插件存储在 MySQL 数据库中")
	fmt.Println("   - 自动加载插件到内存缓存")
	fmt.Println("   - 支持实时更新和热重载")
	fmt.Println("   - 定期检查数据库变化")
	fmt.Println(strings.Repeat("=", 60))

	// 显示插件统计
	if stats, err := ps.dbManager.GetPluginStats(); err == nil {
		fmt.Printf("\n📊 数据库统计:\n")
		if totalPlugins, ok := stats["total_plugins"].(int); ok {
			fmt.Printf("   总插件数: %d\n", totalPlugins)
		}
		if byStatus, ok := stats["by_status"].(map[string]int); ok {
			for status, count := range byStatus {
				fmt.Printf("   %s: %d\n", status, count)
			}
		}
	}

	// 显示缓存统计
	cacheStats := ps.memoryCache.GetCacheStats()
	if totalCached, ok := cacheStats["total_cached_plugins"].(int); ok {
		fmt.Printf("\n💾 内存缓存统计:\n")
		fmt.Printf("   已缓存插件数: %d\n", totalCached)
		if lastCheck, ok := cacheStats["last_check"].(time.Time); ok {
			fmt.Printf("   最后检查时间: %s\n", lastCheck.Format("2006-01-02 15:04:05"))
		}
	}

	fmt.Println()
}

// runInteractiveMode 运行交互模式
func runInteractiveMode(ps *PluginSystem) {
	scanner := bufio.NewScanner(os.Stdin)
	
	showHelp()
	
	for {
		fmt.Print("🔧 请输入命令 (help 查看帮助): ")
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		
		parts := strings.Fields(input)
		command := parts[0]
		
		switch command {
		case "help", "h":
			showHelp()
			
		case "list", "ls":
			listPlugins(ps)
			
		case "call", "c":
			if len(parts) < 3 {
				fmt.Println("❌ 用法: call <插件名> <函数名> [参数...]")
				continue
			}
			callFunction(ps, parts[1], parts[2], parts[3:]...)
			
		case "add", "create":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: add <插件名>")
				continue
			}
			addPluginInteractive(ps, parts[1], scanner)
			
		case "update", "u":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: update <插件名>")
				continue
			}
			updatePlugin(ps, parts[1])
			
		case "delete", "del":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: delete <插件名>")
				continue
			}
			deletePlugin(ps, parts[1])
			
		case "refresh", "r":
			if len(parts) > 1 {
				refreshPlugin(ps, parts[1])
			} else {
				refreshAllPlugins(ps)
			}
			
		case "stats", "stat":
			showStats(ps)
			
		case "check":
			checkUpdates(ps)
			
		case "info", "i":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: info <插件名>")
				continue
			}
			showPluginInfo(ps, parts[1])
			
		case "export":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: export <插件名>")
				continue
			}
			exportPlugin(ps, parts[1])
			
		case "import":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: import <文件路径>")
				continue
			}
			importPlugin(ps, parts[1])
			
		case "exit", "quit", "q":
			fmt.Println("👋 再见!")
			ps.Stop()
			return
			
		default:
			fmt.Printf("❌ 未知命令: %s (输入 help 查看帮助)\n", command)
		}
	}
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println("📖 可用命令:")
	fmt.Println("  help, h              - 显示帮助")
	fmt.Println("  list, ls             - 列出所有插件")
	fmt.Println("  call, c <plugin> <func> [args] - 调用插件函数")
	fmt.Println("  add, create <name>   - 添加新插件")
	fmt.Println("  update, u <name>     - 更新插件")
	fmt.Println("  delete, del <name>   - 删除插件")
	fmt.Println("  refresh, r [name]    - 刷新插件 (不指定名称则刷新全部)")
	fmt.Println("  stats, stat          - 显示系统统计信息")
	fmt.Println("  check                - 手动检查数据库更新")
	fmt.Println("  info, i <name>       - 显示插件详细信息")
	fmt.Println("  export <name>        - 导出插件")
	fmt.Println("  import <file>        - 导入插件")
	fmt.Println("  exit, quit, q        - 退出程序")
	fmt.Println()
}

// listPlugins 列出所有插件
func listPlugins(ps *PluginSystem) {
	plugins := ps.memoryCache.ListPlugins()
	
	if len(plugins) == 0 {
		fmt.Println("📭 暂无已加载的插件")
		return
	}
	
	fmt.Printf("📋 已加载的插件 (共 %d 个):\n", len(plugins))
	for name, cached := range plugins {
		functions := make([]string, 0, len(cached.Functions))
		for funcName := range cached.Functions {
			functions = append(functions, funcName)
		}
		fmt.Printf("  🔹 %s (v%s) - 函数: %v\n", name, cached.Version, functions)
		fmt.Printf("     加载时间: %s, 哈希: %s\n", 
			cached.LoadedAt.Format("2006-01-02 15:04:05"), cached.Hash[:12]+"...")
	}
}

// callFunction 调用插件函数
func callFunction(ps *PluginSystem, pluginName, functionName string, args ...string) {
	// 转换参数
	var convertedArgs []interface{}
	for _, arg := range args {
		// 尝试转换为数字
		if num, err := strconv.Atoi(arg); err == nil {
			convertedArgs = append(convertedArgs, num)
		} else if float, err := strconv.ParseFloat(arg, 64); err == nil {
			convertedArgs = append(convertedArgs, float)
		} else {
			convertedArgs = append(convertedArgs, arg)
		}
	}
	
	start := time.Now()
	result, err := ps.memoryCache.CallFunction(pluginName, functionName, convertedArgs...)
	duration := time.Since(start)
	
	if err != nil {
		fmt.Printf("❌ 调用失败: %v\n", err)
	} else {
		fmt.Printf("✅ 结果: %v (耗时: %v)\n", result, duration)
	}
}

// addPluginInteractive 交互式添加插件
func addPluginInteractive(ps *PluginSystem, name string, scanner *bufio.Scanner) {
	fmt.Printf("📝 添加插件 '%s'\n", name)
	
	fmt.Print("版本号 (默认 1.0.0): ")
	scanner.Scan()
	version := strings.TrimSpace(scanner.Text())
	if version == "" {
		version = "1.0.0"
	}
	
	fmt.Print("描述: ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())
	
	fmt.Println("请输入插件源码 (以 'END' 结束):")
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" {
			break
		}
		lines = append(lines, line)
	}
	sourceCode := strings.Join(lines, "\n")
	
	if err := ps.memoryCache.AddPluginFromSource(name, version, description, sourceCode); err != nil {
		fmt.Printf("❌ 添加插件失败: %v\n", err)
	} else {
		fmt.Printf("✅ 插件 '%s' 添加成功\n", name)
	}
}

// updatePlugin 更新插件
func updatePlugin(ps *PluginSystem, name string) {
	if err := ps.memoryCache.RefreshPlugin(name); err != nil {
		fmt.Printf("❌ 更新插件失败: %v\n", err)
	} else {
		fmt.Printf("✅ 插件 '%s' 更新成功\n", name)
	}
}

// deletePlugin 删除插件
func deletePlugin(ps *PluginSystem, name string) {
	if err := ps.dbManager.DeletePlugin(name); err != nil {
		fmt.Printf("❌ 删除插件失败: %v\n", err)
	} else {
		ps.memoryCache.RemovePlugin(name)
		fmt.Printf("✅ 插件 '%s' 删除成功\n", name)
	}
}

// refreshPlugin 刷新单个插件
func refreshPlugin(ps *PluginSystem, name string) {
	if err := ps.memoryCache.RefreshPlugin(name); err != nil {
		fmt.Printf("❌ 刷新插件失败: %v\n", err)
	} else {
		fmt.Printf("✅ 插件 '%s' 刷新成功\n", name)
	}
}

// refreshAllPlugins 刷新所有插件
func refreshAllPlugins(ps *PluginSystem) {
	if err := ps.memoryCache.ForceRefreshAll(); err != nil {
		fmt.Printf("❌ 刷新所有插件失败: %v\n", err)
	} else {
		fmt.Println("✅ 所有插件刷新成功")
	}
}

// showStats 显示系统统计信息
func showStats(ps *PluginSystem) {
	fmt.Println("📊 系统统计信息:")
	
	// 数据库统计
	if dbStats, err := ps.dbManager.GetPluginStats(); err == nil {
		fmt.Println("\n💾 数据库统计:")
		for key, value := range dbStats {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
	
	// 内存缓存统计
	cacheStats := ps.memoryCache.GetCacheStats()
	fmt.Println("\n🧠 内存缓存统计:")
	for key, value := range cacheStats {
		if key != "plugins" {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

// checkUpdates 手动检查更新
func checkUpdates(ps *PluginSystem) {
	fmt.Println("🔍 检查数据库更新...")
	if err := ps.memoryCache.CheckForUpdates(); err != nil {
		fmt.Printf("❌ 检查更新失败: %v\n", err)
	} else {
		fmt.Println("✅ 更新检查完成")
	}
}

// showPluginInfo 显示插件详细信息
func showPluginInfo(ps *PluginSystem, name string) {
	cached, err := ps.memoryCache.GetPlugin(name)
	if err != nil {
		fmt.Printf("❌ 获取插件信息失败: %v\n", err)
		return
	}
	
	fmt.Printf("📋 插件 '%s' 详细信息:\n", name)
	fmt.Printf("  版本: %s\n", cached.Version)
	fmt.Printf("  哈希: %s\n", cached.Hash)
	fmt.Printf("  加载时间: %s\n", cached.LoadedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  函数数量: %d\n", len(cached.Functions))
	fmt.Printf("  可用函数: %v\n", ps.memoryCache.getFunctionNames(cached.Functions))
	
	if cached.Record != nil {
		fmt.Printf("  描述: %s\n", cached.Record.Description)
		fmt.Printf("  创建时间: %s\n", cached.Record.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  更新时间: %s\n", cached.Record.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
}

// exportPlugin 导出插件
func exportPlugin(ps *PluginSystem, name string) {
	record, err := ps.dbManager.GetPlugin(name)
	if err != nil {
		fmt.Printf("❌ 获取插件失败: %v\n", err)
		return
	}
	
	exportData := map[string]interface{}{
		"name":        record.Name,
		"version":     record.Version,
		"description": record.Description,
		"source_code": record.SourceCode,
		"hash":        record.Hash,
		"exported_at": time.Now(),
	}
	
	filename := fmt.Sprintf("%s_v%s.json", name, record.Version)
	data, _ := json.MarshalIndent(exportData, "", "  ")
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("❌ 导出失败: %v\n", err)
	} else {
		fmt.Printf("✅ 插件已导出到: %s\n", filename)
	}
}

// importPlugin 导入插件
func importPlugin(ps *PluginSystem, filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ 读取文件失败: %v\n", err)
		return
	}
	
	var importData map[string]interface{}
	if err := json.Unmarshal(data, &importData); err != nil {
		fmt.Printf("❌ 解析文件失败: %v\n", err)
		return
	}
	
	name, _ := importData["name"].(string)
	version, _ := importData["version"].(string)
	description, _ := importData["description"].(string)
	sourceCode, _ := importData["source_code"].(string)
	
	if err := ps.memoryCache.AddPluginFromSource(name, version, description, sourceCode); err != nil {
		fmt.Printf("❌ 导入插件失败: %v\n", err)
	} else {
		fmt.Printf("✅ 插件 '%s' 导入成功\n", name)
	}
}