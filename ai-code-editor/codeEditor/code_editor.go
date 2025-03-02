package codeEditor

import (
	codeEditorActions "ai-code-editor/codeEditor/actions"
	codeEditor "ai-code-editor/codeEditor/schemas"
	codeEditorSchemas "ai-code-editor/codeEditor/schemas"
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

func (c *CodeEditor) EditCodeBase(client *ollama.Client, model string, initalPrompt string) {

	// Reading code
	c.LearnFromFiles(client, model, initalPrompt)

	c.EditCode(client, model)
}

func (c *CodeEditor) LearnFromFiles(client *ollama.Client, model string, initialPrompt string) {
	// Create a parser to parse the response
	parser := NewAiResponseParser()

	// Use the new schema object
	expectedFormat := codeEditorSchemas.NewFileRequestSchema()

	var initialReply string = c.SendMessage(client, model, expectedFormat, initialPrompt)

	log.Printf("Initial reply:\n %v", initialReply)

	// Parse the response
	var actions []codeEditorActions.BaseAction = parser.ParseResponse(initialReply)

	if len(actions) == 0 {
		log.Printf("No actions found in response")
		return
	}

	// Execute the actions
	// For 3 times
	var fileContents string = ""
	var reply string = ""

	for i := 0; i < 0; i++ {
		if len(actions) == 0 {
			log.Printf("No actions found in response")
			break
		}

		for _, action := range actions {
			fileContents += "\n\n" + c.ExecuteAction(action)
		}

		reply = c.SendMessage(client, model, "json", fileContents+"\n\nDo you need more context to solve the USER TASK? If you need more files, provide more files to open, if not don't write anything")
		actions = parser.ParseResponse(reply)

		log.Printf("Intermediate reply: %s", reply)
	}

}

func (c *CodeEditor) EditCode(client *ollama.Client, model string) {
	// Create a parser to parse the response
	parser := NewAiResponseParser()

	var prompt string = `
		Edit the code to solve the USER TASK. Use the json structure to write the code.
		You should respond with an array of actions, where each action has a "type" field that is "write_file" only.
		Respond using JSON.
	`

	// Use the new schema object
	expectedFormat := codeEditor.NewEditRequestSchema()

	var reply string = c.SendMessage(client, model, expectedFormat, prompt)

	log.Printf("Edit Reply: %s", reply)

	var actions []codeEditorActions.BaseAction = parser.ParseResponse(reply)

	// Seperate each action by file
	var fileActions map[string][]codeEditorActions.BaseAction = make(map[string][]codeEditorActions.BaseAction)

	// Group actions by file path
	for _, action := range actions {
		if action.GetType() == "write_file" {
			writeAction := action.(*codeEditorActions.EditFileAction)
			path := writeAction.Path
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

	if jsonFormat != "" {
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
		if action.GetType() != "write_file" {
			log.Printf("Error: All actions must be write_file actions")
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
