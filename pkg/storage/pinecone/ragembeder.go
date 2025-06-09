package pineconeutils

import (
	"context"
	"fmt"

	"github.com/farbodsalimi/rago/pkg/rag"
)

var _ rag.Embedder = PineConeRAGEmbedder{}

type PineConeRAGEmbedderConfig struct{}

type PineConeRAGEmbedder struct {
	client *PineconeClient
	config PineConeRAGEmbedderConfig
}

func NewPineConeRAGEmbedder(
	client *PineconeClient,
	config PineConeRAGEmbedderConfig,
) *PineConeRAGEmbedder {
	return &PineConeRAGEmbedder{
		client: client,
		config: config,
	}
}

func (e PineConeRAGEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	res, err := e.client.Embed([]string{text})
	if err != nil {
		return nil, fmt.Errorf("failed to embed text: %w", err)
	}

	if len(res.Data) == 0 {
		return []float32{}, nil
	}

	return *res.Data[0].Values, nil
}
