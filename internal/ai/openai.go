package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

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

	model := openai.ChatModelGPT5_2
	if o.Model != "" {
		model = openai.ChatModel(o.Model)
	}

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(p),
		},
		Seed:  openai.Int(0),
		Model: model,
	})

	var response strings.Builder
	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			response.WriteString(evt.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		return "", stream.Err()
	}

	return response.String(), nil

}
