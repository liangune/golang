package eccrypto

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os/exec"
	"strings"
)

func Md5String(src []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(src)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func FileMd5(path string) string {
	data, err := ioutil.ReadFile(path)
	if nil != err {
		return ""
	}
	return Md5String(data)
}

func FileMd5Shell(path string) string {
	cmdstr := "md5sum " + path
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return ""
	}

	s := out.String()
	if len(s) <= 32 {
		return ""
	}

	res := s[0:32]

	return res
}

func Md5Hash16(src []byte, isLower bool) string {
	md5Ctx := md5.New()
	md5Ctx.Write(src)
	cipherStr := md5Ctx.Sum([]byte(""))

	md5Str := hex.EncodeToString(cipherStr[4:12])
	if !isLower {
		md5Str = strings.ToUpper(md5Str)
	}

	return md5Str
}

func Md5Hash(src []byte, isLower bool) string {
	md5Ctx := md5.New()
	md5Ctx.Write(src)
	cipherStr := md5Ctx.Sum([]byte(""))

	md5Str := hex.EncodeToString(cipherStr)
	if !isLower {
		md5Str = strings.ToUpper(md5Str)
	}

	return md5Str
}
