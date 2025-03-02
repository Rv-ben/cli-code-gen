package services

const (
	BasePrompt = "You are a helpful assistant that can help me with my code. I will provide you with a directory structure and some files. You will then provide me with a list of files that need to be changed. You will then provide me with the code to change the files." +
		"You will only reply to me two actions: Request a file to read in the format <open_file>\n{PATH}\n</open_file>\n and be relative to the current directory" +
		"or Request a file to write in the format <write_file>\n<path>\n<file_contents>\n</write_file>\n" +
		"You will then provide me with the code to change the files. DO NOT PROVIDE ME WITH ANY OTHER FORM OF RESPONSE."
)

type BasePromptProvider struct {
	prompt string
}

func NewBasePromptProvider() *BasePromptProvider {
	return &BasePromptProvider{prompt: BasePrompt}
}

func (b *BasePromptProvider) GetPrompt() string {
	return b.prompt
}
