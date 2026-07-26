package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Claude struct {
	ApiKey       string
	Model        string
	SystemPrompt string
}

func (c Claude) Name() string {
	return "Claude 4.5 Sonnet"
}

func (c Claude) Question(prompt string) (string, error) {
	ApiKey := c.ApiKey
	if ApiKey == "" {
		ApiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if ApiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not configured")
	}
	client := anthropic.NewClient(option.WithAPIKey(ApiKey))

	model := anthropic.Model(anthropic.ModelClaudeSonnet4_5_20250929)
	if c.Model != "" {
		model = anthropic.Model(c.Model)
	}

	params := anthropic.MessageNewParams{
		Model:     model,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	}
	if c.SystemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: c.SystemPrompt},
		}
	}
	stream := client.Messages.NewStreaming(context.TODO(), params)

	message := anthropic.Message{}
	for stream.Next() {
		event := stream.Current()
		err := message.Accumulate(event)
		if err != nil {
			panic(err)
		}

		// Accumulate text without printing during streaming

	}
	if stream.Err() != nil {
		return "", stream.Err()
	}

	if len(message.Content) > 0 {
		return message.Content[0].Text, nil
	}

	return "", nil

}
