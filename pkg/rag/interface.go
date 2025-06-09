package rag

import (
	"context"
)

// Storage interface for vector storage operations
type Storage interface {
	// Store chunks in the vector database
	Store(ctx context.Context, chunks []Chunk) error

	// Search for similar chunks
	Search(ctx context.Context, query Query, embedding []float32) ([]SearchResult, error)

	// Delete chunks by document ID
	Delete(ctx context.Context, documentID string) error
}

// TextSplitter interface for splitting text into chunks
type TextSplitter interface {
	// Split document into chunks
	Split(ctx context.Context, docPath string) ([]Chunk, error)
}

// Embedder interface for generating embeddings
type Embedder interface {
	// Generate embedding for text
	Embed(ctx context.Context, text string) ([]float32, error)
}
