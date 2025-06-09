package main

import (
	"context"
	"fmt"
	"os"

	"github.com/farbodsalimi/rago/pkg/rag"
	"github.com/farbodsalimi/rago/pkg/splitter"
	storage "github.com/farbodsalimi/rago/pkg/storage/pinecone"
)

const (
	ChunkSize      = 300 // size of the chunk is number of characters
	ChunkOverlap   = 30  // overlap is the number of characters that the chunks overlap
	EmbeddingModel = "multilingual-e5-large"
	IndexName      = "your-index-name"
)

func main() {
	ctx := context.Background()

	pc := storage.NewPineconeClient(ctx, storage.PineconeClientConfig{
		ApiKey:         os.Getenv("PINECONE_API_KEY"),
		EmbeddingModel: EmbeddingModel,
	})

	// New Pinecone-based RAG
	namespace := "your-namespace"
	r := rag.NewRAG(
		ctx,
		rag.RAGConfig{
			Storage: storage.NewPineconeRAGStorage(pc, storage.PineconeRAGStorageConfig{
				UserID:    "your-user-id",
				FolderID:  "your-folder-id",
				IndexName: IndexName,
				Namespace: namespace,
			}),
			Splitter: splitter.NewTextSplitter(splitter.TextSplitterConfig{
				ChunkSize:    ChunkSize,
				ChunkOverlap: ChunkOverlap,
			}),
			Embedder:      storage.NewPineconeRAGEmbedder(pc, storage.PineconeRAGEmbedderConfig{}),
			InputFilePath: "./data.txt",
		},
	)

	// Process the input file
	r.ProcessDocument()

	// Query the storage
	res, err := r.Search(rag.Query{
		Text: "What's a RAG system?",
		TopK: 5,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
