package webot

import (
	"fmt"
	"github.com/imroc/req/v3"
)

type WeBot struct {
	client     *req.Client
	webhookURL string
}

type TextMessage struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type MarkdownMessage struct {
	Content string `json:"content"`
}

type Message struct {
	Msgtype  string           `json:"msgtype"`
	Text     *TextMessage     `json:"text,omitempty"`
	Markdown *MarkdownMessage `json:"markdown,omitempty"`
}

type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func New(webhoookURL string) *WeBot {
	return &WeBot{
		client:     req.C(),
		webhookURL: webhoookURL,
	}
}

func (wb *WeBot) Client() *req.Client {
	return wb.client
}

func (wb *WeBot) Send(msg *Message) (resp *Response, err error) {
	resp = &Response{}
	r, err := wb.client.R().
		SetBodyJsonMarshal(msg).
		EnableDumpWithoutRequest().
		SetResult(resp).
		Post(wb.webhookURL)
	if err != nil {
		return
	}
	if !r.IsSuccess() {
		err = fmt.Errorf("bad response:\n%s", r.Dump())
	}
	return
}

func (wb *WeBot) SendMarkdownContent(markdown string) (resp *Response, err error) {
	return wb.SendMarkdown(&MarkdownMessage{
		Content: markdown,
	})
}

func (wb *WeBot) SendMarkdown(markdown *MarkdownMessage) (resp *Response, err error) {
	msg := &Message{Msgtype: "markdown", Markdown: markdown}
	return wb.Send(msg)
}

func (wb *WeBot) SendText(text *TextMessage) (resp *Response, err error) {
	msg := &Message{Msgtype: "text", Text: text}
	return wb.Send(msg)
}

func (wb *WeBot) SendTextContent(text string) (resp *Response, err error) {
	msg := &TextMessage{
		Content: text,
	}
	return wb.SendText(msg)
}

func (wb *WeBot) Debug(debug bool) {
	if debug {
		wb.client.EnableDumpAll().EnableDebugLog().EnableTraceAll()
	} else {
		wb.client.DisableDebugLog().DisableDumpAll().DisableTraceAll()
	}
}
