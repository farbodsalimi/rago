package rag

import (
	"context"
	"time"
)

// RAG is the main RAG class that orchestrates storage, splitting, and embedding
type RAG struct {
	ctx      context.Context
	storage  Storage
	splitter TextSplitter
	embedder Embedder
	config   RAGConfig
}

// RAGConfig holds configuration for the RAG system
type RAGConfig struct {
	Storage       Storage
	Splitter      TextSplitter
	Embedder      Embedder
	BatchSize     int
	MaxRetries    int
	Timeout       time.Duration
	InputFilePath string
}

// Document represents a text document with metadata
type Document struct {
	ID       string
	Content  string
	Metadata map[string]any
}

// Chunk represents a text chunk with its embedding
type Chunk struct {
	DocumentID string
	Content    string
	Embedding  []float32
	Metadata   map[string]any
	Score      float32
}

// Query represents a search query with optional filters
type Query struct {
	Text     string
	TopK     int
	Filters  map[string]any
	Metadata map[string]any
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Chunk Chunk
}
