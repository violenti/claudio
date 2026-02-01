package ai

type Claude struct {
	Token string
}

func (c Claude) Name() string {
	return "Claude 3.5 Sonnet"
}

func (c Claude) Question(p string) (string, error) {
	// Aquí iría el código específico para Anthropic
	return "Question of Claude...", nil
}
