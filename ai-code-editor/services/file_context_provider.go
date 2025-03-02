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

func NewFileContextProvider(files []string) *FileContextProvider {
	return &FileContextProvider{files: files}
}

func (f *FileContextProvider) GetFileContents(currentDir string) string {
	fileContents := ""
	for _, file := range f.files {
		// Convert file path to be relative to current directory
		relPath, err := filepath.Rel(currentDir, filepath.Join(currentDir, file))
		if err != nil {
			log.Printf("Error getting relative path for %s: %v", file, err)
			continue
		}

		content, err := os.ReadFile(relPath)
		if err != nil {
			log.Printf("Error reading file %s: %v", relPath, err)
			continue
		}
		fileContents += fmt.Sprintf("\n<open_file>\n%s\n```%s\n<contents>\n%s\n</contents>\n```\n</open_file>\n",
			relPath, relPath, content)
	}

	return fileContents
}
