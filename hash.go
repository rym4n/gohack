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
func Crc32(str string) (hexValue string) {
	hexValue = fmt.Sprintf("%x", crc32.ChecksumIEEE([]byte(str)))
	return hexValue
}

// Md5 计算字符串的MD5值
func Md5(str string) (hexValue string) {
	strByte := []byte(str)
	hexValue = fmt.Sprintf("%x", md5.Sum(strByte))
	return hexValue
}

/*
HmacSha1 计算字符串的HmacSha1值
@key: 计算hmac的key
@str: 原始字符串
*/
func HmacSha1(key, str string) (hexValue string) {
	mac := hmac.New(sha1.New, []byte(key))
	hexValue = fmt.Sprintf("%x", mac.Sum([]byte(str)))
	return hexValue
}

// Sha256 计算字符串的sha256哈希值
func Sha256(str string) (hexValue string) {
	hexValue = fmt.Sprintf("%x", sha256.New().Sum([]byte(str)))
	return hexValue
}
