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
	userTask := os.Args[2]

	fmt.Printf("Using model: %s\n", model)

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

	defer treeService.Close() // Make sure to clean up resources
	// Get the tree as a string
	treeText := treeService.GetDirectoryString(currentDir)

	// Build complete prompt with base prompt, working directory, and user task
	completePrompt := basePromptProvider.GetPrompt() + "\n\n" +
		"THE CURRENT WORKING DIRECTORY IS: " + currentDir + "\n\n Here is the directory structure, use it to learn about the codebase and open relevant files to solve the USER TASK. " +
		treeText

	var codeEditingService = codeEditor.NewCodeEditor()

	codeEditingService.EditCodeBase(ollamaClient, model, completePrompt, userTask)
}
