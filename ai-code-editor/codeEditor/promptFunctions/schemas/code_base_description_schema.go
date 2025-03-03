package schemas

type CodeBaseDescriptionSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

func NewCodeBaseDescriptionSchema() *CodeBaseDescriptionSchema {
	return &CodeBaseDescriptionSchema{
		Type: "object",
		Properties: map[string]Property{
			"description": {
				Type:        "string",
				Description: "A description of the codebase in a few sentences",
			},
		},
		Required: []string{"description"},
	}
}
