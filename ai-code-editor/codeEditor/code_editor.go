package codeEditor

import (
	codeEditorActions "ai-code-editor/codeEditor/actions"
	codeEditorSchemas "ai-code-editor/codeEditor/promptFunctions/schemas"
	"ai-code-editor/ollama"
	"ai-code-editor/services"
	"fmt"
	"log"
)

type CodeEditor struct {
}

func NewCodeEditor() *CodeEditor {
	return &CodeEditor{}
}

func (c *CodeEditor) EditCodeBase(client *ollama.Client, model string, basePrompt string, userTask string) {

	// Reading code
	c.LearnFromFiles(client, model, basePrompt, userTask)

	c.EditCode(client, model, userTask)
}

func (c *CodeEditor) LearnFromFiles(client *ollama.Client, model string, basePrompt string, userTask string) {
	// Create a parser to parse the response
	parser := NewAiResponseParser()

	// Use the new schema object
	expectedFormat := codeEditorSchemas.NewFileRequestSchema()

	var initialReply string = c.SendMessage(client, model, expectedFormat, basePrompt+"\n\n USER TASK="+userTask+"\n\n  Open at least 1 file that is relevant to the USER TASK. FIND CONTEXT.")

	log.Printf("Initial reply:\n %v", initialReply)

	// Parse the response
	var actions []codeEditorActions.BaseAction = parser.ParseResponse(initialReply)
	// Execute the actions
	// For 3 times
	var fileContents string = ""
	var reply string = ""

	for range 1 {
		if len(actions) == 0 {
			log.Printf("No actions found in response")
			break
		}

		log.Printf("Attempting to execute actions: %v", actions)

		for _, action := range actions {
			fileContents += "\n\n" + c.ExecuteAction(action)
		}

		reply = c.SendMessage(client, model, expectedFormat, fileContents+"\n\nDo you need more context to solve the USER TASK? If you need more files, provide more files to open, return an empty list of actions if you don't need more context. Respond with JSON.")
		actions = parser.ParseResponse(reply)

		log.Printf("Intermediate reply: %s", reply)
	}

}

func (c *CodeEditor) EditCode(client *ollama.Client, model string, userTask string) {
	// Create a parser to parse the response
	parser := NewAiResponseParser()

	var prompt string = `
		Edit code to solve: {{USER TASK}}
		Respond in JSON with edit_file actions only.
		Edit existing code, don't rewrite.

		Example:
		{
			"actions": [{
				"type": "edit_file",
				"path": "file/path",
				"content": "new code",
				"start_line": 1,
				"end_line": 5,
				"action": "replace"
			}]
		}
	`

	prompt = prompt + "\n\n {{USER TASK}}=" + userTask + "\n\n EDIT THE CODE TO SOLVE THE USER TASK. Respond with JSON."

	// Use the new schema object
	expectedFormat := codeEditorSchemas.NewEditRequestSchema()

	var reply string = c.SendMessage(client, model, expectedFormat, prompt)

	log.Printf("Edit Reply: %s", reply)

	var actions []codeEditorActions.BaseAction = parser.ParseResponse(reply)

	// Seperate each action by file
	var fileActions map[string][]codeEditorActions.BaseAction = make(map[string][]codeEditorActions.BaseAction)

	// Group actions by file path
	for _, action := range actions {
		if action.GetType() == "edit_file" {
			editAction := action.(*codeEditorActions.EditFileAction)
			path := editAction.Path
			if _, exists := fileActions[path]; !exists {
				fileActions[path] = make([]codeEditorActions.BaseAction, 0)
			}
			fileActions[path] = append(fileActions[path], action)
		}
	}

	// Execute actions for each file
	for _, actions := range fileActions {
		if len(actions) > 0 {
			c.ExecuteEditFileAction(actions)
		}
	}
}

func (c *CodeEditor) SendMessage(client *ollama.Client, model string, jsonFormat any, prompt string) string {

	req := ollama.ChatRequest{
		Model: model,
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	if jsonFormat != nil {
		req.WithFormat(jsonFormat)
	}

	resp, err := client.ChatCompletion(req)
	if err != nil {
		var errorMessage string = fmt.Errorf("error sending message: %w", err).Error()

		log.Printf("Error: %v\n", errorMessage)

		return ""
	}

	return resp
}

func (c *CodeEditor) ExecuteAction(action codeEditorActions.BaseAction) string {
	log.Printf("Executing action: %v", action.ToString())

	if action.GetType() == "open_file" {
		log.Printf("Executing RequestFileAction")

		fileAction, ok := action.(*codeEditorActions.RequestFileAction)
		if !ok {
			log.Printf("Error: Failed to convert action to RequestFileAction")
			return ""
		}
		fileContextProvider := services.NewFileContextProvider(fileAction.Path)

		log.Printf("File context provider created with path: %s", fileAction.Path)

		fileContents := fileContextProvider.GetFileContents()

		return fileContents
	}

	return ""
}

func (c *CodeEditor) ExecuteEditFileAction(actions []codeEditorActions.BaseAction) {
	if len(actions) == 0 {
		log.Printf("Warning: No actions to execute")
		return
	}

	// Ensure all actions are for the same file
	firstPath := actions[0].(*codeEditorActions.EditFileAction).Path
	for _, action := range actions {
		if action.GetType() != "edit_file" {
			log.Printf("Error: All actions must be edit_file actions")
			return
		}
		editAction := action.(*codeEditorActions.EditFileAction)
		if editAction.Path != firstPath {
			log.Printf("Error: All actions must be for the same file")
			return
		}
	}

	// Convert each action to an EditFileAction
	var editFileActions []codeEditorActions.EditFileAction = make([]codeEditorActions.EditFileAction, len(actions))

	for i, action := range actions {
		editFileActions[i] = *action.(*codeEditorActions.EditFileAction)
	}

	services.EditFile(firstPath, editFileActions)
}
