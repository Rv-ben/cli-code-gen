package main

import (
	"ai-code-editor/config"
	"ai-code-editor/ollama"
	"ai-code-editor/services"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	config := config.Load()
	ollamaClient := ollama.NewClient(config.OllamaBaseURL)

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

	// Read contents of specified files
	fileContents := ""
	for _, file := range files {
		// Convert file path to be relative to current directory
		relPath, err := filepath.Rel(currentDir, filepath.Join(currentDir, file))
		if err != nil {
			log.Printf("Error getting relative path for %s: %v", file, err)
			continue
		}

		content, err := readFile(relPath)
		if err != nil {
			log.Printf("Error reading file %s: %v", relPath, err)
			continue
		}
		fileContents += fmt.Sprintf("\n<open_file>\n%s\n```%s\n%s\n```\n</open_file>\n",
			relPath, relPath, content)
	}

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

	prompt = prompt + "\n\n" + treeText

	handlePrompt(ollamaClient, model, prompt)
}

// Add new helper function to read files
func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func handlePrompt(client *ollama.Client, model string, prompt string) {
	req := ollama.ChatRequest{
		Model: model,
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	resp, err := client.ChatCompletion(req)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Debug: Print response status
	fmt.Printf("Response status: %s\n", resp.Status)

	printResponse(resp)
}

func printResponse(resp *http.Response) {
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var response map[string]interface {
		}

		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			log.Printf("Error parsing response: %v\n", err)
			continue
		}

		if response["message"] != nil {
			fmt.Print(response["message"].(map[string]interface{})["content"].(string))
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading response: %v\n", err)
	}
}
