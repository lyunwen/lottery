package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strconv"
)

//获取唯一ID
//获取唯一ID
func GetUUID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
	return hex.EncodeToString(h.Sum(nil))
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func StringToInt(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	return intValue, err
}
