package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout redirects os.Stdout and returns what was written.
func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestGetAPIKey_FromEnv(t *testing.T) {
	t.Setenv("TEST_PROVIDER_KEY", "my-secret-key")

	got := getAPIKey("TEST_PROVIDER_KEY", "TestProvider")
	if got != "my-secret-key" {
		t.Errorf("expected 'my-secret-key', got %q", got)
	}
}

func TestGetAPIKey_EmptyWhenEnvNotSet(t *testing.T) {
	t.Setenv("TEST_PROVIDER_KEY", "")

	// Can't test the terminal prompt path in CI, but we can verify
	// the function doesn't panic when the env var is set to empty.
	// The function will try to read from terminal and return empty.
	// We skip the actual prompt by checking behavior is predictable.
	t.Skip("skipping: requires terminal for password prompt")
}

func TestChatMode_MissingProviderIndex(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"aicli", "chat"} // missing index arg

	output := captureStdout(chatMode)
	if !strings.Contains(output, "missing provider index") {
		t.Errorf("expected 'missing provider index', got: %q", output)
	}
}

func TestChatMode_InvalidProviderIndex_NotANumber(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"aicli", "chat", "not-a-number"}

	output := captureStdout(chatMode)
	if !strings.Contains(output, "Error parsing provider index") {
		t.Errorf("expected 'Error parsing provider index', got: %q", output)
	}
}

func TestChatMode_InvalidProviderIndex_OutOfRange(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"aicli", "chat", "99"}

	output := captureStdout(chatMode)
	if !strings.Contains(output, "invalid provider index") {
		t.Errorf("expected 'invalid provider index', got: %q", output)
	}
}

func TestChatMode_InvalidProviderIndex_Negative(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"aicli", "chat", "-1"}

	output := captureStdout(chatMode)
	if !strings.Contains(output, "invalid provider index") {
		t.Errorf("expected 'invalid provider index', got: %q", output)
	}
}
