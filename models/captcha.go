package models

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
)

// Captcha 图形验证码
type Captcha struct {
	ID     string
	B64s   string
	Result bool
}

var store = base64Captcha.DefaultMemStore

//  获取验证码
func NewCaptcha() Captcha {
	// 生成默认数字
	driver := base64Captcha.DefaultDriverDigit
	// 生成base64图片
	s := base64Captcha.NewCaptcha(driver, store)

	// 获取
	id, b64s, err := s.Generate()
	if err != nil {
		fmt.Println("Register GetCaptchaPhoto get base64Captcha has err:", err)
		return Captcha{ID: "", B64s: ""}
	}
	return Captcha{ID: id, B64s: b64s}
}

// 验证验证码
func (c *Captcha) Verify() bool {
	if c.ID == "" || c.B64s == "" {
		return false
	}
	// 同时在内存清理掉这个图片
	return store.Verify(c.ID, c.B64s, true)
}
