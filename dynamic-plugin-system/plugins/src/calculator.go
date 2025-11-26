package main

import (
	"fmt"
	"math"
	"strconv"
)

// Calculator 插件 - 数学计算功能

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
	case "power":
		return Power(a, b)
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

// Power 幂运算
func Power(a, b float64) float64 {
	return math.Pow(a, b)
}

// Sqrt 平方根
func Sqrt(x float64) interface{} {
	if x < 0 {
		return "错误: 负数没有实数平方根"
	}
	return math.Sqrt(x)
}

// Factorial 阶乘
func Factorial(n int) interface{} {
	if n < 0 {
		return "错误: 负数没有阶乘"
	}
	if n == 0 || n == 1 {
		return 1
	}
	
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// ParseAndCalculate 解析字符串并计算
func ParseAndCalculate(expression string) interface{} {
	// 简单的表达式解析器，支持 "数字 操作符 数字" 格式
	parts := []rune(expression)
	var num1Str, operatorStr, num2Str string
	var stage int // 0: num1, 1: operator, 2: num2
	
	for _, char := range parts {
		switch {
		case char >= '0' && char <= '9' || char == '.':
			if stage == 0 {
				num1Str += string(char)
			} else if stage == 2 {
				num2Str += string(char)
			}
		case char == '+' || char == '-' || char == '*' || char == '/':
			if stage == 0 {
				stage = 1
				operatorStr = string(char)
				stage = 2
			}
		case char == ' ':
			// 忽略空格
			continue
		}
	}
	
	num1, err1 := strconv.ParseFloat(num1Str, 64)
	num2, err2 := strconv.ParseFloat(num2Str, 64)
	
	if err1 != nil || err2 != nil {
		return "错误: 无法解析数字"
	}
	
	switch operatorStr {
	case "+":
		return Add(num1, num2)
	case "-":
		return Subtract(num1, num2)
	case "*":
		return Multiply(num1, num2)
	case "/":
		return Divide(num1, num2)
	default:
		return "错误: 不支持的操作符"
	}
}

// GetFunctions 返回所有可用函数
func GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Execute":          Execute,
		"Add":              Add,
		"Subtract":         Subtract,
		"Multiply":         Multiply,
		"Divide":           Divide,
		"Power":            Power,
		"Sqrt":             Sqrt,
		"Factorial":        Factorial,
		"ParseAndCalculate": ParseAndCalculate,
	}
}