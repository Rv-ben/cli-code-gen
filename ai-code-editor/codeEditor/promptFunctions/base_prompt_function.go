package promptFunctions

import (
	"ai-code-editor/config"
	"ai-code-editor/ollama"
)

type BasePromptFunction struct {
	Client *ollama.Client
	Model  string
}

func NewBasePromptFunction(model string, config *config.Config) *BasePromptFunction {
	client := ollama.NewClient(config.OllamaBaseURL, true)

	return &BasePromptFunction{
		Client: client,
		Model:  model,
	}
}

func (b *BasePromptFunction) ExecutePrompt(prompt string) (string, error) {
	req := ollama.ChatRequest{
		Model:  b.Model,
		Stream: false,
		Messages: []ollama.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	return b.Client.ChatCompletion(req)
}
