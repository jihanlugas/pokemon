package utils

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
	"unsafe"
)


var regFormatHp *regexp.Regexp


const letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterLen     = 36                   // len(letterBytes)
)

var src = rand.NewSource(time.Now().UnixNano())


func init() {
	regFormatHp = regexp.MustCompile(`(^\+?628)|(^0?8){1}`)
}

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < letterLen {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func Rand4DigitInt() string {
	return strconv.FormatInt((rand.Int63n(8999) + 1000), 10)
}

func Rand6DigitInt() string {
	return strconv.FormatInt((rand.Int63n(899999) + 100000), 10)
}

func FormatPhoneTo62(phone string) string {
	formatPhone := regFormatHp.ReplaceAllString(phone, "628")
	return formatPhone
}
