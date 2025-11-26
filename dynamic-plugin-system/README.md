# Go 动态插件系统

一个强大的 Go 动态插件加载系统，支持热重载和用户输入创建新函数。

## 🚀 功能特性

- **动态加载插件**: 运行时加载 Go 插件 (.so 文件)
- **热重载**: 自动检测文件变化并重新加载插件
- **交互式创建**: 通过用户输入创建新的插件函数
- **多种模板**: 提供简单函数、字符串处理、数学计算等模板
- **函数调用**: 灵活的函数调用机制，支持多种参数类型
- **插件管理**: 列出、加载、重载插件的完整管理功能

## 📁 项目结构

```
dynamic-plugin-system/
├── main.go                    # 主程序
├── file_watcher.go           # 文件监控器
├── README.md                 # 说明文档
└── plugins/
    ├── src/                  # 插件源码目录
    │   ├── calculator.go     # 计算器插件示例
    │   └── textprocessor.go  # 文本处理插件示例
    └── bin/                  # 编译后的插件目录 (.so 文件)
```

## 🛠️ 安装和运行

### 1. 克隆或创建项目

```bash
# 进入项目目录
cd dynamic-plugin-system

# 初始化 Go 模块 (如果需要)
go mod init dynamic-plugin-system
```

### 2. 运行系统

```bash
go run main.go file_watcher.go
```

## 📖 使用指南

### 基本命令

启动系统后，你可以使用以下命令：

```
🔧 请输入命令 (help 查看帮助): help

📖 可用命令:
  help, h          - 显示帮助
  load, l <name>   - 加载插件
  list, ls         - 列出已加载的插件
  call, c <plugin> <func> [args] - 调用插件函数
  create, new <name> - 交互式创建新插件
  reload, r <name> - 重新加载插件
  exit, quit, q    - 退出程序
```

### 示例操作流程

#### 1. 加载现有插件

```bash
# 加载计算器插件
load calculator

# 加载文本处理插件
load textprocessor
```

#### 2. 查看已加载的插件

```bash
list
```

输出示例：
```
📋 已加载的插件:
  - calculator (函数: [Execute Add Subtract Multiply Divide Power Sqrt Factorial ParseAndCalculate])
  - textprocessor (函数: [Execute ToUpper ToLower Reverse WordCount CharCount RemoveSpaces TitleCase IsPalindrome ExtractNumbers RemoveNumbers CountVowels CountConsonants CamelCase SnakeCase KebabCase Truncate GetStatistics])
```

#### 3. 调用插件函数

```bash
# 调用计算器的加法函数
call calculator Add 10 20

# 调用文本处理的大写转换函数
call textprocessor ToUpper "hello world"

# 调用文本处理的反转函数
call textprocessor Reverse "hello"
```

#### 4. 创建新插件

```bash
# 创建新插件
create myPlugin

# 选择模板或输入自定义代码
选择 (1-4): 1  # 选择简单函数模板
```

## 🔧 插件开发

### 插件结构

每个插件都是一个独立的 Go 包，必须包含以下结构：

```go
package main

// 主要函数实现
func Execute() interface{} {
    return "Hello from plugin!"
}

// 其他功能函数
func MyFunction(param string) string {
    return "Processed: " + param
}

// 必需：返回所有可用函数的映射
func GetFunctions() map[string]interface{} {
    return map[string]interface{}{
        "Execute": Execute,
        "MyFunction": MyFunction,
    }
}
```

### 支持的函数类型

插件系统支持以下函数签名：

```go
func() interface{}                    // 无参数
func(interface{}) interface{}         // 单个 interface{} 参数
func(...interface{}) interface{}      // 可变参数
func(string) string                   // 字符串参数和返回值
func(int) int                        // 整数参数和返回值
```

### 示例插件

#### 简单插件示例

```go
package main

import "fmt"

func Execute() interface{} {
    return "Hello from my plugin!"
}

func Greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

func Add(a, b int) int {
    return a + b
}

func GetFunctions() map[string]interface{} {
    return map[string]interface{}{
        "Execute": Execute,
        "Greet":   Greet,
        "Add":     Add,
    }
}
```

## 🔄 热重载功能

系统会自动监控 `plugins/src/` 目录中的 `.go` 文件：

1. **文件修改**: 当你修改插件源码并保存时，系统会自动重新编译和加载插件
2. **实时更新**: 无需重启主程序，插件的更新会立即生效
3. **错误处理**: 如果编译失败，会显示详细的错误信息

### 热重载示例

1. 启动系统并加载插件：
   ```bash
   load calculator
   call calculator Add 5 3
   # 输出: ✅ 结果: 8
   ```

2. 修改 `plugins/src/calculator.go` 中的 Add 函数：
   ```go
   func Add(a, b float64) float64 {
       return a + b + 1  // 添加额外的 1
   }
   ```

3. 保存文件，系统会自动显示：
   ```
   📁 检测到文件变化: plugins/src/calculator.go
   🔄 插件 'calculator' 重新加载成功
   ```

4. 再次调用函数：
   ```bash
   call calculator Add 5 3
   # 输出: ✅ 结果: 9  (现在结果是 9 而不是 8)
   ```

## 🎯 高级功能

### 1. 批量函数调用

你可以连续调用多个函数：

```bash
call textprocessor ToUpper "hello"
call textprocessor Reverse "world"
call calculator Add 10 20
call calculator Multiply 5 6
```

### 2. 复杂参数处理

对于复杂的参数，系统会尝试智能转换：

```bash
# 字符串参数
call textprocessor WordCount "This is a test sentence"

# 数值参数 (自动转换)
call calculator Power 2 3
```

### 3. 插件模板

创建插件时可以选择不同的模板：

- **模板 1**: 简单函数 - 基础的插件结构
- **模板 2**: 字符串处理 - 包含常用字符串操作函数
- **模板 3**: 数学计算 - 包含数学运算函数
- **模板 4**: 自定义代码 - 完全自定义的插件代码

## 🐛 故障排除

### 常见问题

1. **编译错误**
   ```
   ❌ 编译插件失败: # 错误详情
   ```
   - 检查 Go 语法是否正确
   - 确保所有导入的包都可用
   - 验证函数签名是否正确

2. **加载失败**
   ```
   ❌ 加载失败: 插件源文件不存在
   ```
   - 确保插件文件存在于 `plugins/src/` 目录
   - 检查文件名是否正确

3. **函数调用失败**
   ```
   ❌ 调用失败: 函数 'XXX' 在插件 'YYY' 中不存在
   ```
   - 使用 `list` 命令查看可用函数
   - 确保函数名拼写正确
   - 检查函数是否在 `GetFunctions()` 中注册

### 调试技巧

1. **查看插件信息**
   ```bash
   list  # 显示所有已加载的插件和函数
   ```

2. **重新加载插件**
   ```bash
   reload pluginName  # 手动重新加载插件
   ```

3. **检查文件结构**
   ```bash
   # 确保目录结构正确
   ls plugins/src/    # 查看源码文件
   ls plugins/bin/    # 查看编译后的 .so 文件
   ```

## 📝 注意事项

1. **Go 版本**: 确保使用 Go 1.11+ (支持插件功能)
2. **操作系统**: 插件功能在 Linux 和 macOS 上支持最好
3. **并发安全**: 系统使用读写锁确保并发安全
4. **内存管理**: 插件重载时会自动清理旧的插件引用

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

## 📄 许可证

MIT License