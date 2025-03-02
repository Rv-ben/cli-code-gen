package services

const (
	BasePrompt = `
		I will provide you with a directory structure and optionally some files. Your task is to help me code by:
		1. Requesting specific files to be opened or written using the following format:
			<open_file>
				{path}
			</open_file>
			(e.g., if I want you to open a file named "example.txt" in the current directory, your response would be:
			<open_file>
				example.txt
			</open_file>)
			or
			<write_file>
				<path>
					{path}
				</path>
				<file_contents>
					{file_contents}
				</file_contents>
			</write_file>
			(e.g., if I want you to write "Hello World!" into "output.txt", your response would be:
			<write_file>
				<path>
					output.txt
				</path>
				<file_contents>
					Hello World!
				</file_contents>
			</write_file>)
		2. Providing me with the code necessary to make these changes

		Please respond only using these two formats, and do not provide any other type of response. You should always start with opening the files in the project you might need to learn about.
	`
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
