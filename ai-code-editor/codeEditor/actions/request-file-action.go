package codeEditor

type RequestFileAction struct {
	actionName string
	Path       string
}

func NewRequestFileAction(path string) *RequestFileAction {
	return &RequestFileAction{
		actionName: "open_file",
		Path:       path,
	}
}

func (r *RequestFileAction) ToString() string {
	return "<open_file>\n<path>" + r.Path + "\n</open_file>"
}

func (r *RequestFileAction) GetType() string {
	return r.actionName
}
