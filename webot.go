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

type Message struct {
	Msgtype string       `json:"msgtype"`
	Text    *TextMessage `json:"text,omitempty"`
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

func (wb *WeBot) SendText(msg *TextMessage) (resp *Response, err error) {
	resp = &Response{}
	r, err := wb.client.R().
		SetBodyJsonMarshal(&Message{Msgtype: "text", Text: msg}).
		EnableDumpWithoutRequest().
		SetResult(resp).
		Post(wb.webhookURL)
	if err != nil {
		return
	}
	if !r.IsSuccess() {
		err = fmt.Errorf("bad response:\n%s", r.Dump())
	} else {
		fmt.Println("sccess:")
		fmt.Println(r.Dump())
	}
	return
}

func (wb *WeBot) Debug(debug bool) {
	if debug {
		wb.client.EnableDumpAll().EnableDebugLog().EnableTraceAll()
	} else {
		wb.client.DisableDebugLog().DisableDumpAll().DisableTraceAll()
	}
}
