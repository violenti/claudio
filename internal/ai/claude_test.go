package ai

import (
	"os"
	"testing"
)

func TestClaude_Name(t *testing.T) {
	c := Claude{ApiKey: "fake-key"}
	expected := "Claude 4.5 Sonnet"
	if got := c.Name(); got != expected {
		t.Errorf("Name() = %q, want %q", got, expected)
	}
}

func TestClaude_ImplementsProvider(t *testing.T) {
	var _ Provider = Claude{}
}

func TestClaude_Question_NoApiKey(t *testing.T) {
	original := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Setenv("ANTHROPIC_API_KEY", original)

	c := Claude{ApiKey: ""}
	_, err := c.Question("hello")
	if err == nil {
		t.Fatal("expected error when no API key is configured, got nil")
	}
	expected := "ANTHROPIC_API_KEY not configured"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestClaude_Question_UsesEnvKey(t *testing.T) {
	original := os.Getenv("ANTHROPIC_API_KEY")
	os.Setenv("ANTHROPIC_API_KEY", "fake-env-key")
	defer os.Setenv("ANTHROPIC_API_KEY", original)

	c := Claude{ApiKey: ""}
	_, err := c.Question("hello")
	if err != nil && err.Error() == "ANTHROPIC_API_KEY not configured" {
		t.Error("should have used env var, but got 'not configured' error")
	}
}

func TestClaude_Question_UsesStructKey(t *testing.T) {

	c := Claude{ApiKey: "fake-struct-key"}
	_, err := c.Question("hello")
	if err != nil && err.Error() == "ANTHROPIC_API_KEY not configured" {
		t.Error("should have used struct ApiKey, but got 'not configured' error")
	}
}
