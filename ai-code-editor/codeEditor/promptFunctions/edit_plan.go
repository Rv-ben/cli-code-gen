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

Parse out each file that needs to be edited. For each file provide:
1. The full file path
2. A description of what changes need to be made
3. Any constraints or considerations for the edit

Format the response as a JSON array of objects with fields:
- filePath: string
- description: string  
- constraints: string

Only include files that will be edited. Be specific about the changes needed.

Do not include any explanation or additional text, just the JSON array.
`, e.Plan)

	response, err := e.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error parsing edit plan: %v", err)
		return ""
	}

	return response
}
