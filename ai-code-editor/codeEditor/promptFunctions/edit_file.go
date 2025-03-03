package promptFunctions

import (
	"ai-code-editor/config"
	"fmt"
	"log"
)

type EditFile struct {
	SourceFile  string
	FileContent string
	Constraints string
	*BasePromptFunction
}

func NewEditFile(model string, config *config.Config, sourceFile string, fileContent string, constraints string) *EditFile {
	return &EditFile{
		SourceFile:         sourceFile,
		FileContent:        fileContent,
		Constraints:        constraints,
		BasePromptFunction: NewBasePromptFunction(model, config),
	}
}

func (e *EditFile) GetEdit(editRequest string) string {
	prompt := fmt.Sprintf(`
Given this source file: %s
With content:
%s

Make the following edit:
%s

Additional constraints:
%s

Provide the edit in a code block with the file path like so:
'''language:path/to/file
// ... existing code ...
edited code here
// ... existing code ...
'''

Keep the edit minimal and focused. Preserve existing code structure and style.
`, e.SourceFile, e.FileContent, editRequest, e.Constraints)

	response, err := e.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error getting file edit: %v", err)
		return ""
	}

	return response
}
