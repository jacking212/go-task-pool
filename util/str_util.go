package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

const (
	digit           = "0123456789"
	lowerCaseLetter = "abcdefghijklmnopqrstuvwxyz"
	upperCaseLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

//生成随机字符串
func GenRandomString(length int) string {

	return genStrByLenAndBaseStr(digit+lowerCaseLetter+upperCaseLetter, length)
}

// 生成小写字母与数字随机串
func GenRandomDigitLowerLetter(length int) string {
	return genStrByLenAndBaseStr(digit+lowerCaseLetter, length)
}

// 根据传入的基础字符串 base，返回指定长度 length 的随机串
func genStrByLenAndBaseStr(base string, length int) string {

	bytes := []byte(base)
	bytesLen := len(bytes)
	retVal := make([]byte, 0, length)

	randomGen := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	for i := 0; i < length; i++ {
		retVal = append(retVal, bytes[randomGen.Intn(bytesLen)])
	}
	return string(retVal)
}

func MD5(data []byte) string {
	tmp := md5.Sum(data)
	return hex.EncodeToString(tmp[:])
}
