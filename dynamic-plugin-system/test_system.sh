#!/bin/bash

# Go 动态插件系统测试脚本

echo "🚀 开始测试 Go 动态插件系统..."

# 进入项目目录
cd "$(dirname "$0")"

# 检查 Go 版本
echo "📋 检查 Go 版本..."
go version

# 初始化 Go 模块
echo "📦 初始化 Go 模块..."
go mod init dynamic-plugin-system 2>/dev/null || echo "模块已存在"

# 创建必要的目录
echo "📁 创建插件目录..."
mkdir -p plugins/src
mkdir -p plugins/bin

# 编译测试 - 先手动编译一个插件来验证系统
echo "🔨 测试插件编译..."
if [ -f "plugins/src/calculator.go" ]; then
    echo "编译 calculator 插件..."
    go build -buildmode=plugin -o plugins/bin/calculator.so plugins/src/calculator.go
    if [ $? -eq 0 ]; then
        echo "✅ calculator 插件编译成功"
    else
        echo "❌ calculator 插件编译失败"
        exit 1
    fi
fi

if [ -f "plugins/src/textprocessor.go" ]; then
    echo "编译 textprocessor 插件..."
    go build -buildmode=plugin -o plugins/bin/textprocessor.so plugins/src/textprocessor.go
    if [ $? -eq 0 ]; then
        echo "✅ textprocessor 插件编译成功"
    else
        echo "❌ textprocessor 插件编译失败"
        exit 1
    fi
fi

# 检查编译结果
echo "📋 检查编译结果..."
ls -la plugins/bin/

echo ""
echo "🎉 系统测试完成！"
echo ""
echo "🚀 启动动态插件系统："
echo "   go run main.go file_watcher.go"
echo ""
echo "📖 基本使用步骤："
echo "   1. 输入 'load calculator' 加载计算器插件"
echo "   2. 输入 'load textprocessor' 加载文本处理插件"
echo "   3. 输入 'list' 查看已加载的插件"
echo "   4. 输入 'call calculator Add 10 20' 测试函数调用"
echo "   5. 输入 'call textprocessor ToUpper hello' 测试文本处理"
echo ""
echo "🔄 热重载测试："
echo "   1. 修改 plugins/src/calculator.go 中的任意函数"
echo "   2. 保存文件，系统会自动重新加载"
echo ""
echo "✨ 创建新插件："
echo "   1. 输入 'create myplugin' 创建新插件"
echo "   2. 选择模板或输入自定义代码"
echo ""