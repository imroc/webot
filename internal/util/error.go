package util

import (
	"fmt"

	"github.com/imroc/webot/internal/wxbizjsonmsgcrypt"
)

func ConvertCryptError(e *wxbizjsonmsgcrypt.CryptError) error {
	return fmt.Errorf("crypt error(%d): %s", e.ErrCode, e.ErrMsg)
}
