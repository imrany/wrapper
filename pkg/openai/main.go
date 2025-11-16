package openaiwrapper

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	APIKey string
	Model  string
}

func (o *OpenAIClient) GenerateOpenAIContent(ctx context.Context, prompt string) (string, error) {
	client := openai.NewClient(o.APIKey)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: o.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if resp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("empty response")
	}

	return resp.Choices[0].Message.Content, nil
}
