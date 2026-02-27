package entity

type AdapterInfo struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UIConfig    UIConfig `json:"ui_config"`
}

type UIConfig struct {
	Modes        []UIMode      `json:"modes,omitempty"`
	Fields       []FieldConfig `json:"fields,omitempty"`
	SupportsFile bool          `json:"supports_file"`
	FileTypes    []string      `json:"file_types,omitempty"`
}

type UIMode struct {
	Key    string        `json:"key"`
	Label  string        `json:"label"`
	Fields []FieldConfig `json:"fields"`
}

type FieldConfig struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder,omitempty"`
}
