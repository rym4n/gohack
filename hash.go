package gohack

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash/crc32"
)

func Crc32(in string) string {
	return fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(in)))
}

func Md5(str string) string {
	strByte := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(strByte))
}

func HmacSha1(key, str string) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(str))
	return mac.Sum(nil)
}

func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
