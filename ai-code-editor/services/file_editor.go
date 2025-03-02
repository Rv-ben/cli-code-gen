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

	// Sort actions by line number in descending order to avoid offset issues
	// Process replacements first, then insertions
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
		endLine := action.GetEndLine() - 1     // Convert to 0-based index

		// Validate indices
		if startLine < 0 {
			startLine = 0
		}
		if endLine >= len(lines) {
			endLine = len(lines) - 1
		}
		if startLine > endLine {
			continue // Skip invalid range
		}

		// Create new slice with replaced content
		lines = append(lines[:startLine], append(newLines, lines[endLine+1:]...)...)
	}

	// Execute insert actions after
	for _, action := range insertActions {
		newLines := strings.Split(action.Content, "\n")
		insertPos := action.GetStartLine() - 1 // Convert to 0-based index

		// Validate insertion position
		if insertPos < 0 {
			insertPos = 0
		}
		if insertPos > len(lines) {
			insertPos = len(lines)
		}

		// Insert new lines at position
		if insertPos == len(lines) {
			lines = append(lines, newLines...)
		} else {
			lines = append(lines[:insertPos], append(newLines, lines[insertPos:]...)...)
		}
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
