package redisutils

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisClient(ctx context.Context, client *redis.Client) *RedisClient {
	return &RedisClient{ctx: ctx, client: client}
}

func (c RedisClient) CreateIndex(index string) error {
	// delete any index previously created with the name
	c.client.FTDropIndexWithArgs(c.ctx,
		index,
		&redis.FTDropIndexOptions{
			DeleteDocs: true,
		},
	)

	// create the index
	_, err := c.client.FTCreate(c.ctx,
		index,
		&redis.FTCreateOptions{
			OnHash: true,
			Prefix: []any{"doc:"},
		},
		&redis.FieldSchema{
			FieldName: "content",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "genre",
			FieldType: redis.SearchFieldTypeTag,
		},
		&redis.FieldSchema{
			FieldName: "embedding",
			FieldType: redis.SearchFieldTypeVector,
			VectorArgs: &redis.FTVectorArgs{
				HNSWOptions: &redis.FTHNSWOptions{
					Dim:            1536,
					DistanceMetric: "L2",
					Type:           "FLOAT32",
				},
			},
		},
	).Result()

	if err != nil {
		return err
	}

	return nil
}

func (c RedisClient) AddData(sentences, tags []string, embeddings [][]float64) error {
	for i, embedding := range embeddings {
		buffer, err := floatsToBytes(embedding)
		if err != nil {
			return fmt.Errorf("error converting embeddings to bytes: %w", err)
		}

		_, err = c.client.HSet(c.ctx,
			fmt.Sprintf("doc:%v", i),
			map[string]any{
				"content":   sentences[i],
				"genre":     tags[i],
				"embedding": buffer,
			},
		).Result()

		if err != nil {
			return err
		}
	}

	return nil
}
