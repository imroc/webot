package webot

import (
	"bytes"
	"fmt"
	"github.com/imroc/req/v3"
	"strconv"
	"strings"
)

type WeBot struct {
	client     *req.Client
	webhookURL string
	uploadURL  string
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
	File     *FileMessage     `json:"file,omitempty"`
}

type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type UploadResponse struct {
	Response
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

type FileMessage struct {
	MediaId string `json:"media_id"`
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

func (wb *WeBot) getUploadURL() string {
	if wb.uploadURL != "" {
		return wb.uploadURL
	}
	wb.uploadURL = strings.ReplaceAll(wb.webhookURL, "webhook/send", "webhook/upload_media")
	return wb.uploadURL
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
		return
	}
	if resp.Errcode != 0 {
		err = fmt.Errorf(resp.Errmsg)
	}
	return
}

func (wb *WeBot) SendFileContent(filename string, content []byte) (resp *Response, err error) {
	upload, err := wb.Upload(filename, content)
	if err != nil {
		return
	}
	file := &FileMessage{
		MediaId: upload.MediaId,
	}
	return wb.Send(&Message{
		Msgtype: "file",
		File:    file,
	})
}

func (wb *WeBot) Upload(filename string, data []byte) (resp *UploadResponse, err error) {
	resp = &UploadResponse{}
	cd := new(req.ContentDisposition)
	cd.Add("filelength", strconv.Itoa(len(data)))
	r, err := wb.client.R().
		SetFileUpload(req.FileUpload{
			ParamName:               "media",
			FileName:                filename,
			File:                    bytes.NewReader(data),
			ExtraContentDisposition: cd,
		}).EnableDumpWithoutRequest().
		SetQueryParam("type", "file").
		SetResult(resp).
		Post(wb.getUploadURL())
	if err != nil {
		return
	}
	if !r.IsSuccess() {
		err = fmt.Errorf("bad response:\n%s", r.Dump())
		return
	}
	if resp.Errcode != 0 {
		err = fmt.Errorf(resp.Errmsg)
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
