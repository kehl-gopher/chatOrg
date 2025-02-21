package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func (app *application) GenerateEmbeddings(text string) ([]float32, error) {
	resp, err := app.openai.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: []string{text},
	})

	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

func EmbeddingToString(embedding []float32) string {
	embedArray := make([]string, len(embedding))
	for i, v := range embedding {
		embedArray[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("[%s]", strings.Join(embedArray, ","))
}
