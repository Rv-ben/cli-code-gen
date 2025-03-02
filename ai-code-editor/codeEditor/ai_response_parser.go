package codeEditor

import (
	codeEditor "ai-code-editor/codeEditor/actions"
	"encoding/json"
	"log"
	"strings"
)

type AiResponseParser struct {
}

type ActionResponse struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
}

func NewAiResponseParser() *AiResponseParser {
	return &AiResponseParser{}
}

func (a *AiResponseParser) ParseResponse(response string) []codeEditor.BaseAction {
	var actionResponse ActionResponse
	actions := []codeEditor.BaseAction{}

	// Find the first '{' and last '}' to extract the JSON object
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start == -1 || end == -1 || start >= end {
		log.Printf("Invalid JSON format in response")
		return actions
	}

	jsonStr := response[start : end+1]
	err := json.Unmarshal([]byte(jsonStr), &actionResponse)
	if err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		return actions
	}

	for _, action := range actionResponse.Actions {
		switch action.Type {
		case "open_file":
			actions = append(actions, codeEditor.NewRequestFileAction(action.Path))
		case "write_file":
			actions = append(actions, codeEditor.NewEditFileAction(action.Path, action.Content))
		default:
			log.Printf("Unknown action type: %s", action.Type)
		}
	}

	return actions
}
