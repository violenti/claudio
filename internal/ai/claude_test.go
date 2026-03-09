package ai

import (
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
	t.Setenv("ANTHROPIC_API_KEY", "")

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
	t.Setenv("ANTHROPIC_API_KEY", "fake-env-key")

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

func TestClaude_ModelField_DefaultEmpty(t *testing.T) {
	c := Claude{ApiKey: "fake-key"}
	if c.Model != "" {
		t.Errorf("expected empty Model by default, got %q", c.Model)
	}
}

func TestClaude_ModelField_CustomModel(t *testing.T) {
	c := Claude{ApiKey: "fake-key", Model: "claude-3-opus-20240229"}
	if c.Model != "claude-3-opus-20240229" {
		t.Errorf("expected Model = %q, got %q", "claude-3-opus-20240229", c.Model)
	}
}

func TestClaude_Question_NoApiKey_WithModel(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")

	c := Claude{ApiKey: "", Model: "claude-3-opus-20240229"}
	_, err := c.Question("hello")
	if err == nil {
		t.Fatal("expected error when no API key is configured, got nil")
	}
	if err.Error() != "ANTHROPIC_API_KEY not configured" {
		t.Errorf("error = %q, want %q", err.Error(), "ANTHROPIC_API_KEY not configured")
	}
}
