package openrouter

import (
	"context"
	"errors"
	"fmt"

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

type ChatResult struct {
	Answer             string  `json:"answer"`
	Category           string  `json:"category"`
	ShouldCreateTicket bool    `json:"should_create_ticket"`
	Confidence         float64 `json:"confidence"`
}

func (c *OpenRouterClient) Chat(ctx context.Context, systemPrompt string, userPrompt string) (*ChatResult, error) {
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
	// if content == "" {
	// 	return nil, fmt.Errorf("empty assistant content")
	// }

	// var out ChatResult
	// if err := json.Unmarshal([]byte(content), &out); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
