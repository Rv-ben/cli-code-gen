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
)

func main() {
	config := config.Load()
	ollamaClient := ollama.NewClient(config.OllamaBaseURL)

	// Check for required arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: ollama-cli <model> <prompt>")
		fmt.Println("Example: ollama-cli deepseek-r1:7b 'Tell me a joke'")
		os.Exit(1)
	}

	model := os.Args[1]
	prompt := os.Args[2]

	fmt.Printf("Using model: %s\n", model)

	// Example usage of DirectoryTree service
	treeService := services.NewDirectoryTree(
		"    ", // indent
		10,     // maxDepth
		[]string{ // skipPaths
			"node_modules",
			"vendor",
			".git",
		},
	)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return
	}
	treeText, err := treeService.GenerateTree(currentDir)
	if err != nil {
		log.Printf("Error generating directory tree: %v", err)
	}

	prompt = prompt + "\n\n" + treeText

	handlePrompt(ollamaClient, model, prompt)
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
