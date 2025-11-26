package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"sync"
	"time"
)

// MemoryCache 内存缓存管理器
type MemoryCache struct {
	plugins     map[string]*CachedPlugin
	mutex       sync.RWMutex
	tempDir     string
	dbManager   *DatabaseManager
	lastCheck   time.Time
	checkTicker *time.Ticker
	stopChan    chan bool
}

// CachedPlugin 缓存的插件信息
type CachedPlugin struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Hash        string                 `json:"hash"`
	Plugin      *plugin.Plugin         `json:"-"`
	Functions   map[string]interface{} `json:"functions"`
	LoadedAt    time.Time              `json:"loaded_at"`
	TempSoPath  string                 `json:"temp_so_path"`
	Record      *PluginRecord          `json:"-"`
}

// NewMemoryCache 创建新的内存缓存管理器
func NewMemoryCache(dbManager *DatabaseManager) (*MemoryCache, error) {
	// 创建临时目录
	tempDir, err := ioutil.TempDir("", "plugin_cache_")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %v", err)
	}

	mc := &MemoryCache{
		plugins:   make(map[string]*CachedPlugin),
		tempDir:   tempDir,
		dbManager: dbManager,
		lastCheck: time.Now(),
		stopChan:  make(chan bool),
	}

	// 启动定期检查
	mc.startPeriodicCheck()

	return mc, nil
}

// LoadAllPlugins 从数据库加载所有插件到内存
func (mc *MemoryCache) LoadAllPlugins() error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	plugins, err := mc.dbManager.GetAllPlugins()
	if err != nil {
		return fmt.Errorf("从数据库获取插件失败: %v", err)
	}

	fmt.Printf("📦 开始加载 %d 个插件到内存...\n", len(plugins))

	for _, record := range plugins {
		if err := mc.loadPluginToMemory(record); err != nil {
			fmt.Printf("❌ 加载插件 '%s' 失败: %v\n", record.Name, err)
			continue
		}
		fmt.Printf("✅ 插件 '%s' 已加载到内存\n", record.Name)
	}

	fmt.Printf("🎉 内存缓存初始化完成，共加载 %d 个插件\n", len(mc.plugins))
	return nil
}

// loadPluginToMemory 将单个插件加载到内存
func (mc *MemoryCache) loadPluginToMemory(record *PluginRecord) error {
	// 创建临时 .so 文件
	tempSoPath := filepath.Join(mc.tempDir, record.Name+".so")
	
	if err := ioutil.WriteFile(tempSoPath, record.BinaryData, 0755); err != nil {
		return fmt.Errorf("写入临时 .so 文件失败: %v", err)
	}

	// 加载插件
	p, err := plugin.Open(tempSoPath)
	if err != nil {
		return fmt.Errorf("加载插件失败: %v", err)
	}

	// 获取插件函数
	functions := make(map[string]interface{})
	
	// 尝试加载 GetFunctions 函数
	if sym, err := p.Lookup("GetFunctions"); err == nil {
		if getFuncs, ok := sym.(func() map[string]interface{}); ok {
			functions = getFuncs()
		}
	}

	// 如果没有 GetFunctions，尝试加载常见函数
	if len(functions) == 0 {
		commonFunctions := []string{"Execute", "Process", "Handle", "Run", "Calculate", "Transform"}
		for _, funcName := range commonFunctions {
			if sym, err := p.Lookup(funcName); err == nil {
				functions[funcName] = sym
			}
		}
	}

	// 创建缓存条目
	cached := &CachedPlugin{
		Name:       record.Name,
		Version:    record.Version,
		Hash:       record.Hash,
		Plugin:     p,
		Functions:  functions,
		LoadedAt:   time.Now(),
		TempSoPath: tempSoPath,
		Record:     record,
	}

	mc.plugins[record.Name] = cached
	return nil
}

// GetPlugin 从内存缓存获取插件
func (mc *MemoryCache) GetPlugin(name string) (*CachedPlugin, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	cached, exists := mc.plugins[name]
	if !exists {
		return nil, fmt.Errorf("插件 '%s' 未在内存中", name)
	}

	return cached, nil
}

// CallFunction 调用插件函数
func (mc *MemoryCache) CallFunction(pluginName, functionName string, args ...interface{}) (interface{}, error) {
	cached, err := mc.GetPlugin(pluginName)
	if err != nil {
		return nil, err
	}

	fn, exists := cached.Functions[functionName]
	if !exists {
		return nil, fmt.Errorf("函数 '%s' 在插件 '%s' 中不存在", functionName, pluginName)
	}

	// 根据函数类型调用
	switch f := fn.(type) {
	case func() interface{}:
		return f(), nil
	case func(interface{}) interface{}:
		if len(args) > 0 {
			return f(args[0]), nil
		}
		return f(nil), nil
	case func(...interface{}) interface{}:
		return f(args...), nil
	case func(string) string:
		if len(args) > 0 {
			if str, ok := args[0].(string); ok {
				return f(str), nil
			}
		}
		return nil, fmt.Errorf("函数 '%s' 需要字符串参数", functionName)
	case func(int) int:
		if len(args) > 0 {
			if num, ok := args[0].(int); ok {
				return f(num), nil
			}
		}
		return nil, fmt.Errorf("函数 '%s' 需要整数参数", functionName)
	case func(float64) float64:
		if len(args) > 0 {
			if num, ok := args[0].(float64); ok {
				return f(num), nil
			}
		}
		return nil, fmt.Errorf("函数 '%s' 需要浮点数参数", functionName)
	default:
		return nil, fmt.Errorf("不支持的函数类型: %T", fn)
	}
}

// RefreshPlugin 刷新单个插件
func (mc *MemoryCache) RefreshPlugin(name string) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// 从数据库获取最新插件
	record, err := mc.dbManager.GetPlugin(name)
	if err != nil {
		return fmt.Errorf("从数据库获取插件失败: %v", err)
	}

	// 检查是否需要更新
	if cached, exists := mc.plugins[name]; exists {
		if cached.Hash == record.Hash {
			// 哈希相同，无需更新
			return nil
		}
		
		// 清理旧的临时文件
		if cached.TempSoPath != "" {
			os.Remove(cached.TempSoPath)
		}
		
		// 删除旧缓存
		delete(mc.plugins, name)
	}

	// 加载新版本
	if err := mc.loadPluginToMemory(record); err != nil {
		return fmt.Errorf("加载新版本插件失败: %v", err)
	}

	fmt.Printf("🔄 插件 '%s' 已更新到内存缓存 (新哈希: %s)\n", name, record.Hash)
	return nil
}

// CheckForUpdates 检查数据库中的插件更新
func (mc *MemoryCache) CheckForUpdates() error {
	changes, err := mc.dbManager.GetPluginChanges(mc.lastCheck)
	if err != nil {
		return fmt.Errorf("获取插件变更失败: %v", err)
	}

	if len(changes) == 0 {
		return nil
	}

	fmt.Printf("🔍 检测到 %d 个插件变更\n", len(changes))

	for _, change := range changes {
		switch change.Action {
		case "insert", "update":
			if err := mc.RefreshPlugin(change.PluginName); err != nil {
				fmt.Printf("❌ 刷新插件 '%s' 失败: %v\n", change.PluginName, err)
			}
		case "delete":
			mc.RemovePlugin(change.PluginName)
		}
	}

	mc.lastCheck = time.Now()
	return nil
}

// RemovePlugin 从内存缓存中移除插件
func (mc *MemoryCache) RemovePlugin(name string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	if cached, exists := mc.plugins[name]; exists {
		// 清理临时文件
		if cached.TempSoPath != "" {
			os.Remove(cached.TempSoPath)
		}
		
		delete(mc.plugins, name)
		fmt.Printf("🗑️ 插件 '%s' 已从内存缓存中移除\n", name)
	}
}

// ListPlugins 列出内存中的所有插件
func (mc *MemoryCache) ListPlugins() map[string]*CachedPlugin {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	result := make(map[string]*CachedPlugin)
	for name, cached := range mc.plugins {
		result[name] = cached
	}
	
	return result
}

// startPeriodicCheck 启动定期检查
func (mc *MemoryCache) startPeriodicCheck() {
	mc.checkTicker = time.NewTicker(30 * time.Second) // 每30秒检查一次
	
	go func() {
		for {
			select {
			case <-mc.checkTicker.C:
				if err := mc.CheckForUpdates(); err != nil {
					fmt.Printf("❌ 定期检查更新失败: %v\n", err)
				}
			case <-mc.stopChan:
				return
			}
		}
	}()
}

// AddPluginFromSource 从源码添加插件到数据库和内存
func (mc *MemoryCache) AddPluginFromSource(name, version, description, sourceCode string) error {
	// 计算源码哈希
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(sourceCode)))

	// 创建临时源码文件
	tempSrcPath := filepath.Join(mc.tempDir, name+".go")
	if err := ioutil.WriteFile(tempSrcPath, []byte(sourceCode), 0644); err != nil {
		return fmt.Errorf("写入临时源码文件失败: %v", err)
	}
	defer os.Remove(tempSrcPath)

	// 编译插件
	tempSoPath := filepath.Join(mc.tempDir, name+"_build.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", tempSoPath, tempSrcPath)
	cmd.Env = os.Environ()
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("编译插件失败: %s", string(output))
	}
	defer os.Remove(tempSoPath)

	// 读取编译后的二进制数据
	binaryData, err := ioutil.ReadFile(tempSoPath)
	if err != nil {
		return fmt.Errorf("读取编译后的插件失败: %v", err)
	}

	// 创建插件记录
	record := &PluginRecord{
		Name:        name,
		Version:     version,
		Description: description,
		SourceCode:  sourceCode,
		BinaryData:  binaryData,
		Hash:        hash,
		Status:      "active",
	}

	// 保存到数据库
	if err := mc.dbManager.SavePlugin(record); err != nil {
		return fmt.Errorf("保存插件到数据库失败: %v", err)
	}

	// 加载到内存
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	if err := mc.loadPluginToMemory(record); err != nil {
		return fmt.Errorf("加载插件到内存失败: %v", err)
	}

	fmt.Printf("✅ 插件 '%s' 已添加到数据库和内存缓存\n", name)
	return nil
}

// GetCacheStats 获取缓存统计信息
func (mc *MemoryCache) GetCacheStats() map[string]interface{} {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["total_cached_plugins"] = len(mc.plugins)
	stats["temp_directory"] = mc.tempDir
	stats["last_check"] = mc.lastCheck
	
	pluginStats := make(map[string]map[string]interface{})
	for name, cached := range mc.plugins {
		pluginStats[name] = map[string]interface{}{
			"version":        cached.Version,
			"hash":          cached.Hash,
			"loaded_at":     cached.LoadedAt,
			"function_count": len(cached.Functions),
			"functions":     mc.getFunctionNames(cached.Functions),
		}
	}
	stats["plugins"] = pluginStats
	
	return stats
}

// getFunctionNames 获取函数名列表
func (mc *MemoryCache) getFunctionNames(functions map[string]interface{}) []string {
	names := make([]string, 0, len(functions))
	for name := range functions {
		names = append(names, name)
	}
	return names
}

// Stop 停止内存缓存管理器
func (mc *MemoryCache) Stop() {
	if mc.checkTicker != nil {
		mc.checkTicker.Stop()
	}
	
	close(mc.stopChan)
	
	// 清理临时文件
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	for _, cached := range mc.plugins {
		if cached.TempSoPath != "" {
			os.Remove(cached.TempSoPath)
		}
	}
	
	// 删除临时目录
	os.RemoveAll(mc.tempDir)
	
	fmt.Println("🛑 内存缓存管理器已停止")
}

// ForceRefreshAll 强制刷新所有插件
func (mc *MemoryCache) ForceRefreshAll() error {
	fmt.Println("🔄 强制刷新所有插件...")
	
	plugins, err := mc.dbManager.GetAllPlugins()
	if err != nil {
		return fmt.Errorf("获取数据库插件列表失败: %v", err)
	}

	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// 清理所有现有缓存
	for _, cached := range mc.plugins {
		if cached.TempSoPath != "" {
			os.Remove(cached.TempSoPath)
		}
	}
	mc.plugins = make(map[string]*CachedPlugin)

	// 重新加载所有插件
	for _, record := range plugins {
		if err := mc.loadPluginToMemory(record); err != nil {
			fmt.Printf("❌ 重新加载插件 '%s' 失败: %v\n", record.Name, err)
			continue
		}
	}

	fmt.Printf("✅ 已强制刷新 %d 个插件\n", len(mc.plugins))
	return nil
}