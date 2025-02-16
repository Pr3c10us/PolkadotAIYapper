package openai

import (
	"context"
	"github.com/Pr3c10us/boilerplate/internals/domains/llm"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

type Repository struct {
	client *openai.Client
}

func NewOpenAIRepository(client *openai.Client) llm.Repository {
	return &Repository{client: client}
}

func (repo *Repository) Prompt(prompt string) (string, error) {

	chatCompletion, err := repo.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func (repo *Repository) Embed(prompt string) ([]float32, error) {
	response, err := repo.client.Embeddings.New(context.TODO(), openai.EmbeddingNewParams{
		Input:          openai.F[openai.EmbeddingNewParamsInputUnion](shared.UnionString(prompt)),
		Model:          openai.F(openai.EmbeddingModelTextEmbedding3Large),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})
	if err != nil {
		return nil, err
	}
	embedding := response.Data
	return formatEmbedding(embedding[0].Embedding), err
}

func formatEmbedding(embedding []float64) []float32 {
	float32Embedding := make([]float32, len(embedding))
	for i, val := range embedding {
		float32Embedding[i] = float32(val)
	}
	return float32Embedding
}
