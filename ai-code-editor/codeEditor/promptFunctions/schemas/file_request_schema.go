package schemas

type FileRequestSchema struct {
	Type       string `json:"type"`
	Properties struct {
		Actions struct {
			Type  string `json:"type"`
			Items struct {
				Type       string `json:"type"`
				Properties struct {
					Type string `json:"type"`
					Path string `json:"path"`
				} `json:"properties"`
				Required []string `json:"required"`
			} `json:"items"`
		} `json:"actions"`
	} `json:"properties"`
}

func NewFileRequestSchema() *FileRequestSchema {
	schema := &FileRequestSchema{
		Type: "object",
	}
	schema.Properties.Actions.Type = "array"
	schema.Properties.Actions.Items.Type = "object"
	schema.Properties.Actions.Items.Properties.Type = "open_file"
	schema.Properties.Actions.Items.Properties.Path = "path/to/file"
	schema.Properties.Actions.Items.Required = []string{"type", "path"}
	return schema
}
