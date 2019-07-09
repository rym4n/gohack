package gohack

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash/crc32"
)

// Crc32 计算字符串的crc值
func Crc32(in string) string {
	return fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(in)))
}

// Md5 计算字符串的MD5值
func Md5(str string) string {
	strByte := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(strByte))
}

/*
HmacSha1 计算字符串的HmacSha1值
@key: 计算hmac的key
@str: 原始字符串
*/
func HmacSha1(key, str string) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(str))
	return mac.Sum(nil)
}

// Sha256 计算字符串的sha256哈希值
func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
