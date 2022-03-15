package webot

import (
	"encoding/xml"
	"github.com/imroc/webot/internal/tests"
	"testing"
)

func TestText(t *testing.T) {
	var msg CallbackMessage
	data := tests.GetTestFileContent(t, "msg-text.xml")
	err := xml.Unmarshal(data, &msg)
	tests.AssertNoError(t, err)
}
