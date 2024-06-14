package webot

type From struct {
	UserId string `json:"UserId"`
	Name   string `json:"Name"`
	Alias  string `json:"Alias"`
}

type Text struct {
	Content string `json:"Content"`
}

type Image struct {
	ImageUrl string `json:"ImageUrl"`
}

type Event struct {
	EventType string `json:"EventType"`
}

type Attachment struct {
	CallbackId string  `json:"CallbackId"`
	Actions    Actions `json:"Actions"`
}

type Actions struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
	Type  string `json:"Type"`
}

type CallbackMessageItem struct {
	MsgType    string      `json:"MsgType"`
	Text       *Text       `json:"Text,omitempty"`
	Image      *Image      `json:"Image,omitempty"`
	Event      *Event      `json:"Event,omitempty"`
	Attachment *Attachment `json:"Attachment,omitempty"`
}

type CallbackMessage struct {
	CallbackMessageItem
	WebhookUrl     string `json:"WebhookUrl"`
	ChatId         string `json:"ChatId"`
	PostId         string `json:"PostId"`
	ChatType       string `json:"ChatType"`
	GetChatInfoUrl string `json:"GetChatInfoUrl"`
	MsgId          string `json:"MsgId"`
	MsgType        string `json:"MsgType"`
	From           From   `json:"From"`
	AppVersion     string `json:"AppVersion"`
}
