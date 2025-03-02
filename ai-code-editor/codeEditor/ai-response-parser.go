package codeEditor

import (
	codeEditor "ai-code-editor/codeEditor/actions"
	"log"
	"strings"
)

type AiResponseParser struct {
}

func NewAiResponseParser() *AiResponseParser {
	return &AiResponseParser{}
}

func (a *AiResponseParser) ParseResponse(response string) []codeEditor.BaseAction {
	// The response will have a list of actions in the format <open_file> or <write_file>
	// We need to parse the response and return a list of actions, there will possible be text that are not actions, we need to ignore them

	actions := []codeEditor.BaseAction{}
	// Keep track of current action being built
	currentAction := ""
	inAction := false

	// Process response line by line
	for _, line := range strings.Split(response, "\n") {
		// Check if we're starting a new action
		if strings.Contains(line, "<open_file>") || strings.Contains(line, "<write_file>") {
			inAction = true
			currentAction = line + "\n"
			continue
		}

		// Check if we're ending an action
		if strings.Contains(line, "</open_file>") || strings.Contains(line, "</write_file>") {
			currentAction += line
			actions = append(actions, CreateAction(currentAction))
			currentAction = ""
			inAction = false
			continue
		}

		// If we're in an action, add the line
		if inAction {
			currentAction += line + "\n"
		}
	}

	return actions
}

func CreateAction(actionString string) codeEditor.BaseAction {
	lines := strings.Split(actionString, "\n")
	log.Printf("Parsing action with %d lines", len(lines))

	if len(lines) == 0 {
		log.Printf("No lines found in action string")
		return nil
	}

	firstLine := lines[0]
	log.Printf("First line: %s", firstLine)

	switch {
	case strings.Contains(firstLine, "<open_file>"):
		log.Printf("Parsing open_file action")
		// The path will be between <open_file> and </open_file> 3 lines total
		path := lines[1]
		log.Printf("Found path: %s", path)
		return codeEditor.NewRequestFileAction(path)

	case strings.Contains(firstLine, "<write_file>"):
		log.Printf("Parsing write_file action")
		var path, content string
		inContent := false

		for _, line := range lines {
			if strings.Contains(line, "<path>") {
				path = strings.TrimPrefix(line, "<path>")
				path = strings.TrimSpace(path)
				log.Printf("Found path: %s", path)
			}
			if strings.Contains(line, "<file_contents>") {
				log.Printf("Starting content section")
				inContent = true
				continue
			}
			if strings.Contains(line, "</write_file>") {
				log.Printf("Ending write_file action")
				break
			}
			if inContent {
				content += line + "\n"
			}
		}

		log.Printf("Creating EditFileAction with path: %s and content length: %d", path, len(content))
		return codeEditor.NewEditFileAction(path, content)
	}

	log.Printf("No matching action type found")
	return nil
}
