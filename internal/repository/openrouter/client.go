package openrouter

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenRouterClient struct {
	model  string
	client *openai.Client
}

func New(apiKey, model string) *OpenRouterClient {
	c := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)

	return &OpenRouterClient{
		model:  model,
		client: &c,
	}
}

func (c *OpenRouterClient) Chat(ctx context.Context, systemPrompt string, userPrompt string) (*dto.ChatResult, error) {
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("empty response from OpenRouter")
	}

	content := resp.Choices[0].Message.Content
	fmt.Println(content)
	if content == "" {
		return nil, fmt.Errorf("empty assistant content")
	}

	return &dto.ChatResult{Content: content}, nil
}
