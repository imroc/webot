package webot

type From struct {
	UserId string `xml:"UserId"`
	Name   string `xml:"Name"`
	Alias  string `xml:"Alias"`
}

type Text struct {
	Content string `xml:"Content"`
}

type Image struct {
	ImageUrl string `xml:"ImageUrl"`
}

type Event struct {
	EventType string `xml:"EventType"`
}

type Attachment struct {
	CallbackId string  `xml:"CallbackId"`
	Actions    Actions `xml:"Actions"`
}

type Actions struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
	Type  string `xml:"Type"`
}

type CallbackMessageItem struct {
	MsgType    string      `xml:"MsgType"`
	Text       *Text       `xml:"Text,omitempty"`
	Image      *Image      `xml:"Image,omitempty"`
	Event      *Event      `xml:"Event,omitempty"`
	Attachment *Attachment `xml:"Attachment,omitempty"`
}

type CallbackMessage struct {
	CallbackMessageItem
	WebhookUrl     string `xml:"WebhookUrl"`
	ChatId         string `xml:"ChatId"`
	PostId         string `xml:"PostId"`
	ChatType       string `xml:"ChatType"`
	GetChatInfoUrl string `xml:"GetChatInfoUrl"`
	MsgId          string `xml:"MsgId"`
	MsgType        string `json:"MsgType"`
	From           From   `xml:"From"`
	AppVersion     string `xml:"AppVersion"`
}
