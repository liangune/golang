package eccrypto

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

func Hmac(src []byte, key string) string {
	hmac := hmac.New(md5.New, []byte(key))
	hmac.Write(src)
	return hex.EncodeToString(hmac.Sum([]byte("")))
}
