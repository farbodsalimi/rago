package pineconeutils

import (
	"context"
	"fmt"

	"github.com/farbodsalimi/rago/pkg/rag"
)

var _ rag.Embedder = PineconeRAGEmbedder{}

type PineconeRAGEmbedderConfig struct{}

type PineconeRAGEmbedder struct {
	client *PineconeClient
	config PineconeRAGEmbedderConfig
}

func NewPineconeRAGEmbedder(
	client *PineconeClient,
	config PineconeRAGEmbedderConfig,
) *PineconeRAGEmbedder {
	return &PineconeRAGEmbedder{
		client: client,
		config: config,
	}
}

func (e PineconeRAGEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	res, err := e.client.Embed([]string{text})
	if err != nil {
		return nil, fmt.Errorf("failed to embed text: %w", err)
	}

	if len(res.Data) == 0 {
		return []float32{}, nil
	}

	return *res.Data[0].Values, nil
}
