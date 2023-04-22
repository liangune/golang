package eccrypto

import (
	"encoding/base64"
	"net/url"
)

func Base64Decode(in string) ([]byte, error) {
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return out, err
	}

	return out, nil
}

func Base64Encode(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func UrlEncode(in string) string {
	return url.QueryEscape(in)
}

func UrlDecode(in string) (string, error) {
	return url.QueryUnescape(in)
}
