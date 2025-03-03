package main

import (
	promptFunctions "ai-code-editor/codeEditor/promptFunctions"
	"ai-code-editor/config"
	"ai-code-editor/services"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	config := config.Load()

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return
	}

	// Check for required arguments
	if len(os.Args) < 1 {
		fmt.Println("Usage: ollama-cli <prompt> [files...]")
		fmt.Println("Example: ollama-cli 'Fix the bug' file1.go file2.go")
		os.Exit(1)
	}

	userTask := os.Args[1]

	selectedModel := config.MediumModel

	fmt.Printf("Using model: %s\n User task: %s\n", selectedModel, userTask)

	codeBaseDescription := promptFunctions.NewCodeBaseDescription(currentDir, selectedModel, config)
	description := codeBaseDescription.GetDescription()

	fmt.Printf("Codebase description: %s\n", description)

	gainProblemContext := promptFunctions.NewGainProblemContext(selectedModel, config, userTask, description)
	requiredContext := gainProblemContext.GetRequiredContext()

	fmt.Printf("Required context: %s\n", requiredContext)

	directoryTree := services.NewDirectoryTree(currentDir, 10, []string{})
	availableFiles := directoryTree.GetKnownFiles()

	determineFilesToRead := promptFunctions.NewDetermineFilesToRead(selectedModel, config, requiredContext, availableFiles)
	selectedFiles := determineFilesToRead.GetFilesToRead()

	fmt.Printf("Selected files: %v\n", selectedFiles)

	planOfAction := promptFunctions.NewPlanOfAction(selectedModel, config, userTask, requiredContext+"\n"+strings.Join(selectedFiles, "\n"))
	plan := planOfAction.GetPlan(userTask)

	fmt.Printf("Plan of action: %s\n", plan)

	editPlan := promptFunctions.NewEditPlan(selectedModel, config, plan)
	editActions := editPlan.GetEditActions()

	fmt.Printf("Edit actions: %v\n", editActions)

	// Hard code the files to read for now to be the current directory/main.go use the os package
	selectedFiles = []string{"main.go"}

	readFiles := services.NewFileContextProvider(selectedFiles)

	fileContext := readFiles.GetFileContents()

	fmt.Printf("File context: %v\n", fileContext)

	editFile := promptFunctions.NewEditFile(selectedModel, config, selectedFiles[0], fileContext, editActions)
	edit := editFile.GetEdit(userTask)

	fmt.Printf("Edit: %v\n", edit)
}
