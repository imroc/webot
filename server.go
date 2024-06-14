package webot

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imroc/webot/internal/wxbizmsgcrypt"
)

type MessageType string

const (
	MessageTypeText       MessageType = "text"
	MessageTypeAttachment MessageType = "attachment"
	MessageTypeImage      MessageType = "image"
	MessageTypeMixed      MessageType = "mixed"
	MessageTypeEvent      MessageType = "event"
)

type Server struct {
	log             Logger
	wxcpt           *wxbizmsgcrypt.WXBizMsgCrypt
	messageHandlers map[MessageType]MessageHandler
	client          *Client
	token           string
	encodingAeskey  string
	robotName       string
}

func NewServer(client *Client, token, encodingAeskey, robotName string) *Server {
	return &Server{
		token:           token,
		encodingAeskey:  encodingAeskey,
		wxcpt:           wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, "", wxbizmsgcrypt.XmlType),
		log:             createDefaultLogger(),
		client:          client,
		messageHandlers: map[MessageType]MessageHandler{},
	}
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
	Callback       func(body []byte, client *Client) error
	MessageHandler func(msg *CallbackMessage, bot *Client) error
)

func (s *Server) RegisterCallbackHandlerFunc(t MessageType, fn MessageHandler) *Server {
	s.messageHandlers[t] = fn
	return s
}

func (s *Server) Run(addr, path string) error {
	r := gin.Default()
	s.AddGinRoute(r, path)
	return r.Run(addr)
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

func (s *Server) AddGinRoute(e *gin.Engine, path string) {
	e.GET(path, func(c *gin.Context) {
		msgSignature := c.Query("msg_signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echoStr := c.Query("echostr")
		result, err := s.verifyURL(msgSignature, timestamp, nonce, echoStr)
		if err != nil {
			s.log.Errorf("failed to verify url: %s", err.Error())
			c.String(http.StatusBadRequest, "failed to verify url: %s", err.Error())
			return
		}
		c.String(http.StatusOK, string(result))
	})

	e.POST(path, func(c *gin.Context) {
		msgSignature := c.Query("msg_signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			s.log.Errorf("failed to ready body: %s", err.Error())
			c.String(http.StatusBadRequest, "failed to read body: %s", err.Error())
			return
		}
		s.log.Debugf("body received: %s", string(body))
		msg, cryptErr := s.DecryptJsonMsg(msgSignature, timestamp, nonce, body)
		if cryptErr != nil {
			s.log.Errorf("failed to decrypt msg: %v", cryptErr)
			c.String(http.StatusBadRequest, "failed to decrypt msg: %v", cryptErr)
			return
		}

		c.Status(http.StatusOK)

		handle, ok := s.messageHandlers[MessageType(msg.MsgType)]

		if !ok || handle == nil {
			s.log.Debugf("ignore MsgType %q which is not registered", msg.MsgType)
			return
		}
		err = handle(msg, s.client)
		if err != nil {
			s.log.Errorf("failed to handle msg: %s", err.Error())
		}
	})
}
