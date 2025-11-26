# MySQL 动态插件系统

一个基于 MySQL 数据库的 Go 动态插件系统，支持插件存储在数据库中，自动加载到内存缓存，并实现实时更新和热重载功能。

## 🚀 核心特性

### 🎯 解决的问题
- **插件持久化存储**: 插件 (.so 文件) 存储在 MySQL 数据库中，而非文件系统
- **内存缓存机制**: 程序启动时将插件加载到内存，提供高性能访问
- **自动同步更新**: 定期检查数据库变化，自动更新内存中的插件
- **热重载支持**: 插件修改后无需重启程序，自动重新加载最新版本
- **企业级管理**: 完整的插件生命周期管理，支持版本控制和变更追踪

### ✨ 主要功能
- **数据库存储**: 插件源码和二进制文件存储在 MySQL 中
- **内存缓存**: 高性能的内存插件缓存系统
- **实时监控**: 自动检测数据库中的插件变化
- **版本管理**: 支持插件版本控制和历史记录
- **变更追踪**: 详细的插件变更日志
- **导入导出**: 支持插件的导入和导出功能
- **交互管理**: 丰富的命令行交互界面

## 📁 项目结构

```
mysql-plugin-system/
├── main.go              # 主程序 (387行)
├── database.go          # 数据库管理器 (286行)
├── memory_cache.go      # 内存缓存管理器 (329行)
├── go.mod              # Go 模块配置
├── config.json         # 系统配置文件
├── init_database.sql   # 数据库初始化脚本
└── README.md           # 项目说明文档
```

## 🛠️ 安装和配置

### 1. 环境要求

- **Go 1.21+**: 支持插件功能
- **MySQL 5.7+**: 数据库存储
- **Linux/macOS**: 插件功能最佳支持

### 2. 数据库准备

```bash
# 连接到 MySQL
mysql -u root -p

# 执行初始化脚本
source init_database.sql

# 或者手动执行
mysql -u root -p < init_database.sql
```

### 3. 配置文件

编辑 `config.json` 文件：

```json
{
  "database_dsn": "root:password@tcp(localhost:3306)/plugin_system?charset=utf8mb4&parseTime=True&loc=Local",
  "check_interval_seconds": 30,
  "max_cache_size": 100,
  "enable_auto_sync": true
}
```

### 4. 安装依赖

```bash
cd mysql-plugin-system
go mod tidy
```

### 5. 启动系统

```bash
go run main.go database.go memory_cache.go
```

## 📖 使用指南

### 系统启动

启动后会看到类似输出：

```
🚀 MySQL 插件系统启动中...
📦 开始加载 2 个插件到内存...
✅ 插件 'calculator' 已加载到内存
✅ 插件 'textprocessor' 已加载到内存
🎉 内存缓存初始化完成，共加载 2 个插件
✅ 插件系统启动完成

============================================================
🎯 MySQL 动态插件系统
   - 插件存储在 MySQL 数据库中
   - 自动加载插件到内存缓存
   - 支持实时更新和热重载
   - 定期检查数据库变化
============================================================

📊 数据库统计:
   总插件数: 2
   active: 2

💾 内存缓存统计:
   已缓存插件数: 2
   最后检查时间: 2024-01-01 10:00:00
```

### 基本命令

```bash
📖 可用命令:
  help, h              - 显示帮助
  list, ls             - 列出所有插件
  call, c <plugin> <func> [args] - 调用插件函数
  add, create <name>   - 添加新插件
  update, u <name>     - 更新插件
  delete, del <name>   - 删除插件
  refresh, r [name]    - 刷新插件 (不指定名称则刷新全部)
  stats, stat          - 显示系统统计信息
  check                - 手动检查数据库更新
  info, i <name>       - 显示插件详细信息
  export <name>        - 导出插件
  import <file>        - 导入插件
  exit, quit, q        - 退出程序
```

### 使用示例

#### 1. 查看已加载的插件

```bash
🔧 请输入命令: list

📋 已加载的插件 (共 2 个):
  🔹 calculator (v1.0.0) - 函数: [Execute Add Subtract Multiply Divide Sqrt]
     加载时间: 2024-01-01 10:00:00, 哈希: sample_calcul...
  🔹 textprocessor (v1.0.0) - 函数: [Execute ToUpper ToLower Reverse WordCount CountVowels]
     加载时间: 2024-01-01 10:00:00, 哈希: sample_textpr...
```

#### 2. 调用插件函数

```bash
# 调用计算器插件
🔧 请输入命令: call calculator Add 10 20
✅ 结果: 30 (耗时: 123.456µs)

# 调用文本处理插件
🔧 请输入命令: call textprocessor ToUpper hello
✅ 结果: HELLO (耗时: 45.678µs)

# 调用复杂函数
🔧 请输入命令: call calculator Sqrt 16
✅ 结果: 4 (耗时: 67.890µs)
```

#### 3. 添加新插件

```bash
🔧 请输入命令: add greeting
📝 添加插件 'greeting'
版本号 (默认 1.0.0): 1.0.0
描述: 问候插件，提供多语言问候功能
请输入插件源码 (以 'END' 结束):
package main

import "fmt"

func Execute(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

func Greet(name, language string) string {
    switch language {
    case "zh":
        return fmt.Sprintf("你好, %s!", name)
    case "es":
        return fmt.Sprintf("¡Hola, %s!", name)
    case "fr":
        return fmt.Sprintf("Bonjour, %s!", name)
    default:
        return fmt.Sprintf("Hello, %s!", name)
    }
}

func GetFunctions() map[string]interface{} {
    return map[string]interface{}{
        "Execute": Execute,
        "Greet":   Greet,
    }
}
END

✅ 插件 'greeting' 已添加到数据库和内存缓存
✅ 插件 'greeting' 添加成功
```

#### 4. 查看插件详细信息

```bash
🔧 请输入命令: info calculator

📋 插件 'calculator' 详细信息:
  版本: 1.0.0
  哈希: sample_calculator_hash_1234567890abcdef
  加载时间: 2024-01-01 10:00:00
  函数数量: 6
  可用函数: [Execute Add Subtract Multiply Divide Sqrt]
  描述: 基础计算器插件，提供加减乘除等数学运算功能
  创建时间: 2024-01-01 09:00:00
  更新时间: 2024-01-01 09:00:00
```

#### 5. 导出和导入插件

```bash
# 导出插件
🔧 请输入命令: export calculator
✅ 插件已导出到: calculator_v1.0.0.json

# 导入插件
🔧 请输入命令: import calculator_v1.0.0.json
✅ 插件 'calculator' 导入成功
```

## 🔄 热重载演示

### 场景：更新插件功能

1. **当前状态**：
   ```bash
   🔧 请输入命令: call calculator Add 5 3
   ✅ 结果: 8 (耗时: 123µs)
   ```

2. **数据库中更新插件**（通过其他程序或直接 SQL）：
   ```sql
   UPDATE plugins 
   SET source_code = '...(修改后的代码)...', 
       hash = 'new_hash_value',
       version = '1.1.0'
   WHERE name = 'calculator';
   ```

3. **系统自动检测并更新**：
   ```
   🔍 检测到 1 个插件变更
   🔄 插件 'calculator' 已更新到内存缓存 (新哈希: new_hash_value)
   ```

4. **立即使用新功能**：
   ```bash
   🔧 请输入命令: call calculator Add 5 3
   ✅ 结果: 9 (耗时: 134µs)  # 新版本的结果
   ```

## 🏗️ 系统架构

### 数据流架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   MySQL 数据库   │    │   内存缓存管理器   │    │   插件函数调用   │
│                │    │                 │    │                │
│ ┌─────────────┐ │    │ ┌─────────────┐  │    │ ┌─────────────┐ │
│ │ plugins 表  │ │◄──►│ │ CachedPlugin│  │◄──►│ │ 函数执行器   │ │
│ └─────────────┘ │    │ └─────────────┘  │    │ └─────────────┘ │
│ ┌─────────────┐ │    │ ┌─────────────┐  │    │                │
│ │plugin_changes│ │    │ │ 定期检查器   │  │    │                │
│ └─────────────┘ │    │ └─────────────┘  │    │                │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 核心组件

#### 1. DatabaseManager (数据库管理器)
- **职责**: 管理插件在 MySQL 中的存储和检索
- **功能**: 
  - 插件的 CRUD 操作
  - 变更日志记录
  - 统计信息查询
  - 哈希值比较

#### 2. MemoryCache (内存缓存管理器)
- **职责**: 管理插件在内存中的缓存和执行
- **功能**:
  - 插件加载到内存
  - 函数调用执行
  - 定期更新检查
  - 缓存生命周期管理

#### 3. PluginSystem (插件系统)
- **职责**: 协调数据库和内存缓存，提供统一接口
- **功能**:
  - 系统初始化
  - 交互式命令处理
  - 优雅关闭
  - 配置管理

## 🗄️ 数据库设计

### plugins 表结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键，自增 |
| name | VARCHAR(255) | 插件名称，唯一 |
| version | VARCHAR(50) | 插件版本 |
| description | TEXT | 插件描述 |
| source_code | LONGTEXT | 插件源码 |
| binary_data | LONGBLOB | 编译后的 .so 文件 |
| hash | VARCHAR(64) | 插件内容哈希值 |
| status | ENUM | 插件状态 (active/inactive/error) |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### plugin_changes 表结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键，自增 |
| plugin_name | VARCHAR(255) | 插件名称 |
| action | ENUM | 变更动作 (insert/update/delete) |
| old_hash | VARCHAR(64) | 旧哈希值 |
| new_hash | VARCHAR(64) | 新哈希值 |
| changed_at | TIMESTAMP | 变更时间 |

## 🔧 插件开发

### 插件结构要求

每个插件必须是一个有效的 Go 包，包含以下结构：

```go
package main

// 必需：主要功能函数
func Execute(params ...interface{}) interface{} {
    // 插件的主要逻辑
    return result
}

// 可选：其他功能函数
func CustomFunction(param string) string {
    // 自定义功能
    return processedResult
}

// 必需：函数注册
func GetFunctions() map[string]interface{} {
    return map[string]interface{}{
        "Execute": Execute,
        "CustomFunction": CustomFunction,
    }
}
```

### 支持的函数签名

系统支持以下函数类型：

```go
func() interface{}                    // 无参数
func(interface{}) interface{}         // 单个通用参数
func(...interface{}) interface{}      // 可变参数
func(string) string                   // 字符串处理
func(int) int                        // 整数处理
func(float64) float64                // 浮点数处理
```

### 插件示例

#### 数学运算插件

```go
package main

import "math"

func Execute(operation string, a, b float64) interface{} {
    switch operation {
    case "add": return a + b
    case "multiply": return a * b
    case "power": return math.Pow(a, b)
    default: return "Unknown operation"
    }
}

func Add(a, b float64) float64 { return a + b }
func Multiply(a, b float64) float64 { return a * b }
func Power(a, b float64) float64 { return math.Pow(a, b) }

func GetFunctions() map[string]interface{} {
    return map[string]interface{}{
        "Execute": Execute,
        "Add": Add,
        "Multiply": Multiply,
        "Power": Power,
    }
}
```

## 📊 性能特性

### 内存缓存优势

- **快速访问**: 插件加载到内存后，函数调用延迟通常在微秒级别
- **并发安全**: 使用读写锁保护，支持高并发访问
- **智能更新**: 只有当插件哈希值变化时才重新加载
- **资源管理**: 自动清理临时文件和过期缓存

### 数据库优化

- **索引优化**: 在关键字段上建立索引，提高查询性能
- **批量操作**: 支持批量插件操作，减少数据库连接开销
- **变更追踪**: 高效的变更检测机制，避免不必要的更新

## 🔒 安全考虑

### 插件安全

- **代码审查**: 建议对插件源码进行安全审查
- **编译验证**: 只有编译成功的插件才会被加载
- **运行时隔离**: 插件在独立的 goroutine 中执行
- **错误处理**: 完善的错误捕获和恢复机制

### 数据库安全

- **权限控制**: 建议使用专门的数据库用户和权限
- **连接加密**: 支持 SSL/TLS 数据库连接
- **SQL 注入防护**: 使用参数化查询防止 SQL 注入

## 🐛 故障排除

### 常见问题

#### 1. 数据库连接失败
```
❌ 创建数据库管理器失败: dial tcp: connection refused
```
**解决方案**:
- 检查 MySQL 服务是否启动
- 验证 `config.json` 中的连接字符串
- 确认网络连接和防火墙设置

#### 2. 插件编译失败
```
❌ 编译插件失败: package main: syntax error
```
**解决方案**:
- 检查插件源码语法
- 确保包名为 `main`
- 验证所有导入的包是否可用

#### 3. 插件加载失败
```
❌ 加载插件失败: plugin.Open: plugin was built with a different version of package
```
**解决方案**:
- 确保插件和主程序使用相同的 Go 版本编译
- 重新编译插件
- 检查依赖包版本一致性

#### 4. 内存缓存问题
```
❌ 创建临时目录失败: permission denied
```
**解决方案**:
- 检查程序运行用户的权限
- 确保临时目录可写
- 检查磁盘空间是否充足

### 调试技巧

#### 1. 启用详细日志
```bash
# 设置环境变量启用调试模式
export DEBUG=true
go run main.go database.go memory_cache.go
```

#### 2. 检查系统状态
```bash
🔧 请输入命令: stats

📊 系统统计信息:
💾 数据库统计:
  total_plugins: 3
  by_status: map[active:3]
  recent_changes_24h: 5

🧠 内存缓存统计:
  total_cached_plugins: 3
  temp_directory: /tmp/plugin_cache_123456
  last_check: 2024-01-01T10:30:00Z
```

#### 3. 手动检查更新
```bash
🔧 请输入命令: check
🔍 检查数据库更新...
✅ 更新检查完成
```

## 🚀 部署建议

### 生产环境配置

#### 1. 数据库配置
```json
{
  "database_dsn": "plugin_user:secure_password@tcp(db.example.com:3306)/plugin_system?charset=utf8mb4&parseTime=True&loc=Local&tls=true",
  "check_interval_seconds": 60,
  "max_cache_size": 500,
  "enable_auto_sync": true
}
```

#### 2. 系统服务配置
```ini
# /etc/systemd/system/plugin-system.service
[Unit]
Description=MySQL Plugin System
After=network.target mysql.service

[Service]
Type=simple
User=plugin
WorkingDirectory=/opt/plugin-system
ExecStart=/opt/plugin-system/plugin-system
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

#### 3. 监控和日志
- 使用 systemd journal 收集日志
- 配置 Prometheus 监控指标
- 设置告警规则监控系统健康状态

### 高可用部署

- **数据库**: 使用 MySQL 主从复制或集群
- **应用**: 部署多个实例，通过负载均衡器分发请求
- **缓存**: 考虑使用 Redis 作为分布式缓存层

## 🤝 贡献指南

### 开发环境设置

1. **克隆项目**:
   ```bash
   git clone <repository-url>
   cd mysql-plugin-system
   ```

2. **安装依赖**:
   ```bash
   go mod download
   ```

3. **设置测试数据库**:
   ```bash
   mysql -u root -p < init_database.sql
   ```

4. **运行测试**:
   ```bash
   go test ./...
   ```

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加适当的注释和文档
- 编写单元测试

## 📄 许可证

MIT License - 详见 LICENSE 文件

## 📞 支持

如有问题或建议，请：
1. 查看本文档的故障排除部分
2. 提交 GitHub Issue
3. 联系维护团队

---

**注意**: 这是一个企业级的插件系统，适用于需要动态加载和管理插件的复杂应用场景。在生产环境中使用前，请确保进行充分的测试和安全评估。