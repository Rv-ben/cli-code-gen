package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultOllamaEndpoint = "http://localhost:11434"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
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
	}
}

func (c *Client) ChatCompletion(req interface{}) (string, error) {
	// Convert the generic request to our specific format
	jsonData, err := json.Marshal(req)
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

	return responseString, nil
}
