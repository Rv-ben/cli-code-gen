package codeEditor

type EditFileAction struct {
	actionName string
	Path       string
	Content    string
	StartLine  int
	EndLine    int
	Action     string
}

func NewEditFileAction(path string, content string) *EditFileAction {
	return &EditFileAction{
		actionName: "edit_file",
		Path:       path,
		Content:    content,
	}
}

func (e *EditFileAction) ToString() string {
	return "<edit_file>\n<path>" + e.Path + "\n<file_contents>" + e.Content + "\n</edit_file>"
}

func (e *EditFileAction) GetType() string {
	return e.actionName
}

func (e *EditFileAction) GetStartLine() int {
	return e.StartLine
}

func (e *EditFileAction) GetEndLine() int {
	return e.EndLine
}

func (e *EditFileAction) GetAction() string {
	return e.Action
}
