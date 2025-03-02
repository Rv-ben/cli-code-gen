package config

import "os"

type Config struct {
	Port          string
	OllamaBaseURL string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	return &Config{
		Port:          port,
		OllamaBaseURL: ollamaURL,
	}
} 