package promptFunctions

import (
	"ai-code-editor/config"
	"fmt"
	"log"
)

type GainProblemContext struct {
	Task            string
	CodeBaseDesc    string
	RequiredContext string
	*BasePromptFunction
}

func NewGainProblemContext(model string, config *config.Config, task string, codeBaseDesc string) *GainProblemContext {
	return &GainProblemContext{
		Task:               task,
		CodeBaseDesc:       codeBaseDesc,
		RequiredContext:    "",
		BasePromptFunction: NewBasePromptFunction(model, config),
	}
}

func (g *GainProblemContext) GetRequiredContext() string {
	prompt := fmt.Sprintf(`
Given this task: %s

And this codebase description:
%s

List the context and information you would need to solve this task. For example:
- Which files need to be examined
- What specific parts of the code need to be understood
- What additional information might be needed from the user

If you don't know the specific files files needed, just say the general area of the codebase that needs to be examined.

Keep the list concise and focused on gathering the necessary context to solve the task.
`, g.Task, g.CodeBaseDesc)

	response, err := g.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error getting required context: %v", err)
		return ""
	}

	g.RequiredContext = response
	return response
}
