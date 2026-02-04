package ai

type Provider interface {
	Question(prompt string) (string, error)
	Name() string
}
