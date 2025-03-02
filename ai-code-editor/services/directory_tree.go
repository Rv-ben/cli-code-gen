package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DirectoryTree represents a service for generating directory tree structures
type DirectoryTree struct {
	indent    string
	maxDepth  int
	skipPaths []string
}

// FileNode represents a file or directory in the tree
type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Children []FileNode `json:"children,omitempty"`
}

// NewDirectoryTree creates a new DirectoryTree service with default settings
func NewDirectoryTree(indent string, maxDepth int, skipPaths []string) *DirectoryTree {
	if indent == "" {
		indent = "    "
	}
	if maxDepth <= 0 {
		maxDepth = 100 // reasonable default
	}
	return &DirectoryTree{
		indent:    indent,
		maxDepth:  maxDepth,
		skipPaths: skipPaths,
	}
}

// GenerateTree generates a JSON representation of the directory structure
func (dt *DirectoryTree) GenerateTree(rootPath string) (string, error) {
	rootNode, err := dt.buildTree(rootPath, 1)
	if err != nil {
		return "", fmt.Errorf("error generating directory tree: %w", err)
	}

	jsonData, err := json.MarshalIndent(rootNode, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling to JSON: %w", err)
	}

	return string(jsonData), nil
}

// buildTree recursively builds the tree structure
func (dt *DirectoryTree) buildTree(path string, depth int) (*FileNode, error) {
	if depth > dt.maxDepth {
		return nil, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Skip hidden files and directories
	if strings.HasPrefix(info.Name(), ".") {
		return nil, nil
	}

	// Check if path should be skipped
	for _, skipPath := range dt.skipPaths {
		if strings.Contains(path, skipPath) {
			return nil, nil
		}
	}

	node := &FileNode{
		Name:  info.Name(),
		Path:  path,
		IsDir: info.IsDir(),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			childNode, err := dt.buildTree(childPath, depth+1)
			if err != nil {
				return nil, err
			}
			if childNode != nil {
				node.Children = append(node.Children, *childNode)
			}
		}
	}

	return node, nil
}
