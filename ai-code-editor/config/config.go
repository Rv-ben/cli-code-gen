package config

import "os"

type Config struct {
	Port          string
	OllamaBaseURL string
	SmallModel    string
	MediumModel   string
	LargeModel    string
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

	smallModel := os.Getenv("SMALL_MODEL")
	mediumModel := os.Getenv("MEDIUM_MODEL")
	largeModel := os.Getenv("LARGE_MODEL")

	return &Config{
		Port:          port,
		OllamaBaseURL: ollamaURL,
		SmallModel:    smallModel,
		MediumModel:   mediumModel,
		LargeModel:    largeModel,
	}
}
