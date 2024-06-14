package webot

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
)

type Client struct {
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
	Msgtype       string           `json:"msgtype"`
	Text          *TextMessage     `json:"text,omitempty"`
	Markdown      *MarkdownMessage `json:"markdown,omitempty"`
	File          *FileMessage     `json:"file,omitempty"`
	Chatid        string           `json:"chatid,omitempty"`
	PostId        string           `json:"post_id,omitempty"`
	VisibleToUser string           `json:"visible_to_user,omitempty"`
}

type MessageOption func(*Message)

func WithChatId(chatId string) MessageOption {
	return func(msg *Message) {
		msg.Chatid = chatId
	}
}

func WithPostId(postId string) MessageOption {
	return func(msg *Message) {
		msg.PostId = postId
	}
}

func WithVisibleToUser(visibleToUser string) MessageOption {
	return func(msg *Message) {
		msg.VisibleToUser = visibleToUser
	}
}

func WithReplyCallbackMessage(cm *CallbackMessage) MessageOption {
	return func(msg *Message) {
		msg.Chatid = cm.ChatId
		msg.PostId = cm.PostId
	}
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

func NewClient(webhoookURL string) *Client {
	return &Client{
		client:     req.C(),
		webhookURL: webhoookURL,
	}
}

func (wb *Client) Client() *req.Client {
	return wb.client
}

func (wb *Client) getUploadURL() string {
	if wb.uploadURL != "" {
		return wb.uploadURL
	}
	wb.uploadURL = strings.ReplaceAll(wb.webhookURL, "webhook/send", "webhook/upload_media")
	return wb.uploadURL
}

func (wb *Client) Send(msg *Message, opts ...MessageOption) (resp *Response, err error) {
	for _, opt := range opts {
		opt(msg)
	}
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

func (wb *Client) SendFileContent(filename string, content []byte) (resp *Response, err error) {
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

func (wb *Client) Upload(filename string, data []byte) (resp *UploadResponse, err error) {
	resp = &UploadResponse{}
	cd := new(req.ContentDisposition)
	cd.Add("filelength", strconv.Itoa(len(data)))
	r, err := wb.client.R().
		SetFileUpload(req.FileUpload{
			ParamName: "media",
			FileName:  filename,
			GetFileContent: func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(data)), nil
			},
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

func (wb *Client) SendMarkdownContent(markdown string) (resp *Response, err error) {
	return wb.SendMarkdown(&MarkdownMessage{
		Content: markdown,
	})
}

func (wb *Client) SendMarkdown(markdown *MarkdownMessage, opts ...MessageOption) (resp *Response, err error) {
	msg := &Message{Msgtype: "markdown", Markdown: markdown}
	return wb.Send(msg, opts...)
}

func (wb *Client) SendText(text *TextMessage, opts ...MessageOption) (resp *Response, err error) {
	msg := &Message{Msgtype: "text", Text: text}
	return wb.Send(msg, opts...)
}

func (wb *Client) SendTextContent(text string, opts ...MessageOption) (resp *Response, err error) {
	msg := &TextMessage{
		Content: text,
	}
	return wb.SendText(msg, opts...)
}

func (wb *Client) Debug(debug bool) {
	if debug {
		wb.client.EnableDumpAll().EnableDebugLog().EnableTraceAll()
	} else {
		wb.client.DisableDebugLog().DisableDumpAll().DisableTraceAll()
	}
}
