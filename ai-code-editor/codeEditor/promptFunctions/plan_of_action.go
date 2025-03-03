package promptFunctions

import (
	"ai-code-editor/config"
	"fmt"
	"log"
)

type PlanOfAction struct {
	Plan     string
	Context  string
	UserTask string
	*BasePromptFunction
}

func NewPlanOfAction(model string, config *config.Config, userTask string, context string) *PlanOfAction {
	return &PlanOfAction{
		UserTask:           userTask,
		Context:            context,
		BasePromptFunction: NewBasePromptFunction(model, config),
	}
}

func (p *PlanOfAction) GetPlan(request string) string {
	prompt := fmt.Sprintf(`
		Given this context: %s
		Create a plan of action in a few bullet points describing how to implement this request: %s
		Say things open a file or edit a file.
		Keep it short and concise.
		`, p.Context, p.UserTask)

	response, err := p.ExecutePrompt(prompt)
	if err != nil {
		log.Printf("Error getting plan of action: %v", err)
		return ""
	}

	return response
}
