package services

import (
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

// GenerateTree generates a tree-like representation of the directory structure
func (dt *DirectoryTree) GenerateTree(rootPath string) (string, error) {
	var builder strings.Builder

	// Write the root directory name
	rootName := filepath.Base(rootPath)
	builder.WriteString(fmt.Sprintf("|-- %s/\n", rootName))

	err := dt.walkDirectory(rootPath, &builder, "", 1)
	if err != nil {
		return "", fmt.Errorf("error generating directory tree: %w", err)
	}

	return builder.String(), nil
}

// walkDirectory recursively walks through the directory structure
func (dt *DirectoryTree) walkDirectory(path string, builder *strings.Builder, prefix string, depth int) error {
	if depth > dt.maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Sort entries to ensure consistent output
	for i, entry := range entries {
		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Check if path should be skipped
		shouldSkip := false
		for _, skipPath := range dt.skipPaths {
			if strings.Contains(filepath.Join(path, entry.Name()), skipPath) {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}

		// Determine if this is the last item at this level
		isLast := i == len(entries)-1

		// Create the line prefix
		linePrefix := prefix
		if isLast {
			linePrefix += "|-- "
		} else {
			linePrefix += "|-- "
		}

		// Write the entry
		fullPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			builder.WriteString(fmt.Sprintf("%s%s/\n", linePrefix, entry.Name()))

			// Calculate new prefix for children
			newPrefix := prefix
			if isLast {
				newPrefix += dt.indent
			} else {
				newPrefix += dt.indent
			}

			// Recursively process subdirectories
			err := dt.walkDirectory(fullPath, builder, newPrefix, depth+1)
			if err != nil {
				return err
			}
		} else {
			builder.WriteString(fmt.Sprintf("%s%s\n", linePrefix, entry.Name()))
		}
	}

	return nil
}
