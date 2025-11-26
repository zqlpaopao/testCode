package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"time"
)

// PluginManager 管理插件的加载和重载
type PluginManager struct {
	plugins    map[string]*plugin.Plugin
	pluginInfo map[string]*PluginInfo
	mutex      sync.RWMutex
	watcher    *FileWatcher
}

// PluginInfo 存储插件信息
type PluginInfo struct {
	Name     string
	FilePath string
	SoPath   string
	ModTime  time.Time
	Functions map[string]interface{}
}

// FunctionInterface 定义插件函数的通用接口
type FunctionInterface interface {
	Execute(args ...interface{}) (interface{}, error)
	GetName() string
	GetDescription() string
}

// NewPluginManager 创建新的插件管理器
func NewPluginManager() *PluginManager {
	pm := &PluginManager{
		plugins:    make(map[string]*plugin.Plugin),
		pluginInfo: make(map[string]*PluginInfo),
	}
	
	// 创建插件目录
	os.MkdirAll("plugins", 0755)
	os.MkdirAll("plugins/src", 0755)
	os.MkdirAll("plugins/bin", 0755)
	
	// 启动文件监控
	pm.watcher = NewFileWatcher("plugins/src", pm.onFileChanged)
	go pm.watcher.Start()
	
	return pm
}

// LoadPlugin 加载插件
func (pm *PluginManager) LoadPlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	srcPath := filepath.Join("plugins/src", name+".go")
	soPath := filepath.Join("plugins/bin", name+".so")
	
	// 检查源文件是否存在
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("插件源文件不存在: %s", srcPath)
	}
	
	// 编译插件
	if err := pm.compilePlugin(srcPath, soPath); err != nil {
		return fmt.Errorf("编译插件失败: %v", err)
	}
	
	// 加载插件
	p, err := plugin.Open(soPath)
	if err != nil {
		return fmt.Errorf("加载插件失败: %v", err)
	}
	
	// 获取文件信息
	fileInfo, _ := os.Stat(srcPath)
	
	// 存储插件信息
	info := &PluginInfo{
		Name:      name,
		FilePath:  srcPath,
		SoPath:    soPath,
		ModTime:   fileInfo.ModTime(),
		Functions: make(map[string]interface{}),
	}
	
	// 尝试加载标准函数
	pm.loadPluginFunctions(p, info)
	
	pm.plugins[name] = p
	pm.pluginInfo[name] = info
	
	fmt.Printf("✅ 插件 '%s' 加载成功\n", name)
	return nil
}

// compilePlugin 编译插件
func (pm *PluginManager) compilePlugin(srcPath, soPath string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, srcPath)
	cmd.Env = os.Environ()
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("编译错误: %s", string(output))
	}
	
	return nil
}

// loadPluginFunctions 加载插件中的函数
func (pm *PluginManager) loadPluginFunctions(p *plugin.Plugin, info *PluginInfo) {
	// 尝试加载常见的函数名
	commonFunctions := []string{"Execute", "Process", "Handle", "Run", "Calculate", "Transform"}
	
	for _, funcName := range commonFunctions {
		if sym, err := p.Lookup(funcName); err == nil {
			info.Functions[funcName] = sym
		}
	}
	
	// 尝试加载 GetFunctions 函数来获取所有可用函数
	if sym, err := p.Lookup("GetFunctions"); err == nil {
		if getFuncs, ok := sym.(func() map[string]interface{}); ok {
			functions := getFuncs()
			for name, fn := range functions {
				info.Functions[name] = fn
			}
		}
	}
}

// CallFunction 调用插件函数
func (pm *PluginManager) CallFunction(pluginName, functionName string, args ...interface{}) (interface{}, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	info, exists := pm.pluginInfo[pluginName]
	if !exists {
		return nil, fmt.Errorf("插件 '%s' 不存在", pluginName)
	}
	
	fn, exists := info.Functions[functionName]
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
	default:
		return nil, fmt.Errorf("不支持的函数类型: %T", fn)
	}
}

// ReloadPlugin 重新加载插件
func (pm *PluginManager) ReloadPlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	// 删除旧插件
	delete(pm.plugins, name)
	delete(pm.pluginInfo, name)
	
	pm.mutex.Unlock()
	err := pm.LoadPlugin(name)
	pm.mutex.Lock()
	
	if err == nil {
		fmt.Printf("🔄 插件 '%s' 重新加载成功\n", name)
	}
	
	return err
}

// onFileChanged 文件变化回调
func (pm *PluginManager) onFileChanged(filePath string) {
	name := strings.TrimSuffix(filepath.Base(filePath), ".go")
	
	// 检查是否是已加载的插件
	pm.mutex.RLock()
	_, exists := pm.pluginInfo[name]
	pm.mutex.RUnlock()
	
	if exists {
		fmt.Printf("📁 检测到文件变化: %s\n", filePath)
		time.Sleep(100 * time.Millisecond) // 等待文件写入完成
		
		if err := pm.ReloadPlugin(name); err != nil {
			fmt.Printf("❌ 重新加载插件失败: %v\n", err)
		}
	}
}

// ListPlugins 列出所有插件
func (pm *PluginManager) ListPlugins() {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	if len(pm.pluginInfo) == 0 {
		fmt.Println("📭 暂无已加载的插件")
		return
	}
	
	fmt.Println("📋 已加载的插件:")
	for name, info := range pm.pluginInfo {
		fmt.Printf("  - %s (函数: %v)\n", name, pm.getFunctionNames(info))
	}
}

// getFunctionNames 获取函数名列表
func (pm *PluginManager) getFunctionNames(info *PluginInfo) []string {
	names := make([]string, 0, len(info.Functions))
	for name := range info.Functions {
		names = append(names, name)
	}
	return names
}

// CreatePluginFromInput 根据用户输入创建插件
func (pm *PluginManager) CreatePluginFromInput(name, code string) error {
	srcPath := filepath.Join("plugins/src", name+".go")
	
	// 创建插件文件
	file, err := os.Create(srcPath)
	if err != nil {
		return fmt.Errorf("创建插件文件失败: %v", err)
	}
	defer file.Close()
	
	// 写入代码
	_, err = file.WriteString(code)
	if err != nil {
		return fmt.Errorf("写入插件代码失败: %v", err)
	}
	
	fmt.Printf("📝 插件文件已创建: %s\n", srcPath)
	
	// 加载插件
	return pm.LoadPlugin(name)
}

func main() {
	fmt.Println("🚀 动态插件系统启动")
	fmt.Println("支持功能:")
	fmt.Println("  - 动态加载 Go 插件")
	fmt.Println("  - 热重载 (文件变化自动重新加载)")
	fmt.Println("  - 用户输入创建新函数")
	fmt.Println("  - 函数调用和管理")
	fmt.Println()
	
	pm := NewPluginManager()
	scanner := bufio.NewScanner(os.Stdin)
	
	// 显示帮助
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
			
		case "load", "l":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: load <插件名>")
				continue
			}
			if err := pm.LoadPlugin(parts[1]); err != nil {
				fmt.Printf("❌ 加载失败: %v\n", err)
			}
			
		case "list", "ls":
			pm.ListPlugins()
			
		case "call", "c":
			if len(parts) < 3 {
				fmt.Println("❌ 用法: call <插件名> <函数名> [参数...]")
				continue
			}
			
			pluginName := parts[1]
			functionName := parts[2]
			args := make([]interface{}, 0)
			
			// 解析参数
			for i := 3; i < len(parts); i++ {
				args = append(args, parts[i])
			}
			
			result, err := pm.CallFunction(pluginName, functionName, args...)
			if err != nil {
				fmt.Printf("❌ 调用失败: %v\n", err)
			} else {
				fmt.Printf("✅ 结果: %v\n", result)
			}
			
		case "create", "new":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: create <插件名>")
				continue
			}
			createInteractivePlugin(pm, parts[1])
			
		case "reload", "r":
			if len(parts) < 2 {
				fmt.Println("❌ 用法: reload <插件名>")
				continue
			}
			if err := pm.ReloadPlugin(parts[1]); err != nil {
				fmt.Printf("❌ 重载失败: %v\n", err)
			}
			
		case "exit", "quit", "q":
			fmt.Println("👋 再见!")
			return
			
		default:
			fmt.Printf("❌ 未知命令: %s (输入 help 查看帮助)\n", command)
		}
	}
}

func showHelp() {
	fmt.Println("📖 可用命令:")
	fmt.Println("  help, h          - 显示帮助")
	fmt.Println("  load, l <name>   - 加载插件")
	fmt.Println("  list, ls         - 列出已加载的插件")
	fmt.Println("  call, c <plugin> <func> [args] - 调用插件函数")
	fmt.Println("  create, new <name> - 交互式创建新插件")
	fmt.Println("  reload, r <name> - 重新加载插件")
	fmt.Println("  exit, quit, q    - 退出程序")
	fmt.Println()
}

func createInteractivePlugin(pm *PluginManager, name string) {
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Printf("📝 创建插件 '%s'\n", name)
	fmt.Println("请选择插件模板:")
	fmt.Println("  1. 简单函数")
	fmt.Println("  2. 字符串处理")
	fmt.Println("  3. 数学计算")
	fmt.Println("  4. 自定义代码")
	
	fmt.Print("选择 (1-4): ")
	if !scanner.Scan() {
		return
	}
	
	choice := strings.TrimSpace(scanner.Text())
	var code string
	
	switch choice {
	case "1":
		code = generateSimpleTemplate(name)
	case "2":
		code = generateStringTemplate(name)
	case "3":
		code = generateMathTemplate(name)
	case "4":
		fmt.Println("请输入完整的 Go 代码 (以 'END' 结束):")
		var lines []string
		for scanner.Scan() {
			line := scanner.Text()
			if line == "END" {
				break
			}
			lines = append(lines, line)
		}
		code = strings.Join(lines, "\n")
	default:
		fmt.Println("❌ 无效选择")
		return
	}
	
	if err := pm.CreatePluginFromInput(name, code); err != nil {
		fmt.Printf("❌ 创建插件失败: %v\n", err)
	}
}

func generateSimpleTemplate(name string) string {
	return fmt.Sprintf(`package main

// %s 插件 - 简单函数模板

func Execute() interface{} {
	return "Hello from %s plugin!"
}

func Process(input interface{}) interface{} {
	return fmt.Sprintf("Processed: %%v", input)
}

func GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Execute": Execute,
		"Process": Process,
	}
}
`, name, name)
}

func generateStringTemplate(name string) string {
	return fmt.Sprintf(`package main

import (
	"strings"
	"fmt"
)

// %s 插件 - 字符串处理模板

func Execute(input string) string {
	return strings.ToUpper(input)
}

func Reverse(input string) string {
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Count(input string) int {
	return len(input)
}

func GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Execute": Execute,
		"Reverse": Reverse,
		"Count":   Count,
	}
}
`, name)
}

func generateMathTemplate(name string) string {
	return fmt.Sprintf(`package main

import "math"

// %s 插件 - 数学计算模板

func Execute(x float64) float64 {
	return x * x
}

func Add(a, b float64) float64 {
	return a + b
}

func Multiply(a, b float64) float64 {
	return a * b
}

func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

func GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Execute":  Execute,
		"Add":      Add,
		"Multiply": Multiply,
		"Sqrt":     Sqrt,
	}
}
`, name)
}