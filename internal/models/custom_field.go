package models

type CustomFieldEnumOption struct {
	GID     string `json:"gid"`
	Name    string `json:"name"`
	Color   string `json:"color,omitempty"`
	Enabled bool   `json:"enabled"`
}

// CustomField represents both a custom field definition and a per-task custom field value.
// When returned on a Task, it carries the current value fields (TextValue, NumberValue, etc.).
// When returned from workspace or project endpoints, it carries EnumOptions for the field definition.
type CustomField struct {
	GID             string                  `json:"gid"`
	Name            string                  `json:"name"`
	ResourceType    string                  `json:"resource_type,omitempty"`
	Type            string                  `json:"type"` // text, number, enum, multi_enum, date, people
	DisplayValue    string                  `json:"display_value,omitempty"`
	TextValue       *string                 `json:"text_value,omitempty"`
	NumberValue     *float64                `json:"number_value,omitempty"`
	EnumValue       *CustomFieldEnumOption  `json:"enum_value,omitempty"`
	MultiEnumValues []CustomFieldEnumOption `json:"multi_enum_values,omitempty"`
	EnumOptions     []CustomFieldEnumOption `json:"enum_options,omitempty"`
}

func (cf CustomField) GetName() string { return cf.Name }
func (cf CustomField) GetGID() string  { return cf.GID }

type CustomFieldSetting struct {
	GID         string         `json:"gid"`
	CustomField CustomField    `json:"custom_field"`
	IsImportant bool           `json:"is_important,omitempty"`
	Parent      *AsanaResource `json:"parent,omitempty"`
	Project     *AsanaResource `json:"project,omitempty"`
}
