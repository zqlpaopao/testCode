-- MySQL 插件系统数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS plugin_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE plugin_system;

-- 创建插件表
CREATE TABLE IF NOT EXISTS plugins (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE COMMENT '插件名称',
    version VARCHAR(50) NOT NULL DEFAULT '1.0.0' COMMENT '插件版本',
    description TEXT COMMENT '插件描述',
    source_code LONGTEXT NOT NULL COMMENT '插件源码',
    binary_data LONGBLOB COMMENT '编译后的二进制数据',
    hash VARCHAR(64) NOT NULL COMMENT '插件哈希值',
    status ENUM('active', 'inactive', 'error') DEFAULT 'active' COMMENT '插件状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_name (name),
    INDEX idx_status (status),
    INDEX idx_hash (hash),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='插件表';

-- 创建插件变更日志表
CREATE TABLE IF NOT EXISTS plugin_changes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    plugin_name VARCHAR(255) NOT NULL COMMENT '插件名称',
    action ENUM('insert', 'update', 'delete') NOT NULL COMMENT '变更动作',
    old_hash VARCHAR(64) COMMENT '旧哈希值',
    new_hash VARCHAR(64) COMMENT '新哈希值',
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '变更时间',
    INDEX idx_plugin_name (plugin_name),
    INDEX idx_changed_at (changed_at),
    INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='插件变更日志表';

-- 插入示例插件
INSERT INTO plugins (name, version, description, source_code, binary_data, hash, status) VALUES 
(
    'calculator',
    '1.0.0',
    '基础计算器插件，提供加减乘除等数学运算功能',
    'package main

import (
	"fmt"
	"math"
)

// Execute 默认执行函数
func Execute(operation string, a, b float64) interface{} {
	switch operation {
	case "add":
		return Add(a, b)
	case "subtract":
		return Subtract(a, b)
	case "multiply":
		return Multiply(a, b)
	case "divide":
		return Divide(a, b)
	default:
		return fmt.Sprintf("未知操作: %s", operation)
	}
}

// Add 加法
func Add(a, b float64) float64 {
	return a + b
}

// Subtract 减法
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply 乘法
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide 除法
func Divide(a, b float64) interface{} {
	if b == 0 {
		return "错误: 除数不能为零"
	}
	return a / b
}

// Sqrt 平方根
func Sqrt(x float64) interface{} {
	if x < 0 {
		return "错误: 负数没有实数平方根"
	}
	return math.Sqrt(x)
}

// GetFunctions 返回所有可用函数
func GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Execute":  Execute,
		"Add":      Add,
		"Subtract": Subtract,
		"Multiply": Multiply,
		"Divide":   Divide,
		"Sqrt":     Sqrt,
	}
}',
    NULL, -- binary_data 将在首次编译时填充
    'sample_calculator_hash_1234567890abcdef',
    'active'
),
(
    'textprocessor',
    '1.0.0', 
    '文本处理插件，提供字符串操作和文本分析功能',
    'package main

import (
	"strings"
	"unicode"
)

// Execute 默认执行函数
func Execute(text string) interface{} {
	return fmt.Sprintf("处理文本: %s (长度: %d)", text, len(text))
}

// ToUpper 转换为大写
func ToUpper(text string) string {
	return strings.ToUpper(text)
}

// ToLower 转换为小写
func ToLower(text string) string {
	return strings.ToLower(text)
}

// Reverse 反转字符串
func Reverse(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// WordCount 统计单词数量
func WordCount(text string) int {
	words := strings.Fields(text)
	return len(words)
}

// CountVowels 统计元音字母数量
func CountVowels(text string) int {
	vowels := "aeiouAEIOU"
	count := 0
	for _, char := range text {
		if strings.ContainsRune(vowels, char) {
			count++
		}
	}
	return count
}

// GetFunctions 返回所有可用函数
func GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Execute":     Execute,
		"ToUpper":     ToUpper,
		"ToLower":     ToLower,
		"Reverse":     Reverse,
		"WordCount":   WordCount,
		"CountVowels": CountVowels,
	}
}',
    NULL, -- binary_data 将在首次编译时填充
    'sample_textprocessor_hash_abcdef1234567890',
    'active'
);

-- 创建用户和权限（可选）
-- CREATE USER IF NOT EXISTS 'plugin_user'@'localhost' IDENTIFIED BY 'plugin_password';
-- GRANT SELECT, INSERT, UPDATE, DELETE ON plugin_system.* TO 'plugin_user'@'localhost';
-- FLUSH PRIVILEGES;

-- 显示创建结果
SELECT 'Database initialization completed successfully!' as Status;
SELECT COUNT(*) as PluginCount FROM plugins WHERE status = 'active';