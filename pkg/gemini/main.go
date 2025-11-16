package geminiwrapper

import (
	"context"

	"google.golang.org/genai"
)

type GeminiClientConfig struct {
	APIKey string
	Model  string
}

func (g *GeminiClientConfig) GenerateGeminiContent(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  g.APIKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", err
	}

	result, err := client.Models.GenerateContent(
		ctx, g.Model, genai.Text(prompt), nil,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
