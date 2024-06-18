package webot

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type Client struct {
	client     *req.Client
	webhookURL string
	uploadURL  string
}

func NewClient() *Client {
	return &Client{
		client: req.C().SetResultStateCheckFunc(func(resp *req.Response) req.ResultState {
			if errCode := resp.GetHeader("Error-Code"); errCode == "0" {
				return req.SuccessState
			}
			return req.ErrorState
		}).OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if errCode := resp.GetHeader("Error-Code"); errCode == "0" {
				return nil
			}
			resp.Err = fmt.Errorf("Error-Code: %s, Error-Msg: %s", resp.GetHeader("Error-Code"), resp.GetHeader("Error-Msg"))
			return nil
		}),
	}
}

func (client *Client) Client() *req.Client {
	return client.client
}

func (client *Client) SetDumpRequest(dump bool) {
	if dump {
		client.client.EnableDumpAll().EnableDebugLog().EnableTraceAll()
	} else {
		client.client.DisableDebugLog().DisableDumpAll().DisableTraceAll()
	}
}
