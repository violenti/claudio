package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Models map[string]map[string]string `json:"models"`
}

func ReadModels() (Config, error) {
	config := Config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}

	configPath := filepath.Join(homeDir, ".claudio", "config.json")
	configFile, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	return config, err
}
