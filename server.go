package webot

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/imroc/webot/internal/wxbizmsgcrypt"
	"io/ioutil"
	"net/http"
	"strings"
)

type CallbackType string

const (
	CallbackTypeText       CallbackType = "text"
	CallbackTypeAttachment CallbackType = "attachment"
	CallbackTypeImage      CallbackType = "image"
	CallbackTypeMixed      CallbackType = "mixed"
	CallbackTypeEvent      CallbackType = "event"
)

type CallbackServer struct {
	token          string
	encodingAeskey string
	robotName      string
	wxcpt          *wxbizmsgcrypt.WXBizMsgCrypt
	log            Logger
	callbackFuncs  map[CallbackType]CallbackHandlerFunc
	bot            *WeBot
}

func (b *WeBot) NewCallbackServer(token, encodingAeskey, robotName string) *CallbackServer {
	return &CallbackServer{
		token:          token,
		encodingAeskey: encodingAeskey,
		wxcpt:          wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, "", wxbizmsgcrypt.XmlType),
		log:            createDefaultLogger(),
		bot:            b,
		callbackFuncs:  map[CallbackType]CallbackHandlerFunc{},
	}
}

func (s *CallbackServer) VerifyURL(msgSignature, timestamp, nonce, echoStr string) ([]byte, error) {
	result, err := s.wxcpt.VerifyURL(msgSignature, timestamp, nonce, echoStr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CallbackServer) DecryptJsonMsg(msgSignature, timestamp, nonce string, data []byte) (msg *CallbackMessage, err error) {
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

type CallbackHandlerFunc func(msg *CallbackMessage, bot *WeBot) error

func (s *CallbackServer) RegisterCallbackHandlerFunc(t CallbackType, fn CallbackHandlerFunc) *CallbackServer {
	s.callbackFuncs[t] = fn
	return s
}

func (s *CallbackServer) Run(addr, path string) error {
	r := gin.Default()
	s.AddGinRoute(r, path)
	return r.Run(addr)
}

func (s *CallbackServer) cleanContent(content string) string {
	str := strings.Replace(content, "@"+s.robotName, "", -1)
	return strings.TrimSpace(str)
}

func (s *CallbackServer) AddGinRoute(e *gin.Engine, path string) {
	e.GET(path, func(c *gin.Context) {
		msgSignature := c.Query("msg_signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echoStr := c.Query("echostr")
		result, err := s.VerifyURL(msgSignature, timestamp, nonce, echoStr)
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
		body, err := ioutil.ReadAll(c.Request.Body)
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

		handle, ok := s.callbackFuncs[CallbackType(msg.MsgType)]

		if !ok || handle == nil {
			s.log.Debugf("ignore MsgType %q which is not registered", msg.MsgType)
			return
		}
		err = handle(msg, s.bot)
		if err != nil {
			s.log.Errorf("failed to handle msg: %s", err.Error())
		}
	})
}
