package ai

// Provider es el contrato. Si algo tiene estos métodos, es un Provider.
type Provider interface {
	Question(prompt string) (string, error)
	Name() string
}
