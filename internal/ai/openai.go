package ai

type OpenAI struct {
	Token string
}

func (o OpenAI) Name() string {
	return "OpenAI (GPT-4o)"
}

func (o OpenAI) Question(p string) (string, error) {
	return "Question of ChatGPT...", nil
}
