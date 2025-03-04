package promptFunctions

import (
	"ai-code-editor/config"
	"fmt"
	"log"
)

type EditPlan struct {
	Plan string
	*BasePromptFunction
}

type EditAction struct {
	FilePath    string
	Description string
	Constraints string
}

func NewEditPlan(model string, config *config.Config, plan string) *EditPlan {
	return &EditPlan{
		Plan:               plan,
		BasePromptFunction: NewBasePromptFunction(model, config),
	}
}

func (e *EditPlan) GetEditActions() string {
	prompt := fmt.Sprintf(`
Given this plan of action:
%s

Parse out each file that needs to be edited, created, or deleted. For each file provide:
1. The full file path
2. A general description of what changes need to be made to the file
3. Any constraints for the edit so the edit works well with other files and the overall plan, such as "There needs to be a function called train with a string argument called model_name"
4. Any additional files that need to be read to understand the context of the edit
5. The general edit plan for all files

Format the response as a JSON array of objects with fields:
- filePath: string
- description: string
- constraints: string
- additionalFilesForContext: string[]
- generalPlan: string
Only include files that will be edited. Be specific about the changes needed.

Do not include references to your reasoning or any other explanation in the formatted response.
Do not include any explanation or additional text, just the JSON array.
`, e.Plan)

	response, err := e.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error parsing edit plan: %v", err)
		return ""
	}

	return response
}
