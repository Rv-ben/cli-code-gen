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
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Format   any       `json:"format,omitempty"` // Change to any type to support object formats
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultOllamaEndpoint
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 300, // 5 minute timeout
		},
		history: make([]Message, 0),
	}
}

func NewOllamaClient() *Client {
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
	}
}

func (r *ChatRequest) WithFormat(format any) *ChatRequest {
	r.Format = format
	return r
}

func (c *Client) AddMessage(role, content string) {
	c.history = append(c.history, Message{
		Role:    role,
		Content: content,
	})
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

	// Combine history with new messages
	if len(chatReq.Messages) > 0 {
		// Add the new message to history
		c.AddMessage(chatReq.Messages[len(chatReq.Messages)-1].Role,
			chatReq.Messages[len(chatReq.Messages)-1].Content)
	}

	// Use the full history for the request
	chatReq.Messages = c.history

	// Convert the request to JSON
	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

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

	// Set cookie for the next request
	cookie := http.Cookie{
		Name:  "session_id",
		Value: response.Header.Get("Set-Cookie"),
	}
	request.Header.Set("Cookie", cookie.String())

	var responseString string = ""

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		var response map[string]interface {
		}

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

	// After getting a successful response, add it to history
	if responseString != "" {
		c.AddMessage("assistant", responseString)
	}

	return responseString, nil
}
