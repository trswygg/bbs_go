package controllers

import (
	"bbs/models"
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"time"
)

// ClassController Operations about Users
type ClassController struct {
	web.Controller
}
type ClassControllerReply struct {
	Time      int64
	ReplyCode int
	Result    string
	Count     int
	Data      interface{}
}

func (c *ClassController) URLMapping() {
	logs.Trace("[URLMapping] mapping ClassController")
	c.Mapping("class", c.GetClassInfo)
	c.Mapping("class/list", c.QueryClass)
	c.Mapping("class/update", c.UpdateClass)
	c.Mapping("class/create", c.CreateClass)
	c.Mapping("class/by_user", c.QueryClassByUser)
	c.Mapping("class/subscribe", c.SubscribeClass)
}

func (c *ClassController) sendReply(code int, result string, count int, data interface{}) {
	c.Data["json"] = ClassControllerReply{
		Time:      time.Now().Unix(),
		ReplyCode: code,
		Result:    result,
		Count:     count,
		Data:      data,
	}
	err := c.ServeJSON()
	if err != nil {
		logs.Error("sendReply() err:", err.Error())
	}
	return
}

// GetClassInfo 获取class详细信息
// method get
// router /class:id [get]
func (c *ClassController) GetClassInfo() {
	class := new(models.Class)
	requireId := c.GetString("id")
	// get user
	var rows int64
	util.DB.Find(&class, requireId).Count(&rows)
	if rows == 1 {
		c.sendReply(0, "success", 1, class)
		return
	} else {
		c.sendReply(-1, "no such class id:"+requireId, int(rows), nil)
		return
	}
}

// QueryClass 获取class列表
// method get
// router /class/list [get]
func (c *ClassController) QueryClass() {
	var resultArr []models.Class
	util.DB.Model(&models.Class{}).Find(&resultArr)
	c.sendReply(0, "success", len(resultArr), resultArr)
	return
}

// QueryClassByUser 获取class列表
// method get
// router /class/by_user
func (c *ClassController) QueryClassByUser() {
	var classes []models.Class
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	err := util.DB.Model(&models.User{Id: uid}).Association("Classes").Find(&classes)
	if err != nil {
		c.sendReply(-1, "DB error:"+err.Error(), 0, nil)
		return
	}
	c.sendReply(0, "success", len(classes), classes)
	return
}

// SubscribeClass 订阅class
// method get
// router /class/subscribe :id
func (c *ClassController) SubscribeClass() {
	id, err := c.GetInt("id")
	class := models.Class{}
	if err != nil {
		c.sendReply(-1, "wrong parame id:"+err.Error(), 0, nil)
	}
	class.Id = id
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	err = util.DB.Model(&models.User{Id: uid}).Association("Classes").Append(&class)
	if err != nil {
		c.sendReply(-1, "DB error"+err.Error(), 0, nil)
		return
	}
	c.sendReply(0, "succerss", 1, class)
	return

}

// UpdateClass 更新class信息
// method post
// router /class [post]
func (c *ClassController) UpdateClass() {
	panic("Unimplemented Method")
}

// CreateClass 创建class
// method post
func (c *ClassController) CreateClass() {
	panic("Unimplemented Method")
}
