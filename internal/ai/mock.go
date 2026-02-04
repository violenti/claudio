package ai

type MockIA struct{}

func (m MockIA) Name() string { return "IA de Prueba (Gratis)" }
func (m MockIA) Question(p string) (string, error) {
	return "I am 🤖, I don't spend credits. You said: " + p, nil
}
