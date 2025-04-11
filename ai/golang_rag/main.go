package main

import (
	"bufio"
	"fmt"
	"github.com/hantmac/langchaingo-ollama-rag/rag/logger"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
	"golang.org/x/net/context"
	"net/url"
	"os"
	"strings"
)

var (
	collectionName = "langchaingo-ollama-rag"
	qdrantUrl      = "http://localhost:6333"
	ollamaServer   = "http://localhost:11434"
	topk           = 6
	chunkSize      = 6
	chunkOverlap   = 2
)

// GetOllamaEmbedder 获取ollama嵌入器
func getollamaEmbedder() *embeddings.EmbedderImpl {
	// 创建一个新的ollama模型，模型名为"nomic-embed-text:latest"
	ollamaEmbedderModel, err := ollama.New(
		ollama.WithModel("nomic-embed-text:latest"),
		ollama.WithServerURL(ollamaServer))
	if err != nil {
		logger.Fatal("创建ollama模型失败: %v", err)
	}
	// 使用创建的ollama模型创建一个新的嵌入器
	ollamaEmbedder, err := embeddings.NewEmbedder(ollamaEmbedderModel)
	if err != nil {
		logger.Fatal("创建ollama嵌入器失败: %v", err)
	}
	return ollamaEmbedder
}

func getOllamaDeepseek() *ollama.LLM {
	// 创建一个新的ollama模型，模型名为"deepseek-r1:1.5b"
	llm, err := ollama.New(
		ollama.WithModel("deepseek-r1:1.5b"),
		ollama.WithServerURL(ollamaServer))

	if err != nil {
		logger.Fatal("创建ollama模型失败: %v", err)
	}
	return llm
}

// getStore 获取存储对象
func getStore() *qdrant.Store {
	// 解析URL
	qdUrl, err := url.Parse(qdrantUrl)
	if err != nil {
		logger.Fatal("解析URL失败: %v", err)
	}
	// 创建新的qdrant存储
	store, err := qdrant.New(qdrant.WithURL(*qdUrl), // 设置URL
		qdrant.WithAPIKey(""),                     // 设置API密钥
		qdrant.WithCollectionName(collectionName), // 设置集合名称
		qdrant.WithEmbedder(getollamaEmbedder()),  // 设置嵌入器
	)
	if err != nil {
		logger.Fatal("创建qdrant存储失败: %v", err)
	}
	return &store
}

// storeDocs 将文档存储到向量数据库
func storeDocs(docs []schema.Document, store *qdrant.Store) error {
	// 如果文档数组长度大于0
	if len(docs) > 0 {
		// 添加文档到存储
		fmt.Println("添加文档到存储", docs)
		res, err := store.AddDocuments(context.Background(), docs)
		fmt.Println(res)
		if err != nil {
			return err
		}
	}
	return nil
}

// useRetriaver 函数使用检索器
func useRetriaver(store *qdrant.Store, prompt string, topk int) ([]schema.Document, error) {
	// 设置选项向量
	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80), // 设置分数阈值
	}
	// 创建检索器
	retriever := vectorstores.ToRetriever(store, topk, optionsVector...)
	// 搜索
	docRetrieved, err := retriever.GetRelevantDocuments(context.Background(), prompt)
	if err != nil {
		return nil, fmt.Errorf("检索文档失败: %v", err)
	}
	fmt.Println(docRetrieved)
	// 返回检索到的文档
	return docRetrieved, nil
}

// GetAnswer 获取答案
func GetAnswer(ctx context.Context, llm llms.Model, docRetrieved []schema.Document, prompt string) (string, error) {
	// 创建一个新的聊天消息历史记录
	history := memory.NewChatMessageHistory()
	// 将检索到的文档添加到历史记录中
	for _, doc := range docRetrieved {
		history.AddAIMessage(ctx, doc.PageContent)
	}
	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))

	//fmt.Printf("%#v\n", agents.NewConversationalAgent(llm, nil))
	//fmt.Printf("%#v\n", agents.WithMemory(conversation))
	//os.Exit(1)
	// 使用历史记录创建一个新的对话缓冲区
	executor := agents.NewExecutor(
		agents.NewConversationalAgent(llm, nil),
		agents.WithMemory(conversation))
	// 设置链调用选项
	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}
	// 运行链
	res, err := chains.Run(ctx, executor, prompt, options...)
	if err != nil {
		return "", err
	}
	return res, nil
}

func EmbeddingCmd() {
	// 加载文档
	fmt.Println("📂 Loading documents...")
	filepath := "/Users/zhangqiuli24/Desktop/test/testCode/ai/golang_rag/text.txt"

	docs, err := TextToChunks(filepath, chunkSize, chunkOverlap)
	if err != nil {
		fmt.Println("TextToChunks", err)
	}
	err = storeDocs(docs, getStore())
	if err != nil {
		fmt.Println("storeDocs转换失败")
		panic(err)
	}
	fmt.Println("✅ Documents successfully loaded and stored")
}

func TextToChunks(filepath string, chunkSize, chunkOverlap int) ([]schema.Document, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	loader := documentloaders.NewText(file)
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(chunkSize),
		textsplitter.WithChunkOverlap(chunkOverlap),
	)

	docs, err := loader.LoadAndSplit(context.Background(), splitter)
	if err != nil {
		return nil, fmt.Errorf("failed to process text: %w", err)
	}

	return docs, nil
}

func GetUserInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("input error: %w", err)
	}

	return strings.TrimSpace(input), nil
}

func main() {
	EmbeddingCmd()
	llm := getOllamaDeepseek()

	for {
		prompt, err := GetUserInput("\n💬 Ask a question (or type 'exit' to quit)")
		if err != nil {
			fmt.Printf("Input error: %v", err)
		}

		if strings.ToLower(prompt) == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}
		rst, err := useRetriaver(getStore(), prompt, topk)
		if err != nil {
			panic(err)
		}

		answer, err := GetAnswer(context.Background(), llm, rst, prompt)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(answer)
		}
		if err != nil {
			fmt.Printf("Query failed: %v", err)
		}

		fmt.Printf("\n🤖 Answer:\n%s\n", answer)
	}

}
