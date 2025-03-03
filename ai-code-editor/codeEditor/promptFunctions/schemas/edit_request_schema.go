package schemas

type EditRequestSchema struct {
	Type       string `json:"type"`
	Properties struct {
		Actions struct {
			Type  string `json:"type"`
			Items struct {
				Type       string `json:"type"`
				Properties struct {
					Type      string `json:"type"`
					Path      string `json:"path"`
					Content   string `json:"content"`
					StartLine string `json:"start_line"`
					EndLine   string `json:"end_line"`
					Action    string `json:"action"`
				} `json:"properties"`
				Required []string `json:"required"`
			} `json:"items"`
		} `json:"actions"`
	} `json:"properties"`
}

func NewEditRequestSchema() *EditRequestSchema {
	schema := &EditRequestSchema{
		Type: "object",
	}
	schema.Properties.Actions.Type = "array"
	schema.Properties.Actions.Items.Type = "object"
	schema.Properties.Actions.Items.Properties.Type = "write_file"
	schema.Properties.Actions.Items.Properties.Path = "path/to/file"
	schema.Properties.Actions.Items.Properties.Content = "file contents here"
	schema.Properties.Actions.Items.Properties.StartLine = "x"
	schema.Properties.Actions.Items.Properties.EndLine = "y"
	schema.Properties.Actions.Items.Properties.Action = "replace"
	schema.Properties.Actions.Items.Required = []string{"type", "path", "content", "start_line", "end_line", "action"}
	return schema
}
