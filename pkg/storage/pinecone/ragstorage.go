package pineconeutils

import (
	"context"

	"github.com/farbodsalimi/rago/pkg/rag"
)

var _ rag.Storage = PineConeRAGStorage{}

type PineConeRAGStorageConfig struct {
	UserID    string
	FolderID  string
	IndexName string
	Namespace string
}

type PineConeRAGStorage struct {
	client *PineconeClient
	config PineConeRAGStorageConfig
}

func NewPineConeRAGStorage(
	client *PineconeClient,
	config PineConeRAGStorageConfig,
) *PineConeRAGStorage {
	return &PineConeRAGStorage{
		client: client,
		config: config,
	}
}

func (p PineConeRAGStorage) Store(ctx context.Context, chunks []rag.Chunk) error {
	for _, chunk := range chunks {
		err := p.client.Upsert(
			p.config.UserID,
			p.config.FolderID,
			p.config.IndexName,
			p.config.Namespace,
			chunk.Content,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p PineConeRAGStorage) Search(
	ctx context.Context,
	query rag.Query,
	embedding []float32,
) ([]rag.SearchResult, error) {
	results := []rag.SearchResult{}

	values, err := p.client.QueryByVectorValues(
		p.config.UserID,
		p.config.FolderID,
		p.config.IndexName,
		p.config.Namespace,
		query.Text,
		query.TopK,
	)
	if err != nil {
		return results, err
	}

	for _, m := range values.Matches {
		results = append(results, rag.SearchResult{
			Chunk: rag.Chunk{
				Score:      m.Score,
				Content:    m.Vector.Metadata.Fields["content"].GetStringValue(),
				DocumentID: m.Vector.Id,
				Embedding:  *m.Vector.Values,
				Metadata:   m.Vector.Metadata.AsMap(),
			},
		})
	}

	return results, nil
}

func (p PineConeRAGStorage) Delete(ctx context.Context, documentID string) error {
	panic("implement me")
}
