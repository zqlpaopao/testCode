package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"os"
	"strings"
)

// Chat 结构体用于管理聊天会话
type Chat struct {
	llm    llms.LLM
	ctx    context.Context
	reader *bufio.Reader
	// 用于存储对话历史
	history []string
}

// NewChat 创建一个新的聊天会话
func NewChat() (*Chat, error) {
	llm, err := ollama.New(ollama.WithModel("deepseek-r1:1.5b"))
	if err != nil {
		return nil, fmt.Errorf("创建 LLM 失败: %v", err)
	}
	return &Chat{
		llm:     llm,
		ctx:     context.Background(),
		reader:  bufio.NewReader(os.Stdin),
		history: make([]string, 0),
	}, nil
}

// Start 开始交互式聊天
func (c *Chat) Start() error {
	fmt.Println("欢迎使用 LLM 聊天程序！")
	fmt.Println("输入 'exit' 退出")
	fmt.Println("输入 'clear' 清除对话历史")
	fmt.Println("----------------------------------------")
	for {
		// 获取用户输入
		fmt.Print("\nHuman: ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取输入失败: %v", err)
		}
		input = strings.TrimSpace(input)
		// 处理特殊命令
		switch input {
		case "exit":
			fmt.Println("再见！")
			return nil
		case "clear":
			c.history = make([]string, 0)
			fmt.Println("对话历史已清除")
			continue
		}
		// 构建完整的提示词，包含历史记录
		prompt := c.buildPrompt(input)
		// 发送请求并获取响应
		fmt.Print("\nAssistant: ")
		var response strings.Builder
		completion, err := c.llm.Call(c.ctx, prompt,
			llms.WithTemperature(0.8),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fmt.Print(string(chunk))
				response.Write(chunk)
				return nil
			}),
		)
		if err != nil {
			return fmt.Errorf("获取响应失败: %v", err)
		}
		// 打印完整的响应
		fmt.Printf("\n\n完整响应:\n%s\n", completion)
		// 保存对话历史
		c.history = append(c.history,
			fmt.Sprintf("Human: %s", input),
			fmt.Sprintf("Assistant: %s", response.String()),
		)
		fmt.Println("\n----------------------------------------")
	}
}

// buildPrompt 构建包含历史记录的提示词
func (c *Chat) buildPrompt(input string) string {
	var prompt strings.Builder
	// 添加历史记录
	for _, msg := range c.history {
		prompt.WriteString(msg)
		prompt.WriteString("\n")
	}
	// 添加当前输入
	prompt.WriteString(fmt.Sprintf("Human: %s\nAssistant:", input))
	return prompt.String()
}
func main() {
	chat, err := NewChat()
	if err != nil {
		log.Fatal(err)
	}
	if err := chat.Start(); err != nil {
		log.Fatal(err)
	}
}
