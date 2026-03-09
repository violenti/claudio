package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOllamaName(t *testing.T) {
	tests := []struct {
		name     string
		ollama   Ollama
		expected string
	}{
		{
			name:     "custom model",
			ollama:   Ollama{Model: "mistral"},
			expected: "Ollama (mistral)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ollama.Name(); got != tt.expected {
				t.Errorf("Ollama.Name() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOllamaQuestion(t *testing.T) {
	// Mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/generate" {
			t.Errorf("Expected path /api/generate, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Parse request body
		var req OllamaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		if req.Prompt != "test prompt" {
			t.Errorf("Expected prompt 'test prompt', got %s", req.Prompt)
		}

		if req.Model == "" {
			t.Errorf("Expected a model to be set, got empty string")
		}

		if !req.Stream {
			t.Errorf("Expected stream to be true")
		}

		// Send mock streaming response
		w.Header().Set("Content-Type", "application/json")
		responses := []OllamaResponse{
			{Model: req.Model, Response: "Hello ", Done: false},
			{Model: req.Model, Response: "world!", Done: false},
			{Model: req.Model, Response: "", Done: true},
		}

		encoder := json.NewEncoder(w)
		for _, resp := range responses {
			if err := encoder.Encode(resp); err != nil {
				t.Errorf("Error encoding response: %v", err)
			}
		}
	}))
	defer server.Close()

	ollama := Ollama{
		BaseURL: server.URL,
		Model:   "deepseek-r1:8b",
	}

	response, err := ollama.Question("test prompt")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResponse := "Hello world!"
	if response != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, response)
	}
}

func TestOllamaQuestionWithError(t *testing.T) {
	// Mock Ollama server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	ollama := Ollama{
		BaseURL: server.URL,
		Model:   "deepseek-r1:8b",
	}

	_, err := ollama.Question("test prompt")
	if err == nil {
		t.Error("Expected an error, but got none")
	}

	if !strings.Contains(err.Error(), "Ollama API error") {
		t.Errorf("Expected error to contain 'Ollama API error', got %v", err)
	}
}

func TestOllamaQuestionWithOllamaError(t *testing.T) {
	// Mock Ollama server that returns an Ollama-specific error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := OllamaResponse{
			Error: "model not found",
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ollama := Ollama{
		BaseURL: server.URL,
		Model:   "nonexistent",
	}

	_, err := ollama.Question("test prompt")
	if err == nil {
		t.Error("Expected an error, but got none")
	}

	if !strings.Contains(err.Error(), "model not found") {
		t.Errorf("Expected error to contain 'model not found', got %v", err)
	}
}

func TestOllamaEnvironmentVariables(t *testing.T) {
	// Clear env vars — t.Setenv restores originals automatically on cleanup
	t.Setenv("OLLAMA_BASE_URL", "")
	t.Setenv("OLLAMA_MODEL", "")

	// Phase 1: no env vars, BaseURL provided via struct
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			w.Header().Set("Content-Type", "application/json")
			listResponse := OllamaListResponse{
				Models: []OllamaModel{{Name: "deepseek-r1:8b"}},
			}
			_ = json.NewEncoder(w).Encode(listResponse)
			return
		}

		var req OllamaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Error decoding request: %v", err)
		}

		if req.Model == "" {
			t.Errorf("Expected a default model to be set, got empty string")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(OllamaResponse{Model: req.Model, Response: "test", Done: true})
	}))
	defer server.Close()

	ollama := Ollama{BaseURL: server.URL}
	_, err := ollama.Question("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Phase 2: use env vars
	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(OllamaListResponse{
				Models: []OllamaModel{{Name: "mistral"}},
			})
			return
		}

		var req OllamaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Error decoding request: %v", err)
		}

		if req.Model != "mistral" {
			t.Errorf("Expected model from env 'mistral', got %s", req.Model)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(OllamaResponse{Model: req.Model, Response: "test", Done: true})
	}))
	defer server2.Close()

	t.Setenv("OLLAMA_BASE_URL", server2.URL)
	t.Setenv("OLLAMA_MODEL", "mistral")

	ollamaEnv := Ollama{}
	_, err = ollamaEnv.Question("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestOllamaConnectionError(t *testing.T) {
	ollama := Ollama{
		BaseURL: "http://nonexistent-server:11434",
		Model:   "llama2",
	}

	_, err := ollama.Question("test prompt")
	if err == nil {
		t.Error("Expected a connection error, but got none")
	}

	if !strings.Contains(err.Error(), "error making request to Ollama") {
		t.Errorf("Expected error to contain 'error making request to Ollama', got %v", err)
	}
}
