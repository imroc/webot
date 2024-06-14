package webot

type SendMessageType string

const (
	SendMessageTypeText         SendMessageType = "text"
	SendMessageTypeMarkdown     SendMessageType = "markdown"
	SendMessageTypeAttachment   SendMessageType = "attachment"
	SendMessageTypeImage        SendMessageType = "image"
	SendMessageTypeMiniprogram  SendMessageType = "miniprogram"
	SendMessageTypeFile         SendMessageType = "file"
	SendMessageTypeNews         SendMessageType = "news"
	SendMessageTypeTemplateCard SendMessageType = "template_card"
)

type TextMessage struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type MarkdownMessage struct {
	Content string `json:"content"`
}

type FileMessage struct {
	MediaId string `json:"media_id"`
}

type SendMessage struct {
	Msgtype       SendMessageType  `json:"msgtype"`
	Text          *TextMessage     `json:"text,omitempty"`
	Markdown      *MarkdownMessage `json:"markdown,omitempty"`
	File          *FileMessage     `json:"file,omitempty"`
	Chatid        string           `json:"chatid,omitempty"`
	PostId        string           `json:"post_id,omitempty"`
	VisibleToUser string           `json:"visible_to_user,omitempty"`
}
