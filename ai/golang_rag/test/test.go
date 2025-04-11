package main

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

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
)

const (
	collectionName = "langchaingo-ollama-rag"
	qdrantURL      = "http://localhost:6333"
	ollamaServer   = "http://localhost:11434"
	topK           = 5 // 设置检索返回的文档数量
	chunkSize      = 6
	chunkOverlap   = 2
)

type RAGSystem struct {
	store *qdrant.Store
	llm   llms.Model
}

func NewRAGSystem() (*RAGSystem, error) {
	embedder, err := getOllamaEmbedder()
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	store, err := getStore(embedder)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}

	llm, err := getOllamaDeepseek()
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %w", err)
	}

	return &RAGSystem{
		store: store,
		llm:   llm,
	}, nil
}

func getOllamaEmbedder() (*embeddings.EmbedderImpl, error) {
	ollamaEmbedderModel, err := ollama.New(
		ollama.WithModel("nomic-embed-text:latest"),
		ollama.WithServerURL(ollamaServer))
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama model: %w", err)
	}

	ollamaEmbedder, err := embeddings.NewEmbedder(ollamaEmbedderModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	return ollamaEmbedder, nil
}

func getOllamaDeepseek() (*ollama.LLM, error) {
	llm, err := ollama.New(
		ollama.WithModel("deepseek-r1:1.5b"),
		ollama.WithServerURL(ollamaServer))
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama model: %w", err)
	}
	return llm, nil
}

func getStore(embedder *embeddings.EmbedderImpl) (*qdrant.Store, error) {
	qdURL, err := url.Parse(qdrantURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	store, err := qdrant.New(
		qdrant.WithURL(*qdURL),
		qdrant.WithAPIKey(""),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(embedder),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create qdrant store: %w", err)
	}

	return &store, nil
}

func (r *RAGSystem) LoadDocuments(filepath string) error {
	docs, err := TextToChunks(filepath, chunkSize, chunkOverlap)
	if err != nil {
		return fmt.Errorf("failed to process documents: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no documents found in file")
	}

	if _, err := r.store.AddDocuments(context.Background(), docs); err != nil {
		return fmt.Errorf("failed to store documents: %w", err)
	}

	fmt.Println("✅ Documents successfully loaded and stored")
	return nil
}

func (r *RAGSystem) Query(prompt string) (string, error) {
	docs, err := useRetriever(r.store, prompt, topK)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err)
	}

	answer, err := GetAnswer(context.Background(), r.llm, docs, prompt)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	return answer, nil
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

func useRetriever(store *qdrant.Store, prompt string, topK int) ([]schema.Document, error) {
	options := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80),
	}

	retriever := vectorstores.ToRetriever(store, topK, options...)
	docs, err := retriever.GetRelevantDocuments(context.Background(), prompt)
	if err != nil {
		return nil, fmt.Errorf("retrieval error: %w", err)
	}

	return docs, nil
}

func GetAnswer(ctx context.Context, llm llms.Model, docs []schema.Document, prompt string) (string, error) {
	history := memory.NewChatMessageHistory()
	for _, doc := range docs {
		history.AddAIMessage(ctx, doc.PageContent)
	}

	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))
	fmt.Println(conversation)
	executor := agents.NewExecutor(
		agents.NewConversationalAgent(llm, nil),
		agents.WithMemory(conversation))

	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}

	res, err := chains.Run(ctx, executor, prompt, options...)
	if err != nil {
		return "", fmt.Errorf("generation error: %w", err)
	}

	return res, nil
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
	rag, err := NewRAGSystem()
	if err != nil {
		fmt.Printf("Failed to initialize RAG system: %v", err)
	}

	// 加载文档
	fmt.Println("📂 Loading documents...")
	if err := rag.LoadDocuments("/Users/zhangqiuli24/Desktop/test/testCode/ai/golang_rag/text.txt"); err != nil {
		fmt.Println(err)
		os.Exit(1)
		fmt.Printf("Document loading failed: %v", err)
	}

	// 交互式问答
	for {
		prompt, err := GetUserInput("\n💬 Ask a question (or type 'exit' to quit)")
		if err != nil {
			fmt.Printf("Input error: %v", err)
		}

		if strings.ToLower(prompt) == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		answer, err := rag.Query(prompt)
		if err != nil {
			fmt.Printf("Query failed: %v", err)
		}

		fmt.Printf("\n🤖 Answer:\n%s\n", answer)
	}
}
