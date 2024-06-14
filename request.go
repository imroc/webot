package webot

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
)

type Request struct {
	*req.Request
	msg        *SendMessage
	webhookUrl string
}

func (c *Client) NewRequest(webhookUrl string) *Request {
	return &Request{
		Request:    c.client.R(),
		webhookUrl: webhookUrl,
		msg:        &SendMessage{},
	}
}

func (r *Request) Send() (resp *Response, err error) {
	resp = &Response{}
	res, err := r.
		SetBodyJsonMarshal(r.msg).
		EnableDumpWithoutRequest().
		SetResult(resp).
		Post(r.webhookUrl)
	if err != nil {
		return
	}
	if !res.IsSuccess() || resp.Errcode != 0 {
		err = fmt.Errorf("bad response:\n%s", res.Dump())
		return
	}
	return
}

func (r *Request) SendFileContent(filename string, content []byte) (resp *Response, err error) {
	upload, err := r.Upload(filename, content)
	if err != nil {
		return
	}
	file := &FileMessage{
		MediaId: upload.MediaId,
	}
	r.msg.Msgtype = SendMessageTypeFile
	r.msg.File = file
	return r.Send()
}

func (r *Request) Upload(filename string, data []byte) (resp *UploadResponse, err error) {
	uploadUrl := strings.ReplaceAll(r.webhookUrl, "webhook/send", "webhook/upload_media")
	resp = &UploadResponse{}
	cd := new(req.ContentDisposition)
	cd.Add("filelength", strconv.Itoa(len(data)))
	res, err := r.
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
		Post(uploadUrl)
	if err != nil {
		return
	}
	if !res.IsSuccess() || resp.Errcode != 0 {
		err = fmt.Errorf("bad response:\n%s", res.Dump())
		return
	}
	return
}

func (r *Request) SendMarkdownString(markdown string) (resp *Response, err error) {
	return r.SendMarkdown(&MarkdownMessage{
		Content: markdown,
	})
}

func (r *Request) SendMarkdown(markdown *MarkdownMessage) (resp *Response, err error) {
	r.msg.Msgtype = SendMessageTypeMarkdown
	r.msg.Markdown = markdown
	return r.Send()
}

func (r *Request) SendText(text *TextMessage) (resp *Response, err error) {
	r.msg.Msgtype = SendMessageTypeText
	r.msg.Text = text
	return r.Send()
}

func (r *Request) SendTextString(text string) (resp *Response, err error) {
	return r.SendText(&TextMessage{
		Content: text,
	})
}