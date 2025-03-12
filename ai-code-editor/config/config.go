package config

import "os"

type Config struct {
	Port          string
	OllamaBaseURL string
	SmallModel    string
	MediumModel   string
	LargeModel    string
	ChromaURL     string
	EmbedModel    string
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

	chromaURL := os.Getenv("CHROMA_URL")
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}

	embedModel := os.Getenv("EMBED_MODEL")
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}

	return &Config{
		Port:          port,
		OllamaBaseURL: ollamaURL,
		SmallModel:    smallModel,
		MediumModel:   mediumModel,
		LargeModel:    largeModel,
		ChromaURL:     chromaURL,
		EmbedModel:    embedModel,
	}
}
