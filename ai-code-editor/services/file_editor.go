package services

import (
	codeEditor "ai-code-editor/codeEditor/actions"
	"bufio"
	"os"
	"strings"
)

func EditFile(path string, actions []codeEditor.EditFileAction) error {
	// Read the file into memory
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// First handle all replace actions
	replaceActions := make([]codeEditor.EditFileAction, 0)
	insertActions := make([]codeEditor.EditFileAction, 0)

	for _, action := range actions {
		if action.GetAction() == "replace" {
			replaceActions = append(replaceActions, action)
		} else {
			insertActions = append(insertActions, action)
		}
	}

	// Execute replace actions first
	for _, action := range replaceActions {
		// Replace lines between start and end with new content
		newLines := strings.Split(action.Content, "\n")
		startLine := action.GetStartLine() - 1 // Convert to 0-based index
		endLine := action.GetEndLine() - 1

		// Remove old lines and insert new ones
		lines = append(lines[:startLine], append(newLines, lines[endLine+1:]...)...)
	}

	// Execute insert actions after
	for _, action := range insertActions {
		// Insert new lines at specified position
		newLines := strings.Split(action.Content, "\n")
		insertPos := action.GetStartLine() - 1

		// Insert new lines at position
		lines = append(lines[:insertPos], append(newLines, lines[insertPos:]...)...)
	}

	// Write back to file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, line := range lines {
		if i > 0 {
			writer.WriteString("\n")
		}
		writer.WriteString(line)
	}

	return writer.Flush()
}
