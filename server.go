package webot

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/imroc/webot/internal/wxbizmsgcrypt"
)

type Server struct {
	log                       Logger
	wxcpt                     *wxbizmsgcrypt.WXBizMsgCrypt
	client                    *Client
	token                     string
	encodingAeskey            string
	robotName                 string
	messageHandlers           []MessageHandler
	textMessageHandlers       []TextMessageHandler
	imageMessageHandlers      []ImageMessageHandler
	eventMessageHandlers      []EventMessageHandler
	attachmentMessageHandlers []AttachmentMessageHandler
}

func NewServer(token, encodingAeskey, robotName string) *Server {
	return &Server{
		token:          token,
		encodingAeskey: encodingAeskey,
		wxcpt:          wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, "", wxbizmsgcrypt.XmlType),
		log:            createDefaultLogger(),
		client:         NewClient(),
	}
}

func (s *Server) GetClient() *Client {
	return s.client
}

func (s *Server) verifyURL(msgSignature, timestamp, nonce, echoStr string) ([]byte, error) {
	result, err := s.wxcpt.VerifyURL(msgSignature, timestamp, nonce, echoStr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) DecryptJsonMsg(msgSignature, timestamp, nonce string, data []byte) (msg *CallbackMessage, err error) {
	rawMsg, cryptErr := s.wxcpt.DecryptMsg(msgSignature, timestamp, nonce, data)
	if cryptErr != nil {
		return nil, cryptErr
	}
	msg = &CallbackMessage{}
	err = xml.Unmarshal(rawMsg, &msg)
	if err != nil {
		msg = nil
	}
	return
}

type (
	MessageHandler           func(client *Client, msg CallbackMessage) error
	TextMessageHandler       func(client *Client, msg CallbackMessageCommonItem, text Text) error
	ImageMessageHandler      func(client *Client, msg CallbackMessageCommonItem, image Image) error
	EventMessageHandler      func(client *Client, msg CallbackMessageCommonItem, image Event) error
	AttachmentMessageHandler func(client *Client, msg CallbackMessageCommonItem, image Attachment) error
)

func (s *Server) HandleTextMessage(fn TextMessageHandler) *Server {
	s.textMessageHandlers = append(s.textMessageHandlers, fn)
	return s
}

func (s *Server) HandleImageMessage(fn ImageMessageHandler) *Server {
	s.imageMessageHandlers = append(s.imageMessageHandlers, fn)
	return s
}

func (s *Server) HandleEventMessage(fn EventMessageHandler) *Server {
	s.eventMessageHandlers = append(s.eventMessageHandlers, fn)
	return s
}

func (s *Server) HandleAttachmentMessage(fn AttachmentMessageHandler) *Server {
	s.attachmentMessageHandlers = append(s.attachmentMessageHandlers, fn)
	return s
}

func (s *Server) HandleMessage(fn MessageHandler) *Server {
	s.messageHandlers = append(s.messageHandlers, fn)
	return s
}

func (s *Server) cleanContent(content string) string {
	str := strings.ReplaceAll(content, "@"+s.robotName, "")
	return strings.TrimSpace(str)
}

func (s *Server) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	msg_signature := urlQuery.Get("msg_signature")
	timestamp := urlQuery.Get("timestamp")
	nonce := urlQuery.Get("nonce")
	echostr := urlQuery.Get("echostr")

	switch r.Method {
	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.log.Errorf("failed to read body: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		body, cryptErr := s.wxcpt.DecryptMsg(msg_signature, timestamp, nonce, body)
		if cryptErr != nil {
			s.log.Errorf("failed to decrypt message: %s", cryptErr.Error())
			http.Error(w, cryptErr.Error(), http.StatusBadRequest)
			return
		}
		s.log.Infof("received body: \n%s", string(body))
		var msg CallbackMessage
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			s.log.Errorf("failed to unmarshal xml message: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, handler := range s.messageHandlers {
			handler(s.client, msg)
		}
		switch msg.MsgType {
		case CallbackMessageTypeText:
			if msg.Text == nil {
				errMsg := "no text found in text message"
				s.log.Errorf(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
			for _, handler := range s.textMessageHandlers {
				handler(s.client, msg.CallbackMessageCommonItem, *msg.Text)
			}
		case CallbackMessageTypeImage:
			if msg.Image == nil {
				errMsg := "no image found in image message"
				s.log.Errorf(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
			for _, handler := range s.imageMessageHandlers {
				handler(s.client, msg.CallbackMessageCommonItem, *msg.Image)
			}
		case CallbackMessageTypeEvent:
			if msg.Event == nil {
				errMsg := "no event found in event message"
				s.log.Errorf(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
			for _, handler := range s.eventMessageHandlers {
				handler(s.client, msg.CallbackMessageCommonItem, *msg.Event)
			}
		case CallbackMessageTypeAttachment:
			if msg.Attachment == nil {
				errMsg := "no attachment found in attachment message"
				s.log.Errorf(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
			for _, handler := range s.attachmentMessageHandlers {
				handler(s.client, msg.CallbackMessageCommonItem, *msg.Attachment)
			}
		}
	case "GET":
		if echostr != "" {
			echostr, cryptErr := s.verifyURL(msg_signature, timestamp, nonce, echostr)
			if cryptErr != nil {
				s.log.Errorf("failed to verify url: %s", cryptErr.Error())
				http.Error(w, cryptErr.Error(), http.StatusBadRequest)
			} else {
				s.log.Infof("verifyUrl success echostr: %s", echostr)
				w.Write(echostr)
			}
			return
		} else {
			s.log.Errorf("empty echostr in GET request")
			http.Error(w, "empty echostr in GET request", http.StatusBadRequest)
			return
		}
	}
}
