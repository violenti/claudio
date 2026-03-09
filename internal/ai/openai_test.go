package ai

import (
	"testing"
)

func TestOpenAI_Name(t *testing.T) {
	c := OpenAI{Token: "fake-key"}
	expected := "OpenAI (GPT-4o)"
	if got := c.Name(); got != expected {
		t.Errorf("Name() = %q, want %q", got, expected)
	}
}

func TestOpenAI_ImplementsProviders(t *testing.T) {
	var _ Provider = OpenAI{}
}

func TestOpenAI_Question_NoApiKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")

	c := OpenAI{Token: ""}
	_, err := c.Question("hello")
	if err == nil {
		t.Fatal("expected error when no API key is configured, got nil")
	}
	expected := "OPENAI_API_KEY not configured"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestOpenAI_Question_UsesStructKey(t *testing.T) {

	c := OpenAI{Token: "fake-struct-key"}
	_, err := c.Question("hello")
	if err != nil && err.Error() == "OPENAI_API_KEY not configured" {
		t.Error("should have used struct Token, but got 'not configured' error")
	}
}

func TestOpenAI_ModelField_DefaultEmpty(t *testing.T) {
	o := OpenAI{Token: "fake-key"}
	if o.Model != "" {
		t.Errorf("expected empty Model by default, got %q", o.Model)
	}
}

func TestOpenAI_ModelField_CustomModel(t *testing.T) {
	o := OpenAI{Token: "fake-key", Model: "gpt-5-mini"}
	if o.Model != "gpt-5-mini" {
		t.Errorf("expected Model = %q, got %q", "gpt-5-mini", o.Model)
	}
}

func TestOpenAI_Question_NoApiKey_WithModel(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")

	o := OpenAI{Token: "", Model: "gpt-5-mini"}
	_, err := o.Question("hello")
	if err == nil {
		t.Fatal("expected error when no API key is configured, got nil")
	}
	if err.Error() != "OPENAI_API_KEY not configured" {
		t.Errorf("error = %q, want %q", err.Error(), "OPENAI_API_KEY not configured")
	}
}
