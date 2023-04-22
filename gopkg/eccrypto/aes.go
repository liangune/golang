package eccrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

//PKCS7Padding support 16, 24, 32 blockSize padding
//AES-128, AES-192, or AES-256
func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length <= 0 {
		return nil, errors.New("decrypted data src length is 0")
	}
	unpadding := int(src[length-1])
	return src[:(length - unpadding)], nil
}

func AesEncryptCBC(src []byte, key string, iv string) ([]byte, error) {
	// key的长度必须为16, 24或者32
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// get key block size
	blockSize := block.BlockSize()
	// padding
	src = PKCS7Padding(src, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(src))
	// 加密
	blockMode.CryptBlocks(encrypted, src)
	return encrypted, nil
}

func AesDecryptCBC(encrypted []byte, key string, iv string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// get key block size
	//blockSize := block.BlockSize()
	// 解密模式
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	decrypted := make([]byte, len(encrypted))
	// 解密
	blockMode.CryptBlocks(decrypted, encrypted)
	// 去掉填充
	decrypted, err = PKCS7UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

func AesEncryptEBC(src []byte, key string) ([]byte, error) {
	//key的长度必须为16 24 32
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	//padding
	src = PKCS7Padding(src, blockSize)

	encrypted := make([]byte, len(src))
	tmp := make([]byte, blockSize)

	//分组分块加密
	for index := 0; index < len(src); index += blockSize {
		block.Encrypt(tmp, src[index:index+blockSize])
		copy(encrypted, tmp)
	}
	return encrypted, nil
}

func AesDecryptEBC(src []byte, key string) ([]byte, error) {
	//key的长度必须为16 24 32
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	decrypted := make([]byte, len(src))
	tmp := make([]byte, blockSize)

	//分组分块解密
	for index := 0; index < len(src); index += blockSize {
		block.Decrypt(tmp, src[index:index+blockSize])
		copy(decrypted, tmp)
	}
	decrypted, err = PKCS7UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
