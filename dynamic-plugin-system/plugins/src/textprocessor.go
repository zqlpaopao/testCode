package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// TextProcessor 插件 - 文本处理功能

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

// CharCount 统计字符数量
func CharCount(text string) int {
	return len([]rune(text))
}

// RemoveSpaces 移除所有空格
func RemoveSpaces(text string) string {
	return strings.ReplaceAll(text, " ", "")
}

// TitleCase 转换为标题格式
func TitleCase(text string) string {
	return strings.Title(strings.ToLower(text))
}

// IsPalindrome 检查是否为回文
func IsPalindrome(text string) bool {
	// 移除空格并转换为小写
	cleaned := strings.ToLower(strings.ReplaceAll(text, " ", ""))
	return cleaned == Reverse(cleaned)
}

// ExtractNumbers 提取所有数字
func ExtractNumbers(text string) []string {
	re := regexp.MustCompile(`\d+`)
	return re.FindAllString(text, -1)
}

// RemoveNumbers 移除所有数字
func RemoveNumbers(text string) string {
	re := regexp.MustCompile(`\d`)
	return re.ReplaceAllString(text, "")
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

// CountConsonants 统计辅音字母数量
func CountConsonants(text string) int {
	consonantCount := 0
	for _, char := range text {
		if unicode.IsLetter(char) && !strings.ContainsRune("aeiouAEIOU", char) {
			consonantCount++
		}
	}
	return consonantCount
}

// CamelCase 转换为驼峰命名
func CamelCase(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}
	return result
}

// SnakeCase 转换为蛇形命名
func SnakeCase(text string) string {
	words := strings.Fields(text)
	var result []string
	for _, word := range words {
		result = append(result, strings.ToLower(word))
	}
	return strings.Join(result, "_")
}

// KebabCase 转换为短横线命名
func KebabCase(text string) string {
	words := strings.Fields(text)
	var result []string
	for _, word := range words {
		result = append(result, strings.ToLower(word))
	}
	return strings.Join(result, "-")
}

// Truncate 截断文本
func Truncate(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

// ReplaceMultiple 批量替换
func ReplaceMultiple(text string, replacements map[string]string) string {
	result := text
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// GetStatistics 获取文本统计信息
func GetStatistics(text string) map[string]interface{} {
	return map[string]interface{}{
		"length":     len([]rune(text)),
		"words":      WordCount(text),
		"vowels":     CountVowels(text),
		"consonants": CountConsonants(text),
		"lines":      len(strings.Split(text, "\n")),
		"isPalindrome": IsPalindrome(text),
	}
}

// GetFunctions 返回所有可用函数
func GetFunctions() map[string]interface{} {
	return map[string]interface{}{
		"Execute":        Execute,
		"ToUpper":        ToUpper,
		"ToLower":        ToLower,
		"Reverse":        Reverse,
		"WordCount":      WordCount,
		"CharCount":      CharCount,
		"RemoveSpaces":   RemoveSpaces,
		"TitleCase":      TitleCase,
		"IsPalindrome":   IsPalindrome,
		"ExtractNumbers": ExtractNumbers,
		"RemoveNumbers":  RemoveNumbers,
		"CountVowels":    CountVowels,
		"CountConsonants": CountConsonants,
		"CamelCase":      CamelCase,
		"SnakeCase":      SnakeCase,
		"KebabCase":      KebabCase,
		"Truncate":       Truncate,
		"GetStatistics":  GetStatistics,
	}
}