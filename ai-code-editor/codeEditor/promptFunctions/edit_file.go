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

Provide multiple edit file action with:
- The full file path
- The content to be inserted or replaced
- The line numbers where the edit occurs (start and end)
- The edit type (replace, insert, delete)

Format the response as a JSON array of objects with fields:
- filePath: string
- content: string
- startLine: int
- endLine: int
- editType: string (replace, insert, delete)

Keep the code changes focused.
`, e.SourceFile, e.FileContent, editRequest, e.Constraints)

	response, err := e.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error getting file edit: %v", err)
		return ""
	}

	return response
}
