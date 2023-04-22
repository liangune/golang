package eccrypto

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(Md5Hash16([]byte("abc"), true))
	fmt.Println(Md5Hash([]byte("abc"), true))
	fmt.Println(Md5Hash16([]byte("abc"), false))
	fmt.Println(Md5Hash([]byte("abc"), false))
}

func TestAes(t *testing.T) {
	src := []byte("postgres")
	enebc, err := AesEncryptEBC(src, "09542349A30A6AB8")
	fmt.Println(string(enebc), err, Base64Encode(string(enebc)))
	deebc, err := AesDecryptEBC(enebc, "09542349A30A6AB8")
	fmt.Println(string(deebc), err)

	key := "ZAV3141592653510"
	iv := "AHANCKH413566745"

	urlencodebase64encbc := "wMB7tGxuNQ0DYdHS3Xowl14TiVruzr7moTNBDfZjcoX0JkZLM3RiDhjHUxHueGApjDWMHX0VTEJZZBimIJB8qakE9yGKYaHiQJQIIm6JvEiv0VpVTzvqIGo7H8rsMJGFo6lmkPZWUv5wYb61LzknzGBdvlrv1PK2U3rh6MUhT3Xe%2F2WmJ0dGDxuO9oGKLugJO0%2BU003JXAyc36E7wKp8QE4kLMUsruaWC%2BhpcRtkZhrzIhjieS1%2FP9t7IptOtkPcbfOAYk8Xeip25qL3eFDzvPFbYip42zJn9VtKrkm6K3jGSgwDk5OeiyL4m9mJXX0%2FRclXjSe3yTO5FbQEu6IM1jfcjkVa6Y3XUgEPv7GuXc%2BI%2Frl77AZCjlXGx77UkdAI%2FdalAGcYPGmwS820Yjq%2F3u6Av6B9gLcAi%2FY%2FhPRVVM0%3D"
	base64encbc, _ := UrlDecode(urlencodebase64encbc)
	fmt.Println(base64encbc)
	encbc, err := Base64Decode(base64encbc)
	decbc, err := AesDecryptCBC(encbc, key, iv)
	fmt.Println(string(decbc), err)

	encbc1, err := AesEncryptCBC(decbc, key, iv)
	fmt.Println(encbc1, err)
	base64encbc1 := Base64Encode(string(encbc1))
	fmt.Println(base64encbc1)
	urlencodebase64encbc1 := UrlEncode(base64encbc1)
	fmt.Println(urlencodebase64encbc1)
}

func TestHash(t *testing.T) {
	src := []byte("crypto_test")
	key := "key"
	fmt.Println(Hmac(src, key))

	fmt.Println(Sha1(src))
	fmt.Println(Sha256(src))
	fmt.Println(HmacSha1(src, key))
	fmt.Println(HmacSha256(src, key))
}
