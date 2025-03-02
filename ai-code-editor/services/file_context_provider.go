package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileContextProvider struct {
	files []string
}

func NewFileContextProvider(files interface{}) *FileContextProvider {
	switch f := files.(type) {
	case []string:
		return &FileContextProvider{files: f}
	case string:
		return &FileContextProvider{files: []string{f}}
	default:
		return &FileContextProvider{files: []string{}}
	}
}

func (f *FileContextProvider) GetFileContent(path string) (string, error) {
	// Clean the path to remove any duplicate separators
	cleanPath := filepath.Clean(path)

	log.Printf("Getting content for file: %s", cleanPath)
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		log.Printf("Error reading file %s: %v", cleanPath, err)
		return "", fmt.Errorf("error reading file %s: %w", cleanPath, err)
	}

	return string(content), nil
}

func (f *FileContextProvider) GetFileContents() string {
	log.Printf("Getting contents for files: %v", f.files)
	fileContents := ""

	for _, file := range f.files {
		// Clean the path to remove any duplicate separators
		cleanPath := filepath.Clean(file)

		content, err := os.ReadFile(cleanPath)
		if err != nil {
			log.Printf("Error reading file %s: %v", cleanPath, err)
			continue
		}

		fileContents += fmt.Sprintf("\n<File Context>\n%s\n```%s\n%s\n```\n</File Context>\n",
			cleanPath, cleanPath, string(content))
	}

	return fileContents
}
