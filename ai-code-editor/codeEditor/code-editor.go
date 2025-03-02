package codeEditor

import (
	codeEditor "ai-code-editor/codeEditor/actions"
	"ai-code-editor/ollama"
	"ai-code-editor/services"
	"fmt"
	"log"
	"os"
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

	var initialReply string = c.SendMessage(client, model, true, initialPrompt)

	log.Printf("Initial reply:\n %v", initialReply)

	// Parse the response
	var actions []codeEditor.BaseAction = parser.ParseResponse(initialReply)

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

		reply = c.SendMessage(client, model, true, fileContents+"\n\nDo you need more context to solve the USER TASK? If you need more files, provide more files to open, if not don't write anything")
		actions = parser.ParseResponse(reply)

		log.Printf("Intermediate reply: %s", reply)
	}

}

func (c *CodeEditor) EditCode(client *ollama.Client, model string) {
	// Create a parser to parse the response
	parser := NewAiResponseParser()

	var prompt string = `
		Edit the code to solve the USER TASK. Use the json structure to write the code.
		You should try your best to not make drastic changes to the code.
		You should respond with an array of actions, where each action has a "type" field that is either "open_file" or "write_file". 

		Your response should be in the following format and nothing else (EXAMPLE):
		{
			"actions": [
				{
					"type": "write_file",
					"path": "path/to/file",
					"content": "file contents here"
				}
			]
		}
	`

	var reply string = c.SendMessage(client, model, true, prompt)

	log.Printf("Edit Reply: %s", reply)

	var actions []codeEditor.BaseAction = parser.ParseResponse(reply)

	for _, action := range actions {
		c.ExecuteAction(action)
	}
}

func (c *CodeEditor) SendMessage(client *ollama.Client, model string, isJson bool, prompt string) string {

	req := ollama.ChatRequest{
		Model: model,
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	if isJson {
		req.WithFormat("json")
	}

	resp, err := client.ChatCompletion(req)
	if err != nil {
		var errorMessage string = fmt.Errorf("error sending message: %w", err).Error()

		log.Printf("Error: %v\n", errorMessage)

		return ""
	}

	return resp
}

func (c *CodeEditor) ExecuteAction(action codeEditor.BaseAction) string {
	log.Printf("Executing action: %v", action.ToString())

	if action.GetType() == "open_file" {
		log.Printf("Executing RequestFileAction")

		fileAction, ok := action.(*codeEditor.RequestFileAction)
		if !ok {
			log.Printf("Error: Failed to convert action to RequestFileAction")
			return ""
		}
		fileContextProvider := services.NewFileContextProvider(fileAction.Path)

		log.Printf("File context provider created with path: %s", fileAction.Path)

		fileContents := fileContextProvider.GetFileContents()

		return fileContents
	}

	if action.GetType() == "write_file" {
		log.Printf("Executing WriteFileAction")

		fileAction, ok := action.(*codeEditor.EditFileAction)
		if !ok {
			log.Printf("Error: Failed to convert action to EditFileAction")
			return ""
		}

		log.Printf("Writing file: %s", fileAction.Path)

		os.WriteFile(fileAction.Path, []byte(fileAction.Content), 0644)

		return ""
	}

	return ""
}
