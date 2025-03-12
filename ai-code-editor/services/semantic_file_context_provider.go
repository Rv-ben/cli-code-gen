package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SemanticFileContextProvider provides file context based on semantic similarity
type SemanticFileContextProvider struct {
	embeddingService *CodeEmbeddingService
	chunkingService  *CodeChunkingService
	baseDir          string
}

// NewSemanticFileContextProvider creates a new semantic file context provider
func NewSemanticFileContextProvider(embeddingService *CodeEmbeddingService, baseDir string) *SemanticFileContextProvider {
	return &SemanticFileContextProvider{
		embeddingService: embeddingService,
		chunkingService:  NewCodeChunkingService(1500), // Use 1500 as default chunk size
		baseDir:          baseDir,
	}
}

// IndexDirectory indexes all code files in a directory
func (p *SemanticFileContextProvider) IndexDirectory(dir string, extensions []string) error {
	files, err := GetFilesWithExtensions(dir, extensions)
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}

	for _, file := range files {
		// Read file content
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Warning: Failed to read file %s: %v", file, err)
			continue
		}

		// Get relative path for storage
		relPath, err := filepath.Rel(p.baseDir, file)
		if err != nil {
			relPath = file // Use absolute path if relative path fails
		}

		// Store file content using the chunk size from the chunking service
		err = p.embeddingService.StoreCodeChunks(
			relPath,
			string(content),
			p.chunkingService.defaultChunkSize,
			map[string]interface{}{
				"path":      relPath,
				"extension": filepath.Ext(file),
				"type":      "code",
			},
		)
		if err != nil {
			log.Printf("Warning: Failed to store file %s: %v", file, err)
		}
	}

	return nil
}

// GetRelevantFiles returns files relevant to a query
func (p *SemanticFileContextProvider) GetRelevantFiles(query string, limit int) ([]string, error) {
	results, err := p.embeddingService.QuerySimilarCode(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query similar code: %w", err)
	}

	// Extract unique file paths
	uniquePaths := make(map[string]bool)
	for _, result := range results {
		metadata, ok := result["metadata"].(map[string]interface{})
		if !ok {
			continue
		}

		path, ok := metadata["path"].(string)
		if !ok {
			continue
		}

		uniquePaths[path] = true
	}

	// Convert to slice
	paths := make([]string, 0, len(uniquePaths))
	for path := range uniquePaths {
		paths = append(paths, path)
	}

	return paths, nil
}

// GetRelevantContext returns code snippets relevant to a query
func (p *SemanticFileContextProvider) GetRelevantContext(query string, limit int) (map[string]string, error) {
	results, err := p.embeddingService.QuerySimilarCode(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query similar code: %w", err)
	}

	// Group snippets by file
	fileSnippets := make(map[string][]string)
	for _, result := range results {
		metadata, ok := result["metadata"].(map[string]interface{})
		if !ok {
			continue
		}

		path, ok := metadata["path"].(string)
		if !ok {
			continue
		}

		document, ok := result["document"].(string)
		if !ok {
			continue
		}

		fileSnippets[path] = append(fileSnippets[path], document)
	}

	// Combine snippets for each file
	context := make(map[string]string)
	for path, snippets := range fileSnippets {
		context[path] = strings.Join(snippets, "\n\n...\n\n")
	}

	return context, nil
}

// Helper function to get files with specific extensions
func GetFilesWithExtensions(dir string, extensions []string) ([]string, error) {
	var files []string

	// If no extensions provided, use common code extensions
	if len(extensions) == 0 {
		extensions = []string{".go", ".js", ".ts", ".py", ".java", ".c", ".cpp", ".h", ".hpp", ".cs", ".php", ".rb"}
	}

	// Create a map for faster lookup
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has one of the specified extensions
		ext := filepath.Ext(path)
		if extMap[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
