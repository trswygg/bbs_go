package controllers

import (
	"bbs/models"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"math/rand"
	"time"
)

func init() {
	//init rand seed
	rand.Seed(time.Now().UnixNano())
}

// CaptchaController Captcha 图形验证码
type CaptchaController struct {
	web.Controller
}
type CaptchaControllerReply struct {
	Time      int64       `json:"timestamp"`
	ReplyCode int         `json:"code"`
	Result    string      `json:"result"`
	Data      interface{} `json:"data"`
}

func (c *CaptchaController) URLMapping() {
	logs.Trace("[URLMapping] mapping CaptchaController")
	c.Mapping("captcha", c.Get)
	c.Mapping("captcha", c.Post)
}

// Get
// Title GetCaptcha
// Description 获取图形验证码
// Success 200 {object} models.Captcha
// router /captcha/ [get]
func (c *CaptchaController) Get() {
	Captcha := models.NewCaptcha()

	c.Data["json"] = CaptchaControllerReply{
		Time:      time.Now().Unix(),
		ReplyCode: 0,
		Result:    "success",
		Data:      Captcha,
	}
	_ = c.ServeJSON()
}

// Post
// Title VerifyCaptcha
// Description 校验图形验证码
// Param body body models.Captcha true "验证码"
// Success 200 models.Captcha
// router /captcha/ [post]
func (c *CaptchaController) Post() {
	var captcha models.Captcha
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &captcha)
	if captcha.Verify() {
		//c.Data["json"] = models.Captcha{ID: captcha.ID, Result: true}
		c.Data["json"] = CaptchaControllerReply{
			Time:      time.Now().Unix(),
			ReplyCode: 0,
			Result:    "success",
			Data:      models.Captcha{ID: captcha.ID, Result: true},
		}
	} else {
		//c.Data["json"] = models.Captcha{ID: captcha.ID, Result: false}
		c.Data["json"] = CaptchaControllerReply{
			Time:      time.Now().Unix(),
			ReplyCode: 0,
			Result:    "fail",
			Data:      models.Captcha{ID: captcha.ID, Result: false},
		}
	}
	_ = c.ServeJSON()
}
