package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultOllamaEndpoint = "http://localhost:11434"

type Client struct {
	baseURL    string
	httpClient *http.Client
	history    []Message
	stateless  bool // New field to control stateless behavior
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Format   any       `json:"format"` // Change to any type to support object formats
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewClient(baseURL string, stateless bool) *Client {
	if baseURL == "" {
		baseURL = defaultOllamaEndpoint
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 300, // 5 minute timeout
		},
		history:   make([]Message, 0),
		stateless: stateless,
	}
}

func NewOllamaClient(stateless bool) *Client {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "localhost" // default value
	}

	port := os.Getenv("OLLAMA_PORT")
	if port == "" {
		port = "11434" // default value
	}

	baseURL := fmt.Sprintf("http://%s:%s", host, port)

	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		history:    make([]Message, 0),
		stateless:  stateless,
	}
}

func (r *ChatRequest) WithFormat(format any) *ChatRequest {
	r.Format = format
	return r
}

func (c *Client) AddMessage(role, content string) {
	if !c.stateless {
		c.history = append(c.history, Message{
			Role:    role,
			Content: content,
		})
		// For debugging
		log.Printf("Added message to history: %+v", c.history)
	}
}

func (c *Client) ClearHistory() {
	c.history = make([]Message, 0)
}

func (c *Client) ChatCompletion(req interface{}) (string, error) {
	// Convert the generic request to ChatRequest
	chatReq, ok := req.(ChatRequest)
	if !ok {
		// If not already ChatRequest, try to convert from map
		if reqMap, ok := req.(map[string]interface{}); ok {
			chatReq = ChatRequest{
				Model:    reqMap["model"].(string),
				Messages: make([]Message, 0),
			}
			// Only set format if it exists in the request
			if format, exists := reqMap["format"].(string); exists {
				chatReq.Format = format
			}
			if msgs, ok := reqMap["messages"].([]interface{}); ok {
				for _, msg := range msgs {
					if msgMap, ok := msg.(map[string]interface{}); ok {
						chatReq.Messages = append(chatReq.Messages, Message{
							Role:    msgMap["role"].(string),
							Content: msgMap["content"].(string),
						})
					}
				}
			}
		} else {
			return "", fmt.Errorf("invalid request type")
		}
	}

	// Combine history with new messages only if not stateless
	if len(chatReq.Messages) > 0 && !c.stateless {
		// Add the new message to history
		c.AddMessage(chatReq.Messages[len(chatReq.Messages)-1].Role,
			chatReq.Messages[len(chatReq.Messages)-1].Content)
	}

	// Use the full history for the request only if not stateless
	if !c.stateless {
		chatReq.Messages = c.history
	}

	// Convert the request to JSON
	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Request: %+v", chatReq.Messages)

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/chat", c.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("server returned status code %d: %s", response.StatusCode, string(body))
	}

	log.Printf("Response: %+v", response)

	var responseString string = ""

	scanner := bufio.NewScanner(response.Body)
	// Set a larger buffer size to handle longer responses
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		var response map[string]interface{}

		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			return "", fmt.Errorf("error parsing response: %w", err)
		}

		if response["message"] != nil {
			responseString += response["message"].(map[string]interface{})["content"].(string)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// After getting a successful response, add it to history only if not stateless
	if responseString != "" && !c.stateless {
		c.AddMessage("assistant", responseString)
	}

	return responseString, nil
}
