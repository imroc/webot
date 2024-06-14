package webot

import (
	"github.com/imroc/req/v3"
)

type Client struct {
	client     *req.Client
	webhookURL string
	uploadURL  string
}

func NewClient() *Client {
	return &Client{
		client: req.C(),
	}
}

func (client *Client) Client() *req.Client {
	return client.client
}

func (client *Client) Debug(debug bool) {
	if debug {
		client.client.EnableDumpAll().EnableDebugLog().EnableTraceAll()
	} else {
		client.client.DisableDebugLog().DisableDumpAll().DisableTraceAll()
	}
}
