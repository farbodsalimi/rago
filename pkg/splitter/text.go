package splitter

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"

	"github.com/farbodsalimi/rago/pkg/rag"
)

var _ rag.TextSplitter = TextSplitter{}

type TextSplitterConfig struct {
	ChunkSize    int
	ChunkOverlap int
}

type TextSplitter struct {
	ChunkSize    int
	ChunkOverlap int
}

func NewTextSplitter(config TextSplitterConfig) *TextSplitter {
	return &TextSplitter{
		ChunkSize:    config.ChunkSize,
		ChunkOverlap: config.ChunkOverlap,
	}
}

func (t TextSplitter) Split(ct context.Context, txtFilePath string) ([]rag.Chunk, error) {
	f, err := os.Open(txtFilePath)
	if err != nil {
		return nil, err
	}

	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = t.ChunkSize
	split.ChunkOverlap = t.ChunkOverlap

	tl := documentloaders.NewText(f)
	docs, err := tl.LoadAndSplit(context.Background(), split)
	if err != nil {
		return nil, fmt.Errorf("error loading document: %s", err.Error())
	}

	chunks := []rag.Chunk{}
	for _, v := range docs {
		chunks = append(chunks, rag.Chunk{
			DocumentID: uuid.New().String(),
			Content:    v.PageContent,
			Metadata:   v.Metadata,
			Score:      v.Score,
		})
	}

	return chunks, nil
}
