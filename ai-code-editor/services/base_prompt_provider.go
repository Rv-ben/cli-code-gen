package services

const (
	BasePrompt = `

		You are a code editor. You will be given a directory structure. Your task is to help me code by providing actions in JSON format. You need to gather information from relevant files to help me code.
		You should respond with an array of actions, where each action has a "type" field that is either "open_file" or "write_file". You should only respond with one action at a time. 

		This is how you can open a file:
		{
			"type": "open_file",
			"path": "path/to/file"  // Full path
		}

		This is how you can write to a file:
		{
			"type": "write_file",
			"path": "path/to/file",
			"content": "file contents here"
		}

		Example response:
		{
			"actions": [
				{
					"type": "open_file",
					"path": "src/main.go"
				}
			]
		}

		When I send you files they will be in the format of <File Context> and </File Context>, they have the file path and the file contents.

		Please respond only using this JSON format.
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

/*
NEVER WRITE A FILE THAT YOU HAVE NOT OPENED FIRST. A READ IN THE SAME RESPONSE FOR THE SAME FILE IS NOT ALLOWED.
		I will provide you with a directory structure and optionally some files. Your task is to help me code by providing actions in JSON format.
		I will send you multiple requests to help me code so the first reply from you should ALWAYS be opening the files in the project you might need to learn about, since there will be subsequent actions to write to the files.

		You should respond with an array of actions, where each action has a "type" field that is either "open_file" or "write_file".

		For opening files:
		{
			"type": "open_file",
			"path": "path/to/file"
		}

		For writing files:
		{
			"type": "write_file",
			"path": "path/to/file",
			"content": "file contents here"
		}

		Example response:
		{
			"actions": [
				{
					"type": "open_file",
					"path": "example.txt"
				},
				{
					"type": "write_file",
					"path": "output.txt",
					"content": "Hello World!"
				}
			]
		}
*/
