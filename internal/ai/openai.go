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

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(p),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT5_2,
	})
	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			print(evt.Choices[0].Delta.Content)
		}
	}
	println()

	if err := stream.Err(); err != nil {
		return "", stream.Err()
	}

	return "", nil

}
