package services

import (
	"encoding/json"
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

// Add these struct definitions after the existing FileNode struct

type DirectoryStructure struct {
	Directory Directory `json:"directory"`
}

type Directory struct {
	Path  string     `json:"path"`
	Files []FileInfo `json:"files"`
}

type FileInfo struct {
	Path               string   `json:"path"`
	FunctionSignatures []string `json:"functionSignatures"`
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
func (dt *DirectoryTree) GenerateTree(rootPath string) (*DirectoryStructure, error) {
	if dt.tree != nil {
		dt.tree.Close()
	}

	dt.parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language()))

	structure := &DirectoryStructure{
		Directory: Directory{
			Path:  rootPath,
			Files: make([]FileInfo, 0),
		},
	}

	err := dt.buildDirectoryStructure(rootPath, structure)
	if err != nil {
		return nil, fmt.Errorf("error generating directory tree: %w", err)
	}

	return structure, nil
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

func (dt *DirectoryTree) GetDirectoryString(rootPath string) string {
	structure, err := dt.GenerateTree(rootPath)
	if err != nil {
		log.Printf("Error generating tree: %v", err)
		return ""
	}

	jsonBytes, err := json.MarshalIndent(structure, "", "    ")
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return ""
	}

	return string(jsonBytes)
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

// Add new method to build the directory structure
func (dt *DirectoryTree) buildDirectoryStructure(path string, structure *DirectoryStructure) error {

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if strings.HasPrefix(info.Name(), ".") {
		return nil
	}

	for _, skipPath := range dt.skipPaths {
		if strings.Contains(path, skipPath) {
			return nil
		}
	}

	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			childPath := filepath.Join(path, entry.Name())
			err := dt.buildDirectoryStructure(childPath, structure)
			if err != nil {
				return err
			}
		}
	} else {
		if strings.HasSuffix(path, ".go") {
			signatures := dt.GetFunctionSignatures(path)
			fileInfo := FileInfo{
				Path:               path,
				FunctionSignatures: signatures,
			}
			structure.Directory.Files = append(structure.Directory.Files, fileInfo)
		}
	}

	return nil
}

// Add new method to get function signatures
func (dt *DirectoryTree) GetFunctionSignatures(filePath string) []string {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return nil
	}

	tree := dt.parser.Parse(nil, nil)
	if tree == nil {
		log.Printf("Failed to create parse tree for %s", filePath)
		return nil
	}

	tree = dt.parser.Parse(fileContents, nil)
	if tree == nil {
		log.Printf("Failed to parse file %s", filePath)
		return nil
	}

	declarations := dt.GetFunctionDeclarations(tree.RootNode(), fileContents)
	return strings.Split(strings.TrimSpace(declarations), "\n")
}
