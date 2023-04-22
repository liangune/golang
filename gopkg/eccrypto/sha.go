package eccrypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func Sha1(src []byte) string {
	sha1 := sha1.New()
	sha1.Write(src)
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func HmacSha1(src []byte, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write(src)
	return hex.EncodeToString(mac.Sum(nil))
}

func Sha256(src []byte) string {
	sha256 := sha256.New()
	sha256.Write(src)
	return hex.EncodeToString(sha256.Sum([]byte("")))
}

func HmacSha256(src []byte, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(src)
	return hex.EncodeToString(mac.Sum(nil))
}
