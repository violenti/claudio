package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Claude struct {
	Token string
}

func (c Claude) Name() string {
	return "Claude 4.5 Sonnet"
}

func (c Claude) Question(prompt string) (string, error) {

	var AnthropicKey = os.Getenv("ANTHROPIC_API_KEY")

	client := anthropic.NewClient(option.WithAPIKey(AnthropicKey))

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
		Model: anthropic.ModelClaudeSonnet4_5_20250929,
	})
	if err != nil {
		return "", fmt.Errorf("error of API (Anthropic): %w", err)
	}
	if len(message.Content) > 0 {
		return message.Content[0].Text, nil
	}

	return "I received no response from the AI", nil
}
