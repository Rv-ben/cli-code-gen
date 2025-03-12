package main

import (
	"ai-code-editor/config"
	"ai-code-editor/services"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	config := config.Load()

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return
	}

	// Check for required arguments
	if len(os.Args) < 1 {
		fmt.Println("Usage: ollama-cli <prompt> [files...]")
		fmt.Println("Example: ollama-cli 'Fix the bug' file1.go file2.go")
		os.Exit(1)
	}

	userTask := os.Args[1]

	selectedModel := config.LargeModel

	// Initialize the code embedding service
	codeEmbeddingService, err := services.NewCodeEmbeddingService(config, "code_embeddings")
	if err != nil {
		log.Fatalf("Failed to create code embedding service: %v", err)
	}

	// Initialize the semantic file context provider
	semanticContextProvider := services.NewSemanticFileContextProvider(codeEmbeddingService, currentDir)

	// Index the current directory
	fmt.Println("Indexing code files...")
	err = semanticContextProvider.IndexDirectory(currentDir, []string{".go"})
	if err != nil {
		log.Printf("Warning: Error indexing directory: %v", err)
	}

	directoryTree := services.NewDirectoryTree("    ", 10, []string{})

	// Generate tree first to populate known files
	_, err = directoryTree.GenerateTree(currentDir)
	if err != nil {
		log.Printf("Error generating tree: %v", err)
		return
	}

	fmt.Printf("Using model: %s\n User task: %s\n", selectedModel, userTask)

	// Use semantic search to find relevant files
	relevantFiles, err := semanticContextProvider.GetRelevantFiles(userTask, 5)

	// Get relevant context
	relevantContext, err := semanticContextProvider.GetRelevantContext(userTask, 10)

	fmt.Printf("Relevant context: %v\n", relevantContext)
	fmt.Printf("Relevant files: %v\n", relevantFiles)
}

// Helper function to merge file lists without duplicates
func mergeFileLists(list1, list2 []string) []string {
	uniqueFiles := make(map[string]bool)

	// Add all files from both lists to the map
	for _, file := range list1 {
		uniqueFiles[file] = true
	}

	for _, file := range list2 {
		uniqueFiles[file] = true
	}

	// Convert map keys to slice
	result := make([]string, 0, len(uniqueFiles))
	for file := range uniqueFiles {
		result = append(result, file)
	}

	return result
}
