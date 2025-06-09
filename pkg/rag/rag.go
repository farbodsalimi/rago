package rag

import "context"

// NewRAG creates a new RAG instance with the provided components
func NewRAG(ctx context.Context, config RAGConfig) *RAG {
	return &RAG{
		ctx:      ctx,
		storage:  config.Storage,
		splitter: config.Splitter,
		embedder: config.Embedder,
		config:   config,
	}
}

// ProcessDocument processes and stores documents in the RAG system
func (r *RAG) ProcessDocument() error {
	chunks, err := r.splitter.Split(r.ctx, r.config.InputFilePath)
	if err != nil {
		return err
	}
	return r.storage.Store(r.ctx, chunks)
}

// Search performs similarity search and returns relevant chunks
func (r *RAG) Search(query Query) ([]SearchResult, error) {
	embeddings, err := r.embedder.Embed(r.ctx, query.Text)
	if err != nil {
		return []SearchResult{}, err
	}
	return r.storage.Search(r.ctx, query, embeddings)
}

// DeleteDocument removes a document and its chunks
func (r *RAG) DeleteDocument(documentID string) error {
	return r.storage.Delete(r.ctx, documentID)
}
