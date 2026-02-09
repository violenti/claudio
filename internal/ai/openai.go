package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
)

type OpenAI struct {
	Token string
	Model string
}

func (o OpenAI) Name() string {
	return "OpenAI (GPT-4o)"
}

func (o OpenAI) Question(p string) (string, error) {
	Token := o.Token
	if Token == "" {
		Token = os.Getenv("OPENAI_API_KEY")
	}
	if Token == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not configured")
	}

	client := openai.NewClient()
	ctx := context.Background()

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(p),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT5_2,
	}

	completion, err := client.Chat.Completions.New(ctx, params)

	if err != nil {
		return "", fmt.Errorf("error of API (OpenIA): %w", err)
	}
	if len(completion.Choices) > 0 {
		return completion.Choices[0].Message.Content, nil
	}

	return "I received no response from the AI", nil
}
