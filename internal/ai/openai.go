package ai

type OpenAI struct {
	APIKey string
}

func (o OpenAI) Name() string {
	return "OpenAI (GPT-4o)"
}

func (o OpenAI) Question(p string) (string, error) {
	// Aquí iría el código de http.Post a api.openai.com
	return "Question of ChatGPT...", nil
}
