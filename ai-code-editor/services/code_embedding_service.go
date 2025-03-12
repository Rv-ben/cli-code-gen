package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"ai-code-editor/config"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	ollama "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
)

// CodeEmbeddingService handles embedding and retrieving code using Chroma and Ollama
type CodeEmbeddingService struct {
	chromaClient      *chromago.Client
	collectionName    string
	collection        *chromago.Collection
	chromaURL         string
	embeddingFunction *ollama.OllamaEmbeddingFunction
	chunkingService   *CodeChunkingService
}

// NewCodeEmbeddingService creates a new instance of the code embedding service
func NewCodeEmbeddingService(cfg *config.Config, collectionName string) (*CodeEmbeddingService, error) {
	chromaURL := cfg.ChromaURL
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}

	// Create Chroma client
	chromaClient, err := chromago.NewClient(chromaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Chroma client: %w", err)
	}

	ef, err := ollama.NewOllamaEmbeddingFunction(ollama.WithBaseURL("http://127.0.0.1:11434"), ollama.WithModel("nomic-embed-text"))
	if err != nil {
		fmt.Printf("Error creating Ollama embedding function: %s \n", err)
	}

	// Create or get collection
	collection, err := chromaClient.NewCollection(
		context.Background(),
		collection.WithName(collectionName),
		collection.WithCreateIfNotExist(true),
		collection.WithEmbeddingFunction(ef),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return &CodeEmbeddingService{
		chromaClient:    chromaClient,
		collectionName:  collectionName,
		collection:      collection,
		chromaURL:       chromaURL,
		chunkingService: NewCodeChunkingService(1000), // Default 1000 chunk size
	}, nil
}

// StoreCode embeds and stores code in the vector database
func (s *CodeEmbeddingService) StoreCode(filePath, code string, metadata map[string]interface{}) error {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	// Add file path to metadata if not present
	if _, exists := metadata["path"]; !exists {
		metadata["path"] = filePath
	}

	// Create a unique ID for this code snippet
	id := fmt.Sprintf("%s-%d", strings.ReplaceAll(filePath, "/", "-"), len(code))

	// Use direct API call instead of the collection.Add method
	endpoint := fmt.Sprintf("%s/api/v1/collections/%s/add", s.chromaURL, s.collectionName)

	// Prepare request body
	reqBody := map[string]interface{}{
		"ids":       []string{id},
		"metadatas": []map[string]interface{}{metadata},
		"documents": []string{code},
	}

	// Get embedding from Ollama
	embedding, err := s.embeddingFunction.EmbedDocuments(context.Background(), []string{code})
	if err != nil {
		return fmt.Errorf("failed to get embedding: %w", err)
	}

	// Add embedding to request
	reqBody["embeddings"] = [][]float32{*embedding[0].GetFloat32()}

	// Send request
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to add to collection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to add to collection: %s", string(body))
	}

	log.Printf("Stored code from %s in vector database", filePath)
	return nil
}

// StoreCodeChunks splits code into chunks and stores them
func (s *CodeEmbeddingService) StoreCodeChunks(filePath, code string, chunkSize int, metadata map[string]interface{}) error {
	// Use the chunking service to split the code and get metadata for each chunk
	chunks, chunkMetadatas := s.chunkingService.ChunkCodeWithMetadata(code, metadata, chunkSize)

	// Store each chunk
	for i, chunk := range chunks {
		err := s.StoreCode(
			fmt.Sprintf("%s-chunk-%d", filePath, i),
			chunk,
			chunkMetadatas[i],
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// QuerySimilarCode finds similar code based on a query
func (s *CodeEmbeddingService) QuerySimilarCode(query string, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 5 // Default limit
	}

	// Query parameters
	queryParams := map[string]interface{}{
		"where": map[string]interface{}{}, // No filtering
	}
	includeParams := map[string]interface{}{
		"documents": true,
		"metadatas": true,
		"distances": true,
	}

	// Query the collection
	results, err := s.collection.Query(
		context.Background(),
		[]string{query}, // Query text
		int32(limit),
		queryParams,
		includeParams,
		nil, // No query enum needed
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	// Format results
	formattedResults := make([]map[string]interface{}, 0)
	for i, doc := range results.Documents {
		if i >= len(results.Metadatas) || i >= len(results.Distances) {
			break
		}

		formattedResults = append(formattedResults, map[string]interface{}{
			"document": doc,
			"metadata": results.Metadatas[i],
			"distance": results.Distances[i],
		})
	}

	return formattedResults, nil
}
