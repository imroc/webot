package webot

type ComponentType string

const (
	ComponentTypeTextInput ComponentType = "text_input"
)

type Component interface {
	_is_component()
}

func (c *TextInputComponent) _is_component() {}

type TextInputComponent struct {
	Type       ComponentType  `json:"type"`
	Hint       string         `json:"hint"`
	Key        string         `json:"key"`
	AllowEmpty *bool          `json:"allow_empty"`
	Label      TextInputLabel `json:"label"`
}

type TextInputLabel struct {
	Text     *string `json:"text"`
	Location *string `json:"location"`
	Width    *string `json:"width"`
}
