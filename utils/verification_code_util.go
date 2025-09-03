package utils

import (
	"math/rand"
	"time"
)

var CaptchaStore = make(map[string]string)

// 生成包含字母和数字的6位验证码
func GenerateCaptcha() string {
	// 字符集：包含大写字母、小写字母和数字
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// 使用新的随机数生成方式，替代已弃用的rand.Seed()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成6位验证码
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
