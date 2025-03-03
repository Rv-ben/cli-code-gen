package promptFunctions

import (
	"ai-code-editor/config"
	"ai-code-editor/services"
	"fmt"
	"log"
)

type CodeBaseDescription struct {
	Path            string
	Description     string
	ProjectTreeJson string
	*BasePromptFunction
}

func NewCodeBaseDescription(path string, model string, config *config.Config) *CodeBaseDescription {
	projectTreeJson := services.NewDirectoryTree(
		"    ", // indent
		10,     // maxDepth
		[]string{ // skipPaths
			"node_modules",
			"vendor",
			".git",
		},
	).GetDirectoryString(path)

	return &CodeBaseDescription{
		Path:               path,
		Description:        "",
		ProjectTreeJson:    projectTreeJson,
		BasePromptFunction: NewBasePromptFunction(model, config),
	}
}

func (c *CodeBaseDescription) GetDescription() string {
	prompt := fmt.Sprintf(`
Describe the codebase in a few sentences, keep it short and concise, give a short list of important files and directories. Identify things like main entry points, and programming languages used. The codebase is located at
	%s. 
Here is the directory structure: 
	%s 
	`, c.Path, c.ProjectTreeJson)

	response, err := c.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error getting codebase description: %v", err)
		return ""
	}

	return response
}
