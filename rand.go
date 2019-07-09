package gohack

import (
	"math/rand"
	"time"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// NewRand 实例化随机数结构体，源为时间微秒
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandString 产生指定长度随机字符串, letters 默认为纯数字 "0123456789"
func RandString(letters string, n int) string {
	if letters == "" {
		letters = "0123456789"
	}
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// Shuffle 打乱 []int 的元素顺序
func Shuffle(vals []int) *[]int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return &ret
}
