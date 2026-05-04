package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Ollama struct {
	BaseURL string
	Model   string
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

type OllamaModel struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
}

type OllamaListResponse struct {
	Models []OllamaModel `json:"models"`
}

func (o Ollama) Name() string {
	model := o.Model
	if model == "" {
		model = os.Getenv("OLLAMA_MODEL")
		if model == "" {
			baseURL := o.BaseURL
			if baseURL == "" {
				baseURL = os.Getenv("OLLAMA_BASE_URL")
				if baseURL == "" {
					baseURL = "http://localhost:11434"
				}
			}
			availableModel, err := o.getAvailableModel(baseURL)
			if err != nil {
				model = "auto"
			} else {
				model = availableModel
			}
		}
	}
	return fmt.Sprintf("Ollama (%s)", model)
}

func (o Ollama) getAvailableModel(baseURL string) (string, error) {
	url := strings.TrimSuffix(baseURL, "/") + "/api/tags"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error getting available models: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting models (status %d)", resp.StatusCode)
	}

	var listResponse OllamaListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		return "", fmt.Errorf("error decoding models response: %v", err)
	}

	if len(listResponse.Models) == 0 {
		return "", fmt.Errorf("no models available in Ollama")
	}

	// Return the first available model
	return listResponse.Models[0].Name, nil
}

func (o Ollama) Question(prompt string) (string, error) {
	baseURL := o.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("OLLAMA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
	}

	model := o.Model
	if model == "" {
		model = os.Getenv("OLLAMA_MODEL")
		if model == "" {
			// Try to get the first available model
			availableModel, err := o.getAvailableModel(baseURL)
			if err != nil {
				return "", fmt.Errorf("no model specified and couldn't get available models: %v", err)
			}
			model = availableModel
		}
	}

	url := strings.TrimSuffix(baseURL, "/") + "/api/generate"

	requestBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error making request to Ollama: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)

	for {
		var response OllamaResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("error decoding response: %v", err)
		}

		if response.Error != "" {
			return "", fmt.Errorf("Ollama error: %s", response.Error)
		}

		fullResponse.WriteString(response.Response)

		if response.Done {
			break
		}
	}

	return fullResponse.String(), nil
}
