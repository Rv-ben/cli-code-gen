package main

import (
	"ai-code-editor/codeEditor"
	"ai-code-editor/config"
	"ai-code-editor/ollama"
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
	ollamaClient := ollama.NewClient(config.OllamaBaseURL)

	basePromptProvider := services.NewBasePromptProvider()

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return
	}

	// Check for required arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: ollama-cli <model> <prompt> [files...]")
		fmt.Println("Example: ollama-cli deepseek-r1:7b 'Fix the bug' file1.go file2.go")
		os.Exit(1)
	}

	model := os.Args[1]
	prompt := os.Args[2]
	files := os.Args[3:] // Get all remaining arguments as files

	fmt.Printf("Using model: %s\n", model)

	fileContextProvider := services.NewFileContextProvider(files)
	fileContents := fileContextProvider.GetFileContents(currentDir)

	// Add file contents to prompt if any files were read
	if fileContents != "" {
		prompt = prompt + "\n" + fileContents
	}

	// Generate directory tree
	treeService := services.NewDirectoryTree(
		"    ", // indent
		10,     // maxDepth
		[]string{ // skipPaths
			"node_modules",
			"vendor",
			".git",
		},
	)

	treeText, err := treeService.GenerateTree(currentDir)
	if err != nil {
		log.Printf("Error generating directory tree: %v", err)
	}

	log.Printf("Working directory: %s", currentDir)

	prompt = prompt + "\n\n" + treeText
	completePrompt := basePromptProvider.GetPrompt() + "\n" + "THE CURRENT WORKING DIRECTORY IS: " + currentDir + "\n" + prompt

	codeEditor := codeEditor.NewCodeEditor()
	codeEditor.EditCodeBase(ollamaClient, model, completePrompt)
}
