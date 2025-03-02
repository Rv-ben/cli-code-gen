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

	// Create a parser to parse the response
	parser := NewAiResponseParser()

	var initialReply string = c.SendMessage(client, model, true, initalPrompt)

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

	for i := 0; i < 3; i++ {
		if len(actions) == 0 {
			log.Printf("No actions found in response")
			break
		}

		for _, action := range actions {
			fileContents += "\n\n" + c.ExecuteAction(action)
		}

		log.Printf("File contents: %s", fileContents)

		reply = c.SendMessage(client, model, true, fileContents+"\n\nDo you need more context to solve the USER TASK? If you need more files, provide more files to open, if not don't write anything")
		actions = parser.ParseResponse(reply)

		log.Printf("Intermediate reply: %s", reply)
	}

	var learnedReply string = c.SendMessage(client, model, false, "Tell me what you learned from the files, you don't need to write in json format")

	log.Printf("Learned reply: %s", learnedReply)
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

		wd, err := os.Getwd()
		if err != nil {
			log.Printf("Error getting working directory: %v", err)
			return ""
		}

		fileContents := fileContextProvider.GetFileContents(wd)

		return fileContents
	}

	return ""
}
