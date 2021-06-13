package controllers

import (
	"bbs/models"
	"bbs/util"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"time"
)

type TribeController struct {
	web.Controller
}
type TribeControllerReply struct {
	Time      int64
	ReplyCode int
	Result    string
	Count     int
	Data      interface{}
}

func (c *TribeController) URLMapping() {
	logs.Trace("[URLMapping] mapping PostController")
	c.Mapping("tribe", c.GetTribeInfo)
	c.Mapping("tribe/list", c.QueryTribes)
	c.Mapping("tribe/by_user", c.GetTribesByUser)
	c.Mapping("tribe/subscribe", c.SubscribeTribe)
	c.Mapping("tribe/update", c.UpdateTribe)
	c.Mapping("tribe/create", c.CreateTribe)
	c.Mapping("tribe/delete", c.DeleteTribe)

}

func (c *TribeController) sendReply(code int, result string, count int, data interface{}) {
	c.Data["json"] = TribeControllerReply{
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

// GetTribeInfo 获取TribeInfo
// param id
// router /tribe :id
func (c *TribeController) GetTribeInfo() {

}

// QueryTribes 批量查询
// param class_id
// router /tribe/list
func (c *TribeController) QueryTribes() {
	id, err := c.GetInt("id")
	if err != nil {
		c.sendReply(-1, "err in parse int :"+err.Error(), 0, nil)
		return
	}
	var PostsCount int64
	// id,created_at,name,disc,creator_id,class_id
	var tribes []models.TribeQry
	if db := util.DB.Model(&models.Tribe{}).Select("id,created_at,name,disc,creator_id,class_id").Where("class_id = ?", id).Find(&tribes); db.Error != nil {
		c.sendReply(-1, "err on db :"+db.Error.Error(), 0, nil)
		return
	}
	for i := range tribes {
		PostsCount = 0
		util.DB.Model(&models.Post{}).Where("tribe_id = ?", tribes[i].Id).Count(&PostsCount)
		tribes[i].PostsCount = int(PostsCount)
	}
	util.DB.Model(&models.Post{}).Where("tribe_id = ?")
	c.sendReply(0, "success", len(tribes), tribes)
	return
}

// GetTribesByUser 获取用户关注的Tribe
// param id
// router /tribe :id
func (c *TribeController) GetTribesByUser() {

}

// SubscribeTribe 订阅Tribe
// method get
// router /tribe/subscribe :id
func (c *TribeController) SubscribeTribe() {
	id, err := c.GetInt("id")
	tribe := models.Tribe{}
	if err != nil {
		c.sendReply(-1, "wrong parame id:"+err.Error(), 0, nil)
	}
	tribe.Id = id
	uid, b := models.GetUserIdByKey(c.GetString("AccessKey"))
	if !b {
		c.sendReply(-1, "you should login first", 0, nil)
		return
	}
	err = util.DB.Model(&models.User{Id: uid}).Association("Tribes").Append(&tribe)
	if err != nil {
		c.sendReply(-1, "DB error"+err.Error(), 0, nil)
		return
	}
	c.sendReply(0, "succerss", 1, tribe)
	return

}

// CUID

// UpdateTribe 获取TribeInfo
// param id
// router /tribe :id
func (c *TribeController) UpdateTribe() {

}

// CreateTribe 创建Tribe
// param id
// router /tribe :id
func (c *TribeController) CreateTribe() {

}

// DeleteTribe 删除Tribe
// param id
// router /tribe :id
func (c *TribeController) DeleteTribe() {

}
