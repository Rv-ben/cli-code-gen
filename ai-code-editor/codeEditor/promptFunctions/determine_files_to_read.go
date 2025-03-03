package promptFunctions

import (
	"ai-code-editor/config"
	"fmt"
	"log"
	"strings"
)

type DetermineFilesToRead struct {
	Context        string
	AvailableFiles []string
	SelectedFiles  []string
	*BasePromptFunction
}

func NewDetermineFilesToRead(model string, config *config.Config, context string, availableFiles []string) *DetermineFilesToRead {
	return &DetermineFilesToRead{
		Context:            context,
		AvailableFiles:     availableFiles,
		SelectedFiles:      []string{},
		BasePromptFunction: NewBasePromptFunction(model, config),
	}
}

func (d *DetermineFilesToRead) GetFilesToRead() []string {
	filesStr := strings.Join(d.AvailableFiles, "\n")

	prompt := fmt.Sprintf(`
Given this context about what we need to understand:
%s

And this list of available files:
%s

List only the file FULL paths (one per line) that would be most relevant to examine, based on the context.
Only include files that are actually in the list above.
Keep the list minimal - only include files that are directly relevant. Order the files by relevance.
Do not include wildcard paths like * or **.
Do not include any explanation or additional text.
`, d.Context, filesStr)

	response, err := d.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error determining files to read: %v", err)
		return []string{}
	}

	// Clean up the response and convert to slice
	files := strings.Split(strings.TrimSpace(response), "\n")

	// Filter out any empty lines
	var cleanFiles []string
	for _, file := range files {
		if trimmed := strings.TrimSpace(file); trimmed != "" {
			cleanFiles = append(cleanFiles, trimmed)
		}
	}

	d.SelectedFiles = cleanFiles
	return cleanFiles
}
