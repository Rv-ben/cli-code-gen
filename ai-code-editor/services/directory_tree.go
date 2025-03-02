package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

// DirectoryTree represents a service for generating directory tree structures
type DirectoryTree struct {
	indent          string
	maxDepth        int
	skipPaths       []string
	parser          *tree_sitter.Parser
	tree            *tree_sitter.Tree
	directoryString string
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

	parser := tree_sitter.NewParser()

	return &DirectoryTree{
		indent:    indent,
		maxDepth:  maxDepth,
		skipPaths: skipPaths,
		parser:    parser,
	}
}

// GenerateTree generates a tree-sitter tree representation of the directory structure
func (dt *DirectoryTree) GenerateTree(rootPath string) (string, error) {
	if dt.tree != nil {
		dt.tree.Close()
	}

	dt.parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language()))

	contents, err := dt.buildDirectoryText(rootPath, 0)
	if err != nil {
		return "", fmt.Errorf("error generating directory tree: %w", err)
	}

	dt.directoryString = contents

	return contents, nil
}

// Close releases resources associated with the DirectoryTree
func (dt *DirectoryTree) Close() {
	if dt.tree != nil {
		dt.tree.Close()
		dt.tree = nil
	}
	if dt.parser != nil {
		dt.parser.Close()
		dt.parser = nil
	}
}

// buildDirectoryText recursively builds a text representation of the directory structure
func (dt *DirectoryTree) buildDirectoryText(path string, depth int) (string, error) {
	if depth > dt.maxDepth {
		return "", nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(info.Name(), ".") {
		return "", nil
	}

	for _, skipPath := range dt.skipPaths {
		if strings.Contains(path, skipPath) {
			return "", nil
		}
	}

	var sb strings.Builder

	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}

		for _, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			childText, err := dt.buildDirectoryText(childPath, depth+1)
			if err != nil {
				return "", err
			}
			if childText != "" {
				sb.WriteString(childText)
			}
		}
	} else {
		sb.WriteString("-------\n")
		sb.WriteString(path + "\n")
		sb.WriteString(dt.GetTree(path) + "\n")
	}

	return sb.String(), nil
}

func (dt *DirectoryTree) GetDirectoryString(rootPath string) string {
	dt.GenerateTree(rootPath)

	return dt.directoryString
}

func (dt *DirectoryTree) GetTree(filePath string) string {
	log.Printf("Getting tree for file: %s", filePath)

	// Only parse .go files
	if !strings.HasSuffix(filePath, ".go") {
		log.Printf("Skipping non-Go file: %s", filePath)
		return ""
	}

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return ""
	}

	tree := dt.parser.Parse(nil, nil)
	if tree == nil {
		log.Printf("Failed to create parse tree for %s", filePath)
		return ""
	}

	tree = dt.parser.Parse(fileContents, nil)
	if tree == nil {
		log.Printf("Failed to parse file %s", filePath)
		return ""
	}

	return dt.GetFunctionDeclarations(tree.RootNode(), fileContents)
}

func (dt *DirectoryTree) GetFunctionDeclarations(node *tree_sitter.Node, fileContents []byte) string {
	cursor := node.Walk()
	defer cursor.Close()

	var declarations strings.Builder

	// Move the cursor to the first child
	cursor.GotoFirstChild()

	// Traverse the AST
	for {
		log.Printf("Cursor node: %s", cursor.Node().Kind())
		if cursor.Node().Kind() == "method_declaration" || cursor.Node().Kind() == "function_declaration" {
			funcNode := cursor.Node()

			// Get function name
			nameNode := funcNode.ChildByFieldName("name")
			if nameNode != nil {
				declarations.WriteString("func ")

				// Handle receiver for methods
				if cursor.Node().Kind() == "method_declaration" {
					receiverNode := funcNode.ChildByFieldName("receiver")
					if receiverNode != nil {
						declarations.WriteString("(")
						declarations.WriteString(receiverNode.Utf8Text(fileContents))
						declarations.WriteString(") ")
					}
				}

				declarations.WriteString(nameNode.Utf8Text(fileContents))

				// Get parameters
				paramListNode := funcNode.ChildByFieldName("parameters")
				if paramListNode != nil {
					declarations.WriteString(paramListNode.Utf8Text(fileContents))
				}

				// Get return type if exists
				resultNode := funcNode.ChildByFieldName("result")
				if resultNode != nil {
					declarations.WriteString(" ")
					declarations.WriteString(resultNode.Utf8Text(fileContents))
				}

				declarations.WriteString("\n")
			}
		}

		if !cursor.GotoNextSibling() {
			break
		}
	}

	return declarations.String()
}
