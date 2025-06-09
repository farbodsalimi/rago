package pineconeutils

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/pinecone-io/go-pinecone/v3/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
)

type PineconeClient struct {
	ctx            context.Context
	client         *pinecone.Client
	embeddingModel string
}

type PineconeClientConfig struct {
	ApiKey         string
	EmbeddingModel string
}

func NewPineconeClient(
	ctx context.Context,
	config PineconeClientConfig,
) *PineconeClient {
	client, err := pinecone.NewClient(pinecone.NewClientParams{ApiKey: config.ApiKey})
	if err != nil {
		log.Fatalf("Failed to create Pinecone client: %v", err)
	}
	return &PineconeClient{ctx: ctx, client: client, embeddingModel: config.EmbeddingModel}
}

func (c PineconeClient) CreateIndex(indexName string) (*pinecone.Index, error) {
	metric := pinecone.Cosine
	dimension := int32(1024)
	idx, err := c.client.CreateServerlessIndex(c.ctx, &pinecone.CreateServerlessIndexRequest{
		Name:      indexName,
		Cloud:     pinecone.Aws,
		Region:    "us-east-1",
		Metric:    &metric,
		Dimension: &dimension,
		Tags:      &pinecone.IndexTags{"environment": "development"},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating index (may already exist): %v", err)
	} else {
		fmt.Printf("Successfully created serverless index: %s", idx.Name)
	}
	return idx, nil
}

func (c PineconeClient) Upsert(
	userID, folderID, indexName, namespace string,
	document string,
) error {
	idxConnection, err := c.getIndexConnection(indexName, namespace)
	if err != nil {
		return err
	}

	embeddings, err := c.Embed([]string{document})
	if err != nil {
		return err
	}

	metadataMap := map[string]any{
		"userId":   userID,
		"folderId": folderID,
		"content":  document,
	}
	metadata, err := structpb.NewStruct(metadataMap)
	if err != nil {
		log.Fatalf("Failed to create metadata: %v", err)
	}

	vectors := []*pinecone.Vector{}
	for _, embedding := range embeddings.Data {
		vectors = append(vectors, &pinecone.Vector{
			Id:       uuid.New().String(),
			Values:   embedding.Values,
			Metadata: metadata,
		})
	}

	count, err := idxConnection.UpsertVectors(c.ctx, vectors)
	if err != nil {
		return fmt.Errorf("failed to upsert vectors: %s", err.Error())
	} else {
		fmt.Printf("successfully upserted %d vector(s)\n", count)
	}
	return nil
}

func (c PineconeClient) Embed(documents []string) (*pinecone.EmbedResponse, error) {
	queryParameters := pinecone.EmbedParameters{"input_type": "passage", "truncate": "END"}
	res, err := c.client.Inference.Embed(c.ctx, &pinecone.EmbedRequest{
		Model:      c.embeddingModel,
		TextInputs: documents,
		Parameters: queryParameters,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %v", err)
	}
	return res, nil
}

func (c PineconeClient) QueryByVectorValues(
	userID, folderID, indexName, namespace, query string,
	topK int,
) (*pinecone.QueryVectorsResponse, error) {
	embeddings, err := c.Embed([]string{query})
	if err != nil {
		log.Fatalf("failed to add embeddings to pinecone: %v\n\n", err)
	}

	idxConnection, err := c.getIndexConnection(indexName, namespace)
	if err != nil {
		return nil, err
	}

	metadataMap := map[string]any{"folderId": map[string]any{"$eq": folderID}}
	metadataFilter, err := structpb.NewStruct(metadataMap)
	if err != nil {
		log.Fatalf("Failed to create metadataFilter: %v", err)
	}

	res, err := idxConnection.QueryByVectorValues(c.ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:          *embeddings.Data[0].Values,
		TopK:            uint32(topK),
		MetadataFilter:  metadataFilter,
		IncludeValues:   true,
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error encountered when querying by vector: %s", err.Error())
	}

	return res, nil
}

func (c PineconeClient) getIndexConnection(
	indexName string,
	namespace string,
) (*pinecone.IndexConnection, error) {
	idx, err := c.client.DescribeIndex(c.ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to describe index \"%v\": %v", indexName, err)
	}

	idxConnection, err := c.client.Index(pinecone.NewIndexConnParams{
		Host:      idx.Host,
		Namespace: namespace,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create IndexConnection for Host: %v Namespace: %s with err: %w",
			idx.Host,
			namespace,
			err,
		)
	}

	return idxConnection, nil
}
