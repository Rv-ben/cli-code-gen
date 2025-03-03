package services

import (
	codeEditor "ai-code-editor/codeEditor/actions"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary test file
func createTempFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}

// Helper function to read file content
func readFile(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	return string(content)
}

func TestEditFile_Replace(t *testing.T) {
	initialContent := "line1\nline2\nline3\nline4\nline5"
	tmpFile := createTempFile(t, initialContent)

	actions := []codeEditor.EditFileAction{
		{
			Action:    "replace",
			StartLine: 2,
			EndLine:   4,
			Content:   "new line2\nnew line3",
		},
	}

	err := EditFile(tmpFile, actions)
	if err != nil {
		t.Errorf("EditFile failed: %v", err)
	}

	expected := "line1\nnew line2\nnew line3\nline5"
	result := readFile(t, tmpFile)

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestEditFile_Insert(t *testing.T) {
	initialContent := "line1\nline2\nline3"
	tmpFile := createTempFile(t, initialContent)

	actions := []codeEditor.EditFileAction{
		{
			Action:    "insert",
			StartLine: 2,
			Content:   "new line",
		},
	}

	err := EditFile(tmpFile, actions)
	if err != nil {
		t.Errorf("EditFile failed: %v", err)
	}

	expected := "line1\nnew line\nline2\nline3"
	result := readFile(t, tmpFile)

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestEditFile_MultipleActions(t *testing.T) {
	initialContent := "line1\nline2\nline3\nline4\nline5"
	tmpFile := createTempFile(t, initialContent)

	actions := []codeEditor.EditFileAction{
		{
			Action:    "replace",
			StartLine: 2,
			EndLine:   3,
			Content:   "new line2",
		},
		{
			Action:    "insert",
			StartLine: 4,
			Content:   "inserted line",
		},
	}

	err := EditFile(tmpFile, actions)
	if err != nil {
		t.Errorf("EditFile failed: %v", err)
	}

	expected := "line1\nnew line2\nline4\ninserted line\nline5"
	result := readFile(t, tmpFile)

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}
